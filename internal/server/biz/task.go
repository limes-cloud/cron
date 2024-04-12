package biz

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/limes-cloud/kratosx"
	ktypes "github.com/limes-cloud/kratosx/types"

	clientV1 "github.com/limes-cloud/cron/api/client/v1"
	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type TaskFactory interface {
	DrySpec(s string) bool
	AddCron(id uint32, spec string) error
	UpdateCron(id uint32, spec string) error
	RemoveCron(id uint32) error
}

type TaskRepo interface {
	GetSpecs(ctx kratosx.Context) map[uint32]string
	AddTask(ctx kratosx.Context, in *Task) (uint32, error)
	GetTask(ctx kratosx.Context, id uint32) (*Task, error)
	GetWorkerByUuid(ctx kratosx.Context, uuid string) (*Worker, error)
	PageTask(ctx kratosx.Context, req *PageTaskRequest) ([]*Task, uint32, error)
	UpdateTask(ctx kratosx.Context, c *Task) error
	DeleteTask(ctx kratosx.Context, id uint32) error
	UpdateTaskStatus(ctx kratosx.Context, id uint32, status string) error
}

type PageTaskRequest struct {
	Page     uint32
	PageSize uint32
	Tag      *string
	Status   *string
	Name     *string
}

type Task struct {
	ktypes.BaseModel
	GroupId       uint32  `json:"group_id"`
	Name          string  `json:"name"`
	Tag           string  `json:"tag"`
	Spec          string  `json:"spec"`
	SelectType    string  `json:"select_type"`
	SelectValue   string  `json:"select_value"`
	WorkerGroupId *uint32 `json:"worker_group_id"`
	WorkerId      *uint32 `json:"worker_id"`
	ExecType      string  `json:"exec_type"`
	ExecValue     string  `json:"exec_value"`
	ExpectCode    uint32  `json:"expect_code"`
	RetryCount    uint32  `json:"retry_count"`
	RetryWaitTime uint32  `json:"retry_wait_time"`
	MaxExecTime   uint32  `json:"max_exec_time"`
	Status        string  `json:"status"`
	Description   string  `json:"description"`
	// 后续新增告警
}

type TaskUseCase struct {
	config  *conf.Config
	repo    TaskRepo
	factory TaskFactory
}

const (
	TaskEnabled   = "enabled"
	TaskDisabled  = "disabled"
	ExecTypeGroup = "group"
)

// NewTaskUseCase 创建UseCase实体
func NewTaskUseCase(config *conf.Config, repo TaskRepo, factory TaskFactory) *TaskUseCase {
	return &TaskUseCase{config: config, repo: repo, factory: factory}
}

// PageTask 获取分页任务
func (u *TaskUseCase) PageTask(ctx kratosx.Context, req *PageTaskRequest) ([]*Task, uint32, error) {
	task, total, err := u.repo.PageTask(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return task, total, nil
}

// AddTask 添加任务
func (u *TaskUseCase) AddTask(ctx kratosx.Context, task *Task) (uint32, error) {
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
func (u *TaskUseCase) UpdateTask(ctx kratosx.Context, task *Task) error {
	old, err := u.repo.GetTask(ctx, task.ID)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	if old.Status != TaskDisabled {
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
func (u *TaskUseCase) DeleteTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if task.Status != TaskDisabled {
		return errors.DeleteNotDisabledTask()
	}

	if err := u.repo.DeleteTask(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// CancelTask 取消执行任务
func (u *TaskUseCase) CancelTask(ctx kratosx.Context, uuid string) error {
	worker, err := u.repo.GetWorkerByUuid(ctx, uuid)
	if err != nil {
		return errors.NotFound()
	}

	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(worker.IP))
	if err != nil {
		return errors.WorkerNotAvailable()
	}
	client := clientV1.NewServiceClient(conn)
	_, err = client.CancelExecTask(ctx, &clientV1.CancelExecTaskRequest{Uuid: uuid})
	return err
}

// EnableTask 启动指定任务
func (u *TaskUseCase) EnableTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if task.Status != TaskDisabled {
		return errors.EnableNotDisabledTask()
	}

	if u.factory.DrySpec(task.Spec) {
		return errors.CronSpec()
	}

	if err := u.repo.UpdateTaskStatus(ctx, id, TaskEnabled); err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if err := u.factory.AddCron(id, task.Spec); err != nil {
		return errors.EnableCronFormat(err.Error())
	}

	return nil
}

// DisableTask 禁用指定任务
func (u *TaskUseCase) DisableTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if task.Status != TaskEnabled {
		return errors.DisableNotEnableTask()
	}

	if err := u.repo.UpdateTaskStatus(ctx, id, TaskDisabled); err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if err := u.factory.RemoveCron(id); err != nil {
		return errors.DisableCronFormat(err.Error())
	}

	return nil
}
