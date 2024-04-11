package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/config"
	_ "go.uber.org/automaxprocs"

	v1 "github.com/limes-cloud/cron/api/client/v1"
	"github.com/limes-cloud/cron/internal/client/conf"
	"github.com/limes-cloud/cron/internal/client/service"
)

func main() {
	path := flag.String("conf", "internal/client/conf/config.yaml", "config path, eg: -conf config.yaml")
	app := kratosx.New(
		kratosx.Config(file.NewSource(*path)),
		kratosx.RegistrarServer(RegisterServer),
	)

	if err := app.Run(); err != nil {
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

	srv := service.New(cfg)
	v1.RegisterServiceServer(gs, srv)
}
