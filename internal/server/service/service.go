package service

import (
	v1 "github.com/limes-cloud/cron/api/server/v1"
	"github.com/limes-cloud/cron/internal/server/biz"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/data"
	"github.com/limes-cloud/cron/internal/server/factory/task"
	"github.com/limes-cloud/cron/internal/server/factory/worker"
)

type Service struct {
	v1.UnimplementedServiceServer
	conf   *conf.Config
	task   *biz.TaskUseCase
	worker *biz.WorkerUseCase
	log    *biz.LogUseCase
}

func New(conf *conf.Config) *Service {
	repo := struct {
		biz.TaskRepo
		biz.LogRepo
		biz.WorkerRepo
	}{
		TaskRepo:   data.NewTaskRepo(),
		LogRepo:    data.NewLogRepo(),
		WorkerRepo: data.NewWorkerRepo(),
	}
	return &Service{
		conf:   conf,
		task:   biz.NewTaskUseCase(conf, repo.TaskRepo, task.New(repo)),
		worker: biz.NewWorkerUseCase(conf, repo.WorkerRepo, worker.New()),
		log:    biz.NewLogUseCase(conf, repo.LogRepo),
	}
}
