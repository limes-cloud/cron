package task

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type UseCase struct {
	config  *conf.Config
	repo    Repo
	factory Factory
}

const (
	StatusRunning       = "running"
	StatusEnabled       = "enabled"
	StatusDisabled      = "disabled"
	ExecTypeGroupWorker = "group"
	ExecTypeWorker      = "node"
)

// NewUseCase 创建UseCase实体
func NewUseCase(config *conf.Config, repo Repo) *UseCase {
	return &UseCase{config: config, repo: repo}
}

// PageTaskGroup 获取分页任务分组
func (u *UseCase) PageTaskGroup(ctx kratosx.Context, req *PageTaskGroupRequest) ([]*TaskGroup, uint32, error) {
	tg, total, err := u.repo.PageTaskGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return tg, total, nil
}

// AddTaskGroup 添加任务分组
func (u *UseCase) AddTaskGroup(ctx kratosx.Context, tg *TaskGroup) (uint32, error) {
	id, err := u.repo.AddTaskGroup(ctx, tg)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateTaskGroup 删除指定任务分组
func (u *UseCase) UpdateTaskGroup(ctx kratosx.Context, tg *TaskGroup) error {
	if err := u.repo.UpdateTaskGroup(ctx, tg); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteTaskGroup 删除指定任务分组
func (u *UseCase) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	if err := u.repo.DeleteTaskGroup(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// PageTask 获取分页任务
func (u *UseCase) PageTask(ctx kratosx.Context, req *PageTaskRequest) ([]*Task, uint32, error) {
	task, total, err := u.repo.PageTask(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return task, total, nil
}

// AddTask 添加任务
func (u *UseCase) AddTask(ctx kratosx.Context, task *Task) (uint32, error) {
	if u.factory.DrySpec(task.Spec) {
		return 0, errors.CronSpec()
	}
	id, err := u.repo.AddTask(ctx, task)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateTask 删除指定任务
func (u *UseCase) UpdateTask(ctx kratosx.Context, task *Task) error {
	old, err := u.repo.GetTask(ctx, task.ID)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	if old.Status != StatusDisabled {
		return errors.UpdateNotEnableTask()
	}

	if u.factory.DrySpec(task.Spec) {
		return errors.CronSpec()
	}
	if err := u.repo.UpdateTask(ctx, task); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteTask 删除指定任务
func (u *UseCase) DeleteTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if task.Status != StatusDisabled {
		return errors.DeleteNotDisabledTask()
	}

	if err := u.repo.DeleteTask(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// CancelTask 取消执行任务
func (u *UseCase) CancelTask(ctx kratosx.Context, id uint32) error {
	if err := u.repo.CancelTask(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// EnableTask 启动指定任务
func (u *UseCase) EnableTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if task.Status != StatusDisabled {
		return errors.EnableNotDisabledTask()
	}

	if u.factory.DrySpec(task.Spec) {
		return errors.CronSpec()
	}

	if err := u.repo.UpdateTaskStatus(ctx, id, StatusEnabled); err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if err := u.factory.AddCron(id, task.Spec); err != nil {
		return errors.EnableCronFormat(err.Error())
	}

	return nil
}

// DisableTask 禁用指定任务
func (u *UseCase) DisableTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if task.Status != StatusEnabled {
		return errors.DisableNotEnableTask()
	}

	if err := u.repo.UpdateTaskStatus(ctx, id, StatusDisabled); err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if err := u.factory.DeleteCron(id); err != nil {
		return errors.DisableCronFormat(err.Error())
	}

	return nil
}
