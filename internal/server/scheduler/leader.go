package scheduler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// LeaderStore 选主存储（保证多实例下同一时刻只有一个 leader）
type LeaderStore interface {
	// Acquire 尝试成为 leader（原子）；成为返回 true，他人持有返回 false
	Acquire(ctx context.Context, instanceID string, ttl time.Duration) (bool, error)
	// Renew 续期 leader 身份（仅持有者有效）；返回是否仍持有
	Renew(ctx context.Context, instanceID string, ttl time.Duration) (bool, error)
	// Current 返回当前 leader 实例 ID（无则空串）
	Current(ctx context.Context) (string, error)
}

// RedisLeaderStore LeaderStore 的 Redis 实现（cron:leader，10s TTL）
type RedisLeaderStore struct {
	rdb *redis.Client
}

const leaderKey = "cron:leader"

// renewLeaderScript 续期：仅当持有者是本人时延长
const renewLeaderScript = `
if redis.call('get', KEYS[1]) == ARGV[1] then
    return redis.call('pexpire', KEYS[1], ARGV[2])
end
return 0`

// NewRedisLeaderStore 创建选主存储
func NewRedisLeaderStore(rdb *redis.Client) *RedisLeaderStore {
	return &RedisLeaderStore{rdb: rdb}
}

func (s *RedisLeaderStore) Acquire(ctx context.Context, instanceID string, ttl time.Duration) (bool, error) {
	return s.rdb.SetNX(ctx, leaderKey, instanceID, ttl).Result()
}

func (s *RedisLeaderStore) Renew(ctx context.Context, instanceID string, ttl time.Duration) (bool, error) {
	n, err := s.rdb.Eval(ctx, renewLeaderScript, []string{leaderKey}, instanceID, ttl.Milliseconds()).Int64()
	return n == 1, err
}

func (s *RedisLeaderStore) Current(ctx context.Context) (string, error) {
	return s.rdb.Get(ctx, leaderKey).Result()
}

// Leader 选主组件：Run 阻塞式维护 leader 身份，切换时回调
type Leader struct {
	store      LeaderStore
	instanceID string
	ttl        time.Duration // leader 身份 TTL（默认 10s）
	interval   time.Duration // 抢主/续期检查间隔（默认 3s）
}

// NewLeader 创建选主组件
func NewLeader(store LeaderStore, instanceID string) *Leader {
	return &Leader{
		store:      store,
		instanceID: instanceID,
		ttl:        10 * time.Second,
		interval:   3 * time.Second,
	}
}

// Run 阻塞式选主循环：
//   - 抢主成功 → onBecome（启动调度器）→ 进入续期循环
//   - 续期失败 / ctx 取消 → onLose（停止调度器）→ 回到抢主
//
// 返回时（ctx 取消）当前实例不再是 leader
func (l *Leader) Run(ctx context.Context, onBecome func(), onLose func()) error {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	isLeader := false
	become := func() {
		isLeader = true
		leaderSwitchTotal.WithLabelValues(l.instanceID).Inc()
		onBecome()
	}
	lose := func() {
		if isLeader {
			isLeader = false
			onLose()
		}
	}

	for {
		select {
		case <-ctx.Done():
			lose()
			return nil
		case <-ticker.C:
			if isLeader {
				// 续期：失败（锁被他人抢走/过期）= 失去 leader
				ok, err := l.store.Renew(ctx, l.instanceID, l.ttl)
				if err != nil || !ok {
					lose()
				}
				continue
			}
			// 抢主：成功即成为 leader（故障接管窗口 ≤ ttl + interval）
			ok, err := l.store.Acquire(ctx, l.instanceID, l.ttl)
			if err == nil && ok {
				become()
			}
		}
	}
}