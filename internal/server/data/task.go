package data

import (
	"fmt"

	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"
	"google.golang.org/protobuf/proto"

	biz "github.com/limes-cloud/cron/internal/server/biz/task"
	"github.com/limes-cloud/cron/internal/server/data/model"
)

type taskRepo struct {
}

func NewTaskRepo() biz.Repo {
	return &taskRepo{}
}

// ToTaskGroupEntity model转entity
func (r taskRepo) ToTaskGroupEntity(m *model.TaskGroup) *biz.TaskGroup {
	e := &biz.TaskGroup{}
	_ = valx.Transform(m, e)
	return e
}

// ToTaskGroupModel entity转model
func (r taskRepo) ToTaskGroupModel(e *biz.TaskGroup) *model.TaskGroup {
	m := &model.TaskGroup{}
	_ = valx.Transform(e, m)
	return m
}

// GetTaskGroup 获取指定的数据
func (r taskRepo) GetTaskGroup(ctx kratosx.Context, id uint32) (*biz.TaskGroup, error) {
	var (
		m  = model.TaskGroup{}
		fs = []string{"*"}
	)
	db := ctx.DB().Select(fs)
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	return r.ToTaskGroupEntity(&m), nil
}

// ListTaskGroup 获取列表
func (r taskRepo) ListTaskGroup(ctx kratosx.Context, req *biz.ListTaskGroupRequest) ([]*biz.TaskGroup, uint32, error) {
	var (
		bs    []*biz.TaskGroup
		ms    []*model.TaskGroup
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(model.TaskGroup{}).Select(fs)

	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	if req.OrderBy == nil || *req.OrderBy == "" {
		req.OrderBy = proto.String("id")
	}
	if req.Order == nil || *req.Order == "" {
		req.Order = proto.String("asc")
	}
	db = db.Order(fmt.Sprintf("%s %s", *req.OrderBy, *req.Order))
	if *req.OrderBy != "id" {
		db = db.Order("id asc")
	}

	if err := db.Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	for _, m := range ms {
		bs = append(bs, r.ToTaskGroupEntity(m))
	}
	return bs, uint32(total), nil
}

// CreateTaskGroup 创建数据
func (r taskRepo) CreateTaskGroup(ctx kratosx.Context, req *biz.TaskGroup) (uint32, error) {
	m := r.ToTaskGroupModel(req)
	return m.Id, ctx.DB().Create(m).Error
}

// UpdateTaskGroup 更新数据
func (r taskRepo) UpdateTaskGroup(ctx kratosx.Context, req *biz.TaskGroup) error {
	return ctx.DB().Updates(r.ToTaskGroupModel(req)).Error
}

// DeleteTaskGroup 删除数据
func (r taskRepo) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&model.TaskGroup{}).Error
}

// ToTaskEntity model转entity
func (r taskRepo) ToTaskEntity(m *model.Task) *biz.Task {
	e := &biz.Task{}
	_ = valx.Transform(m, e)
	return e
}

// ToTaskModel entity转model
func (r taskRepo) ToTaskModel(e *biz.Task) *model.Task {
	m := &model.Task{}
	_ = valx.Transform(e, m)
	return m
}

// GetTask 获取指定的数据
func (r taskRepo) GetTask(ctx kratosx.Context, id uint32) (*biz.Task, error) {
	var (
		m  = model.Task{}
		fs = []string{"*"}
	)
	db := ctx.DB().Select(fs).Preload("Group").Preload("Worker").Preload("WorkerGroup")
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	return r.ToTaskEntity(&m), nil
}

// ListTask 获取列表
func (r taskRepo) ListTask(ctx kratosx.Context, req *biz.ListTaskRequest) ([]*biz.Task, uint32, error) {
	var (
		bs    []*biz.Task
		ms    []*model.Task
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(model.Task{}).Select(fs)

	if req.GroupId != nil {
		db = db.Where("group_id = ?", *req.GroupId)
	}
	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}
	if req.Tag != nil {
		db = db.Where("tag = ?", *req.Tag)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	if req.OrderBy == nil || *req.OrderBy == "" {
		req.OrderBy = proto.String("id")
	}
	if req.Order == nil || *req.Order == "" {
		req.Order = proto.String("asc")
	}
	db = db.Order(fmt.Sprintf("%s %s", *req.OrderBy, *req.Order))
	if *req.OrderBy != "id" {
		db = db.Order("id asc")
	}

	if err := db.Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	for _, m := range ms {
		bs = append(bs, r.ToTaskEntity(m))
	}
	return bs, uint32(total), nil
}

// CreateTask 创建数据
func (r taskRepo) CreateTask(ctx kratosx.Context, req *biz.Task) (uint32, error) {
	m := r.ToTaskModel(req)
	return m.Id, ctx.DB().Create(m).Error
}

// UpdateTask 更新数据
func (r taskRepo) UpdateTask(ctx kratosx.Context, req *biz.Task) error {
	return ctx.DB().Updates(r.ToTaskModel(req)).Error
}

// UpdateTaskStatus 更新数据状态
func (r taskRepo) UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error {
	return ctx.DB().Model(model.Task{}).Where("id=?", id).Update("status", status).Error
}

// DeleteTask 删除数据
func (r taskRepo) DeleteTask(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&model.Task{}).Error
}

// GetSpecs 获取当前时间表达式
func (r taskRepo) GetSpecs(ctx kratosx.Context) map[uint32]string {
	var (
		list []*model.Task
		m    = map[uint32]string{}
	)
	ctx.DB().Model(model.Task{}).Where("status=true").Find(&list)
	for _, item := range list {
		m[item.Id] = item.Spec
	}
	return m
}
