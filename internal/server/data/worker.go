package data

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"
	"google.golang.org/protobuf/proto"

	biz "github.com/limes-cloud/cron/internal/server/biz/worker"
	"github.com/limes-cloud/cron/internal/server/data/model"
)

type workerRepo struct {
}

func NewWorkerRepo() biz.Repo {
	return &workerRepo{}
}

// ToWorkerGroupEntity model转entity
func (r workerRepo) ToWorkerGroupEntity(m *model.WorkerGroup) *biz.WorkerGroup {
	e := &biz.WorkerGroup{}
	_ = valx.Transform(m, e)
	return e
}

// ToWorkerGroupModel entity转model
func (r workerRepo) ToWorkerGroupModel(e *biz.WorkerGroup) *model.WorkerGroup {
	m := &model.WorkerGroup{}
	_ = valx.Transform(e, m)
	return m
}

// GetWorkerGroup 获取指定的数据
func (r workerRepo) GetWorkerGroup(ctx kratosx.Context, id uint32) (*biz.WorkerGroup, error) {
	var (
		m  = model.WorkerGroup{}
		fs = []string{"*"}
	)
	db := ctx.DB().Select(fs)
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	return r.ToWorkerGroupEntity(&m), nil
}

// ListWorkerGroup 获取列表
func (r workerRepo) ListWorkerGroup(ctx kratosx.Context, req *biz.ListWorkerGroupRequest) ([]*biz.WorkerGroup, uint32, error) {
	var (
		bs    []*biz.WorkerGroup
		ms    []*model.WorkerGroup
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(model.WorkerGroup{}).Select(fs)

	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	if req.OrderBy == nil || *req.OrderBy == "" {
		req.OrderBy = proto.String("id")
	}
	if req.Order == nil || *req.Order == "" {
		req.Order = proto.String("asc")
	}
	db = db.Order(fmt.Sprintf("%s %s", *req.OrderBy, *req.Order))
	if *req.OrderBy != "id" {
		db = db.Order("id asc")
	}

	if err := db.Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	for _, m := range ms {
		bs = append(bs, r.ToWorkerGroupEntity(m))
	}
	return bs, uint32(total), nil
}

// CreateWorkerGroup 创建数据
func (r workerRepo) CreateWorkerGroup(ctx kratosx.Context, req *biz.WorkerGroup) (uint32, error) {
	m := r.ToWorkerGroupModel(req)
	return m.Id, ctx.DB().Create(m).Error
}

// UpdateWorkerGroup 更新数据
func (r workerRepo) UpdateWorkerGroup(ctx kratosx.Context, req *biz.WorkerGroup) error {
	return ctx.DB().Updates(r.ToWorkerGroupModel(req)).Error
}

// DeleteWorkerGroup 删除数据
func (r workerRepo) DeleteWorkerGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id=?", id).Delete(&model.WorkerGroup{}).Error
}

// ToWorkerEntity model转entity
func (r workerRepo) ToWorkerEntity(m *model.Worker) *biz.Worker {
	e := &biz.Worker{}
	_ = valx.Transform(m, e)
	return e
}

// ToWorkerModel entity转model
func (r workerRepo) ToWorkerModel(e *biz.Worker) *model.Worker {
	m := &model.Worker{}
	_ = valx.Transform(e, m)
	return m
}

// GetWorkerByIp 获取指定数据
func (r workerRepo) GetWorkerByIp(ctx kratosx.Context, ip string) (*biz.Worker, error) {
	var (
		m  = model.Worker{}
		fs = []string{"*"}
	)
	db := ctx.DB().Select(fs)
	if err := db.Where("ip = ?", ip).First(&m).Error; err != nil {
		return nil, err
	}

	return r.ToWorkerEntity(&m), nil
}

// GetWorker 获取指定的数据
func (r workerRepo) GetWorker(ctx kratosx.Context, id uint32) (*biz.Worker, error) {
	var (
		m  = model.Worker{}
		fs = []string{"*"}
	)
	db := ctx.DB().Select(fs)
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	return r.ToWorkerEntity(&m), nil
}

// ListWorker 获取列表
func (r workerRepo) ListWorker(ctx kratosx.Context, req *biz.ListWorkerRequest) ([]*biz.Worker, uint32, error) {
	var (
		bs    []*biz.Worker
		ms    []*model.Worker
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(model.Worker{}).Preload("Group").Select(fs)

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

	if req.OrderBy == nil || *req.OrderBy == "" {
		req.OrderBy = proto.String("id")
	}
	if req.Order == nil || *req.Order == "" {
		req.Order = proto.String("asc")
	}
	db = db.Order(fmt.Sprintf("%s %s", *req.OrderBy, *req.Order))
	if *req.OrderBy != "id" {
		db = db.Order("id asc")
	}

	if err := db.Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	for _, m := range ms {
		bs = append(bs, r.ToWorkerEntity(m))
	}
	return bs, uint32(total), nil
}

// CreateWorker 创建数据
func (r workerRepo) CreateWorker(ctx kratosx.Context, req *biz.Worker) (uint32, error) {
	m := r.ToWorkerModel(req)
	return m.Id, ctx.DB().Create(m).Error
}

// UpdateWorker 更新数据
func (r workerRepo) UpdateWorker(ctx kratosx.Context, req *biz.Worker) error {
	return ctx.DB().Updates(r.ToWorkerModel(req)).Error
}

// UpdateWorkerStatus 更新数据状态
func (r workerRepo) UpdateWorkerStatus(ctx kratosx.Context, id uint32, status bool) error {
	return ctx.DB().Model(model.Worker{}).Where("id=?", id).Update("status", status).Error
}

// DeleteWorker 删除数据
func (r workerRepo) DeleteWorker(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id in ?", id).Delete(&model.Worker{}).Error
}

func (r workerRepo) GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*biz.Worker, error) {
	var list []*biz.Worker
	if err := ctx.DB().Where("group_id=? and status=true", id).Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("not exist worker")
	}
	x := rand.Intn(len(list))
	return list[x], nil
}
