package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/library/pool"

	v1 "github.com/limes-cloud/cron/api/client/v1"
)

const (
	allow         = "allow"
	maxCheckCount = 10

	publish     = "cron:worker:publish"
	queuePrefix = "cron:worker:queue:"
)

var (
	_f   *Factory
	once sync.Once
)

type Factory struct {
	cip   chan string
	redis *redis.Client
	pub   *redis.PubSub
}

func New() *Factory {
	once.Do(func() {
		ctx := kratosx.MustContext(context.Background())
		_f = &Factory{
			cip:   make(chan string, maxCheckCount+1),
			redis: ctx.Redis(),
			pub:   ctx.Redis().Subscribe(ctx, publish),
		}
		_f.watch(ctx)
		_f.check(ctx)
	})
	return _f
}

// watch 监听接收发送的检测信号
func (w *Factory) watch(ctx kratosx.Context) {
	_ = ctx.Go(pool.AddRunner(func() error {
		for {
			msg, err := w.pub.ReceiveMessage(ctx)
			if err != nil {
				return err
			}
			w.cip <- msg.Payload
		}
	}))
}

// check 并发监听需要ip
func (w *Factory) check(ctx kratosx.Context) {
	for i := 0; i < maxCheckCount; i++ {
		_ = ctx.Go(pool.AddRunner(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ip := <-w.cip:
					status := allow
					if err := w.ping(ctx, ip); err != nil {
						status = err.Error()
					}
					return w.redis.RPush(ctx, queuePrefix+ip, status).Err()
				}
			}
		}))
	}
}

// healthy 进行节点的健康检查
func (w *Factory) ping(ctx kratosx.Context, ip string) error {
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(ip))
	if err != nil {
		return err
	}
	_, err = v1.NewServiceClient(conn).Healthy(ctx, nil)
	return err
}

// CheckIP 检查指定的ip是否健康
func (w *Factory) CheckIP(ctx kratosx.Context, ip string) error {
	count, err := w.redis.Publish(ctx, publish, ip).Result()
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
				res, err := w.redis.BLPop(ctx, 1*time.Second, queuePrefix+ip).Result()
				if err != nil || len(res) == 0 {
					continue
				}
				ct++
				if res[0] != allow {
					ec <- errors.New(res[0])
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
