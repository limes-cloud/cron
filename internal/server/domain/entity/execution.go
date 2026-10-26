package entity

// Execution 执行意图：一次触发产生一条记录，贯穿 pending → running → 终态
type Execution struct {
	App         string `json:"app"         gorm:"column:app"`                 // 所属应用（X-App-Key，随任务落库快照）
	Id          string `json:"id"          gorm:"primaryKey;column:id"`       // 执行意图 ID（uuid，幂等键之一）
	TaskId      uint32 `json:"taskId"      gorm:"column:task_id"`             // 所属任务
	TriggerAt   int64  `json:"triggerAt"   gorm:"column:trigger_at"`          // 触发时刻（整周期对齐，Unix 秒）；与 TaskId 组成唯一键
	TriggerType string `json:"triggerType" gorm:"column:trigger_type"`        // 触发方式：cron | manual | misfire
	WorkerId    *uint32 `json:"workerId"   gorm:"column:worker_id"`           // 实际执行节点；NULL = 未认领
	LogUuid     string  `json:"logUuid"    gorm:"column:log_uuid"`            // 关联执行日志（log.uuid），执行时回填
	Executor    string `json:"executor"    gorm:"column:executor"`            // 执行体标识（instance_id:goroutine）；空 = 未认领
	Status      string `json:"status"      gorm:"column:status"`              // 状态：pending | running | orphan | success | fail | cancel
	HeartbeatAt *int64 `json:"heartbeatAt" gorm:"column:heartbeat_at"`        // 最近心跳（Unix 秒）；孤儿检测依据
	StartedAt   *int64 `json:"startedAt"   gorm:"column:started_at"`          // 开始执行时间（Unix 秒）
	FinishedAt  *int64 `json:"finishedAt"  gorm:"column:finished_at"`         // 结束时间（Unix 秒）
	RetryCount  uint32 `json:"retryCount"  gorm:"column:retry_count"`         // 连接重试次数
	LastError   string `json:"lastError"   gorm:"column:last_error"`          // 最近一次错误（终态失败时填充）
	CreatedAt   int64  `json:"createdAt"   gorm:"column:created_at"`          // 创建时间（Unix 秒）
	UpdatedAt   int64  `json:"updatedAt"   gorm:"column:updated_at"`          // 修改时间（Unix 秒）
}

// 执行状态常量（与 proto/业务层共用语义）
const (
	ExecPending = "pending" // 已落库未认领
	ExecRunning = "running" // 执行中（心跳维持）
	ExecOrphan  = "orphan"  // 心跳超时被 Watchdog 标记（P2 使用）
	ExecSkip    = "skip"    // 错过帧被丢弃（misfire=skip 或补帧后旧帧）
	ExecSuccess = "success" // 执行成功
	ExecFail    = "fail"    // 执行失败（含超时/取消后异常）
	ExecCancel  = "cancel"  // 主动取消
)

// 触发方式常量
const (
	TriggerCron    = "cron"    // 定时触发
	TriggerManual  = "manual"  // 手动执行
	TriggerMisfire = "misfire" // 错过补偿重投（P2 使用）
)

// IsFinal 是否为终态
func (e *Execution) IsFinal() bool {
	return e.Status == ExecSuccess || e.Status == ExecFail || e.Status == ExecCancel || e.Status == ExecSkip
}