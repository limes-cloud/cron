package dbs

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/limes-cloud/kratosx"
	"gorm.io/gorm"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Log struct {
}

var (
	logIns  *Log
	logOnce sync.Once
)

func NewLog() *Log {
	logOnce.Do(func() {
		logIns = &Log{}
	})
	return logIns
}

// GetLogStatusByUuid 获取指定uuid的状态
func (l *Log) GetLogStatusByUuid(ctx kratosx.Context, uuid string) (string, error) {
	var (
		log = entity.Log{}
		fs  = []string{"status"}
	)
	return log.Status, ctx.DB().Select(fs).First(&log, "uuid = ?", uuid).Error
}

// GetLogByUuid 获取指定的数据
func (l *Log) GetLogByUuid(ctx kratosx.Context, uuid string) (*entity.Log, error) {
	var (
		log = entity.Log{}
		fs  = []string{"*"}
	)
	return &log, ctx.DB().Select(fs).First(&log, "uuid = ?", uuid).Error
}

// GetLog 获取指定的数据
func (l *Log) GetLog(ctx kratosx.Context, id uint32) (*entity.Log, error) {
	var (
		log = entity.Log{}
		fs  = []string{"*"}
	)
	return &log, ctx.DB().Select(fs).First(&log, id).Error
}

// ListLog 获取列表
func (l *Log) ListLog(ctx kratosx.Context, req *types.ListLogRequest) ([]*entity.Log, uint32, error) {
	var (
		list  []*entity.Log
		total int64
		fs    = []string{"id", "uuid", "worker_id", "task_id", "start_at", "end_at", "status"}
	)

	db := ctx.DB().Model(entity.Log{}).Select(fs).Where("task_id=?", req.TaskId)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).Order("id desc")
	return list, uint32(total), db.Find(&list).Error
}

// CreateLog 创建数据
func (l *Log) CreateLog(ctx kratosx.Context, log *entity.Log) (uint32, error) {
	return log.Id, ctx.DB().Create(log).Error
}

// UpdateLog 更新数据
func (l *Log) UpdateLog(ctx kratosx.Context, log *entity.Log) error {
	return ctx.DB().Updates(log).Error
}

// DeleteLog 删除数据
func (l *Log) DeleteLog(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.Log{}).Error
}

// AppendLogContent 追加指定的uuid的日志内容
func (l *Log) AppendLogContent(ctx kratosx.Context, uuid string, content string) error {
	return ctx.DB().Model(&entity.Log{}).
		Where("uuid=?", uuid).
		UpdateColumn("content", gorm.Expr("CONCAT(content,?)", ","+content)).Error
}

// UpdateLogStatus 更新任务状态
func (l *Log) UpdateLogStatus(ctx kratosx.Context, uuid string, status string) error {
	log := &entity.Log{
		Status: status,
		EndAt:  time.Now().Unix(),
	}
	return ctx.DB().Where("uuid=?", uuid).Updates(log).Error
}

// GetTargetIpByUuid 获取指定uuid
func (l *Log) GetTargetIpByUuid(ctx kratosx.Context, uuid string) (string, error) {
	var (
		log    entity.Log
		worker entity.Worker
	)
	if err := ctx.DB().First(&log, "uuid=?", uuid).Error; err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(log.WorkerSnapshot), &worker); err != nil {
		return "", err
	}
	return worker.Ip, nil
}
