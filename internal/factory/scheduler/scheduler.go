package scheduler

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/limes-cloud/kratosx/pkg/lock"
)

const (
	scheduleKey      = "cron:schedule:lock"
	scheduleQueueKey = "cron:schedule:queue"
	sleep            = time.Second
)

type scheduler struct {
	redis *redis.Client
}

func Run(redis *redis.Client) {
	lock.New(redis, scheduleKey)
}

func (s *scheduler) Run(ctx context.Context) {
	for {
		// 获取调度权限，这里是进行固定延迟1秒，这样由于启动时间的先后顺序问题，某个主机会固定获取调度权限。
		if s.Acquire(ctx) {
			s.Scheduler(ctx)
		}
		time.Sleep(sleep)
	}
}

// Acquire 获取调度的权限
func (s *scheduler) Acquire(ctx context.Context) bool {
	is, _ := s.redis.SetNX(ctx, scheduleKey, 1, sleep).Result()
	return is
}

func (s *scheduler) Scheduler(ctx context.Context) {
	// 1 获取当前时间之前的所有可执行的任务ID列表
	ids := s.GetTaskIds(ctx)

	// 2 将任务id放进待执行的队列
	for _, id := range ids {
		s.Push(ctx, id)
	}
}

func (s *scheduler) GetTaskIds(ctx context.Context) []int32 {
	return []int32{}
}

func (s *scheduler) Push(ctx context.Context, id int32) {
	s.redis.RPush(ctx, scheduleQueueKey, id)
}
