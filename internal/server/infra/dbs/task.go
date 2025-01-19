package dbs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/pkg/cron"
	"github.com/limes-cloud/cron/internal/server/types"
)

const (
	repairSleep = 180

	taskLockPrefix   = "cron:task:lock:"
	subscribeChannel = "cron:subscribe:channel"

	add    = "add"
	remove = "remove"
	update = "update"
)

type Task struct {
	rdb       *redis.Client
	pubSub    *redis.PubSub
	store     *store
	cron      *cron.Cron
	ctx       kratosx.Context
	wg        sync.WaitGroup
	closer    atomic.Bool
	scheduler func(ctx kratosx.Context, id uint32, spec string, force bool) error
	once      sync.Once
}

type message struct {
	ID   uint32 `json:"id"`
	Type string `json:"type"`
	Spec string `json:"spec"`
}

var (
	taskIns  *Task
	taskOnce sync.Once
)

func NewTask() *Task {
	taskOnce.Do(func() {
		ctx := kratosx.MustContext(context.Background())
		location, _ := time.LoadLocation("Asia/Shanghai")

		taskIns = &Task{
			rdb:    ctx.Redis(),
			pubSub: ctx.Redis().Subscribe(ctx, subscribeChannel),
			store:  &store{rdb: ctx.Redis()},
			cron: cron.New(
				cron.WithSeconds(),
				cron.WithLocation(location),
				cron.WithLogger(ctx.Logger()),
			),
			ctx: ctx,
		}
	})
	return taskIns
}

// GetTaskGroup 获取指定的数据
func (t *Task) GetTaskGroup(ctx kratosx.Context, id uint32) (*entity.TaskGroup, error) {
	var (
		taskGroup = entity.TaskGroup{}
		fs        = []string{"*"}
	)
	return &taskGroup, ctx.DB().Select(fs).First(&taskGroup, id).Error
}

// ListTaskGroup 获取列表
func (t *Task) ListTaskGroup(ctx kratosx.Context, req *types.ListTaskGroupRequest) ([]*entity.TaskGroup, uint32, error) {
	var (
		list  []*entity.TaskGroup
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.TaskGroup{}).Select(fs)
	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	return list, uint32(total), db.Find(&list).Error
}

// CreateTaskGroup 创建数据
func (t *Task) CreateTaskGroup(ctx kratosx.Context, tg *entity.TaskGroup) (uint32, error) {
	return tg.Id, ctx.DB().Create(tg).Error
}

// UpdateTaskGroup 更新数据
func (t *Task) UpdateTaskGroup(ctx kratosx.Context, tg *entity.TaskGroup) error {
	return ctx.DB().Updates(tg).Error
}

// DeleteTaskGroup 删除数据
func (t *Task) DeleteTaskGroup(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.TaskGroup{}).Error
}

// GetTask 获取指定的数据
func (t *Task) GetTask(ctx kratosx.Context, id uint32) (*entity.Task, error) {
	var (
		task = entity.Task{}
		fs   = []string{"*"}
	)
	return &task, ctx.DB().Select(fs).
		Preload("Group").
		Preload("Worker").
		Preload("WorkerGroup").
		First(&task, id).Error
}

// ListTask 获取列表
func (t *Task) ListTask(ctx kratosx.Context, req *types.ListTaskRequest) ([]*entity.Task, uint32, error) {
	var (
		list  []*entity.Task
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.Task{}).Select(fs)
	if req.GroupId != nil {
		db = db.Where("group_id = ?", *req.GroupId)
	}
	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}
	if req.Tag != nil {
		db = db.Where("tag = ?", *req.Tag)
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

// CreateTask 创建数据
func (t *Task) CreateTask(ctx kratosx.Context, task *entity.Task) (uint32, error) {
	return task.Id, ctx.DB().Create(task).Error
}

// UpdateTask 更新数据
func (t *Task) UpdateTask(ctx kratosx.Context, task *entity.Task) error {
	return ctx.DB().Updates(task).Error
}

// UpdateTaskStatus 更新数据状态
func (t *Task) UpdateTaskStatus(ctx kratosx.Context, id uint32, status bool) error {
	return ctx.DB().Model(entity.Task{}).Where("id=?", id).Update("status", status).Error
}

// DeleteTask 删除数据
func (t *Task) DeleteTask(ctx kratosx.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.Task{}).Error
}

// AllTaskSpecs 获取当前启用的任务的所有定时时间表达式
func (t *Task) AllTaskSpecs(ctx kratosx.Context) map[uint32]string {
	var (
		list []*entity.Task
		m    = map[uint32]string{}
	)
	ctx.DB().Model(entity.Task{}).
		Where("status=true").
		// 提前加载需要启动的定时任务 -2*repairSleep 防止存在加载空隙2️而没被触发
		Where("start is null or start <= ?", time.Now().Unix()-2*repairSleep).
		Where("end is null or end <= ?", time.Now().Unix()).
		Find(&list)
	for _, item := range list {
		m[item.Id] = item.Spec
	}
	return m
}

// StartCron 启动定时调度程序
func (t *Task) StartCron(scheduler func(ctx kratosx.Context, id uint32, spec string, force bool) error) {
	t.once.Do(func() {
		// 注册调度程序
		t.scheduler = scheduler

		// 启动定时任务
		t.cron.Start()

		// 启动坚定调度指令队列
		t.watch()

		// 启动修复数据
		t.startRepair()
	})
}

// ValidateSpec 检验定时任务表达式
func (t *Task) ValidateSpec(s string) error {
	return t.cron.ValidateSpec(s)
}

// AddCron 添加定时任务
func (t *Task) AddCron(id uint32, spec string) error {
	if t.store.get(id) != "" {
		return errors.New("task id exist")
	}
	if err := t.cron.ValidateSpec(spec); err != nil {
		return err
	}
	return t.publish(t.ctx, &message{ID: id, Spec: spec, Type: add})
}

// CloseCron 关闭定时任务
func (t *Task) CloseCron() {
	t.closer.Store(true)
	t.wg.Wait()
}

// RemoveCron 移除定时任务
func (t *Task) RemoveCron(id uint32) error {
	if t.store.get(id) == "" {
		return nil
	}
	return t.publish(t.ctx, &message{ID: id, Type: remove})
}

// UpdateCron 更新定时任务
func (t *Task) UpdateCron(id uint32, spec string) error {
	if t.store.get(id) == "" {
		return errors.New("id not exist")
	}
	if err := t.cron.ValidateSpec(spec); err != nil {
		return err
	}
	return t.publish(t.ctx, &message{ID: id, Spec: spec, Type: update})
}

// publish 发送定时任务变更信息
func (t *Task) publish(ctx kratosx.Context, msg *message) error {
	b, _ := json.Marshal(msg)
	count, err := t.rdb.Publish(ctx, subscribeChannel, string(b)).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("not exist subscribe")
	}

	ctx.Logger().Info("add publish count ", count)
	return nil
}

// parseMsg 解析信息数据
func (t *Task) parseMsg(msg string) (*message, error) {
	res := &message{}
	if err := json.Unmarshal([]byte(msg), res); err != nil {
		return nil, err
	}
	if res.ID == 0 {
		return nil, errors.New("msg format error")
	}
	return res, nil
}

// watch 监听定时任务变更信息
func (t *Task) watch() {
	go func() {
		ctx := kratosx.MustContext(context.Background())

		for {
			source, err := t.pubSub.ReceiveMessage(ctx)
			if err != nil {
				t.ctx.Logger().Errorw("msg", "receive message fail", "err", err.Error())
			}

			msg, err := t.parseMsg(source.Payload)
			if err != nil {
				t.ctx.Logger().Errorw("msg", "parse message fail", "err", err.Error())
			}

			t.exec(msg)
		}
	}()
}

// exec 指定定时任务指令
func (t *Task) exec(msg *message) {
	switch msg.Type {
	case add:
		t.store.set(msg.ID, msg.Spec)
		_ = t.cron.AddFunc(msg.ID, msg.Spec, t.lockScheduler(msg.ID, msg.Spec))
	case remove:
		t.store.delete(msg.ID)
		t.cron.Remove(msg.ID)
	case update:
		t.store.set(msg.ID, msg.Spec)
		t.cron.Remove(msg.ID)
		_ = t.cron.AddFunc(msg.ID, msg.Spec, t.lockScheduler(msg.ID, msg.Spec))
	}
}

// lockScheduler 分布式调度任务
func (t *Task) lockScheduler(id uint32, spec string) func() {
	return func() {
		// 如果已经关闭服务则停止再调度新的任务
		if t.closer.Load() {
			return
		}

		// 获取服务的下次执行时间，上锁
		entry := t.cron.Entry(id)
		expire := entry.Next.Unix() - entry.Prev.Unix() - 1
		if expire < 0 {
			expire = 1
		}
		if expire > 60 {
			expire = 60
		}
		ep := time.Duration(expire)*time.Second - 100*time.Millisecond

		// 获取锁执行
		key := fmt.Sprintf("%s%d", taskLockPrefix, id)
		if is, _ := t.rdb.SetNX(context.Background(), key, 1, ep).Result(); is {
			// 添加任务正在执行数量
			t.wg.Add(1)
			defer t.wg.Done()

			// 调度任务
			if err := t.scheduler(t.ctx, id, spec, false); err != nil {
				t.ctx.Logger().Errorw("msg", "exec task error", "err", err.Error())
			}
		}
	}
}

// startRepair 启动异步修复
func (t *Task) startRepair() {
	go func() {
		ctx := kratosx.MustContext(context.Background())
		for {
			t.repair(ctx)
			time.Sleep(repairSleep * time.Second)
		}
	}()
}

// repair 从数据库中同步最新的定时任务列表进行修复到内存
// 主要是为了防止再监听redis的消息队列的时候，造成的消息丢失问题导致的内存与数据库定时表达式不一致的问题
// 所以每隔一段时间同步一下数据库与内存的数据
// 当然再具体调度的时候，也需要再具体看看表达式是否还存在，不存在则不需要再调度
func (t *Task) repair(ctx kratosx.Context) {
	specs := t.cron.Specs()
	source := t.AllTaskSpecs(ctx)

	for k, v := range source {
		val, ok := specs[k]
		if !ok {
			_ = t.cron.AddFunc(k, v, t.lockScheduler(k, v))
			continue
		}

		if v != val {
			t.cron.Remove(k)
			_ = t.cron.AddFunc(k, v, t.lockScheduler(k, v))
			continue
		}
	}

	for k := range specs {
		if _, ok := source[k]; !ok {
			t.cron.Remove(k)
		}
	}
}
