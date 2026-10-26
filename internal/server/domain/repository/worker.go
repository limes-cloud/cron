package repository

import (
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Worker interface {
	// GetWorkerGroup 获取指定的节点分组
	GetWorkerGroup(ctx kratosx.Context, id uint32) (*entity.WorkerGroup, error)

	// ListWorkerGroup 获取节点分组列表
	ListWorkerGroup(ctx kratosx.Context, req *types.ListWorkerGroupRequest) ([]*entity.WorkerGroup, uint32, error)

	// CreateWorkerGroup 创建节点分组
	CreateWorkerGroup(ctx kratosx.Context, req *entity.WorkerGroup) (uint32, error)

	// UpdateWorkerGroup 更新节点分组
	UpdateWorkerGroup(ctx kratosx.Context, req *entity.WorkerGroup) error

	// DeleteWorkerGroup 删除节点分组
	DeleteWorkerGroup(ctx kratosx.Context, id uint32) error

	// GetWorker 获取指定的节点信息
	GetWorker(ctx kratosx.Context, id uint32) (*entity.Worker, error)

	// ListWorker 获取节点信息列表
	ListWorker(ctx kratosx.Context, req *types.ListWorkerRequest) ([]*entity.Worker, uint32, error)

	// CreateWorker 创建节点信息
	CreateWorker(ctx kratosx.Context, req *entity.Worker) (uint32, error)

	// UpdateWorker 更新节点信息
	UpdateWorker(ctx kratosx.Context, req *entity.Worker) error

	// UpdateWorkerStatus 更新节点信息状态
	UpdateWorkerStatus(ctx kratosx.Context, id uint32, status bool) error

	// DeleteWorker 删除节点信息
	DeleteWorker(ctx kratosx.Context, id uint32) error

	// GetWorkerByIp 获取指定的节点信息
	GetWorkerByIp(ctx kratosx.Context, ip string) (*entity.Worker, error)

	// GetWorkerByGroupId 获取随机节点（仅返回授权给 app 的节点；无授权记录的节点对任何 app 可用）
	GetWorkerByGroupId(ctx kratosx.Context, app string, id uint32) (*entity.Worker, error)

	// ListWorkerApps 获取节点的可用应用白名单（空 = 所有应用可用）
	ListWorkerApps(ctx kratosx.Context, workerID uint32) ([]string, error)

	// ReplaceWorkerApps 替换节点的可用应用白名单（事务；空列表 = 清空授权 = 所有应用可用）
	ReplaceWorkerApps(ctx kratosx.Context, workerID uint32, apps []string) error

	// ListEnabledWorkers 获取启用的节点（探测候选，含授权白名单）
	ListEnabledWorkers(ctx kratosx.Context) ([]*entity.Worker, error)

	// UpdateProbeResult 更新节点探测结果（last_check_at 恒刷新；成功时刷新 last_heartbeat 并清空 last_error）
	UpdateProbeResult(ctx kratosx.Context, id uint32, checkAt int64, heartbeatAt *int64, errMsg string) error

}
