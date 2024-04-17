package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/util"

	clientV1 "github.com/limes-cloud/cron/api/client/v1"
	apierrors "github.com/limes-cloud/cron/api/errors"
	"github.com/limes-cloud/cron/internal/server/biz"
	"github.com/limes-cloud/cron/internal/server/pkg/cron"
)

const (
	repairSleep = 180
	maxWaitTime = 10

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
	biz.WorkerRepo
	biz.TaskRepo
	biz.LogRepo
}

type Factory struct {
	repo   Repo
	store  *store
	rdb    *redis.Client
	pubSub *redis.PubSub
	cron   *cron.Cron
	log    *log.Helper
	ctx    kratosx.Context
	wg     sync.WaitGroup
	closer atomic.Bool
}

func GlobalFactory() *Factory {
	return _f
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
			ctx:   ctx,
			log:   ctx.Logger(),
			wg:    sync.WaitGroup{},
			store: &store{rdb: ctx.Redis()},
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

func (t *Factory) DrySpec(s string) bool {
	return t.cron.ValidateSpec(s) == nil
}

func (t *Factory) AddCron(id uint32, spec string) error {
	if t.store.get(id) != "" {
		return errors.New("id exist")
	}
	if err := t.cron.ValidateSpec(spec); err != nil {
		return err
	}
	return t.publish(t.ctx, &message{ID: id, Spec: spec, Type: add})
}

func (t *Factory) Close() {
	t.closer.Store(true)
	t.wg.Wait()
}

func (t *Factory) RemoveCron(id uint32) error {
	if t.store.get(id) == "" {
		return nil
	}
	return t.publish(t.ctx, &message{ID: id, Type: remove})
}

func (t *Factory) UpdateCron(id uint32, spec string) error {
	if t.store.get(id) == "" {
		return errors.New("id not exist")
	}
	if err := t.cron.ValidateSpec(spec); err != nil {
		return err
	}
	return t.publish(t.ctx, &message{ID: id, Spec: spec, Type: update})
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
		if !t.closer.Load() {
			return
		}

		expire := 2 * time.Second
		if strings.Contains(spec, "0/") {
			entry := t.cron.Entry(id)
			expire = time.Duration(entry.Next.Unix()-time.Now().Unix()) * time.Second
		}

		key := fmt.Sprintf("%s%d", taskLockPrefix, id)
		if is, _ := t.rdb.SetNX(context.Background(), key, 1, expire).Result(); is {
			t.wg.Add(1)
			defer t.wg.Done()
			if err := t.scheduler(id); err != nil {
				t.log.Errorw("exec task error", err.Error())
				return
			}
		}
	}
}

func (t *Factory) Scheduler(id uint32) error {
	return t.scheduler(id)
}

func (t *Factory) genUid() string {
	return util.MD5ToUpper([]byte(uuid.NewString()))
}

func (t *Factory) scheduler(id uint32) error {
	t.log.Infof("start scheduler task：%d", id)
	tk, err := t.repo.GetTask(t.ctx, id)
	if err != nil {
		return err
	}
	if tk.Status == biz.TaskDisabled {
		return errors.New("task is disabled")
	}

	var (
		wk *biz.Worker
	)
	if tk.ExecType == biz.ExecTypeGroup {
		wk, err = t.repo.GetWorkerByGroupId(t.ctx, *tk.WorkerGroupId)
	} else {
		wk, err = t.repo.GetWorker(t.ctx, *tk.WorkerId)
	}
	if err != nil {
		return err
	}

	uid := t.genUid()
	if _, err := t.repo.AddLog(t.ctx, t.makeLog(uid, tk, wk)); err != nil {
		return err
	}

	// 连接worker
	ctx := context.Background()
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(wk.IP))
	if err != nil {
		return err
	}

	client := clientV1.NewServiceClient(conn)
	req := &clientV1.ExecTaskRequest{
		Id:            tk.ID,
		Type:          tk.ExecType,
		Value:         tk.ExecValue,
		ExpectCode:    *tk.ExpectCode,
		RetryCount:    *tk.RetryCount,
		RetryWaitTime: *tk.RetryWaitTime,
		MaxExecTime:   *tk.MaxExecTime,
		Uuid:          uid,
	}
	stream, err := client.ExecTask(ctx, req)
	if err != nil {
		return err
	}

	var (
		rer error
		msg *clientV1.ExecTaskReply
	)
	for {
		msg, rer = stream.Recv()
		if rer != nil {
			if rer == io.EOF || apierrors.IsExecTaskFail(rer) {
				if rer == io.EOF {
					rer = nil
				}
				break
			}

			for t.repo.TaskIsRunning(t.ctx, uid) {
				if stream, err = client.ExecTask(ctx, req); err == nil {
					break
				}
				wt := rand.Intn(maxWaitTime)
				time.Sleep(time.Duration(wt) * time.Second)
			}
		}

		b, _ := json.Marshal(msg)
		if err := t.repo.AppendLogContent(t.ctx, uid, string(b)); err != nil {
			t.log.Errorw("append log error", err.Error())
		}
	}
	if err := t.repo.UpdateLogStatus(t.ctx, uid, rer); err != nil {
		t.log.Errorw("update log status error", err.Error())
	}
	return rer
}

func (t *Factory) makeLog(uuid string, task *biz.Task, worker *biz.Worker) *biz.Log {
	tb, _ := json.Marshal(task)
	wb, _ := json.Marshal(worker)
	start := clientV1.ExecTaskReply{
		Type:    "info",
		Content: "start scheduler",
		Time:    uint32(time.Now().Unix()),
	}
	b, _ := json.Marshal(start)
	return &biz.Log{
		Uuid:           uuid,
		TaskId:         task.ID,
		TaskSnapshot:   string(tb),
		WorkerId:       worker.ID,
		WorkerSnapshot: string(wb),
		Content:        string(b),
		Start:          time.Now().Unix(),
		Status:         biz.ExecRunning,
	}
}

func (t *Factory) startRepair() {
	go func() {
		ctx := kratosx.MustContext(context.Background())
		for {
			t.repair(ctx)
			time.Sleep(repairSleep * time.Second)
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