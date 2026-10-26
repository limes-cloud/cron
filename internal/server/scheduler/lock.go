package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock 执行标记锁：同一任务同时只允许一个执行体执行。
// 语义：锁 = “该任务正在执行中”（value 为 execution_id），TTL 由心跳续期维持，
// 与 cron 周期、任务执行时长彻底解耦。
type Lock interface {
	// Acquire 取得执行权；锁被他人持有返回 (false, nil)
	Acquire(ctx context.Context, taskID uint32, executionID string, ttl time.Duration) (bool, error)
	// Renew 续期（仅当 value 是本人时生效，防止他人锁被延长）
	Renew(ctx context.Context, taskID uint32, executionID string, ttl time.Duration) error
	// Release 释放（仅当 value 是本人时删除，防止误删他执行体的锁）
	Release(ctx context.Context, taskID uint32, executionID string) error
}

// RedisLock Lock 的 Redis 实现
type RedisLock struct {
	rdb *redis.Client
}

const (
	lockKeyPrefix = "cron:exec:lock:"

	// renewScript 续期：仅当持有者是本人时延长过期时间
	renewScript = `
if redis.call('get', KEYS[1]) == ARGV[1] then
    return redis.call('pexpire', KEYS[1], ARGV[2])
end
return 0`

	// releaseScript 释放：仅当持有者是本人时删除
	releaseScript = `
if redis.call('get', KEYS[1]) == ARGV[1] then
    return redis.call('del', KEYS[1])
end
return 0`
)

// NewRedisLock 创建 Redis 执行标记锁
func NewRedisLock(rdb *redis.Client) *RedisLock {
	return &RedisLock{rdb: rdb}
}

func (l *RedisLock) key(taskID uint32) string {
	return fmt.Sprintf("%s%d", lockKeyPrefix, taskID)
}

func (l *RedisLock) Acquire(ctx context.Context, taskID uint32, executionID string, ttl time.Duration) (bool, error) {
	return l.rdb.SetNX(ctx, l.key(taskID), executionID, ttl).Result()
}

func (l *RedisLock) Renew(ctx context.Context, taskID uint32, executionID string, ttl time.Duration) error {
	return l.rdb.Eval(ctx, renewScript, []string{l.key(taskID)}, executionID, ttl.Milliseconds()).Err()
}

func (l *RedisLock) Release(ctx context.Context, taskID uint32, executionID string) error {
	return l.rdb.Eval(ctx, releaseScript, []string{l.key(taskID)}, executionID).Err()
}