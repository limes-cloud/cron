package scheduler

import (
	"context"
	"time"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
)

const (
	watchdogInterval = 30 * time.Second // 孤儿检测周期
	reconcilePeriod  = 5 * time.Minute  // 调度表对账周期
)

// watchdogOnce 孤儿检测：running 且心跳过期（>2 个心跳周期）→ 标 orphan → 重投。
// 由 watchdogLoop 周期调用（仅 leader 实例运行）。
func (s *Scheduler) watchdogOnce() {
	ctx := kratosx.MustContext(context.Background())
	heartbeatBefore := time.Now().Unix() - 2*int64(s.cfg.Heartbeat.Seconds())

	runnings, err := s.cfg.ExecRepo.ListRunning(ctx, heartbeatBefore, 100)
	if err != nil {
		s.cfg.Logger.Warnw("msg", "watchdog: list running fail", "err", err.Error())
		return
	}
	for _, ex := range runnings {
		if err := s.cfg.ExecRepo.MarkOrphan(ctx, ex.Id); err != nil {
			s.cfg.Logger.Warnw("msg", "watchdog: mark orphan fail", "execution", ex.Id, "err", err.Error())
			continue
		}
		s.cfg.Logger.Infow("msg", "watchdog: orphan detected, redispatch", "execution", ex.Id, "task", ex.TaskId)
		// 重投：复用原 execution（trigger_at 幂等键不变），execute 中 MarkRunning 重置状态
		s.wg.Add(1)
		go func(ex *entity.Execution) {
			defer s.wg.Done()
			s.execute(ctx, ex)
		}(ex)
	}
}

// watchdogLoop 孤儿检测循环（随 s.ctx 退出）
func (s *Scheduler) watchdogLoop() {
	ticker := time.NewTicker(watchdogInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.watchdogOnce()
		}
	}
}

// reconcileOnce 以 DB 为权威对账调度表（任务增删改/状态变更的最终一致入口）
func (s *Scheduler) reconcileOnce() {
	ctx := kratosx.MustContext(context.Background())
	specs := s.cfg.TaskRepo.AllTaskSpecs(ctx)
	s.cfg.Trigger.Sync(specs, s.wrap)
	s.cfg.Logger.Infow("msg", "reconcile done", "tasks", len(specs))
}

// reconcileLoop 对账循环（随 s.ctx 退出）
func (s *Scheduler) reconcileLoop() {
	ticker := time.NewTicker(reconcilePeriod)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.reconcileOnce()
		}
	}
}

// handleMisfire 执行权被占用时的错过帧处理：
// 本帧统一标 skip 终态（记录为错过帧）；非 skip 策略由执行完成方的 refillMisfire 补执行最近一帧。
func (s *Scheduler) handleMisfire(ctx kratosx.Context, task *entity.Task, ex *entity.Execution) {
	s.cfg.Logger.Infow("msg", "misfire: task already running, frame skipped", "task", task.Id, "execution", ex.Id)
	misfireTotal.WithLabelValues(taskIDLabel(task.Id), task.Misfire).Inc()
	_ = s.cfg.ExecRepo.FinishExecution(ctx, ex.Id, entity.ExecSkip, time.Now().Unix(), "misfire: task already running")
}

// refillMisfire 执行完成后补偿：若本次执行期间出现过错过帧（skip）且任务策略允许，补执行一帧。
// LATEST 语义：只补最近一帧（丢弃中间帧），防错过堆积雪崩。
// 频率保护：同一任务 1 分钟内最多补一帧——执行时长 > 周期时，
// 补帧执行期间仍会产生新错过帧，若不限频将形成"执行→补帧→执行→补帧"的永续循环。
func (s *Scheduler) refillMisfire(ctx kratosx.Context, task *entity.Task, ex *entity.Execution) {
	if task.Misfire == entity.MisfireSkip {
		return
	}
	cur, err := s.cfg.ExecRepo.GetExecution(ctx, ex.Id)
	if err != nil || cur == nil || cur.StartedAt == nil {
		return
	}
	count, err := s.cfg.ExecRepo.CountSkippedSince(ctx, task.Id, *cur.StartedAt)
	if err != nil || count <= 0 {
		return
	}
	// 频率保护：1 分钟内不重复补帧
	s.refillMu.Lock()
	if last, ok := s.lastRefill[task.Id]; ok && time.Since(last) < time.Minute {
		s.refillMu.Unlock()
		s.cfg.Logger.Infow("msg", "misfire refill skipped (rate limit)", "task", task.Id)
		return
	}
	s.lastRefill[task.Id] = time.Now()
	s.refillMu.Unlock()

	s.cfg.Logger.Infow("msg", "misfire refill after execution", "task", task.Id, "skipped", count)
	// 补一帧：新建执行意图（trigger_at=当前时刻，与历史帧不冲突），走正常 dispatch 互斥
	s.dispatch(ctx, task.Id, task.Spec, false, entity.TriggerMisfire)
}