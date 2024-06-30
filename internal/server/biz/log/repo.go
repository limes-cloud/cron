package log

import (
	"github.com/limes-cloud/kratosx"
)

type Repo interface {
	// GetLog 获取指定的日志信息
	GetLog(ctx kratosx.Context, id uint32) (*Log, error)

	// ListLog 获取日志信息列表
	ListLog(ctx kratosx.Context, req *ListLogRequest) ([]*Log, uint32, error)

	// CreateLog 创建日志信息
	CreateLog(ctx kratosx.Context, req *Log) (uint32, error)

	// DeleteLog 删除日志信息
	DeleteLog(ctx kratosx.Context, ids []uint32) (uint32, error)

	// IsRunning 任务是否还在执行
	IsRunning(ctx kratosx.Context, uuid string) bool

	// AppendLogContent 追加日志内容
	AppendLogContent(ctx kratosx.Context, uuid string, c string) error

	// UpdateLogStatus 更新日志状态
	UpdateLogStatus(ctx kratosx.Context, uuid string, err error) error

	// CancelTaskByUUID 取消任务通过uuid
	CancelTaskByUUID(ctx kratosx.Context, uuid string) error

	// GetTargetIpByUuid 获取指定任务的执行节点ip
	GetTargetIpByUuid(ctx kratosx.Context, uuid string) (string, error)
}
