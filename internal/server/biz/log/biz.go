package log

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/conf"
)

type UseCase struct {
	config *conf.Config
	repo   Repo
}

const (
	StatusEnabled   = "enabled"
	StatusDisabling = "disabling"
	StatusDisabled  = "disabled"
)

// NewUseCase 创建UseCase实体
func NewUseCase(config *conf.Config, repo Repo) *UseCase {
	return &UseCase{config: config, repo: repo}
}

// GetLog 获取分页工作节点
func (u *UseCase) GetLog(ctx kratosx.Context, req *PageLogRequest) ([]*Log, uint32, error) {
	worker, total, err := u.repo.PageLog(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return worker, total, nil
}

// PageLog 获取分页工作节点
func (u *UseCase) PageLog(ctx kratosx.Context, req *PageLogRequest) ([]*Log, uint32, error) {
	worker, total, err := u.repo.PageLog(ctx, req)
	if err != nil {
		return nil, 0, errors.NotFound()
	}
	return worker, total, nil
}
