package scheduler

import (
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

// Trigger 触发源：维护任务的调度表，到点回调。
// 当前实现为 robfig/cron/v3；接口化以便替换（如基于队列的触发源）。
type Trigger interface {
	// Start 启动触发循环（幂等；Stop 后再次 Start 可重新运行）
	Start()
	// Add 注册任务（同 taskID 已存在时先移除再注册，幂等）
	Add(taskID uint32, spec string, fn func()) error
	// Remove 移除任务
	Remove(taskID uint32)
	// Update 更新任务表达式（同 Add 的幂等语义）
	Update(taskID uint32, spec string, fn func()) error
	// Sync 以期望的 spec 集合为准全量同步（对账用；多出的移除，缺的注册）
	Sync(specs map[uint32]string, wrapFn func(taskID uint32, spec string) func())
	// ValidateSpec 校验 cron 表达式（秒级）
	ValidateSpec(spec string) error
	// Stop 停止触发循环（幂等；不再产生新触发，在办任务由 Scheduler 管理）
	Stop()
}

// CronTrigger Trigger 的 robfig/cron/v3 实现
type CronTrigger struct {
	c       *cron.Cron
	parser  cron.Parser
	entries map[uint32]cron.EntryID // taskID → cron EntryID（cron 的 EntryID 自增，需映射）
	mu      sync.Mutex
}

var (
	// secondsParser 秒级表达式解析器（与原 pkg/cron 的 WithSeconds 语义一致）
	secondsParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

// recoverLogger 将 cron 的日志/panic 适配到 kratos logger
type recoverLogger struct {
	h *log.Helper
}

func (l recoverLogger) Info(msg string, keysAndValues ...interface{}) {
	l.h.Infow(append([]interface{}{"msg", msg}, keysAndValues...)...)
}

func (l recoverLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.h.Errorw(append([]interface{}{"msg", msg, "err", err.Error()}, keysAndValues...)...)
}

// NewCronTrigger 创建触发源（Asia/Shanghai 时区；启用 Recover 防单任务 panic 崩进程）
func NewCronTrigger(logger *log.Helper) *CronTrigger {
	location, _ := time.LoadLocation("Asia/Shanghai")
	if logger == nil {
		logger = log.NewHelper(log.DefaultLogger)
	}
	c := cron.New(
		cron.WithParser(secondsParser),
		cron.WithLocation(location),
		cron.WithChain(cron.Recover(recoverLogger{h: logger})),
	)
	return &CronTrigger{
		c:       c,
		parser:  secondsParser,
		entries: make(map[uint32]cron.EntryID),
	}
}

// Start 启动触发循环（幂等；由 Scheduler.Start 显式调用，非 leader 实例不启动）
func (t *CronTrigger) Start() {
	t.c.Start()
}

func (t *CronTrigger) Add(taskID uint32, spec string, fn func()) error {
	id, err := t.c.AddFunc(spec, fn)
	if err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	// 幂等：同 taskID 已存在则先移除旧条目
	if old, ok := t.entries[taskID]; ok {
		t.c.Remove(old)
	}
	t.entries[taskID] = id
	return nil
}

func (t *CronTrigger) Remove(taskID uint32) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if old, ok := t.entries[taskID]; ok {
		t.c.Remove(old)
		delete(t.entries, taskID)
	}
}

func (t *CronTrigger) Update(taskID uint32, spec string, fn func()) error {
	return t.Add(taskID, spec, fn)
}

func (t *CronTrigger) ValidateSpec(spec string) error {
	_, err := t.parser.Parse(spec)
	return err
}

func (t *CronTrigger) Stop() {
	// cron.Stop 返回等 job 结束的 ctx，这里不阻塞（等办任务由 Scheduler.wg 管理）
	t.c.Stop()
}

// Sync 以期望的 spec 集合为准全量同步（对账用；多出的移除，缺的注册）
func (t *CronTrigger) Sync(specs map[uint32]string, wrapFn func(taskID uint32, spec string) func()) {
	// 新增/更新（Add 内部幂等：同 ID 先移除）
	for id, spec := range specs {
		_ = t.Add(id, spec, wrapFn(id, spec))
	}
	// 移除已不存在
	t.mu.Lock()
	for id := range t.entries {
		if _, ok := specs[id]; !ok {
			t.c.Remove(t.entries[id])
			delete(t.entries, id)
		}
	}
	t.mu.Unlock()
}