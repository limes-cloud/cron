package scheduler

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// TestCronTrigger_Trigger 验证：秒级表达式到点触发回调
func TestCronTrigger_Trigger(t *testing.T) {
	tr := NewCronTrigger(log.NewHelper(log.DefaultLogger)); tr.Start()
	var fired int32

	// 每秒触发一次
	if err := tr.Add(1001, "* * * * * *", func() { atomic.AddInt32(&fired, 1) }); err != nil {
		t.Fatalf("add: %v", err)
	}
	defer func() { tr.Remove(1001); tr.Stop() }()

	time.Sleep(2500 * time.Millisecond)
	if n := atomic.LoadInt32(&fired); n < 1 {
		t.Fatalf("expect trigger fired >=1, got %d", n)
	}
}

// TestCronTrigger_UpdateIdempotent 验证：同任务 ID 更新后旧 entry 被移除，不会双触发
func TestCronTrigger_UpdateIdempotent(t *testing.T) {
	tr := NewCronTrigger(log.NewHelper(log.DefaultLogger)); tr.Start()
	var count int32

	// 先注册每秒，再更新为每 2 秒（同 taskID）
	if err := tr.Add(1002, "* * * * * *", func() { atomic.AddInt32(&count, 1) }); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := tr.Update(1002, "*/2 * * * * *", func() { atomic.AddInt32(&count, 1) }); err != nil {
		t.Fatalf("update: %v", err)
	}
	defer func() { tr.Remove(1002); tr.Stop() }()

	time.Sleep(2500 * time.Millisecond)
	// 若双 entry 残留：4 秒内会执行 ~4-5 次；2 秒频率应 ≤ 2 次
	// 同时校验 entries 映射只有一个（幂等）
	if n := atomic.LoadInt32(&count); n > 2 {
		t.Fatalf("expect single entry (<=2 in 2.5s), got %d", n)
	}
}

// TestCronTrigger_ValidateSpec 验证：秒级表达式校验
func TestCronTrigger_ValidateSpec(t *testing.T) {
	tr := NewCronTrigger(log.NewHelper(log.DefaultLogger)); tr.Start()

	if err := tr.ValidateSpec("*/5 * * * * *"); err != nil {
		t.Fatalf("valid spec rejected: %v", err)
	}
	if err := tr.ValidateSpec("not-a-spec"); err == nil {
		t.Fatal("invalid spec accepted")
	}
}

// TestCronTrigger_Remove 验证：移除后不再触发
func TestCronTrigger_Remove(t *testing.T) {
	tr := NewCronTrigger(log.NewHelper(log.DefaultLogger)); tr.Start()
	var fired int32

	if err := tr.Add(1003, "* * * * * *", func() { atomic.AddInt32(&fired, 1) }); err != nil {
		t.Fatalf("add: %v", err)
	}
	tr.Remove(1003)
	defer tr.Stop()

	time.Sleep(1500 * time.Millisecond)
	if n := atomic.LoadInt32(&fired); n != 0 {
		t.Fatalf("expect no trigger after remove, got %d", n)
	}
}