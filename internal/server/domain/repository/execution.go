package repository

import (
	"errors"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

// ErrExecutionExists 执行意图已存在（创建幂等命中的统一错误）
var ErrExecutionExists = errors.New("execution already exists")

// Execution 执行意图仓储接口
type Execution interface {
	// CreateExecution 创建执行意图（幂等：同 (task_id, trigger_at) 已存在时返回 entity.ErrExecutionExists）
	CreateExecution(ctx kratosx.Context, e *entity.Execution) error

	// GetExecution 根据 ID 查询执行意图
	GetExecution(ctx kratosx.Context, id string) (*entity.Execution, error)

	// AttachLog 回填关联的执行日志（log.uuid）
	AttachLog(ctx kratosx.Context, id string, logUUID string) error

	// ListExecution 查询执行历史（按 app+task 过滤，分页）
	ListExecution(ctx kratosx.Context, req *types.ListExecutionRequest) ([]*entity.Execution, uint32, error)

	// MarkRunning 认领执行：写 worker/executor/startedAt/status=running/心跳
	MarkRunning(ctx kratosx.Context, id string, workerID *uint32, executor string) error

	// UpdateHeartbeat 刷新心跳时间（执行体续期）
	UpdateHeartbeat(ctx kratosx.Context, id string) error

	// MarkOrphan 标记孤儿执行（心跳超时，P2 由 Watchdog 调用）
	MarkOrphan(ctx kratosx.Context, id string) error

	// FinishExecution 写入终态（success/fail/cancel）并记录结束时间与错误
	FinishExecution(ctx kratosx.Context, id string, status string, finishedAt int64, lastError string) error

	// ListPending 查询待执行（启动恢复用）：created_at 早于 before 且仍 pending
	ListPending(ctx kratosx.Context, before int64, limit int) ([]*entity.Execution, error)

	// ListRunning 查询执行中（启动恢复/孤儿扫描用）：心跳早于 heartbeatBefore
	ListRunning(ctx kratosx.Context, heartbeatBefore int64, limit int) ([]*entity.Execution, error)

	// CountSkippedSince 查询任务自 since 起被跳过的帧数（misfire 补帧判断用）
	CountSkippedSince(ctx kratosx.Context, taskID uint32, since int64) (int64, error)
}