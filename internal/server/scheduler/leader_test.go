package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// TestLeader_ElectionAndTakeover 验证：A 先选主；A 退出后 B 接管（故障接管）
func TestLeader_ElectionAndTakeover(t *testing.T) {
	rdb := newTestRedis(t)
	store := NewRedisLeaderStore(rdb)
	ctx := context.Background()

	// 清理选主 key
	_ = rdb.Del(ctx, leaderKey).Err()
	defer func() { _ = rdb.Del(ctx, leaderKey).Err() }()

	var aBecome, aLose, bBecome int32

	// 实例 A：先启动（快速心跳便于测试）
	leaderA := NewLeader(store, "instance-a")
	leaderA.interval = 100 * time.Millisecond
	leaderA.ttl = 500 * time.Millisecond
	ctxA, cancelA := context.WithCancel(ctx)
	doneA := make(chan struct{})
	go func() {
		defer close(doneA)
		_ = leaderA.Run(ctxA, func() { atomic.AddInt32(&aBecome, 1) }, func() { atomic.AddInt32(&aLose, 1) })
	}()

	// 等 A 成为 leader
	waitFor(t, 2*time.Second, func() bool { return atomic.LoadInt32(&aBecome) == 1 })
	if cur, _ := store.Current(ctx); cur != "instance-a" {
		t.Fatalf("expect leader instance-a, got %q", cur)
	}

	// 实例 B：启动，此时 A 持有锁，B 不能成为 leader（互斥）
	leaderB := NewLeader(store, "instance-b")
	leaderB.interval = 100 * time.Millisecond
	leaderB.ttl = 500 * time.Millisecond
	ctxB, cancelB := context.WithCancel(ctx)
	defer cancelB()
	go func() { _ = leaderB.Run(ctxB, func() { atomic.AddInt32(&bBecome, 1) }, func() {}) }()

	time.Sleep(500 * time.Millisecond)
	if n := atomic.LoadInt32(&bBecome); n != 0 {
		t.Fatalf("expect B not leader while A alive, got become=%d", n)
	}

	// A 退出（模拟故障）→ B 应在 ttl（500ms）+ interval 内接管
	cancelA()
	<-doneA // 等 A 的 Run 退出（清理）
	waitFor(t, 3*time.Second, func() bool { return atomic.LoadInt32(&bBecome) == 1 })
	if cur, _ := store.Current(ctx); cur != "instance-b" {
		t.Fatalf("expect leader instance-b after takeover, got %q", cur)
	}

	// A 的 onLose 应被调用（失去 leader 时）
	if n := atomic.LoadInt32(&aLose); n != 1 {
		t.Fatalf("expect A lose once, got %d", n)
	}
}

// TestLeader_RenewKeepAlive 验证：leader 持续续期，ttl 内保持身份
func TestLeader_RenewKeepAlive(t *testing.T) {
	rdb := newTestRedis(t)
	store := NewRedisLeaderStore(rdb)
	ctx := context.Background()

	_ = rdb.Del(ctx, leaderKey).Err()
	defer func() { _ = rdb.Del(ctx, leaderKey).Err() }()

	leader := NewLeader(store, "instance-a")
	leader.interval = 100 * time.Millisecond
	leader.ttl = 300 * time.Millisecond
	ctxA, cancelA := context.WithCancel(ctx)
	defer cancelA()

	var become int32
	go func() { _ = leader.Run(ctxA, func() { atomic.AddInt32(&become, 1) }, func() {}) }()

	waitFor(t, 2*time.Second, func() bool { return atomic.LoadInt32(&become) == 1 })
	// 等 2 个 ttl：若有续期则仍持有；无续期早已过期
	time.Sleep(600 * time.Millisecond)
	if cur, _ := store.Current(ctx); cur != "instance-a" {
		t.Fatalf("expect leader still instance-a after renew, got %q", cur)
	}
}

// waitFor 轮询等待条件成立（超时失败）
func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("waitFor timeout")
}