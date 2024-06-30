package task

import (
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/crypto"
	"google.golang.org/protobuf/proto"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type UseCase struct {
	conf    *conf.Config
	repo    Repo
	factory Factory
}

func NewUseCase(config *conf.Config, repo Repo, factory Factory) *UseCase {
	return &UseCase{conf: config, repo: repo, factory: factory}
}

// GetTaskGroup 获取指定的任务分组
func (u *UseCase) GetTaskGroup(ctx kratosx.Context, req *GetTaskGroupRequest) (*TaskGroup, error) {
	var (
		res *TaskGroup
		err error
	)

	if req.Id != nil {
		res, err = u.repo.GetTaskGroup(ctx, *req.Id)
	} else {
		return nil, errors.ParamsError()
	}

	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListTaskGroup 获取任务分组列表
func (u *UseCase) ListTaskGroup(ctx kratosx.Context, req *ListTaskGroupRequest) ([]*TaskGroup, uint32, error) {
	list, total, err := u.repo.ListTaskGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateTaskGroup 创建任务分组
func (u *UseCase) CreateTaskGroup(ctx kratosx.Context, req *TaskGroup) (uint32, error) {
	id, err := u.repo.CreateTaskGroup(ctx, req)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateTaskGroup 更新任务分组
func (u *UseCase) UpdateTaskGroup(ctx kratosx.Context, req *TaskGroup) error {
	if err := u.repo.UpdateTaskGroup(ctx, req); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteTaskGroup 删除任务分组
func (u *UseCase) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	err := u.repo.DeleteTaskGroup(ctx, id)
	if err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetTask 获取指定的任务信息
func (u *UseCase) GetTask(ctx kratosx.Context, req *GetTaskRequest) (*Task, error) {
	var (
		res *Task
		err error
	)

	if req.Id != nil {
		res, err = u.repo.GetTask(ctx, *req.Id)
	} else {
		return nil, errors.ParamsError()
	}

	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListTask 获取任务信息列表
func (u *UseCase) ListTask(ctx kratosx.Context, req *ListTaskRequest) ([]*Task, uint32, error) {
	list, total, err := u.repo.ListTask(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateTask 创建任务信息
func (u *UseCase) CreateTask(ctx kratosx.Context, req *Task) (uint32, error) {
	if !u.factory.DrySpec(req.Spec) {
		return 0, errors.CronSpecError()
	}
	req.Version = crypto.MD5ToUpper([]byte(req.Spec + req.ExecType + req.ExecValue))
	req.Status = proto.Bool(false)

	id, err := u.repo.CreateTask(ctx, req)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateTask 更新任务信息
func (u *UseCase) UpdateTask(ctx kratosx.Context, req *Task) error {
	task, err := u.repo.GetTask(ctx, req.Id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if task.Status != nil && *task.Status {
		return errors.UpdateNotDisableTaskError()
	}

	if !u.factory.DrySpec(req.Spec) {
		return errors.CronSpecError()
	}

	if err := u.repo.UpdateTask(ctx, req); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// UpdateTaskStatus 更新任务信息状态
func (u *UseCase) UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}

	return ctx.Transaction(func(ctx kratosx.Context) error {
		if err := u.repo.UpdateTaskStatus(ctx, id, status); err != nil {
			return errors.DatabaseError(err.Error())
		}
		if status {
			if err := u.factory.AddCron(id, task.Spec); err != nil {
				return errors.EnableCronError(err.Error())
			}
		} else {
			if err := u.factory.RemoveCron(id); err != nil {
				return errors.DisableCronError(err.Error())
			}
		}

		return nil
	})
}

// DeleteTask 删除任务信息
func (u *UseCase) DeleteTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}

	if task.Status != nil && *task.Status {
		return errors.DeleteNotDisabledTaskError()
	}

	if err := u.repo.DeleteTask(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// CancelExecTask 取消执行任务
func (u *UseCase) CancelExecTask(ctx kratosx.Context, uuid string) error {
	if err := u.factory.CancelExec(ctx, uuid); err != nil {
		return errors.CancelTaskError(err.Error())
	}
	return nil
}

// ExecTask 执行任务
func (u *UseCase) ExecTask(ctx kratosx.Context, id uint32) error {
	go func() {
		if err := u.factory.Scheduler(id, true); err != nil {
			ctx.Logger().Errorw("exec task error", err.Error())
		}
	}()
	return nil
}
