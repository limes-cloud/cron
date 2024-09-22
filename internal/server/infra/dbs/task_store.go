package dbs

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type store struct {
	rdb *redis.Client
}

const (
	key = "cron:task:hash:spec"
)

func (s *store) tid(id uint32) string {
	return strconv.Itoa(int(id))
}

func (s *store) set(id uint32, spec string) {
	s.rdb.HSet(context.Background(), key, s.tid(id), spec)
}

func (s *store) get(id uint32) string {
	res, _ := s.rdb.HGet(context.Background(), key, s.tid(id)).Result()
	return res
}

func (s *store) delete(id uint32) {
	s.rdb.HDel(context.Background(), key, s.tid(id))
}
