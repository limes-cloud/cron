package biz

import (
	"github.com/limes-cloud/kratosx"
	ktypes "github.com/limes-cloud/kratosx/types"

	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type WorkerFactory interface {
	CheckIP(ctx kratosx.Context, ip string) error
}

type WorkerRepo interface {
	AddWorkerGroup(ctx kratosx.Context, in *WorkerGroup) (uint32, error)
	PageWorkerGroup(ctx kratosx.Context, req *PageWorkerGroupRequest) ([]*WorkerGroup, uint32, error)
	UpdateWorkerGroup(ctx kratosx.Context, c *WorkerGroup) error
	DeleteWorkerGroup(ctx kratosx.Context, id uint32) error

	GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*Worker, error)
	AddWorker(ctx kratosx.Context, in *Worker) (uint32, error)
	GetWorker(ctx kratosx.Context, id uint32) (*Worker, error)
	GetWorkersByTag(ctx kratosx.Context, tag string) ([]*Worker, error)
	PageWorker(ctx kratosx.Context, req *PageWorkerRequest) ([]*Worker, uint32, error)
	UpdateWorker(ctx kratosx.Context, c *Worker) error
	DeleteWorker(ctx kratosx.Context, id uint32) error
}

type WorkerGroup struct {
	ktypes.BaseModel
	Name        string
	Description string `json:"description"`
}

type Worker struct {
	ktypes.BaseModel
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Tag         string `json:"tag"`
	Status      string `json:"status"`
	StopDesc    string `json:"stop_desc"`
	Description string `json:"description"`
}

type PageWorkerGroupRequest struct {
	Page     uint32
	PageSize uint32
}

type PageWorkerRequest struct {
	Page     uint32
	PageSize uint32
	Tag      *string
	Status   *string
	IP       *string
	Name     *string
}

type WorkerUseCase struct {
	config  *conf.Config
	repo    WorkerRepo
	factory WorkerFactory
}

const (
	WorkerEnabled  = "enabled"
	WorkerDisabled = "disabled"
)

// NewWorkerUseCase 创建UseCase实体
func NewWorkerUseCase(config *conf.Config, repo WorkerRepo, factory WorkerFactory) *WorkerUseCase {
	return &WorkerUseCase{config: config, repo: repo, factory: factory}
}

// PageWorkerGroup 获取分页任务分组
func (u *WorkerUseCase) PageWorkerGroup(ctx kratosx.Context, req *PageWorkerGroupRequest) ([]*WorkerGroup, uint32, error) {
	tg, total, err := u.repo.PageWorkerGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return tg, total, nil
}

// AddWorkerGroup 添加任务分组
func (u *WorkerUseCase) AddWorkerGroup(ctx kratosx.Context, tg *WorkerGroup) (uint32, error) {
	id, err := u.repo.AddWorkerGroup(ctx, tg)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateWorkerGroup 删除指定任务分组
func (u *WorkerUseCase) UpdateWorkerGroup(ctx kratosx.Context, tg *WorkerGroup) error {
	if err := u.repo.UpdateWorkerGroup(ctx, tg); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteWorkerGroup 删除指定任务分组
func (u *WorkerUseCase) DeleteWorkerGroup(ctx kratosx.Context, id uint32) error {
	if err := u.repo.DeleteWorkerGroup(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// PageWorker 获取分页工作节点
func (u *WorkerUseCase) PageWorker(ctx kratosx.Context, req *PageWorkerRequest) ([]*Worker, uint32, error) {
	worker, total, err := u.repo.PageWorker(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return worker, total, nil
}

// AddWorker 添加工作节点
func (u *WorkerUseCase) AddWorker(ctx kratosx.Context, worker *Worker) (uint32, error) {
	id, err := u.repo.AddWorker(ctx, worker)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateWorker 删除指定工作节点
func (u *WorkerUseCase) UpdateWorker(ctx kratosx.Context, worker *Worker) error {
	if worker.Status != WorkerDisabled {
		worker.IP = ""
	}
	if err := u.repo.UpdateWorker(ctx, worker); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteWorker 删除指定工作节点
func (u *WorkerUseCase) DeleteWorker(ctx kratosx.Context, id uint32) error {
	worker, err := u.repo.GetWorker(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	if worker.Status != WorkerDisabled {
		return errors.DeleteNotDisabledWorker()
	}
	if err := u.repo.DeleteWorker(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// EnableWorker 启动指定工作节点
func (u *WorkerUseCase) EnableWorker(ctx kratosx.Context, id uint32) error {
	worker, err := u.repo.GetWorker(ctx, id)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}

	if worker.Status != WorkerDisabled {
		return errors.EnableNotDisabledWorker()
	}

	if err := u.factory.CheckIP(ctx, worker.IP); err != nil {
		return err
	}

	worker.Status = WorkerEnabled
	if err := u.repo.UpdateWorker(ctx, worker); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DisableWorker 禁用指定工作节点
func (u *WorkerUseCase) DisableWorker(ctx kratosx.Context, id uint32) error {
	worker := Worker{
		BaseModel: ktypes.BaseModel{ID: id},
		Status:    WorkerDisabled,
	}
	if err := u.repo.UpdateWorker(ctx, &worker); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}
