package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// newTestRedis 连接本地测试 Redis（环境不可用时跳过）
func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skipf("redis not available: %v", err)
	}
	return rdb
}

// TestRedisLock_AcquireNoReentry 验证：锁被持有后，其他执行体无法获取（互斥）
func TestRedisLock_AcquireNoReentry(t *testing.T) {
	rdb := newTestRedis(t)
	lock := NewRedisLock(rdb)
	ctx := context.Background()
	const taskID = 91001

	// 清理
	_ = rdb.Del(ctx, lock.key(taskID)).Err()
	defer func() { _ = rdb.Del(ctx, lock.key(taskID)).Err() }()

	// 第一次获取成功
	ok, err := lock.Acquire(ctx, taskID, "exec-1", 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("first acquire: ok=%v err=%v", ok, err)
	}
	// 第二次（另一个执行体）获取失败 —— 同一任务不并发
	ok, err = lock.Acquire(ctx, taskID, "exec-2", 30*time.Second)
	if err != nil {
		t.Fatalf("second acquire err: %v", err)
	}
	if ok {
		t.Fatal("expect second acquire fail, got ok")
	}
}

// TestRedisLock_RenewOnlyOwner 验证：仅持有者可续期；他人续期无效
func TestRedisLock_RenewOnlyOwner(t *testing.T) {
	rdb := newTestRedis(t)
	lock := NewRedisLock(rdb)
	ctx := context.Background()
	const taskID = 91002

	_ = rdb.Del(ctx, lock.key(taskID)).Err()
	defer func() { _ = rdb.Del(ctx, lock.key(taskID)).Err() }()

	if ok, _ := lock.Acquire(ctx, taskID, "exec-a", 30*time.Second); !ok {
		t.Fatal("acquire fail")
	}
	// 他人续期：不改变锁值（用未持有者尝试后，锁值仍为本人）
	if err := lock.Renew(ctx, taskID, "exec-b", 30*time.Second); err != nil {
		t.Fatalf("renew by other: %v", err)
	}
	// 本人续期成功
	if err := lock.Renew(ctx, taskID, "exec-a", 30*time.Second); err != nil {
		t.Fatalf("renew by owner: %v", err)
	}
	// 续期后锁仍然存在且归本人
	val, err := rdb.Get(ctx, lock.key(taskID)).Result()
	if err != nil || val != "exec-a" {
		t.Fatalf("lock value want exec-a, got %q err=%v", val, err)
	}
}

// TestRedisLock_ReleaseOnlyOwner 验证：仅持有者可释放；他人释放后锁仍在
func TestRedisLock_ReleaseOnlyOwner(t *testing.T) {
	rdb := newTestRedis(t)
	lock := NewRedisLock(rdb)
	ctx := context.Background()
	const taskID = 91003

	_ = rdb.Del(ctx, lock.key(taskID)).Err()
	defer func() { _ = rdb.Del(ctx, lock.key(taskID)).Err() }()

	if ok, _ := lock.Acquire(ctx, taskID, "exec-x", 30*time.Second); !ok {
		t.Fatal("acquire fail")
	}
	// 他人释放：无效
	if err := lock.Release(ctx, taskID, "exec-y"); err != nil {
		t.Fatalf("release by other: %v", err)
	}
	if n, _ := rdb.Exists(ctx, lock.key(taskID)).Result(); n != 1 {
		t.Fatal("expect lock still exists after other release")
	}
	// 本人释放：成功
	if err := lock.Release(ctx, taskID, "exec-x"); err != nil {
		t.Fatalf("release by owner: %v", err)
	}
	if n, _ := rdb.Exists(ctx, lock.key(taskID)).Result(); n != 0 {
		t.Fatal("expect lock removed after owner release")
	}
}

// TestRedisLock_TTLExpire 验证：锁超时自动过期（无心跳场景）
func TestRedisLock_TTLExpire(t *testing.T) {
	rdb := newTestRedis(t)
	lock := NewRedisLock(rdb)
	ctx := context.Background()
	const taskID = 91004

	_ = rdb.Del(ctx, lock.key(taskID)).Err()
	defer func() { _ = rdb.Del(ctx, lock.key(taskID)).Err() }()

	if ok, _ := lock.Acquire(ctx, taskID, "exec-t", 300*time.Millisecond); !ok {
		t.Fatal("acquire fail")
	}
	time.Sleep(500 * time.Millisecond)
	// 锁过期后可重新获取（崩溃恢复的前提）
	ok, err := lock.Acquire(ctx, taskID, "exec-t2", 30*time.Second)
	if err != nil || !ok {
		t.Fatalf("re-acquire after ttl: ok=%v err=%v", ok, err)
	}
}