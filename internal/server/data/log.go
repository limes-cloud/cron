package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	biz "github.com/limes-cloud/cron/internal/server/biz/log"
	"github.com/limes-cloud/cron/internal/server/data/model"
)

type logRepo struct {
}

func NewLogRepo() biz.Repo {
	return &logRepo{}
}

// ToLogEntity model转entity
func (r logRepo) ToLogEntity(m *model.Log) *biz.Log {
	e := &biz.Log{}
	_ = valx.Transform(m, e)
	return e
}

// ToLogModel entity转model
func (r logRepo) ToLogModel(e *biz.Log) *model.Log {
	m := &model.Log{}
	_ = valx.Transform(e, m)
	return m
}

// GetLog 获取指定的数据
func (r logRepo) GetLog(ctx kratosx.Context, id uint32) (*biz.Log, error) {
	var (
		m  = model.Log{}
		fs = []string{"*"}
	)
	db := ctx.DB().Select(fs)
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	return r.ToLogEntity(&m), nil
}

// ListLog 获取列表
func (r logRepo) ListLog(ctx kratosx.Context, req *biz.ListLogRequest) ([]*biz.Log, uint32, error) {
	var (
		bs    []*biz.Log
		ms    []*model.Log
		total int64
		fs    = []string{"id", "uuid", "worker_id", "task_id", "start_at", "end_at", "status"}
	)

	db := ctx.DB().Model(model.Log{}).Select(fs).Where("task_id=?", req.TaskId)

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
		bs = append(bs, r.ToLogEntity(m))
	}
	return bs, uint32(total), nil
}

// CreateLog 创建数据
func (r logRepo) CreateLog(ctx kratosx.Context, req *biz.Log) (uint32, error) {
	m := r.ToLogModel(req)
	return m.Id, ctx.DB().Create(m).Error
}

// UpdateLog 更新数据
func (r logRepo) UpdateLog(ctx kratosx.Context, req *biz.Log) error {
	return ctx.DB().Updates(r.ToLogModel(req)).Error
}

// DeleteLog 删除数据
func (r logRepo) DeleteLog(ctx kratosx.Context, ids []uint32) (uint32, error) {
	db := ctx.DB().Where("id in ?", ids).Delete(&model.Log{})
	return uint32(db.RowsAffected), db.Error
}

func (r logRepo) IsRunning(ctx kratosx.Context, uuid string) bool {
	var status string
	if err := ctx.DB().Model(&model.Log{}).
		Select("status").
		Where("uuid=?", uuid).
		Scan(&status).Error; err != nil {
		return false
	}
	return status == biz.ExecRunning
}

func (r logRepo) AppendLogContent(ctx kratosx.Context, uuid string, content string) error {
	return ctx.DB().Model(&model.Log{}).
		Where("uuid=?", uuid).
		UpdateColumn("content", gorm.Expr("CONCAT(content,?)", ","+content)).Error
}

func (r logRepo) UpdateLogStatus(ctx kratosx.Context, uuid string, err error) error {
	status := biz.ExecSuccess
	if err != nil {
		status = biz.ExecFail
	}
	log := &model.Log{
		Status: status,
		EndAt:  time.Now().Unix(),
	}
	return ctx.DB().Where("uuid=?", uuid).Updates(log).Error
}

func (t logRepo) CancelTaskByUUID(ctx kratosx.Context, uuid string) error {
	return ctx.DB().Model(biz.Log{}).Where("uuid=?", uuid).UpdateColumn("status", biz.ExecCancel).Error
}

func (r logRepo) GetTargetIpByUuid(ctx kratosx.Context, uuid string) (string, error) {
	var (
		log    biz.Log
		worker model.Worker
	)
	if err := ctx.DB().First(&log, "uuid=?", uuid).Error; err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(log.WorkerSnapshot), &worker); err != nil {
		return "", err
	}
	return worker.Ip, nil
}
