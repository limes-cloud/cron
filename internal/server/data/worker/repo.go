package worker

import (
	"github.com/limes-cloud/kratosx"

	biz "github.com/limes-cloud/cron/internal/server/biz/worker"
)

type repo struct {
}

func NewRepo() biz.Repo {
	return &repo{}
}

func (r repo) AddWorker(ctx kratosx.Context, in *biz.Worker) (uint32, error) {
	in.Status = biz.StatusDisabled
	return in.ID, ctx.DB().Create(in).Error
}

func (r repo) GetWorker(ctx kratosx.Context, id uint32) (*biz.Worker, error) {
	w := biz.Worker{}
	return &w, ctx.DB().First(&w).Error
}

func (r repo) GetWorkersByTag(ctx kratosx.Context, tag string) ([]*biz.Worker, error) {
	var ws []*biz.Worker
	return ws, ctx.DB().Where("tag=?", tag).Find(&ws).Error
}

func (r repo) PageWorker(ctx kratosx.Context, req *biz.PageWorkerRequest) ([]*biz.Worker, uint32, error) {
	var list []*biz.Worker
	var total int64
	db := ctx.DB().Model(biz.Worker{})
	if req.Tag != nil {
		db.Where("tag=?", *req.Tag)
	}
	if req.Status != nil {
		db.Where("status=?", *req.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, uint32(total), err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	return list, uint32(total), db.Find(&list).Error
}

func (r repo) UpdateWorker(ctx kratosx.Context, in *biz.Worker) error {
	in.Status = ""
	return ctx.DB().Updates(in).Error
}

func (r repo) DeleteWorker(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Delete(biz.Worker{}, id).Error
}

func (r repo) EnableWorker(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.StatusEnabled}).Error
}

func (r repo) UpdateWorkerRunning(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.StatusEnabled}).Error
}

func (r repo) UpdateWorkerStopping(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.StatusEnabled}).Error
}

func (r repo) DisableWorker(ctx kratosx.Context, id uint32, desc string) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.StatusDisabled, StopDesc: desc}).Error
}
