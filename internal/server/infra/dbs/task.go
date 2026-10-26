package dbs

import (
	"sync"
	"time"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

// Task 任务仓储（纯数据访问；调度逻辑已抽离到 scheduler 组件）
type Task struct{}

var (
	taskIns  *Task
	taskOnce sync.Once
)

func NewTask() *Task {
	taskOnce.Do(func() {
		taskIns = &Task{}
	})
	return taskIns
}

// GetTaskGroup 获取指定的数据
func (t *Task) GetTaskGroup(ctx kratosx.Context, id uint32) (*entity.TaskGroup, error) {
	var (
		taskGroup = entity.TaskGroup{}
		fs        = []string{"*"}
	)
	return &taskGroup, ctx.DB().Select(fs).First(&taskGroup, id).Error
}

// ListTaskGroup 获取列表
func (t *Task) ListTaskGroup(ctx kratosx.Context, req *types.ListTaskGroupRequest) ([]*entity.TaskGroup, uint32, error) {
	var (
		list  []*entity.TaskGroup
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.TaskGroup{}).Select(fs)
	if req.App != "" {
		db = db.Where("app = ?", req.App)
	}
	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	return list, uint32(total), db.Find(&list).Error
}

// CreateTaskGroup 创建数据
func (t *Task) CreateTaskGroup(ctx kratosx.Context, tg *entity.TaskGroup) (uint32, error) {
	return tg.Id, ctx.DB().Create(tg).Error
}

// UpdateTaskGroup 更新数据
func (t *Task) UpdateTaskGroup(ctx kratosx.Context, tg *entity.TaskGroup) error {
	return ctx.DB().Updates(tg).Error
}

// DeleteTaskGroup 删除数据
func (t *Task) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.TaskGroup{}).Error
}

// GetTask 获取指定的数据
func (t *Task) GetTask(ctx kratosx.Context, id uint32) (*entity.Task, error) {
	var (
		task = entity.Task{}
		fs   = []string{"*"}
	)
	return &task, ctx.DB().Select(fs).
		Preload("Group").
		Preload("Worker").
		Preload("WorkerGroup").
		First(&task, id).Error
}

// ListTask 获取列表
func (t *Task) ListTask(ctx kratosx.Context, req *types.ListTaskRequest) ([]*entity.Task, uint32, error) {
	var (
		list  []*entity.Task
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.Task{}).Select(fs)
	if req.App != "" {
		db = db.Where("app = ?", req.App)
	}
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
	return list, uint32(total), db.Find(&list).Error
}

// CreateTask 创建数据
func (t *Task) CreateTask(ctx kratosx.Context, task *entity.Task) (uint32, error) {
	return task.Id, ctx.DB().Create(task).Error
}

// UpdateTask 更新数据
func (t *Task) UpdateTask(ctx kratosx.Context, task *entity.Task) error {
	return ctx.DB().Updates(task).Error
}

// UpdateTaskStatus 更新数据状态
func (t *Task) UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error {
	return ctx.DB().Model(entity.Task{}).Where("id=?", id).Update("status", status).Error
}

// DeleteTask 删除数据
func (t *Task) DeleteTask(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.Task{}).Error
}

// AllTaskSpecs 获取当前启用的任务的所有定时时间表达式（scheduler 启动加载用）
func (t *Task) AllTaskSpecs(ctx kratosx.Context) map[uint32]string {
	var (
		list []*entity.Task
		m    = map[uint32]string{}
	)
	ctx.DB().Model(entity.Task{}).
		Where("status=true").
		Where("start is null or start <= ?", time.Now().Unix()).
		Where("end is null or end >= ?", time.Now().Unix()).
		Find(&list)
	for _, item := range list {
		m[item.Id] = item.Spec
	}
	return m
}