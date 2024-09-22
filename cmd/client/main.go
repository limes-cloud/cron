package main

import (
	"context"
	"flag"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/config"
	_ "go.uber.org/automaxprocs"

	"github.com/limes-cloud/cron/internal/client/app"
	"github.com/limes-cloud/cron/internal/client/conf"
)

func main() {
	path := flag.String("conf", "internal/client/conf/conf.yaml", "config path, eg: -conf config.yaml")
	server := kratosx.New(
		kratosx.Config(file.NewSource(*path)),
		kratosx.RegistrarServer(RegisterServer),
		kratosx.Options(kratos.BeforeStop(func(ctx context.Context) error {
			return nil
		})),
	)

	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func RegisterServer(c config.Config, hs *http.Server, gs *grpc.Server) {
	cfg := &conf.Config{}
	c.ScanWatch("business", func(value config.Value) {
		if err := value.Scan(cfg); err != nil {
			log.Error("business 配置变更失败:" + err.Error())
		} else {
			log.Info("business 配置变更成功")
		}
	})

	app.New(cfg, hs, gs)
}
