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
	AllWorkerGroup(ctx kratosx.Context) ([]*WorkerGroup, error)
	UpdateWorkerGroup(ctx kratosx.Context, c *WorkerGroup) error
	DeleteWorkerGroup(ctx kratosx.Context, id uint32) error

	GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*Worker, error)
	AddWorker(ctx kratosx.Context, in *Worker) (uint32, error)
	GetWorker(ctx kratosx.Context, id uint32) (*Worker, error)
	PageWorker(ctx kratosx.Context, req *PageWorkerRequest) ([]*Worker, uint32, error)
	UpdateWorker(ctx kratosx.Context, c *Worker) error
	DeleteWorker(ctx kratosx.Context, id uint32) error
	EnableWorker(ctx kratosx.Context, id uint32) error
	DisableWorker(ctx kratosx.Context, id uint32) error
}

type WorkerGroup struct {
	ktypes.BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Worker struct {
	ktypes.BaseModel
	GroupId     *uint32      `json:"group_id"`
	Name        string       `json:"name"`
	IP          string       `json:"ip"`
	Status      *bool        `json:"status"`
	Description string       `json:"description"`
	Group       *WorkerGroup `json:"group"`
}

type PageWorkerGroupRequest struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"page_size"`
}

type PageWorkerRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"page_size"`
	Name     *string `json:"name"`
	GroupId  *uint32 `json:"group_id"`
}

type WorkerUseCase struct {
	config  *conf.Config
	repo    WorkerRepo
	factory WorkerFactory
}

const (
	WorkerRunning  = "running"
	WorkerEnabled  = true
	WorkerDisabled = false
)

// NewWorkerUseCase 创建UseCase实体
func NewWorkerUseCase(config *conf.Config, repo WorkerRepo, factory WorkerFactory) *WorkerUseCase {
	return &WorkerUseCase{config: config, repo: repo, factory: factory}
}

// AllWorkerGroup 获取全部节点分组
func (u *WorkerUseCase) AllWorkerGroup(ctx kratosx.Context) ([]*WorkerGroup, error) {
	tg, err := u.repo.AllWorkerGroup(ctx)
	if err != nil {
		return nil, errors.NotFound()
	}
	return tg, nil
}

// AddWorkerGroup 添加节点分组
func (u *WorkerUseCase) AddWorkerGroup(ctx kratosx.Context, tg *WorkerGroup) (uint32, error) {
	id, err := u.repo.AddWorkerGroup(ctx, tg)
	if err != nil {
		return 0, errors.DatabaseFormat(err.Error())
	}
	return id, nil
}

// UpdateWorkerGroup 删除指定节点分组
func (u *WorkerUseCase) UpdateWorkerGroup(ctx kratosx.Context, tg *WorkerGroup) error {
	if err := u.repo.UpdateWorkerGroup(ctx, tg); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DeleteWorkerGroup 删除指定节点分组
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

// UpdateWorker 更新指定工作节点
func (u *WorkerUseCase) UpdateWorker(ctx kratosx.Context, worker *Worker) error {
	old, err := u.repo.GetWorker(ctx, worker.ID)
	if err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	if *old.Status != WorkerDisabled {
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
	if *worker.Status != WorkerDisabled {
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

	if *worker.Status != WorkerDisabled {
		return errors.EnableNotDisabledWorker()
	}

	if err := u.factory.CheckIP(ctx, worker.IP); err != nil {
		ctx.Logger().Errorw("check ip error", err.Error())
		return errors.WorkerNotAvailable()
	}

	if err := u.repo.EnableWorker(ctx, worker.ID); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}

// DisableWorker 禁用指定工作节点
func (u *WorkerUseCase) DisableWorker(ctx kratosx.Context, id uint32) error {
	if err := u.repo.DisableWorker(ctx, id); err != nil {
		return errors.DatabaseFormat(err.Error())
	}
	return nil
}
