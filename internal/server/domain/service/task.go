package service

import (
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/crypto"
	"google.golang.org/protobuf/proto"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/scheduler"
	"github.com/limes-cloud/cron/internal/server/types"
)

// Task 任务服务：CRUD + 执行触发（调度主流程已下沉到 scheduler 组件）
// app 纬度隔离：所有 API 操作的 app 来自 X-App-Key（app 层传递），跨应用访问按 404 处理
type Task struct {
	repo  repository.Task
	log   repository.Log
	sched *scheduler.Scheduler
}

func NewTask(
	repo repository.Task,
	logRepo repository.Log,
	sched *scheduler.Scheduler,
) *Task {
	return &Task{
		repo:  repo,
		log:   logRepo,
		sched: sched,
	}
}

// checkApp 校验资源所属应用与请求应用一致（防跨应用越权）；任一为空时跳过（存量兼容）
func checkApp(reqApp, resApp string) error {
	if reqApp != "" && resApp != "" && reqApp != resApp {
		return errors.GetError("not found")
	}
	return nil
}

// GetTaskGroup 获取指定的任务分组
func (t *Task) GetTaskGroup(ctx kratosx.Context, app string, id uint32) (*entity.TaskGroup, error) {
	res, err := t.repo.GetTaskGroup(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	if err := checkApp(app, res.App); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTaskGroup 获取任务分组列表（按 app 隔离）
func (t *Task) ListTaskGroup(ctx kratosx.Context, req *types.ListTaskGroupRequest) ([]*entity.TaskGroup, uint32, error) {
	return t.repo.ListTaskGroup(ctx, req)
}

// CreateTaskGroup 创建任务分组（app 落库）
func (t *Task) CreateTaskGroup(ctx kratosx.Context, app string, tg *entity.TaskGroup) (uint32, error) {
	tg.App = app
	id, err := t.repo.CreateTaskGroup(ctx, tg)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateTaskGroup 更新任务分组（校验归属）
func (t *Task) UpdateTaskGroup(ctx kratosx.Context, app string, tg *entity.TaskGroup) error {
	old, err := t.repo.GetTaskGroup(ctx, tg.Id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, old.App); err != nil {
		return err
	}
	if err := t.repo.UpdateTaskGroup(ctx, tg); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteTaskGroup 删除任务分组（校验归属）
func (t *Task) DeleteTaskGroup(ctx kratosx.Context, app string, id uint32) error {
	old, err := t.repo.GetTaskGroup(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, old.App); err != nil {
		return err
	}
	if err := t.repo.DeleteTaskGroup(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetTask 获取指定的任务信息
func (t *Task) GetTask(ctx kratosx.Context, app string, id uint32) (*entity.Task, error) {
	res, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	if err := checkApp(app, res.App); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTask 获取任务信息列表（按 app 隔离）
func (t *Task) ListTask(ctx kratosx.Context, req *types.ListTaskRequest) ([]*entity.Task, uint32, error) {
	return t.repo.ListTask(ctx, req)
}

// CreateTask 创建任务信息（app 落库）
func (t *Task) CreateTask(ctx kratosx.Context, app string, task *entity.Task) (uint32, error) {
	// 验证表达式是否正确
	if err := t.sched.ValidateSpec(task.Spec); err != nil {
		return 0, errors.CronSpecError(err.Error())
	}
	task.App = app
	task.Version = crypto.MD5ToUpper([]byte(task.Spec + task.ExecType + task.ExecValue))
	task.Status = proto.Bool(false)

	id, err := t.repo.CreateTask(ctx, task)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateTask 更新任务信息（校验归属；约定：仅 disabled 状态可更新，无需同步触发源）
func (t *Task) UpdateTask(ctx kratosx.Context, app string, task *entity.Task) error {
	oldTask, err := t.repo.GetTask(ctx, task.Id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, oldTask.App); err != nil {
		return err
	}
	if oldTask.Status != nil && *oldTask.Status {
		return errors.UpdateNotDisableTaskError()
	}

	if err := t.sched.ValidateSpec(task.Spec); err != nil {
		return errors.CronSpecError(err.Error())
	}

	if err := t.repo.UpdateTask(ctx, task); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// UpdateTaskStatus 更新任务信息状态（校验归属；DB 提交后同步触发源）
func (t *Task) UpdateTaskStatus(ctx kratosx.Context, app string, id uint32, status bool) error {
	task, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, task.App); err != nil {
		return err
	}

	if err := t.repo.UpdateTaskStatus(ctx, id, status); err != nil {
		return errors.UpdateError(err.Error())
	}
	if status {
		if err := t.sched.Add(id, task.Spec); err != nil {
			return errors.EnableCronError(err.Error())
		}
	} else {
		t.sched.Remove(id)
	}
	return nil
}

// DeleteTask 删除任务信息（校验归属）
func (t *Task) DeleteTask(ctx kratosx.Context, app string, id uint32) error {
	task, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, task.App); err != nil {
		return err
	}

	if task.Status != nil && *task.Status {
		return errors.DeleteNotDisabledTaskError()
	}

	if err := t.repo.DeleteTask(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	t.sched.Remove(id)
	return nil
}

// CancelExecTask 取消执行任务（校验归属：跨应用取消按 404）
func (t *Task) CancelExecTask(ctx kratosx.Context, app string, uuid string) error {
	lg, err := t.log.GetLogByUuid(ctx, uuid)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, lg.App); err != nil {
		return err
	}
	if err := t.sched.CancelExec(ctx, uuid); err != nil {
		return errors.CancelTaskError(err.Error())
	}
	return nil
}

// ExecTask 执行任务（手动触发：校验归属；与定时执行共用执行标记锁）
func (t *Task) ExecTask(ctx kratosx.Context, app string, id uint32) error {
	task, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if err := checkApp(app, task.App); err != nil {
		return err
	}
	return t.sched.Trigger(ctx, id)
}