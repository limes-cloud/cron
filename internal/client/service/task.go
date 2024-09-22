package service

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/client/conf"
)

type TaskFactory interface {
	ExecTask(ctx kratosx.Context, task *ExecTaskRequest, fn ExecTaskReplyFunc) error
	CancelExecTask(uuid string)
}

type ExecTaskRequest struct {
	Id            uint32
	Uuid          string
	Type          string
	Value         string
	ExpectCode    uint32
	RetryCount    uint32
	RetryWaitTime uint32
	MaxExecTime   uint32
}

type ExecTaskReply struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    uint32 `json:"time"`
}

type ExecTaskReplyFunc func(*ExecTaskReply) error

type Task struct {
	conf    *conf.Config
	factory TaskFactory
}

func NewTask(config *conf.Config, factory TaskFactory) *Task {
	return &Task{conf: config, factory: factory}
}

func (uc *Task) ExecTask(ctx kratosx.Context, task *ExecTaskRequest, fn ExecTaskReplyFunc) error {
	if err := uc.factory.ExecTask(ctx, task, fn); err != nil {
		return errors.ExecTaskFailError(err.Error())
	}
	return nil
}

func (uc *Task) CancelExecTask(uuid string) {
	uc.factory.CancelExecTask(uuid)
}
