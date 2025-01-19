package service

import (
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/crypto"
	"google.golang.org/protobuf/proto"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Task struct {
	conf       *conf.Config
	repo       repository.Task
	worker     repository.Worker
	log        repository.Log
	taskClient repository.TaskClient
}

const (
	maxWaitTime   = 600
	maxRetryCount = 5
)

const (
	ExecTypeGroup = "group"
)

const (
	ExecRunning = "running"
	ExecFail    = "fail"
	ExecCancel  = "cancel"
	ExecSuccess = "success"
)

func NewTask(
	conf *conf.Config,
	repo repository.Task,
	worker repository.Worker,
	log repository.Log,
	taskClient repository.TaskClient,

) *Task {
	task := &Task{
		conf:       conf,
		repo:       repo,
		worker:     worker,
		log:        log,
		taskClient: taskClient,
	}
	repo.StartCron(task.Scheduler)
	return task
}

// GetTaskGroup 获取指定的任务分组
func (t *Task) GetTaskGroup(ctx kratosx.Context, id uint32) (*entity.TaskGroup, error) {
	res, err := t.repo.GetTaskGroup(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListTaskGroup 获取任务分组列表
func (t *Task) ListTaskGroup(ctx kratosx.Context, req *types.ListTaskGroupRequest) ([]*entity.TaskGroup, uint32, error) {
	list, total, err := t.repo.ListTaskGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateTaskGroup 创建任务分组
func (t *Task) CreateTaskGroup(ctx kratosx.Context, tg *entity.TaskGroup) (uint32, error) {
	id, err := t.repo.CreateTaskGroup(ctx, tg)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateTaskGroup 更新任务分组
func (t *Task) UpdateTaskGroup(ctx kratosx.Context, tg *entity.TaskGroup) error {
	if err := t.repo.UpdateTaskGroup(ctx, tg); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteTaskGroup 删除任务分组
func (t *Task) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	err := t.repo.DeleteTaskGroup(ctx, id)
	if err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetTask 获取指定的任务信息
func (t *Task) GetTask(ctx kratosx.Context, id uint32) (*entity.Task, error) {
	res, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListTask 获取任务信息列表
func (t *Task) ListTask(ctx kratosx.Context, req *types.ListTaskRequest) ([]*entity.Task, uint32, error) {
	list, total, err := t.repo.ListTask(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateTask 创建任务信息
func (t *Task) CreateTask(ctx kratosx.Context, task *entity.Task) (uint32, error) {
	// 验证表达式是否正确
	if err := t.repo.ValidateSpec(task.Spec); err != nil {
		return 0, errors.CronSpecError(err.Error())
	}
	task.Version = crypto.MD5ToUpper([]byte(task.Spec + task.ExecType + task.ExecValue))
	task.Status = proto.Bool(false)

	id, err := t.repo.CreateTask(ctx, task)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateTask 更新任务信息
func (t *Task) UpdateTask(ctx kratosx.Context, task *entity.Task) error {
	oldTask, err := t.repo.GetTask(ctx, task.Id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if oldTask.Status != nil && *oldTask.Status {
		return errors.UpdateNotDisableTaskError()
	}

	if err := t.repo.ValidateSpec(task.Spec); err != nil {
		return errors.CronSpecError(err.Error())
	}

	if err := t.repo.UpdateTask(ctx, task); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// UpdateTaskStatus 更新任务信息状态
func (t *Task) UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error {
	task, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}

	return ctx.Transaction(func(ctx kratosx.Context) error {
		if err := t.repo.UpdateTaskStatus(ctx, id, status); err != nil {
			return errors.DatabaseError(err.Error())
		}
		if status {
			if err := t.repo.AddCron(id, task.Spec); err != nil {
				return errors.EnableCronError(err.Error())
			}
		} else {
			if err := t.repo.RemoveCron(id); err != nil {
				return errors.DisableCronError(err.Error())
			}
		}

		return nil
	})
}

// DeleteTask 删除任务信息
func (t *Task) DeleteTask(ctx kratosx.Context, id uint32) error {
	task, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}

	if task.Status != nil && *task.Status {
		return errors.DeleteNotDisabledTaskError()
	}

	if err := t.repo.DeleteTask(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// CancelExecTask 取消执行任务
func (t *Task) CancelExecTask(ctx kratosx.Context, uuid string) error {
	// 从日志获取执行节点
	log, err := t.log.GetLogByUuid(ctx, uuid)
	if err != nil {
		return errors.CancelTaskError(err.Error())
	}

	worker := log.GetWorker()
	if err := t.taskClient.CancelExec(ctx, &types.CancelExecRequest{
		Uuid: uuid,
		IP:   worker.Ip,
		Ak:   worker.Ak,
		Sk:   worker.Sk,
	}); err != nil {
		return errors.CancelTaskError(err.Error())
	}
	return nil
}

// ExecTask 执行任务
func (t *Task) ExecTask(ctx kratosx.Context, id uint32) error {
	go func() {
		if err := t.Scheduler(ctx.Clone(), id, "", true); err != nil {
			ctx.Logger().Errorw("exec task error", err.Error())
		}
	}()
	return nil
}

func (t *Task) Scheduler(ctx kratosx.Context, id uint32, spec string, force bool) error {
	ctx.Logger().Infow("msg", "start scheduler task", "id", id)
	defer ctx.Logger().Infow("msg", "end scheduler task", "id", id)

	// 获取当前调度的任务
	tk, err := t.repo.GetTask(ctx, id)
	if err != nil {
		return err
	}

	// 不是强制调度的情况下，需要判断任务是否正常以及任务的表达式是否符合当前预期
	if !force {
		if tk.Status == nil || !*tk.Status {
			return errors.ExecTaskFailError("task is disabled")
		}
		if tk.Spec != spec {
			return errors.ExecTaskFailError("task is abnormal")
		}
	}

	// 获取需要调度到的机器节点
	var wk *entity.Worker
	if tk.ExecType == ExecTypeGroup {
		wk, err = t.worker.GetWorkerByGroupId(ctx, *tk.WorkerGroupId)
	} else {
		wk, err = t.worker.GetWorker(ctx, *tk.WorkerId)
	}
	if err != nil {
		return errors.ExecTaskFailError("msg", "get worker error", "err", err.Error())
	}

	// 生成调用日志
	uid := crypto.MD5ToUpper([]byte(uuid.NewString()))
	if _, err := t.log.CreateLog(ctx, t.makeLog(uid, tk, wk)); err != nil {
		return errors.ExecTaskFailError("msg", "create task log error", "err", err.Error())
	}

	// 下发任务
	var (
		rer    error
		count  = 0
		status = ExecFail
	)

	// 消息处理函数
	recvFunc := func(msg *types.ExecTaskLog) {
		b, _ := json.Marshal(msg)
		if err := t.log.AppendLogContent(ctx, uid, string(b)); err != nil {
			ctx.Logger().Errorw("msg", "append log error", "err", err.Error())
		}
	}

	for {
		rer = t.taskClient.ExecTask(ctx, &types.ExecTaskRequest{
			Id:            tk.Id,
			Uuid:          uid,
			IP:            wk.Ip,
			Ak:            wk.Ak,
			Sk:            wk.Sk,
			ExecType:      tk.ExecType,
			ExecValue:     tk.ExecValue,
			ExpectCode:    tk.ExpectCode,
			RetryCount:    tk.RetryCount,
			RetryWaitTime: tk.RetryWaitTime,
			MaxExecTime:   tk.MaxExecTime,
		}, recvFunc)
		if rer != nil {
			// 这里除了这两种错误之外还有可能存在其他的连接错误
			// 但是只允许这两种错误的情况下才断开连接
			if rer == io.EOF || errors.IsExecTaskFailError(rer) {
				break
			}

			// 出现非正常的意外错误
			// 获取当前任务的状态，这里重新获取是因为我们可能手动停止任务。
			curStatus, err := t.log.GetLogStatusByUuid(ctx, uid)
			if err != nil {
				ctx.Logger().Warnw("msg", "get task log status error", "err", err.Error())
			}

			// 如果不是正在运行中的，则退出
			if curStatus != ExecRunning {
				rer = errors.ExecTaskFailError("task is not running")
				break
			}

			// 判断是否超过最大的重试次数
			if count < maxRetryCount {
				count++
				// 计算重试等待时间
				wait := maxWaitTime / maxRetryCount * count
				time.Sleep(time.Duration(wait) * time.Second)
				continue
			}
			break
		}
		break
	}

	// 处理执行结果
	if rer == io.EOF || rer == nil {
		rer = nil
		status = ExecSuccess
	}

	// 更新结果到日志
	if err := t.log.UpdateLogStatus(ctx, uid, status); err != nil {
		ctx.Logger().Warnw("msg", "update log status error", "err", err.Error())
	}
	return rer
}

func (t *Task) makeLog(uuid string, task *entity.Task, worker *entity.Worker) *entity.Log {
	tb, _ := json.Marshal(task)
	wb, _ := json.Marshal(worker)
	content := types.ExecTaskLog{
		Type:    "info",
		Content: "start scheduler",
		Time:    time.Now().Unix(),
	}
	cb, _ := json.Marshal(content)
	return &entity.Log{
		Uuid:           uuid,
		TaskId:         task.Id,
		TaskSnapshot:   string(tb),
		WorkerId:       worker.Id,
		WorkerSnapshot: string(wb),
		Content:        string(cb),
		StartAt:        time.Now().Unix(),
		Status:         ExecRunning,
	}
}
