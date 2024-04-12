package biz

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/client/conf"
)

type TaskFactory interface {
	ExecTask(ctx kratosx.Context, task *Task, fn ExecTaskReplyFunc) error
	CancelExecTask(uuid string)
}

type ExecTaskReply struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    uint32 `json:"time"`
}

type ExecTaskReplyFunc func(*ExecTaskReply) error

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

type TaskUseCase struct {
	config  *conf.Config
	factory TaskFactory
}

func NewTaskUseCase(config *conf.Config, factory TaskFactory) *TaskUseCase {
	return &TaskUseCase{config: config, factory: factory}
}

func (uc *TaskUseCase) ExecTask(ctx kratosx.Context, task *Task, fn ExecTaskReplyFunc) error {
	if err := uc.factory.ExecTask(ctx, task, fn); err != nil {
		return errors.ExecTaskFailFormat(err.Error())
	}
	return nil
}

func (uc *TaskUseCase) CancelExecTask(uuid string) {
	uc.factory.CancelExecTask(uuid)
}
