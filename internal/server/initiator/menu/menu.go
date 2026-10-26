package menu

import (
	"context"

	"github.com/limes-cloud/kratosx"
	managerMenu "github.com/limes-cloud/manager/api/menu"
)

const registerModeAppend = "append"

// Register 注册 cron 菜单到 manager（幂等：append 模式，服务启动后执行）
func Register(ctx kratosx.Context) {
	conn, err := ctx.GrpcConn("Manager")
	if err != nil {
		ctx.Logger().Errorw("msg", "get manager grpc conn error", "err", err.Error())
		return
	}

	reply, err := managerMenu.NewMenuClient(conn).RegisterMenu(context.Background(), &managerMenu.RegisterMenuRequest{
		AppKeyword: appKeyword,
		Mode:       registerModeAppend,
		Menus:      menus,
	})
	if err != nil {
		ctx.Logger().Errorw("msg", "register cron menu error", "err", err.Error())
		return
	}

	ctx.Logger().Infow(
		"msg", "register cron menu success",
		"created", reply.Created,
		"updated", reply.Updated,
		"deleted", reply.Deleted,
	)
}