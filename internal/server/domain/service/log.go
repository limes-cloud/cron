package service

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Log struct {
	repo     repository.Log
	execRepo repository.Execution
}

func NewLog(
	repo repository.Log,
	execRepo repository.Execution,
) *Log {
	return &Log{
		repo:     repo,
		execRepo: execRepo,
	}
}

// GetLog 获取指定的日志信息
func (l *Log) GetLog(ctx kratosx.Context, id uint32) (*entity.Log, error) {
	log, err := l.repo.GetLog(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return log, nil
}

// ListExecution 获取执行历史列表（按 app 隔离）
func (l *Log) ListExecution(ctx kratosx.Context, req *types.ListExecutionRequest) ([]*entity.Execution, uint32, error) {
	return l.execRepo.ListExecution(ctx, req)
}

// ListLog 获取日志信息列表
func (l *Log) ListLog(ctx kratosx.Context, req *types.ListLogRequest) ([]*entity.Log, uint32, error) {
	list, total, err := l.repo.ListLog(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}
