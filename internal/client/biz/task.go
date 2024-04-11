package biz

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/client/conf"
)

type Task struct {
	Id            uint32
	Uuid          string
	Type          string
	Value         string
	ExpectCode    uint32
	RetryCount    uint32
	RetryWaitTime uint32
	MaxExecTime   uint32
}

type TaskFactory interface {
	ExecTask(ctx kratosx.Context, task *Task, fn ExecTaskReplyFunc) error
	CancelExecTask(uuid string)
}

type TaskUseCase struct {
	config  *conf.Config
	factory TaskFactory
}

func NewTaskUseCase(config *conf.Config, factory TaskFactory) *TaskUseCase {
	return &TaskUseCase{config: config, factory: factory}
}

func (uc *TaskUseCase) ExecTask(ctx kratosx.Context, task *Task, fn ExecTaskReplyFunc) error {
	return uc.factory.ExecTask(ctx, task, fn)
}

func (uc *TaskUseCase) CancelExecTask(uuid string) {
	uc.factory.CancelExecTask(uuid)
}
