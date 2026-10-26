package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/crypto"

	cronErrors "github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// execTypeGroup 选机方式：按节点分组随机（与执行类型 exec_type 不同维度）
const execTypeGroup = "group"

// Config 调度器配置（app 层组装注入）
type Config struct {
	TaskRepo   repository.Task       // 任务仓储（查询/校验）
	WorkerRepo repository.Worker     // 节点仓储（选机）
	LogRepo    repository.Log        // 执行日志（快照 + 流式回传）
	ExecRepo   repository.Execution  // 执行意图（状态机载体）
	TaskClient repository.TaskClient // 执行器（gRPC 下发）
	Lock       Lock                  // 执行标记锁
	Trigger    Trigger               // 触发源（robfig/cron/v3）
	Logger     *log.Helper           // 日志
	InstanceID string                // 执行体标识（实例唯一）

	LockTTL   time.Duration // 执行锁 TTL（默认 30s）
	Heartbeat time.Duration // 心跳间隔（默认 10s）
	MaxRetry  int           // 连接重试次数（默认 5）
	RetryWait int           // 每次重试等待基准秒（默认 120）
}

// Scheduler 调度内核：触发源 + 执行意图状态机 + 执行标记锁 + 执行器
// 单实例部署即保证：触发不丢、执行不重、崩溃可恢复；多实例由 Leader 保证只有 leader 运行调度器。
type Scheduler struct {
	cfg *Config

	mu       sync.Mutex
	started  bool
	ctx      context.Context   // Start 时创建，Stop 时取消（watchdog/reconcile/心跳跟随）
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	stopping atomic.Bool

	// misfire 补帧频率保护：执行时长 > 周期时避免"执行完→补帧→执行中又错过→再补"的级联风暴
	refillMu   sync.Mutex
	lastRefill map[uint32]time.Time
}

var (
	defaultInst *Scheduler
	defaultOnce sync.Once
)

// New 创建调度器（不启动；由 Start/Stop 管理生命周期，支持 Leader 切换反复启停）
func New(cfg *Config) *Scheduler {
	if cfg.LockTTL <= 0 {
		cfg.LockTTL = 30 * time.Second
	}
	if cfg.Heartbeat <= 0 {
		cfg.Heartbeat = 10 * time.Second
	}
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 5
	}
	if cfg.RetryWait <= 0 {
		cfg.RetryWait = 120
	}
	if cfg.Logger == nil {
		cfg.Logger = log.NewHelper(log.DefaultLogger)
	}
	return &Scheduler{
		cfg:        cfg,
		lastRefill: make(map[uint32]time.Time),
	}
}

// Init 注册全局调度器实例（幂等；**不启动**——由 Leader 选举决定何时 Start，
// 保证多实例下只有 leader 运行调度器；单实例场景下 Leader 必然胜出，行为等价于启动）
func Init(cfg *Config) *Scheduler {
	defaultOnce.Do(func() {
		defaultInst = New(cfg)
	})
	return defaultInst
}

// Stop 停止全局调度器（幂等；main.go BeforeStop 调用）
func Stop() {
	if defaultInst != nil {
		defaultInst.Stop()
	}
}

// Start 启动调度器（幂等；支持 Leader 切换后再次启动）：
// 启动触发循环 → 以 DB 为权威加载任务 → 崩溃恢复 → 启动 watchdog/reconcile 循环
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.started = true
	s.stopping.Store(false)
	s.mu.Unlock()

	// 1. 启动触发循环
	s.cfg.Trigger.Start()

	// 2. 加载启用任务到触发源
	ctx := kratosx.MustContext(context.Background())
	specs := s.cfg.TaskRepo.AllTaskSpecs(ctx)
	s.cfg.Trigger.Sync(specs, s.wrap)

	// 3. 崩溃恢复（pending 重投 + running 心跳过期重投）
	s.recover()

	// 4. 后台循环：孤儿检测 + 对账 + 节点探活（随 s.ctx 退出）
	go s.watchdogLoop()
	go s.reconcileLoop()
	go s.workerProbeLoop()

	s.cfg.Logger.Infow("msg", "scheduler started", "tasks", len(specs), "instance", s.cfg.InstanceID)
}

// Stop 优雅停止（幂等）：停触发循环 → 取消后台循环/心跳 → 等在办执行（带超时兜底）
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	s.started = false
	s.stopping.Store(true)
	s.mu.Unlock()

	s.cfg.Trigger.Stop()
	s.cancel()
	done := make(chan struct{})
	go func() { s.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(60 * time.Second):
		s.cfg.Logger.Warnw("msg", "scheduler stop timeout, force exit")
	}
}

// InstanceID 返回当前实例标识（Leader 选举用）
func (s *Scheduler) InstanceID() string {
	return s.cfg.InstanceID
}

// Trigger 手动触发（force：跳过状态/表达式校验，与定时执行共用锁互斥）
func (s *Scheduler) Trigger(ctx kratosx.Context, taskID uint32) error {
	s.dispatch(ctx, taskID, "", true, entity.TriggerManual)
	return nil
}

// Add 任务启用/创建后注册到触发源（幂等）
func (s *Scheduler) Add(taskID uint32, spec string) error {
	return s.cfg.Trigger.Add(taskID, spec, s.wrap(taskID, spec))
}

// Remove 任务禁用/删除后从触发源移除
func (s *Scheduler) Remove(taskID uint32) {
	s.cfg.Trigger.Remove(taskID)
}

// Update 任务更新后同步触发源（P0-3 修复：消除更新后静默停摆）
func (s *Scheduler) Update(taskID uint32, spec string) error {
	return s.cfg.Trigger.Update(taskID, spec, s.wrap(taskID, spec))
}

// ValidateSpec 校验 cron 表达式（秒级）
func (s *Scheduler) ValidateSpec(spec string) error {
	return s.cfg.Trigger.ValidateSpec(spec)
}

// CancelExec 取消执行：从日志快照定位节点后下发取消指令
func (s *Scheduler) CancelExec(ctx kratosx.Context, uuid string) error {
	lg, err := s.cfg.LogRepo.GetLogByUuid(ctx, uuid)
	if err != nil {
		return err
	}
	worker := lg.GetWorker()
	return s.cfg.TaskClient.CancelExec(ctx, &types.CancelExecRequest{
		Uuid: uuid,
		IP:   worker.Ip,
		Ak:   worker.Ak,
		Sk:   worker.Sk,
	})
}

// wrap 生成触发回调：到点异步调度，不阻塞 cron 调度循环
func (s *Scheduler) wrap(taskID uint32, spec string) func() {
	return func() {
		if s.stopping.Load() {
			return
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.dispatch(kratosx.MustContext(context.Background()), taskID, spec, false, entity.TriggerCron)
		}()
	}
}

// dispatch 产生一次执行意图并异步执行（幂等：同 (task_id, trigger_at) 只生效一次）
func (s *Scheduler) dispatch(ctx kratosx.Context, taskID uint32, spec string, force bool, triggerType string) {
	// 1. 查任务
	task, err := s.cfg.TaskRepo.GetTask(ctx, taskID)
	if err != nil {
		s.cfg.Logger.Warnw("msg", "dispatch: get task fail", "id", taskID, "err", err.Error())
		return
	}
	// 2. 校验（非 force）：状态启用 + 表达式与当前一致（防串任务）
	if !force {
		if task.Status == nil || !*task.Status {
			s.cfg.Logger.Infow("msg", "dispatch: task disabled", "id", taskID)
			return
		}
		if task.Spec != spec {
			s.cfg.Logger.Infow("msg", "dispatch: task spec changed", "id", taskID)
			return
		}
	}

	// 3. 生成执行意图（唯一键 (task_id, trigger_at) 幂等，防同帧重复触发）
	ex := &entity.Execution{
		App:         task.App,
		Id:          uuid.New().String(),
		TaskId:      taskID,
		TriggerAt:   time.Now().Unix(),
		TriggerType: triggerType,
		Status:      entity.ExecPending,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	if err := s.cfg.ExecRepo.CreateExecution(ctx, ex); err != nil {
		if errors.Is(err, repository.ErrExecutionExists) {
			s.cfg.Logger.Infow("msg", "dispatch: execution exists, skip", "task", taskID, "triggerAt", ex.TriggerAt)
			return
		}
		s.cfg.Logger.Errorw("msg", "dispatch: create execution fail", "task", taskID, "err", err.Error())
		return
	}
	executionTotal.WithLabelValues(taskIDLabel(taskID), triggerType, "created").Inc()

	// 4. 异步执行（Clone 脱离请求生命周期：手动触发场景请求结束会取消 ctx）
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.execute(ctx.Clone(), ex)
	}()
}

// execute 一次执行的完整流程：选机 → 抢执行标记锁 → 认领+心跳 → 下发 → 终态
func (s *Scheduler) execute(ctx kratosx.Context, ex *entity.Execution) {
	start := time.Now()
	// 1. 查询任务与选机
	task, err := s.cfg.TaskRepo.GetTask(ctx, ex.TaskId)
	if err != nil {
		s.finishFail(ex.Id, "get task fail: "+err.Error())
		s.recordExecution(ex, entity.ExecFail, start)
		return
	}
	wk, err := s.pickWorker(ctx, task)
	if err != nil {
		s.finishFail(ex.Id, "pick worker fail: "+err.Error())
		s.recordExecution(ex, entity.ExecFail, start)
		return
	}

	// 2. 抢执行标记锁（锁被其他执行体持有 = 该任务正在执行 → 本帧按 misfire 策略处理）
	ok, err := s.cfg.Lock.Acquire(ctx, task.Id, ex.Id, s.cfg.LockTTL)
	if err != nil {
		s.finishFail(ex.Id, "acquire lock fail: "+err.Error())
		return
	}
	if !ok {
		s.handleMisfire(ctx, task, ex)
		return
	}
	// 释放锁 + 补帧（任何退出路径执行；补帧需在锁释放之后，避免自锁）
	defer func() {
		s.releaseLockSafe(ex.Id, task.Id)
		s.refillMisfire(ctx, task, ex)
	}()

	// 3. 认领执行 + 启动心跳
	executorID := fmt.Sprintf("%s:%s", s.cfg.InstanceID, ex.Id[:8])
	if err := s.cfg.ExecRepo.MarkRunning(ctx, ex.Id, &wk.Id, executorID); err != nil {
		s.finishFail(ex.Id, "mark running fail: "+err.Error())
		return
	}
	stopHeartbeat := s.startHeartbeat(ex.Id)
	defer stopHeartbeat()

	// 4. 生成执行日志（快照 + 首条），流式回传
	logUUID := crypto.MD5ToUpper([]byte(uuid.NewString()))
	if _, err := s.cfg.LogRepo.CreateLog(ctx, s.makeLog(logUUID, task, wk)); err != nil {
		s.finishFail(ex.Id, "create task log fail: "+err.Error())
		return
	}
	// 回填执行历史与日志的关联（log_uuid），点开执行历史可查看该次日志
	if err := s.cfg.ExecRepo.AttachLog(ctx, ex.Id, logUUID); err != nil {
		s.cfg.Logger.Warnw("msg", "attach log to execution fail", "execution", ex.Id, "err", err.Error())
	}
	recvFunc := func(msg *types.ExecTaskLog) {
		b, _ := json.Marshal(msg)
		if err := s.cfg.LogRepo.AppendLogContent(ctx, logUUID, string(b)); err != nil {
			s.cfg.Logger.Errorw("msg", "append log error", "err", err.Error())
		}
	}

	// 5. 下发执行（P1 重试语义：连接错误按退避重试；业务失败立即终止）
	status, lastErr := s.execWithRetry(ctx, task, wk, ex, logUUID, recvFunc)

	// 6. 写终态
	if err := s.cfg.ExecRepo.FinishExecution(ctx, ex.Id, status, time.Now().Unix(), lastErr); err != nil {
		s.cfg.Logger.Errorw("msg", "finish execution fail", "id", ex.Id, "err", err.Error())
	}
	if err := s.cfg.LogRepo.UpdateLogStatus(ctx, logUUID, status); err != nil {
		s.cfg.Logger.Warnw("msg", "update log status error", "err", err.Error())
	}
	s.recordExecution(ex, status, start)
}

// recordExecution 执行意图指标打点（终态 + 耗时）
func (s *Scheduler) recordExecution(ex *entity.Execution, status string, start time.Time) {
	executionTotal.WithLabelValues(taskIDLabel(ex.TaskId), ex.TriggerType, status).Inc()
	executionDuration.WithLabelValues(taskIDLabel(ex.TaskId), status).Observe(time.Since(start).Seconds())
}

// execWithRetry 执行任务并处理连接重试；返回 (终态, 错误描述)
// failover：group 路由模式下重连会重新选机（可能换到组内其他节点），避免重试固定打同一台坏节点
func (s *Scheduler) execWithRetry(ctx kratosx.Context, task *entity.Task, wk *entity.Worker, ex *entity.Execution, logUUID string, recvFunc func(msg *types.ExecTaskLog)) (string, string) {
	// 连接错误重试（业务错误/EOF 立即结束）
	for count := 0; ; count++ {
		curWk := wk
		// 重试且任务支持多机路由 → 重新选机（failover）
		if count > 0 && s.isMultiNode(task) {
			re, err := s.pickWorker(ctx, task)
			if err != nil {
				s.cfg.Logger.Warnw("msg", "failover: pick worker fail", "task", task.Id, "err", err.Error())
			} else {
				curWk = re
				if curWk.Id != wk.Id {
					s.cfg.Logger.Infow("msg", "failover to another worker", "task", task.Id, "worker", curWk.Ip)
				}
			}
		}
		rer := s.cfg.TaskClient.ExecTask(ctx, &types.ExecTaskRequest{
			Id:            task.Id,
			Uuid:          logUUID,
			IP:            curWk.Ip,
			Ak:            curWk.Ak,
			Sk:            curWk.Sk,
			ExecType:      task.ExecType,
			ExecValue:     task.ExecValue,
			ExpectCode:    task.ExpectCode,
			RetryCount:    task.RetryCount,
			RetryWaitTime: task.RetryWaitTime,
			MaxExecTime:   task.MaxExecTime,
		}, recvFunc)
		if rer == nil || errors.Is(rer, io.EOF) || isCanceled(rer) {
			if rer == nil || errors.Is(rer, io.EOF) {
				return entity.ExecSuccess, ""
			}
			// 执行被取消（client 端取消/停机）：不重试，直接终止
			return entity.ExecCancel, "context canceled"
		}

		// 业务错误：立即终止（client 已按 RetryCount 重试过任务体）
		if cronErrors.IsExecTaskFailError(rer) {
			return entity.ExecFail, rer.Error()
		}

		// 连接错误：查当前状态（可能已被取消），未取消则退避重连
		cur, err := s.cfg.ExecRepo.GetExecution(ctx, ex.Id)
		if err != nil || cur == nil || cur.Status != entity.ExecRunning {
			return entity.ExecFail, "task is not running"
		}
		if count < s.cfg.MaxRetry {
			wait := time.Duration(s.cfg.RetryWait/(s.cfg.MaxRetry+1)*(count+1)) * time.Second
			s.cfg.Logger.Warnw("msg", "retry connect worker", "task", task.Id, "count", count+1, "wait", wait.String(), "err", rer.Error())
			select {
			case <-time.After(wait):
			case <-s.ctx.Done():
				return entity.ExecCancel, "scheduler stopped"
			}
			continue
		}
		return entity.ExecFail, "max retry exceeded: " + rer.Error()
	}
}

// isMultiNode 判断任务是否支持多机路由（group 模式；旧数据按 exec_type 回退）
func (s *Scheduler) isMultiNode(task *entity.Task) bool {
	if task.RouteType == entity.RouteGroup {
		return true
	}
	return task.RouteType == "" && task.ExecType == execTypeGroup
}

// isCanceled 判断错误是否为执行取消（gRPC codes.Canceled 或 context.Canceled）
func isCanceled(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return true
	}
	return status.Code(err) == codes.Canceled
}

// pickWorker 选机：按 RouteType 路由（group 随机 / single 指定）；旧数据兼容：route_type 为空时按 exec_type 回退
func (s *Scheduler) pickWorker(ctx kratosx.Context, task *entity.Task) (*entity.Worker, error) {
	route := task.RouteType
	if route == "" {
		// 兼容迁移前数据：exec_type == "group" 表示按组路由
		if task.ExecType == execTypeGroup {
			route = entity.RouteGroup
		} else {
			route = entity.RouteSingle
		}
	}
	if route == entity.RouteGroup {
		if task.WorkerGroupId == nil {
			return nil, errors.New("worker_group_id is nil")
		}
		return s.cfg.WorkerRepo.GetWorkerByGroupId(ctx, task.App, *task.WorkerGroupId)
	}
	if task.WorkerId == nil {
		return nil, errors.New("worker_id is nil")
	}
	return s.cfg.WorkerRepo.GetWorker(ctx, *task.WorkerId)
}

// startHeartbeat 启动心跳续期（后台 goroutine，随 Stop 取消）
func (s *Scheduler) startHeartbeat(exID string) func() {
	stop := make(chan struct{})
	var once sync.Once
	go func() {
		ticker := time.NewTicker(s.cfg.Heartbeat)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := s.cfg.ExecRepo.UpdateHeartbeat(kratosx.MustContext(s.ctx), exID); err != nil {
					s.cfg.Logger.Warnw("msg", "heartbeat fail", "execution", exID, "err", err.Error())
				}
			case <-s.ctx.Done():
				return
			case <-stop:
				return
			}
		}
	}()
	return func() { once.Do(func() { close(stop) }) }
}

// releaseLockSafe 释放执行锁（带超时兜底，防止请求 ctx 取消导致锁残留）
func (s *Scheduler) releaseLockSafe(exID string, taskID uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.cfg.Lock.Release(ctx, taskID, exID); err != nil {
		s.cfg.Logger.Warnw("msg", "release lock fail", "task", taskID, "err", err.Error())
	}
}

// finishFail 执行前置失败统一收口（写终态 + 日志）
func (s *Scheduler) finishFail(exID string, lastErr string) {
	if err := s.cfg.ExecRepo.FinishExecution(kratosx.MustContext(context.Background()), exID, entity.ExecFail, time.Now().Unix(), lastErr); err != nil {
		s.cfg.Logger.Warnw("msg", "finish fail execution", "id", exID, "err", err.Error())
	}
}

// recover 崩溃恢复：pending 重投 + running 心跳过期标孤儿重投（P2 由 Watchdog 覆盖运行时场景）
func (s *Scheduler) recover() {
	ctx := kratosx.MustContext(context.Background())
	now := time.Now().Unix()

	// pending 超时（> 1 分钟）未认领：视为上次进程崩溃遗留，重投
	if pendings, err := s.cfg.ExecRepo.ListPending(ctx, now-int64(time.Minute.Seconds()), 100); err == nil {
		for _, ex := range pendings {
			s.cfg.Logger.Infow("msg", "recover pending execution", "execution", ex.Id, "task", ex.TaskId)
			go s.execute(ctx, ex)
		}
	}

	// running 心跳过期（> 2 个心跳周期）：标孤儿后重投
	heartbeatBefore := now - 2*int64(s.cfg.Heartbeat.Seconds())
	if runnings, err := s.cfg.ExecRepo.ListRunning(ctx, heartbeatBefore, 100); err == nil {
		for _, ex := range runnings {
			if err := s.cfg.ExecRepo.MarkOrphan(ctx, ex.Id); err != nil {
				s.cfg.Logger.Warnw("msg", "mark orphan fail", "execution", ex.Id, "err", err.Error())
				continue
			}
			s.cfg.Logger.Infow("msg", "recover orphan execution", "execution", ex.Id, "task", ex.TaskId)
			go s.execute(ctx, ex)
		}
	}
}

// makeLog 生成一次执行的日志记录（含任务/节点快照），与原逻辑一致
func (s *Scheduler) makeLog(uuid string, task *entity.Task, worker *entity.Worker) *entity.Log {
	tb, _ := json.Marshal(task)
	wb, _ := json.Marshal(worker)
	content := types.ExecTaskLog{
		Type:    "info",
		Content: "start scheduler",
		Time:    time.Now().Unix(),
	}
	cb, _ := json.Marshal(content)
	return &entity.Log{
		App:            task.App,
		Uuid:           uuid,
		TaskId:         task.Id,
		TaskSnapshot:   string(tb),
		WorkerId:       worker.Id,
		WorkerSnapshot: string(wb),
		Content:        string(cb),
		StartAt:        time.Now().Unix(),
		Status:         entity.ExecRunning,
	}
}