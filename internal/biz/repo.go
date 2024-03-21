package task

import "github.com/limes-cloud/kratosx"

type Repo interface {
	GetGroupByID(ctx kratosx.Context, id uint32) (*Group, error)
	GetGroupByKeyword(ctx kratosx.Context, keyword string) (*Group, error)
	PageGroup(ctx kratosx.Context, req *PageGroupRequest) ([]*Group, uint32, error)
	AddGroup(ctx kratosx.Context, c *Group) (uint32, error)
	UpdateGroup(ctx kratosx.Context, c *Group) error
	DeleteGroup(ctx kratosx.Context, uint322 uint32) error

	GetTaskByID(ctx kratosx.Context, id uint32) (*Task, error)
	PageTask(ctx kratosx.Context, req *PageTaskRequest) ([]*Task, uint32, error)
	AddTask(ctx kratosx.Context, c *Task) (uint32, error)
	UpdateTask(ctx kratosx.Context, c *Task) error
	DeleteTask(ctx kratosx.Context, uint322 uint32) error
}
