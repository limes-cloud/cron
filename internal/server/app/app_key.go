package app

import (
	"context"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/api/cron/errors"
)

// currentApp 从请求上下文取应用标识（X-App-Key），缺失即拒绝（app 纬度隔离的入口）。
func currentApp(c context.Context) (string, error) {
	app := kratosx.MustContext(c).AppKey()
	if app == "" {
		return "", errors.ParamsError("缺少应用标识（X-App-Key）")
	}
	return app, nil
}