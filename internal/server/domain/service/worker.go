package service

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Worker struct {
	conf       *conf.Config
	repo       repository.Worker
	taskClient repository.TaskClient
}

func NewWorker(
	conf *conf.Config,
	repo repository.Worker,
	taskClient repository.TaskClient,
) *Worker {
	worker := &Worker{
		conf:       conf,
		repo:       repo,
		taskClient: taskClient,
	}
	repo.RegistryCheckIP(
		kratosx.MustContext(context.Background()),
		func(ctx kratosx.Context, req *types.CheckWorkerRequest) error {
			return taskClient.CheckHealthy(ctx, req)
		},
	)
	return worker
}

// GetWorkerGroup 获取指定的节点分组
func (w *Worker) GetWorkerGroup(ctx kratosx.Context, id uint32) (*entity.WorkerGroup, error) {
	res, err := w.repo.GetWorkerGroup(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListWorkerGroup 获取节点分组列表
func (w *Worker) ListWorkerGroup(ctx kratosx.Context, req *types.ListWorkerGroupRequest) ([]*entity.WorkerGroup, uint32, error) {
	list, total, err := w.repo.ListWorkerGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateWorkerGroup 创建节点分组
func (w *Worker) CreateWorkerGroup(ctx kratosx.Context, wg *entity.WorkerGroup) (uint32, error) {
	id, err := w.repo.CreateWorkerGroup(ctx, wg)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateWorkerGroup 更新节点分组
func (w *Worker) UpdateWorkerGroup(ctx kratosx.Context, wg *entity.WorkerGroup) error {
	if err := w.repo.UpdateWorkerGroup(ctx, wg); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteWorkerGroup 删除节点分组
func (w *Worker) DeleteWorkerGroup(ctx kratosx.Context, id uint32) error {
	if err := w.repo.DeleteWorkerGroup(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetWorker 获取指定的节点信息
func (w *Worker) GetWorker(ctx kratosx.Context, req *types.GetWorkerRequest) (*entity.Worker, error) {
	var (
		res *entity.Worker
		err error
	)

	if req.Id != nil {
		res, err = w.repo.GetWorker(ctx, *req.Id)
	} else if req.Ip != nil {
		res, err = w.repo.GetWorkerByIp(ctx, *req.Ip)
	} else {
		return nil, errors.ParamsError()
	}

	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListWorker 获取节点信息列表
func (w *Worker) ListWorker(ctx kratosx.Context, req *types.ListWorkerRequest) ([]*entity.Worker, uint32, error) {
	list, total, err := w.repo.ListWorker(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateWorker 创建节点信息
func (w *Worker) CreateWorker(ctx kratosx.Context, worker *entity.Worker) (uint32, error) {
	worker.Status = proto.Bool(false)
	id, err := w.repo.CreateWorker(ctx, worker)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateWorker 更新节点信息
func (w *Worker) UpdateWorker(ctx kratosx.Context, worker *entity.Worker) error {
	if err := w.repo.UpdateWorker(ctx, worker); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// UpdateWorkerStatus 更新节点信息状态
func (w *Worker) UpdateWorkerStatus(ctx kratosx.Context, id uint32, status bool) error {
	if status {
		worker, err := w.repo.GetWorker(ctx, id)
		if err != nil {
			return errors.GetError(err.Error())
		}

		if err := w.repo.CheckIP(ctx, &types.CheckWorkerRequest{
			IP: worker.Ip,
			Ak: worker.Ak,
			Sk: worker.Sk,
		}); err != nil {
			ctx.Logger().Errorw("check ip error", err.Error())
			return errors.WorkerNotAvailableError()
		}
	}

	if err := w.repo.UpdateWorkerStatus(ctx, id, status); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteWorker 删除节点信息
func (w *Worker) DeleteWorker(ctx kratosx.Context, id uint32) error {
	worker, err := w.repo.GetWorker(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if worker.Status != nil && *worker.Status {
		return errors.DeleteNotDisabledWorkerError()
	}

	if err := w.repo.DeleteWorker(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}
