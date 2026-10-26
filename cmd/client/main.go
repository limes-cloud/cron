package main

import (
	"flag"
	"log"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/config"
	_ "go.uber.org/automaxprocs"

	"github.com/limes-cloud/cron/internal/client/app"
	"github.com/limes-cloud/cron/internal/client/conf"
)

func main() {
	path := flag.String("conf", "internal/client/conf/conf.yaml", "config path, eg: -conf config.yaml")
	flag.Parse()

	server := kratosx.New(
		kratosx.WithConfigSource(file.NewSource(*path)),
		kratosx.WithConfigWatch(func(w config.Watcher) {
			w("business", func(value config.Value) {
				cfg := &conf.Config{}
				if err := value.Scan(cfg); err != nil {
					log.Printf("business 配置变更失败: %s", err.Error())
					return
				}
				conf.Set(cfg) // 热更新 client 执行配置（Shell 等）
			})
		}),
		kratosx.WithRegistrarServer(app.Register),
	)

	if err := server.App().Run(); err != nil {
		log.Fatal(err.Error())
	}
}