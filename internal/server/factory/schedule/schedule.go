package schedule

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/limes-cloud/kratosx"
)

const (
	scheduleLockKey  = "cron:schedule:lock"
	scheduleQueueKey = "cron:schedule:queue"
	sleep            = time.Second
)

type Scheduler interface {
	Schedule(ctx kratosx.Context)
}

type scheduler struct {
}

func New(redis *redis.Client) Scheduler {
	return &scheduler{}
}

func (s *scheduler) Schedule(ctx kratosx.Context) {}
