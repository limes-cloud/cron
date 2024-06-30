package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/crypto"

	clientV1 "github.com/limes-cloud/cron/api/cron/client/v1"
	apierrors "github.com/limes-cloud/cron/api/cron/errors"
	"github.com/limes-cloud/cron/internal/server/biz/log"
	"github.com/limes-cloud/cron/internal/server/biz/task"
	"github.com/limes-cloud/cron/internal/server/biz/worker"
	"github.com/limes-cloud/cron/internal/server/pkg/cron"
)

const (
	repairSleep   = 180
	maxWaitTime   = 600
	maxRetryCount = 5

	taskLockPrefix   = "cron:task:lock:"
	subscribeChannel = "cron:subscribe:channel"

	add    = "add"
	remove = "remove"
	update = "update"

	ExecTypeGroup = "group"
)

var (
	once sync.Once
	_f   *Factory
)

type Repo struct {
	Task   task.Repo
	Log    log.Repo
	Worker worker.Repo
}

type Factory struct {
	repo   *Repo
	store  *store
	rdb    *redis.Client
	pubSub *redis.PubSub
	cron   *cron.Cron
	log    *klog.Helper
	ctx    kratosx.Context
	wg     sync.WaitGroup
	closer atomic.Bool
}

func GlobalFactory() *Factory {
	return _f
}

func Init(repo *Repo) *Factory {
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
		if t.closer.Load() {
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
			if err := t.scheduler(id, false); err != nil {
				t.log.Errorw("exec task error", err.Error())
				return
			}
		}
	}
}

func (t *Factory) Scheduler(id uint32, force bool) error {
	return t.scheduler(id, force)
}

func (t *Factory) genUid() string {
	return crypto.MD5ToUpper([]byte(uuid.NewString()))
}

func (t *Factory) scheduler(id uint32, force bool) error {
	t.log.Infof("start scheduler task：%d", id)
	tk, err := t.repo.Task.GetTask(t.ctx, id)
	if err != nil {
		return err
	}
	if !force {
		if tk.Status == nil || !*tk.Status {
			return errors.New("task is disabled")
		}
	}

	var (
		wk *worker.Worker
	)
	if tk.ExecType == ExecTypeGroup {
		wk, err = t.repo.Worker.GetWorkerByGroupId(t.ctx, *tk.WorkerGroupId)
	} else {
		wk, err = t.repo.Worker.GetWorker(t.ctx, *tk.WorkerId)
	}
	if err != nil {
		return err
	}

	uid := t.genUid()
	if _, err := t.repo.Log.CreateLog(t.ctx, t.makeLog(uid, tk, wk)); err != nil {
		return err
	}

	// 连接worker
	ctx := context.Background()
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(wk.Ip))
	if err != nil {
		return err
	}

	client := clientV1.NewClientClient(conn)
	req := &clientV1.ExecTaskRequest{
		Id:            tk.Id,
		Type:          tk.ExecType,
		Value:         tk.ExecValue,
		ExpectCode:    tk.ExpectCode,
		RetryCount:    tk.RetryCount,
		RetryWaitTime: tk.RetryWaitTime,
		MaxExecTime:   tk.MaxExecTime,
		Uuid:          uid,
	}
	stream, err := client.ExecTask(ctx, req)
	if err != nil {
		return err
	}

	var (
		rer   error
		msg   *clientV1.ExecTaskReply
		count = 0
	)
	for {
		msg, rer = stream.Recv()
		if rer != nil {
			if rer == io.EOF || apierrors.IsExecTaskFailError(rer) {
				if rer == io.EOF {
					rer = nil
				}
				break
			}

			for t.repo.Log.IsRunning(t.ctx, uid) {
				if count < maxRetryCount {
					count++
					wait := maxWaitTime / maxRetryCount * count
					if stream, err = client.ExecTask(ctx, req); err == nil {
						break
					}
					time.Sleep(time.Duration(wait) * time.Second)
				} else {
					rer = errors.New("exceeded retry attempts")
					if err := t.repo.Log.UpdateLogStatus(t.ctx, uid, rer); err != nil {
						t.log.Errorw("update log status error", err.Error())
					}
					return rer
				}
			}
		}

		b, _ := json.Marshal(msg)
		if err := t.repo.Log.AppendLogContent(t.ctx, uid, string(b)); err != nil {
			t.log.Errorw("append log error", err.Error())
		}
	}
	if err := t.repo.Log.UpdateLogStatus(t.ctx, uid, rer); err != nil {
		t.log.Errorw("update log status error", err.Error())
	}
	return rer
}

func (t *Factory) makeLog(uuid string, task *task.Task, worker *worker.Worker) *log.Log {
	tb, _ := json.Marshal(task)
	wb, _ := json.Marshal(worker)
	start := clientV1.ExecTaskReply{
		Type:    "info",
		Content: "start scheduler",
		Time:    uint32(time.Now().Unix()),
	}
	b, _ := json.Marshal(start)
	return &log.Log{
		Uuid:           uuid,
		TaskId:         task.Id,
		TaskSnapshot:   string(tb),
		WorkerId:       worker.Id,
		WorkerSnapshot: string(wb),
		Content:        string(b),
		StartAt:        time.Now().Unix(),
		Status:         log.ExecRunning,
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
	source := t.repo.Task.GetSpecs(ctx)

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

func (t *Factory) CancelExec(ctx kratosx.Context, uuid string) error {
	ip, err := t.repo.Log.GetTargetIpByUuid(ctx, uuid)
	if err != nil {
		return err
	}
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(ip))
	if err != nil {
		return err
	}
	client := clientV1.NewClientClient(conn)
	if _, err = client.CancelExecTask(ctx, &clientV1.CancelExecTaskRequest{Uuid: uuid}); err != nil {
		return err
	}
	return t.repo.Log.CancelTaskByUUID(ctx, uuid)
}
