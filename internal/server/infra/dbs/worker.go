package dbs

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Worker struct{}

var (
	workerIns  *Worker
	workerOnce sync.Once
)

func NewWorker() *Worker {
	workerOnce.Do(func() {
		workerIns = &Worker{}
	})
	return workerIns
}

// GetWorkerGroup 获取指定的数据
func (w *Worker) GetWorkerGroup(ctx kratosx.Context, id uint32) (*entity.WorkerGroup, error) {
	var (
		wg = entity.WorkerGroup{}
		fs = []string{"*"}
	)
	return &wg, ctx.DB().Select(fs).First(&wg, id).Error
}

// ListWorkerGroup 获取列表
func (w *Worker) ListWorkerGroup(ctx kratosx.Context, req *types.ListWorkerGroupRequest) ([]*entity.WorkerGroup, uint32, error) {
	var (
		list  []*entity.WorkerGroup
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.WorkerGroup{}).Select(fs)

	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	return list, uint32(total), db.Find(&list).Error
}

// CreateWorkerGroup 创建数据
func (w *Worker) CreateWorkerGroup(ctx kratosx.Context, wg *entity.WorkerGroup) (uint32, error) {
	return wg.Id, ctx.DB().Create(wg).Error
}

// UpdateWorkerGroup 更新数据
func (w *Worker) UpdateWorkerGroup(ctx kratosx.Context, wg *entity.WorkerGroup) error {
	return ctx.DB().Updates(wg).Error
}

// DeleteWorkerGroup 删除数据
func (w *Worker) DeleteWorkerGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Delete(&entity.WorkerGroup{}).Error
}

// GetWorkerByIp 获取指定数据
func (w *Worker) GetWorkerByIp(ctx kratosx.Context, ip string) (*entity.Worker, error) {
	var (
		worker = entity.Worker{}
		fs     = []string{"*"}
	)
	return &worker, ctx.DB().Select(fs).Where("ip = ?", ip).First(&worker).Error
}

// GetWorker 获取指定的数据
func (w *Worker) GetWorker(ctx kratosx.Context, id uint32) (*entity.Worker, error) {
	var (
		worker = entity.Worker{}
		fs     = []string{"*"}
	)
	if err := ctx.DB().Select(fs).First(&worker, id).Error; err != nil {
		return nil, err
	}
	w.loadApps(ctx, []*entity.Worker{&worker})
	return &worker, nil
}

// ListWorker 获取列表
func (w *Worker) ListWorker(ctx kratosx.Context, req *types.ListWorkerRequest) ([]*entity.Worker, uint32, error) {
	var (
		list  []*entity.Worker
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.Worker{}).Preload("Group").Select(fs)

	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}
	if req.Ip != nil {
		db = db.Where("ip = ?", *req.Ip)
	}
	if req.GroupId != nil {
		db = db.Where("group_id = ?", *req.GroupId)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	w.loadApps(ctx, list)
	return list, uint32(total), nil
}

// CreateWorker 创建数据
func (w *Worker) CreateWorker(ctx kratosx.Context, worker *entity.Worker) (uint32, error) {
	return worker.Id, ctx.DB().Create(worker).Error
}

// UpdateWorker 更新数据
func (w *Worker) UpdateWorker(ctx kratosx.Context, worker *entity.Worker) error {
	return ctx.DB().Updates(worker).Error
}

// UpdateWorkerStatus 更新数据状态
func (w *Worker) UpdateWorkerStatus(ctx kratosx.Context, id uint32, status bool) error {
	return ctx.DB().Model(entity.Worker{}).Where("id=?", id).Update("status", status).Error
}

// DeleteWorker 删除数据
func (w *Worker) DeleteWorker(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id in ?", id).Delete(&entity.Worker{}).Error
}

// GetWorkerByGroupId 按 app 授权随机选机：仅返回"授权给该 app（或无授权记录=全可用）"且启用的节点
func (w *Worker) GetWorkerByGroupId(ctx kratosx.Context, app string, id uint32) (*entity.Worker, error) {
	var list []*entity.Worker
	// 无授权记录（wa.id IS NULL）或授权命中该 app 的节点均可被使用
	db := ctx.DB().
		Joins("LEFT JOIN worker_app wa ON wa.worker_id = worker.id AND wa.app = ? AND wa.id IS NOT NULL", app).
		Joins("LEFT JOIN (SELECT worker_id, COUNT(*) AS cnt FROM worker_app GROUP BY worker_id) wac ON wac.worker_id = worker.id").
		Where("worker.group_id=? AND worker.status=true", id).
		Where("wa.worker_id IS NOT NULL OR wac.cnt IS NULL").
		Where("worker.last_heartbeat IS NULL OR worker.last_heartbeat > ?", time.Now().Unix()-90)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("not exist worker authorized for app " + app)
	}
	x := rand.Intn(len(list))
	return list[x], nil
}

// ListWorkerApps 获取节点的可用应用白名单（空 = 所有应用可用）
func (w *Worker) ListWorkerApps(ctx kratosx.Context, workerID uint32) ([]string, error) {
	var list []*entity.WorkerApp
	err := ctx.DB().
		Select("app").
		Where("worker_id = ?", workerID).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	apps := make([]string, 0, len(list))
	for _, wa := range list {
		apps = append(apps, wa.App)
	}
	return apps, nil
}

// ReplaceWorkerApps 替换节点的可用应用白名单（全量替换；空列表 = 清空授权 = 所有应用可用）
func (w *Worker) ReplaceWorkerApps(ctx kratosx.Context, workerID uint32, apps []string) error {
	// 去重
	seen := make(map[string]struct{}, len(apps))
	unique := make([]string, 0, len(apps))
	for _, app := range apps {
		if app == "" {
			continue
		}
		if _, ok := seen[app]; ok {
			continue
		}
		seen[app] = struct{}{}
		unique = append(unique, app)
	}

	return ctx.Transaction(func(tx kratosx.Context) error {
		if err := tx.DB().Where("worker_id = ?", workerID).Delete(&entity.WorkerApp{}).Error; err != nil {
			return err
		}
		for _, app := range unique {
			wa := &entity.WorkerApp{
				WorkerId: workerID,
				App:      app,
			}
			if err := tx.DB().Create(wa).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// loadApps 批量填充节点的授权白名单
func (w *Worker) loadApps(ctx kratosx.Context, workers []*entity.Worker) {
	if len(workers) == 0 {
		return
	}
	ids := make([]uint32, 0, len(workers))
	for _, wk := range workers {
		ids = append(ids, wk.Id)
	}
	var list []*entity.WorkerApp
	if err := ctx.DB().Where("worker_id IN ?", ids).Find(&list).Error; err != nil {
		return
	}
	m := make(map[uint32][]string, len(workers))
	for _, wa := range list {
		m[wa.WorkerId] = append(m[wa.WorkerId], wa.App)
	}
	for _, wk := range workers {
		wk.Apps = m[wk.Id]
	}
}

// ListEnabledWorkers 获取启用的节点（探测候选）
func (w *Worker) ListEnabledWorkers(ctx kratosx.Context) ([]*entity.Worker, error) {
	var list []*entity.Worker
	if err := ctx.DB().Where("status = ?", true).Find(&list).Error; err != nil {
		return nil, err
	}
	w.loadApps(ctx, list)
	return list, nil
}

// UpdateProbeResult 更新节点探测结果：last_check_at 恒刷新；
// 成功（heartbeatAt != nil）时刷新 last_heartbeat 并清空 last_error；失败记录 last_error
// 故障回溯：last_check_at 与 last_heartbeat 的差距即连续不可用时长
func (w *Worker) UpdateProbeResult(ctx kratosx.Context, id uint32, checkAt int64, heartbeatAt *int64, errMsg string) error {
	updates := map[string]interface{}{
		"last_check_at": checkAt,
	}
	if heartbeatAt != nil {
		updates["last_heartbeat"] = *heartbeatAt
		updates["last_error"] = nil
	} else {
		updates["last_error"] = errMsg
	}
	return ctx.DB().Model(&entity.Worker{}).Where("id = ?", id).Updates(updates).Error
}
