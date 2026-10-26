package dbs

import (
	"strings"
	"sync"
	"time"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/types"
)

// Execution 执行意图仓储实现
type Execution struct{}

var (
	executionIns  *Execution
	executionOnce sync.Once
)

func NewExecution() *Execution {
	executionOnce.Do(func() {
		executionIns = &Execution{}
	})
	return executionIns
}

// CreateExecution 创建执行意图（幂等：唯一键 (task_id, trigger_at) 冲突视为已存在）
func (e *Execution) CreateExecution(ctx kratosx.Context, ex *entity.Execution) error {
	err := ctx.DB().Create(ex).Error
	if err == nil {
		return nil
	}
	if isDuplicateKey(err) {
		return repository.ErrExecutionExists
	}
	return err
}

// GetExecution 根据 ID 查询执行意图
func (e *Execution) GetExecution(ctx kratosx.Context, id string) (*entity.Execution, error) {
	var (
		ex = entity.Execution{}
		fs = []string{"*"}
	)
	return &ex, ctx.DB().Select(fs).First(&ex, "id = ?", id).Error
}

// MarkRunning 认领执行
func (e *Execution) MarkRunning(ctx kratosx.Context, id string, workerID *uint32, executor string) error {
	now := time.Now().Unix()
	return ctx.DB().Model(&entity.Execution{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       entity.ExecRunning,
			"worker_id":    workerID,
			"executor":     executor,
			"started_at":   now,
			"heartbeat_at": now,
		}).Error
}

// UpdateHeartbeat 刷新心跳时间
func (e *Execution) UpdateHeartbeat(ctx kratosx.Context, id string) error {
	return ctx.DB().Model(&entity.Execution{}).
		Where("id = ?", id).
		Update("heartbeat_at", time.Now().Unix()).Error
}

// MarkOrphan 标记孤儿执行
func (e *Execution) MarkOrphan(ctx kratosx.Context, id string) error {
	return ctx.DB().Model(&entity.Execution{}).
		Where("id = ? AND status = ?", id, entity.ExecRunning).
		Update("status", entity.ExecOrphan).Error
}

// FinishExecution 写入终态
func (e *Execution) FinishExecution(ctx kratosx.Context, id string, status string, finishedAt int64, lastError string) error {
	return ctx.DB().Model(&entity.Execution{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"finished_at": finishedAt,
			"last_error":  lastError,
		}).Error
}

// ListPending 查询待执行（启动恢复用）
func (e *Execution) ListPending(ctx kratosx.Context, before int64, limit int) ([]*entity.Execution, error) {
	var list []*entity.Execution
	err := ctx.DB().Model(&entity.Execution{}).
		Where("status = ?", entity.ExecPending).
		Where("created_at <= ?", before).
		Order("created_at ASC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// ListRunning 查询执行中（启动恢复/孤儿扫描用）
func (e *Execution) ListRunning(ctx kratosx.Context, heartbeatBefore int64, limit int) ([]*entity.Execution, error) {
	var list []*entity.Execution
	err := ctx.DB().Model(&entity.Execution{}).
		Where("status = ?", entity.ExecRunning).
		Where("heartbeat_at <= ?", heartbeatBefore).
		Order("heartbeat_at ASC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// CountSkippedSince 查询任务自 since 起被跳过的帧数
func (e *Execution) CountSkippedSince(ctx kratosx.Context, taskID uint32, since int64) (int64, error) {
	var count int64
	err := ctx.DB().Model(&entity.Execution{}).
		Where("task_id = ?", taskID).
		Where("status = ?", entity.ExecSkip).
		Where("finished_at >= ?", since).
		Count(&count).Error
	return count, err
}

// isDuplicateKey 判断是否为 MySQL 唯一键冲突（1062 / Duplicate entry）
func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "1062") || strings.Contains(msg, "Duplicate entry")
}
// AttachLog 回填关联的执行日志（log.uuid）
func (e *Execution) AttachLog(ctx kratosx.Context, id string, logUUID string) error {
	return ctx.DB().Model(&entity.Execution{}).
		Where("id = ?", id).
		Update("log_uuid", logUUID).Error
}

// ListExecution 查询执行历史（按 app+task 过滤，分页）
func (e *Execution) ListExecution(ctx kratosx.Context, req *types.ListExecutionRequest) ([]*entity.Execution, uint32, error) {
	var (
		list  []*entity.Execution
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(&entity.Execution{}).Select(fs)
	if req.App != "" {
		db = db.Where("app = ?", req.App)
	}
	if req.TaskId > 0 {
		db = db.Where("task_id = ?", req.TaskId)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.TriggerType != nil {
		db = db.Where("trigger_type = ?", *req.TriggerType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).Order("trigger_at DESC")
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, uint32(total), nil
}
