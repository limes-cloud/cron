package repository

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/types"
)

type TaskClient interface {
	// ExecTask 执行任务
	ExecTask(ctx kratosx.Context, req *types.ExecTaskRequest, recv func(msg *types.ExecTaskLog)) error

	// CancelExec 取消执行任务
	CancelExec(ctx kratosx.Context, req *types.CancelExecRequest) error

	// CheckHealthy 检查节点健康状况 error
	CheckHealthy(ctx kratosx.Context, req *types.CheckWorkerRequest) error
}
