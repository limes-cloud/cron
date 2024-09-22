package service

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/repository"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Log struct {
	conf *conf.Config
	repo repository.Log
}

func NewLog(
	conf *conf.Config,
	repo repository.Log,
) *Log {
	return &Log{
		conf: conf,
		repo: repo,
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

// ListLog 获取日志信息列表
func (l *Log) ListLog(ctx kratosx.Context, req *types.ListLogRequest) ([]*entity.Log, uint32, error) {
	list, total, err := l.repo.ListLog(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}
