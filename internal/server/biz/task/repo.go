package task

import "github.com/limes-cloud/kratosx"

type Repo interface {
	AddTaskGroup(ctx kratosx.Context, in *TaskGroup) (uint32, error)
	PageTaskGroup(ctx kratosx.Context, req *PageTaskGroupRequest) ([]*TaskGroup, uint32, error)
	UpdateTaskGroup(ctx kratosx.Context, c *TaskGroup) error
	DeleteTaskGroup(ctx kratosx.Context, id uint32) error

	AddTask(ctx kratosx.Context, in *Task) (uint32, error)
	GetTask(ctx kratosx.Context, id uint32) (*Task, error)
	GetTasksByTag(ctx kratosx.Context, tag string) ([]*Task, error)
	PageTask(ctx kratosx.Context, req *PageTaskRequest) ([]*Task, uint32, error)
	UpdateTask(ctx kratosx.Context, c *Task) error
	DeleteTask(ctx kratosx.Context, id uint32) error
	UpdateTaskStatus(ctx kratosx.Context, id uint32, status string) error
	CancelTask(ctx kratosx.Context, id uint32) error
}
