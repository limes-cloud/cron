package service

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/data"
	tf "github.com/limes-cloud/cron/internal/server/factory/task"
)

type registryFunc func(c *conf.Config, hs *http.Server, gs *grpc.Server)

var registries []registryFunc

func register(fn registryFunc) {
	registries = append(registries, fn)
}

func initFactory() {
	tf.Init(&tf.Repo{
		Task:   data.NewTaskRepo(),
		Log:    data.NewLogRepo(),
		Worker: data.NewWorkerRepo(),
	})
}

func New(c *conf.Config, hs *http.Server, gs *grpc.Server) {
	initFactory()
	for _, registry := range registries {
		registry(c, hs, gs)
	}
}
