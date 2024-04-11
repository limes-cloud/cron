package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/biz/task"
	"github.com/limes-cloud/cron/internal/server/biz/worker"
	"github.com/limes-cloud/cron/internal/server/pkg/cron"
)

const (
	taskLockPrefix   = "cron:task:lock:"
	subscribeChannel = "cron:subscribe:channel"

	add    = "add"
	remove = "remove"
	update = "update"
)

var (
	once sync.Once
	_f   *Factory
)

type Repo interface {
	GetSpecs(ctx kratosx.Context) map[uint32]string
	GetTask(ctx kratosx.Context, id uint32) (*task.Task, error)
	GetWorkerByGroupId(ctx kratosx.Context, id uint32) (*worker.Worker, error)
	GetWorker(ctx kratosx.Context, id uint32) (*worker.Worker, error)
}

type Factory struct {
	repo   Repo
	store  *store
	rdb    *redis.Client
	pubSub *redis.PubSub
	cron   *cron.Cron
	log    *log.Helper
	ctx    kratosx.Context
}

func New(repo Repo) *Factory {
	once.Do(func() {
		ctx := kratosx.MustContext(context.Background())
		location, _ := time.LoadLocation("Asia/Shanghai")
		_f = &Factory{
			repo:   repo,
			rdb:    ctx.Redis(),
			pubSub: ctx.Redis().Subscribe(ctx, subscribeChannel),
			cron: cron.New(
				cron.WithSeconds(),
				cron.WithLocation(location),
				cron.WithLogger(ctx.Logger()),
			),
			ctx: ctx,
			log: ctx.Logger(),
		}
		_f.start()
	})
	return _f
}

type message struct {
	ID   uint32 `json:"id"`
	Type string `json:"type"`
	Spec string `json:"spec"`
}

func (t *Factory) AddCron(ctx kratosx.Context, id uint32, spec string) error {
	if t.store.get(id) != "" {
		return errors.New("id exist")
	}
	if err := t.cron.ValidateSpec(spec); err != nil {
		return err
	}
	return t.publish(ctx, &message{ID: id, Spec: spec, Type: add})
}

func (t *Factory) RemoveCron(ctx kratosx.Context, id uint32) error {
	if t.store.get(id) == "" {
		return nil
	}
	return t.publish(ctx, &message{ID: id, Type: remove})
}

func (t *Factory) UpdateCron(ctx kratosx.Context, id uint32, spec string) error {
	if t.store.get(id) == "" {
		return errors.New("id not exist")
	}
	if err := t.cron.ValidateSpec(spec); err != nil {
		return err
	}
	return t.publish(ctx, &message{ID: id, Spec: spec, Type: update})
}

func (t *Factory) start() {
	t.cron.Start()
	t.watch()
	t.startRepair()
}

func (t *Factory) publish(ctx kratosx.Context, msg *message) error {
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

func (t *Factory) parseMsg(msg string) (*message, error) {
	res := &message{}
	if err := json.Unmarshal([]byte(msg), res); err != nil {
		return nil, err
	}
	if res.ID == 0 {
		return nil, errors.New("msg format error")
	}
	return res, nil
}

func (t *Factory) watch() {
	go func() {
		ctx := kratosx.MustContext(context.Background())

		for {
			source, err := t.pubSub.ReceiveMessage(ctx)
			if err != nil {
				t.log.Error("receive message fail ", err.Error())
			}

			msg, err := t.parseMsg(source.Payload)
			if err != nil {
				t.log.Error("parse message fail ", err.Error())
			}

			t.exec(msg)
		}
	}()
}

func (t *Factory) exec(msg *message) {
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

func (t *Factory) lockScheduler(id uint32, spec string) func() {
	return func() {
		expire := 2 * time.Second
		if strings.Contains(spec, "0/") {
			entry := t.cron.Entry(id)
			expire = time.Duration(entry.Next.Unix()-time.Now().Unix()) * time.Second
		}

		key := fmt.Sprintf("%s%d", taskLockPrefix, id)
		if is, _ := t.rdb.SetNX(context.Background(), key, 1, expire).Result(); is {
			// 放回重新执行的队列中进行执行
			if err := t.scheduler(id); err != nil {
			}

		}
	}
}

func (t *Factory) scheduler(id uint32) error {
	t.log.Infof("开始进行调度任务：%d", id)
	tk, err := t.repo.GetTask(t.ctx, id)
	if err != nil {
		return err
	}

	var (
		wk *worker.Worker
	)
	if tk.ExecType == task.ExecTypeGroupWorker {
		wk, err = t.repo.GetWorkerByGroupId(t.ctx, *tk.WorkerGroupId)
	} else {
		wk, err = t.repo.GetWorker(t.ctx, *tk.WorkerId)
	}
	if err != nil {
		return err
	}

	// 连接worker
	return nil
	// 发送执行任务
}

func (t *Factory) startRepair() {
	go func() {
		ctx := kratosx.MustContext(context.Background())
		for {
			t.repair(ctx)
			time.Sleep(10 * time.Second)
		}
	}()
}

func (t *Factory) repair(ctx kratosx.Context) {
	specs := t.cron.Specs()
	source := t.repo.GetSpecs(ctx)

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
