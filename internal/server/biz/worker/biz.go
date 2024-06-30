package worker

import (
	"github.com/limes-cloud/kratosx"
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

// GetWorkerGroup 获取指定的节点分组
func (u *UseCase) GetWorkerGroup(ctx kratosx.Context, req *GetWorkerGroupRequest) (*WorkerGroup, error) {
	var (
		res *WorkerGroup
		err error
	)

	if req.Id != nil {
		res, err = u.repo.GetWorkerGroup(ctx, *req.Id)
	} else {
		return nil, errors.ParamsError()
	}

	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListWorkerGroup 获取节点分组列表
func (u *UseCase) ListWorkerGroup(ctx kratosx.Context, req *ListWorkerGroupRequest) ([]*WorkerGroup, uint32, error) {
	list, total, err := u.repo.ListWorkerGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateWorkerGroup 创建节点分组
func (u *UseCase) CreateWorkerGroup(ctx kratosx.Context, req *WorkerGroup) (uint32, error) {
	id, err := u.repo.CreateWorkerGroup(ctx, req)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateWorkerGroup 更新节点分组
func (u *UseCase) UpdateWorkerGroup(ctx kratosx.Context, req *WorkerGroup) error {
	if err := u.repo.UpdateWorkerGroup(ctx, req); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteWorkerGroup 删除节点分组
func (u *UseCase) DeleteWorkerGroup(ctx kratosx.Context, id uint32) error {
	if err := u.repo.DeleteWorkerGroup(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetWorker 获取指定的节点信息
func (u *UseCase) GetWorker(ctx kratosx.Context, req *GetWorkerRequest) (*Worker, error) {
	var (
		res *Worker
		err error
	)

	if req.Id != nil {
		res, err = u.repo.GetWorker(ctx, *req.Id)
	} else if req.Ip != nil {
		res, err = u.repo.GetWorkerByIp(ctx, *req.Ip)
	} else {
		return nil, errors.ParamsError()
	}

	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListWorker 获取节点信息列表
func (u *UseCase) ListWorker(ctx kratosx.Context, req *ListWorkerRequest) ([]*Worker, uint32, error) {
	list, total, err := u.repo.ListWorker(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateWorker 创建节点信息
func (u *UseCase) CreateWorker(ctx kratosx.Context, req *Worker) (uint32, error) {
	req.Status = proto.Bool(false)
	id, err := u.repo.CreateWorker(ctx, req)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateWorker 更新节点信息
func (u *UseCase) UpdateWorker(ctx kratosx.Context, req *Worker) error {
	worker, err := u.repo.GetWorker(ctx, req.Id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if worker.Status != nil && *worker.Status {
		req.Status = nil
	}
	if err := u.repo.UpdateWorker(ctx, req); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// UpdateWorkerStatus 更新节点信息状态
func (u *UseCase) UpdateWorkerStatus(ctx kratosx.Context, id uint32, status bool) error {
	if status {
		worker, err := u.repo.GetWorker(ctx, id)
		if err != nil {
			return errors.GetError(err.Error())
		}

		if err := u.factory.CheckIP(ctx, worker.Ip); err != nil {
			ctx.Logger().Errorw("check ip error", err.Error())
			return errors.WorkerNotAvailableError()
		}
	}

	if err := u.repo.UpdateWorkerStatus(ctx, id, status); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteWorker 删除节点信息
func (u *UseCase) DeleteWorker(ctx kratosx.Context, id uint32) error {
	worker, err := u.repo.GetWorker(ctx, id)
	if err != nil {
		return errors.GetError(err.Error())
	}
	if worker.Status != nil && *worker.Status {
		return errors.DeleteNotDisabledWorkerError()
	}

	if err := u.repo.DeleteWorker(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}
