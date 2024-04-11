package service

import (
	"context"

	"github.com/limes-cloud/kratosx"

	v1 "github.com/limes-cloud/cron/api/server/v1"
	biz "github.com/limes-cloud/cron/internal/server/biz/worker"
	"github.com/limes-cloud/cron/internal/server/conf"
	data "github.com/limes-cloud/cron/internal/server/data/worker"
)

type WorkerService struct {
	v1.UnimplementedServiceServer
	uc   *biz.UseCase
	conf *conf.Config
}

func NewWorkerService(conf *conf.Config) *WorkerService {
	return &WorkerService{
		conf: conf,
		uc:   biz.NewUseCase(conf, data.NewRepo()),
	}
}

func (s *WorkerService) AddWorker(ctx context.Context, in *v1.AuthRequest) (*v1.AuthReply, error) {
	return s.uc.AddWorker(kratosx.MustContext(ctx), in)
}
