package worker

import (
	"github.com/limes-cloud/kratosx"
)

type Repo interface {
	// GetWorkerGroup 获取指定的节点分组
	GetWorkerGroup(ctx kratosx.Context, id uint32) (*WorkerGroup, error)

	// ListWorkerGroup 获取节点分组列表
	ListWorkerGroup(ctx kratosx.Context, req *ListWorkerGroupRequest) ([]*WorkerGroup, uint32, error)

	// CreateWorkerGroup 创建节点分组
	CreateWorkerGroup(ctx kratosx.Context, req *WorkerGroup) (uint32, error)

	// UpdateWorkerGroup 更新节点分组
	UpdateWorkerGroup(ctx kratosx.Context, req *WorkerGroup) error

	// DeleteWorkerGroup 删除节点分组
	DeleteWorkerGroup(ctx kratosx.Context, id uint32) error

	// GetWorker 获取指定的节点信息
	GetWorker(ctx kratosx.Context, id uint32) (*Worker, error)

	// ListWorker 获取节点信息列表
	ListWorker(ctx kratosx.Context, req *ListWorkerRequest) ([]*Worker, uint32, error)

	// CreateWorker 创建节点信息
	CreateWorker(ctx kratosx.Context, req *Worker) (uint32, error)

	// UpdateWorker 更新节点信息
	UpdateWorker(ctx kratosx.Context, req *Worker) error

	// UpdateWorkerStatus 更新节点信息状态
	UpdateWorkerStatus(ctx kratosx.Context, id uint32, status bool) error

	// DeleteWorker 删除节点信息
	DeleteWorker(ctx kratosx.Context, id uint32) error

	// GetWorkerByIp 获取指定的节点信息
	GetWorkerByIp(ctx kratosx.Context, ip string) (*Worker, error)

	// GetWorkerByGroupId 获取随机节点
	GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*Worker, error)
}
