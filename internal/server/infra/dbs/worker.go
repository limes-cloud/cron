package dbs

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/library/pool"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/types"
)

const (
	allow         = "allow"
	maxCheckCount = 10

	publish     = "cron:worker:publish"
	queuePrefix = "cron:worker:queue:"
)

type Worker struct {
	cip   chan string
	redis *redis.Client
	pub   *redis.PubSub
}

var (
	workerIns  *Worker
	workerOnce sync.Once
)

func NewWorker() *Worker {
	workerOnce.Do(func() {
		ctx := kratosx.MustContext(context.Background())
		workerIns = &Worker{
			cip:   make(chan string, maxCheckCount+1),
			redis: ctx.Redis(),
			pub:   ctx.Redis().Subscribe(ctx, publish),
		}
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
	return &worker, ctx.DB().Select(fs).First(&worker, id).Error
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
	return list, uint32(total), db.Find(&list).Error
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

func (w *Worker) GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*entity.Worker, error) {
	var list []*entity.Worker
	if err := ctx.DB().Where("group_id=? and status=true", id).Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("not exist worker")
	}
	x := rand.Intn(len(list))
	return list[x], nil
}

// RegistryCheckIP 注册监听接收发送的检测信号
func (w *Worker) RegistryCheckIP(ctx kratosx.Context, check func(kratosx.Context, *types.CheckWorkerRequest) error) {
	_ = ctx.Go(pool.AddRunner(func() error {
		for {
			msg, err := w.pub.ReceiveMessage(ctx)
			if err != nil {
				return err
			}
			w.cip <- msg.Payload
		}
	}))

	for i := 0; i < maxCheckCount; i++ {
		_ = ctx.Go(pool.AddRunner(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case data := <-w.cip:
					var info types.CheckWorkerRequest
					_ = json.Unmarshal([]byte(data), &info)
					status := allow
					if err := check(ctx, &info); err != nil {
						status = err.Error()
					}
					return w.redis.RPush(ctx, queuePrefix+info.IP, status).Err()
				}
			}
		}))
	}
}

// CheckIP 检查指定的ip是否健康
func (w *Worker) CheckIP(ctx kratosx.Context, req *types.CheckWorkerRequest) error {
	data, _ := json.Marshal(req)
	count, err := w.redis.Publish(ctx, publish, string(data)).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("not exist checker")
	}
	ec := make(chan error, 1)
	go func() {
		ct := 0
		for {
			select {
			case <-ctx.Done():
				ec <- ctx.Err()
				break
			default:
				res, err := w.redis.BLPop(ctx, 1*time.Second, queuePrefix+req.IP).Result()
				if err != nil || len(res) != 2 {
					continue
				}
				ct++
				if res[1] != allow {
					ec <- errors.New(res[1])
					break
				}
				if ct >= int(count) {
					ec <- nil
					break
				}
			}
		}
	}()
	return <-ec
}
