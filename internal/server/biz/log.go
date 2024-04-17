package biz

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type LogRepo interface {
	AddLog(ctx kratosx.Context, in *Log) (uint32, error)
	GetLog(ctx kratosx.Context, id uint32) (*Log, error)
	PageLog(ctx kratosx.Context, req *PageLogRequest) ([]*Log, uint32, error)
	AppendLogContent(ctx kratosx.Context, uuid string, c string) error
	UpdateLogStatus(ctx kratosx.Context, uuid string, err error) error
	TaskIsRunning(ctx kratosx.Context, uuid string) bool
}

type PageLogRequest struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"page_size"`
	TaskId   uint32 `json:"task_id"`
}

type Log struct {
	ID             uint32 `json:"id"`
	Uuid           string `json:"uuid"`
	WorkerId       uint32 `json:"worker_id"`
	WorkerSnapshot string `json:"worker_snapshot"`
	TaskId         uint32 `json:"task_id"`
	TaskSnapshot   string `json:"task_snapshot"`
	Start          int64  `json:"start"`
	End            int64  `json:"end"`
	Content        string `json:"content"`
	Status         string `json:"status"`
}

type LogUseCase struct {
	config *conf.Config
	repo   LogRepo
}

const (
	ExecRunning = "running"
	ExecFail    = "fail"
	ExecCancel  = "cancel"
	ExecSuccess = "success"
)

// NewLogUseCase 创建UseCase实体
func NewLogUseCase(config *conf.Config, repo LogRepo) *LogUseCase {
	return &LogUseCase{config: config, repo: repo}
}

// GetLog 获取分页工作节点
func (u *LogUseCase) GetLog(ctx kratosx.Context, id uint32) (*Log, error) {
	log, err := u.repo.GetLog(ctx, id)
	if err != nil {
		return nil, errors.NotFound()
	}
	return log, nil
}

// PageLog 获取分页工作节点
func (u *LogUseCase) PageLog(ctx kratosx.Context, req *PageLogRequest) ([]*Log, uint32, error) {
	log, total, err := u.repo.PageLog(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return log, total, nil
}
