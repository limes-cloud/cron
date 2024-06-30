package task

import (
	"github.com/limes-cloud/kratosx"
)

type Repo interface {
	// GetTaskGroup 获取指定的任务分组
	GetTaskGroup(ctx kratosx.Context, id uint32) (*TaskGroup, error)

	// ListTaskGroup 获取任务分组列表
	ListTaskGroup(ctx kratosx.Context, req *ListTaskGroupRequest) ([]*TaskGroup, uint32, error)

	// CreateTaskGroup 创建任务分组
	CreateTaskGroup(ctx kratosx.Context, req *TaskGroup) (uint32, error)

	// UpdateTaskGroup 更新任务分组
	UpdateTaskGroup(ctx kratosx.Context, req *TaskGroup) error

	// DeleteTaskGroup 删除任务分组
	DeleteTaskGroup(ctx kratosx.Context, id uint32) error

	// GetTask 获取指定的任务信息
	GetTask(ctx kratosx.Context, id uint32) (*Task, error)

	// ListTask 获取任务信息列表
	ListTask(ctx kratosx.Context, req *ListTaskRequest) ([]*Task, uint32, error)

	// CreateTask 创建任务信息
	CreateTask(ctx kratosx.Context, req *Task) (uint32, error)

	// UpdateTask 更新任务信息
	UpdateTask(ctx kratosx.Context, req *Task) error

	// UpdateTaskStatus 更新任务信息状态
	UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error

	// DeleteTask 删除任务信息
	DeleteTask(ctx kratosx.Context, id uint32) error

	// GetSpecs 获取所有的表达式
	GetSpecs(ctx kratosx.Context) map[uint32]string
}
