package repository

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Task interface {
	// GetTaskGroup 获取指定的任务分组
	GetTaskGroup(ctx kratosx.Context, id uint32) (*entity.TaskGroup, error)

	// ListTaskGroup 获取任务分组列表
	ListTaskGroup(ctx kratosx.Context, req *types.ListTaskGroupRequest) ([]*entity.TaskGroup, uint32, error)

	// CreateTaskGroup 创建任务分组
	CreateTaskGroup(ctx kratosx.Context, req *entity.TaskGroup) (uint32, error)

	// UpdateTaskGroup 更新任务分组
	UpdateTaskGroup(ctx kratosx.Context, req *entity.TaskGroup) error

	// DeleteTaskGroup 删除任务分组
	DeleteTaskGroup(ctx kratosx.Context, id uint32) error

	// GetTask 获取指定的任务信息
	GetTask(ctx kratosx.Context, id uint32) (*entity.Task, error)

	// ListTask 获取任务信息列表
	ListTask(ctx kratosx.Context, req *types.ListTaskRequest) ([]*entity.Task, uint32, error)

	// CreateTask 创建任务信息
	CreateTask(ctx kratosx.Context, req *entity.Task) (uint32, error)

	// UpdateTask 更新任务信息
	UpdateTask(ctx kratosx.Context, req *entity.Task) error

	// UpdateTaskStatus 更新任务状态信息
	UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error

	// DeleteTask 删除任务信息
	DeleteTask(ctx kratosx.Context, id uint32) error

	// AllTaskSpecs 获取所有的表达式
	AllTaskSpecs(ctx kratosx.Context) map[uint32]string

	// StartCron 启动定时任务
	StartCron(scheduler func(ctx kratosx.Context, id uint32, spec string, force bool) error)

	// ValidateSpec 验证定时任务表达式
	ValidateSpec(s string) error

	// AddCron 添加定时任务
	AddCron(id uint32, spec string) error

	// RemoveCron 删除定时任务
	RemoveCron(id uint32) error

	// UpdateCron 更新定时任务
	UpdateCron(id uint32, spec string) error

	// CloseCron 优雅关闭停止定时任务
	CloseCron()
}
