package data

import (
	"github.com/limes-cloud/kratosx"

	biz "github.com/limes-cloud/cron/internal/server/biz"
)

type taskRepo struct {
}

func NewTaskRepo() biz.TaskRepo {
	return &taskRepo{}
}

func (t taskRepo) AddTaskGroup(ctx kratosx.Context, in *biz.TaskGroup) (uint32, error) {
	return in.ID, ctx.DB().Create(in).Error
}

func (t taskRepo) AllTaskGroup(ctx kratosx.Context) ([]*biz.TaskGroup, error) {
	var list []*biz.TaskGroup
	return list, ctx.DB().Model(biz.TaskGroup{}).Find(&list).Error
}

func (t taskRepo) UpdateTaskGroup(ctx kratosx.Context, c *biz.TaskGroup) error {
	return ctx.DB().Updates(c).Error
}

func (t taskRepo) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Delete(&biz.TaskGroup{}).Error
}

func (t taskRepo) AddTask(ctx kratosx.Context, in *biz.Task) (uint32, error) {
	in.Status = biz.TaskDisabled
	return in.ID, ctx.DB().Create(in).Error
}

func (t taskRepo) GetTask(ctx kratosx.Context, id uint32) (*biz.Task, error) {
	var task biz.Task
	return &task, ctx.DB().Where("id=?", id).First(&task).Error
}

func (t taskRepo) PageTask(ctx kratosx.Context, req *biz.PageTaskRequest) ([]*biz.Task, uint32, error) {
	var list []*biz.Task
	var total int64
	db := ctx.DB().Model(biz.Task{}).Preload("Group")
	if req.Tag != nil {
		db = db.Where("tag=?", *req.Tag)
	}
	if req.Status != nil {
		db = db.Where("status=?", *req.Status)
	}
	if req.Name != nil {
		db = db.Where("name like ?", *req.Name+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, uint32(total), err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	return list, uint32(total), db.Find(&list).Error
}

func (t taskRepo) UpdateTask(ctx kratosx.Context, c *biz.Task) error {
	return ctx.DB().Updates(c).Error
}

func (t taskRepo) DeleteTask(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Delete(&biz.Task{}).Error
}

func (t taskRepo) UpdateTaskStatus(ctx kratosx.Context, id uint32, status *bool) error {
	return ctx.DB().Model(&biz.Task{}).Where("id", id).Update("status", *status).Error
}

func (t taskRepo) GetSpecs(ctx kratosx.Context) map[uint32]string {
	var (
		list []*biz.Task
		m    = map[uint32]string{}
	)
	ctx.DB().Model(biz.Task{}).Where("status=?", biz.TaskEnabled).Find(&list)
	for _, item := range list {
		m[item.ID] = item.Spec
	}
	return m
}

func (t taskRepo) GetWorkerByUuid(ctx kratosx.Context, uuid string) (*biz.Worker, error) {
	var (
		log    biz.Log
		worker biz.Worker
	)
	if err := ctx.DB().First(&log, "uuid=?", uuid).Error; err != nil {
		return nil, err
	}

	return &worker, ctx.DB().First(&worker, "id=?", log.TaskId).Error
}

func (t taskRepo) CancelTask(ctx kratosx.Context, uuid string) error {
	return ctx.DB().Model(biz.Log{}).Where("uuid=?", uuid).UpdateColumn("status", biz.ExecCancel).Error
}