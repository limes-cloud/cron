package data

import (
	"errors"
	"math/rand"

	"github.com/limes-cloud/kratosx"

	biz "github.com/limes-cloud/cron/internal/server/biz"
)

type workerRepo struct {
}

func NewWorkerRepo() biz.WorkerRepo {
	return &workerRepo{}
}

func (r workerRepo) AddWorkerGroup(ctx kratosx.Context, in *biz.WorkerGroup) (uint32, error) {
	return in.ID, ctx.DB().Create(in).Error
}

func (r workerRepo) PageWorkerGroup(ctx kratosx.Context, req *biz.PageWorkerGroupRequest) ([]*biz.WorkerGroup, uint32, error) {
	var list []*biz.WorkerGroup
	var total int64
	db := ctx.DB().Model(biz.WorkerGroup{})

	if err := db.Count(&total).Error; err != nil {
		return nil, uint32(total), err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	return list, uint32(total), db.Find(&list).Error
}

func (r workerRepo) UpdateWorkerGroup(ctx kratosx.Context, c *biz.WorkerGroup) error {
	return ctx.DB().Updates(c).Error
}

func (r workerRepo) DeleteWorkerGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Delete(&biz.WorkerGroup{}).Error
}

func (r workerRepo) AddWorker(ctx kratosx.Context, in *biz.Worker) (uint32, error) {
	in.Status = biz.WorkerDisabled
	return in.ID, ctx.DB().Create(in).Error
}

func (r workerRepo) GetWorker(ctx kratosx.Context, id uint32) (*biz.Worker, error) {
	w := biz.Worker{}
	return &w, ctx.DB().First(&w).Error
}

func (r workerRepo) GetWorkersByTag(ctx kratosx.Context, tag string) ([]*biz.Worker, error) {
	var ws []*biz.Worker
	return ws, ctx.DB().Where("tag=?", tag).Find(&ws).Error
}

func (r workerRepo) PageWorker(ctx kratosx.Context, req *biz.PageWorkerRequest) ([]*biz.Worker, uint32, error) {
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

func (r workerRepo) UpdateWorker(ctx kratosx.Context, in *biz.Worker) error {
	in.Status = ""
	return ctx.DB().Updates(in).Error
}

func (r workerRepo) DeleteWorker(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Delete(biz.Worker{}, id).Error
}

func (r workerRepo) EnableWorker(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.WorkerEnabled}).Error
}

func (r workerRepo) UpdateWorkerRunning(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.WorkerEnabled}).Error
}

func (r workerRepo) UpdateWorkerStopping(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.WorkerEnabled}).Error
}

func (r workerRepo) DisableWorker(ctx kratosx.Context, id uint32, desc string) error {
	return ctx.DB().Where("id=?", id).Updates(biz.Worker{Status: biz.WorkerDisabled, StopDesc: desc}).Error
}

func (r workerRepo) GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*biz.Worker, error) {
	var list []*biz.Worker
	if err := ctx.DB().Where("group_id=?", id).Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("not exist worker")
	}
	x := rand.Intn(len(list))
	return list[x], nil
}
