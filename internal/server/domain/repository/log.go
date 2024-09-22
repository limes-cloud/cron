package repository

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Log interface {
	// GetLogStatusByUuid 获取指定的日志的状态
	GetLogStatusByUuid(ctx kratosx.Context, uid string) (string, error)

	// GetLog 获取指定的日志信息
	GetLog(ctx kratosx.Context, id uint32) (*entity.Log, error)

	// GetLogByUuid 获取指定的日志信息
	GetLogByUuid(ctx kratosx.Context, uuid string) (*entity.Log, error)

	// ListLog 获取日志信息列表
	ListLog(ctx kratosx.Context, req *types.ListLogRequest) ([]*entity.Log, uint32, error)

	// CreateLog 创建日志信息
	CreateLog(ctx kratosx.Context, req *entity.Log) (uint32, error)

	// DeleteLog 删除日志信息
	DeleteLog(ctx kratosx.Context, id uint32) error

	// AppendLogContent 追加日志内容
	AppendLogContent(ctx kratosx.Context, uuid string, c string) error

	// UpdateLogStatus 更新日志状态
	UpdateLogStatus(ctx kratosx.Context, uuid string, status string) error

	// GetTargetIpByUuid 获取指定任务的执行节点ip
	GetTargetIpByUuid(ctx kratosx.Context, uuid string) (string, error)
}
