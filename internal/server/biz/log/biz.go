package log

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type UseCase struct {
	conf *conf.Config
	repo Repo
}

func NewUseCase(config *conf.Config, repo Repo) *UseCase {
	return &UseCase{conf: config, repo: repo}
}

// GetLog 获取指定的日志信息
func (u *UseCase) GetLog(ctx kratosx.Context, req *GetLogRequest) (*Log, error) {
	var (
		res *Log
		err error
	)

	if req.Id != nil {
		res, err = u.repo.GetLog(ctx, *req.Id)
	} else {
		return nil, errors.ParamsError()
	}

	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListLog 获取日志信息列表
func (u *UseCase) ListLog(ctx kratosx.Context, req *ListLogRequest) ([]*Log, uint32, error) {
	list, total, err := u.repo.ListLog(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateLog 创建日志信息
func (u *UseCase) CreateLog(ctx kratosx.Context, req *Log) (uint32, error) {
	id, err := u.repo.CreateLog(ctx, req)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// DeleteLog 删除日志信息
func (u *UseCase) DeleteLog(ctx kratosx.Context, ids []uint32) (uint32, error) {
	total, err := u.repo.DeleteLog(ctx, ids)
	if err != nil {
		return 0, errors.DeleteError(err.Error())
	}
	return total, nil
}
