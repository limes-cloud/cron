package worker

import (
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/types"

	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type UseCase struct {
	config  *conf.Config
	repo    Repo
	factory Factory
}

const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"
)

// NewUseCase 创建UseCase实体
func NewUseCase(config *conf.Config, repo Repo, factory Factory) *UseCase {
	return &UseCase{config: config, repo: repo, factory: factory}
}

// PageWorker 获取分页工作节点
func (u *UseCase) PageWorker(ctx kratosx.Context, req *PageWorkerRequest) ([]*Worker, uint32, error) {
	worker, total, err := u.repo.PageWorker(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return worker, total, nil
}

// AddWorker 添加工作节点
func (u *UseCase) AddWorker(ctx kratosx.Context, worker *Worker) (uint32, error) {
	id, err := u.repo.AddWorker(ctx, worker)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateWorker 删除指定工作节点
func (u *UseCase) UpdateWorker(ctx kratosx.Context, worker *Worker) error {
	if worker.Status != StatusDisabled {
		worker.IP = ""
	}
	if err := u.repo.UpdateWorker(ctx, worker); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteWorker 删除指定工作节点
func (u *UseCase) DeleteWorker(ctx kratosx.Context, id uint32) error {
	worker, err := u.repo.GetWorker(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	if worker.Status != StatusDisabled {
		return errors.DeleteNotDisabledWorker()
	}
	if err := u.repo.DeleteWorker(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// EnableWorker 启动指定工作节点
func (u *UseCase) EnableWorker(ctx kratosx.Context, id uint32) error {
	worker, err := u.repo.GetWorker(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if worker.Status != StatusDisabled {
		return errors.EnableNotDisabledWorker()
	}

	if err := u.factory.CheckIP(ctx, worker.IP); err != nil {
		return err
	}

	worker.Status = StatusEnabled
	if err := u.repo.UpdateWorker(ctx, worker); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DisableWorker 禁用指定工作节点
func (u *UseCase) DisableWorker(ctx kratosx.Context, id uint32) error {
	worker := Worker{
		BaseModel: types.BaseModel{ID: id},
		Status:    StatusDisabled,
	}
	if err := u.repo.UpdateWorker(ctx, &worker); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}
