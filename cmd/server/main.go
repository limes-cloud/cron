package main

import (
	"context"
	"flag"
	"log"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/limes-cloud/configure/api/configure/client"
	"github.com/limes-cloud/kratosx"
	_ "go.uber.org/automaxprocs"

	"github.com/limes-cloud/cron/internal/server/app"
	"github.com/limes-cloud/cron/internal/server/initiator/menu"
	"github.com/limes-cloud/cron/internal/server/scheduler"
)

func main() {
	// 本地启动：-conf conf/config.yaml；生产：走配置中心（client.NewFromEnv()）
	path := flag.String("conf", "", "config path, eg: -conf conf/config.yaml (default: configuration center)")
	flag.Parse()

	server := kratosx.New(
		kratosx.WithConfigSource(configSource(path)),
		kratosx.WithRegistrarServer(app.Register),
	)

	ctx := kratosx.MustContext(context.Background())
	// 优雅停机：停止调度器（触发源 + 在办执行等待，见 scheduler.Stop）
	ctx.BeforeStop("stop-scheduler", scheduler.Stop)
	// 注册 cron 菜单到 manager（BeforeStart 回调，manager 不可达时仅日志）
	menu.Register(ctx)

	if err := server.App().Run(); err != nil {
		log.Println("run service fail", err.Error())
	}
}

// configSource 配置源：-conf 指定文件，否则配置中心
func configSource(path *string) kratosconfig.Source {
	if *path != "" {
		return file.NewSource(*path)
	}
	return client.NewFromEnv()
}