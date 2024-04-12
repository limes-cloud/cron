package data

import (
	"github.com/limes-cloud/kratosx"
	"gorm.io/gorm"

	biz "github.com/limes-cloud/cron/internal/server/biz"
)

type logRepo struct {
}

func NewLogRepo() biz.LogRepo {
	return &logRepo{}
}

func (t logRepo) AddLog(ctx kratosx.Context, in *biz.Log) (uint32, error) {
	return in.ID, ctx.DB().Create(in).Error
}

func (t logRepo) GetLog(ctx kratosx.Context, id uint32) (*biz.Log, error) {
	var log biz.Log
	return &log, ctx.DB().Where("id=?", id).First(&log).Error
}

func (t logRepo) PageLog(ctx kratosx.Context, req *biz.PageLogRequest) ([]*biz.Log, uint32, error) {
	var list []*biz.Log
	var total int64
	db := ctx.DB().Model(biz.Log{}).Where("task_id=?", req.TaskId)

	if err := db.Count(&total).Error; err != nil {
		return nil, uint32(total), err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	return list, uint32(total), db.Find(&list).Error
}

func (t logRepo) AppendLogContent(ctx kratosx.Context, uuid string, c string) error {
	return ctx.DB().Model(&biz.Log{}).Where("uuid=?", uuid).UpdateColumn("content", gorm.Expr("content+?", c)).Error
}

func (t logRepo) UpdateLogStatus(ctx kratosx.Context, uuid string, err error) error {
	return ctx.DB().Model(&biz.Log{}).Where("uuid=?", uuid).UpdateColumn("status", err == nil).Error
}

func (t logRepo) TaskIsRunning(ctx kratosx.Context, uuid string) bool {
	var log biz.Log
	if err := ctx.DB().Model(&biz.Log{}).Where("uuid=?", uuid).First(&log); err != nil {
		return false
	}
	return log.Status == biz.ExecRunning
}
