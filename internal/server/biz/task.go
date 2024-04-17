package biz

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/util"
	ktypes "github.com/limes-cloud/kratosx/types"
	"google.golang.org/protobuf/proto"

	clientV1 "github.com/limes-cloud/cron/api/client/v1"
	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type TaskFactory interface {
	DrySpec(s string) bool
	AddCron(id uint32, spec string) error
	UpdateCron(id uint32, spec string) error
	RemoveCron(id uint32) error
	Scheduler(id uint32) error
}

type TaskRepo interface {
	AddTaskGroup(ctx kratosx.Context, in *TaskGroup) (uint32, error)
	AllTaskGroup(ctx kratosx.Context) ([]*TaskGroup, error)
	UpdateTaskGroup(ctx kratosx.Context, c *TaskGroup) error
	DeleteTaskGroup(ctx kratosx.Context, id uint32) error

	GetSpecs(ctx kratosx.Context) map[uint32]string
	AddTask(ctx kratosx.Context, in *Task) (uint32, error)
	GetTask(ctx kratosx.Context, id uint32) (*Task, error)
	GetWorkerByUuid(ctx kratosx.Context, uuid string) (*Worker, error)
	PageTask(ctx kratosx.Context, req *PageTaskRequest) ([]*Task, uint32, error)
	UpdateTask(ctx kratosx.Context, c *Task) error
	DeleteTask(ctx kratosx.Context, id uint32) error
	CancelTask(ctx kratosx.Context, uuid string) error
	UpdateTaskStatus(ctx kratosx.Context, id uint32, status *bool) error
}

type PageTaskRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"page_size"`
	Tag      *string `json:"tag"`
	Status   *string `json:"status"`
	Name     *string `json:"name"`
}

type TaskGroup struct {
	ktypes.BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Task struct {
	ktypes.BaseModel
	GroupId       uint32     `json:"group_id"`
	Name          string     `json:"name"`
	Tag           string     `json:"tag"`
	Spec          string     `json:"spec"`
	SelectType    string     `json:"select_type"`
	WorkerGroupId *uint32    `json:"worker_group_id"`
	WorkerId      *uint32    `json:"worker_id"`
	ExecType      string     `json:"exec_type"`
	ExecValue     string     `json:"exec_value"`
	ExpectCode    *uint32    `json:"expect_code"`
	RetryCount    *uint32    `json:"retry_count"`
	RetryWaitTime *uint32    `json:"retry_wait_time"`
	MaxExecTime   *uint32    `json:"max_exec_time"`
	Status        *bool      `json:"status"`
	Version       string     `json:"version"`
	Description   string     `json:"description"`
	Group         *TaskGroup `json:"group"`
}

type TaskUseCase struct {
	config  *conf.Config
	repo    TaskRepo
	factory TaskFactory
}

var (
	TaskEnabled  = proto.Bool(true)
	TaskDisabled = proto.Bool(false)
)

const (
	ExecTypeGroup = "group"
)

// NewTaskUseCase 创建UseCase实体
func NewTaskUseCase(config *conf.Config, repo TaskRepo, factory TaskFactory) *TaskUseCase {
	return &TaskUseCase{config: config, repo: repo, factory: factory}
}

// AllTaskGroup 获取全部任务分组
func (u *TaskUseCase) AllTaskGroup(ctx kratosx.Context) ([]*TaskGroup, error) {
	tg, err := u.repo.AllTaskGroup(ctx)
	if err != nil {
		return nil, errors.NotFound()
	}
	return tg, nil
}

// AddTaskGroup 添加任务分组
func (u *TaskUseCase) AddTaskGroup(ctx kratosx.Context, tg *TaskGroup) (uint32, error) {
	id, err := u.repo.AddTaskGroup(ctx, tg)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateTaskGroup 删除指定任务分组
func (u *TaskUseCase) UpdateTaskGroup(ctx kratosx.Context, tg *TaskGroup) error {
	if err := u.repo.UpdateTaskGroup(ctx, tg); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteTaskGroup 删除指定任务分组
func (u *TaskUseCase) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	if err := u.repo.DeleteTaskGroup(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
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
	if !u.factory.DrySpec(task.Spec) {
		return 0, errors.CronSpec()
	}
	task.Version = util.MD5ToUpper([]byte(task.Spec + task.ExecType + task.ExecValue))
	task.Status = TaskDisabled
	id, err := u.repo.AddTask(ctx, task)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateTask 修改指定任务
func (u *TaskUseCase) UpdateTask(ctx kratosx.Context, task *Task) error {
	old, err := u.repo.GetTask(ctx, task.ID)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	if *old.Status != *TaskDisabled {
		return errors.UpdateNotEnableTask()
	}

	if !u.factory.DrySpec(task.Spec) {
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

	if *task.Status != *TaskDisabled {
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
	if _, err = client.CancelExecTask(ctx, &clientV1.CancelExecTaskRequest{Uuid: uuid}); err != nil {
		return errors.WorkerNotAvailable()
	}
	if err := u.repo.CancelTask(ctx, uuid); err != nil {
		ctx.Logger().Errorw("cancel task db error", err.Error())
	}

	return err
}

// ExecTask 执行任务
func (u *TaskUseCase) ExecTask(ctx kratosx.Context, id uint32) error {
	go func() {
		if err := u.factory.Scheduler(id); err != nil {
			ctx.Logger().Errorw("exec task error", err.Error())
		}
	}()
	return nil
}

// EnableTask 启动指定任务
func (u *TaskUseCase) EnableTask(ctx kratosx.Context, id uint32) error {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if *task.Status != *TaskDisabled {
		return errors.EnableNotDisabledTask()
	}

	if !u.factory.DrySpec(task.Spec) {
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

	if *task.Status != *TaskEnabled {
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
