package scheduler

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// 调度指标（自注册到默认 registry；若配置中心未启用 /metrics 端点则仅注册不暴露）
var (
	// executionTotal 执行意图计数：{task_id, trigger_type, status}
	executionTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cron",
		Name:      "execution_total",
		Help:      "任务执行意图总数（按任务/触发方式/终态）",
	}, []string{"task_id", "trigger_type", "status"})

	// executionDuration 执行耗时（秒）
	executionDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cron",
		Name:      "execution_duration_seconds",
		Help:      "任务执行耗时（秒）",
		Buckets:   []float64{1, 5, 15, 60, 300, 900, 3600},
	}, []string{"task_id", "status"})

	// misfireTotal 错过帧计数：{task_id, policy}
	misfireTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cron",
		Name:      "misfire_total",
		Help:      "错过帧计数（按任务/策略）",
	}, []string{"task_id", "policy"})

	// leaderSwitchTotal 选主切换计数：{instance_id}
	leaderSwitchTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cron",
		Name:      "leader_switch_total",
		Help:      "选主切换次数（按实例）",
	}, []string{"instance_id"})
)

func init() {
	// 容错注册：与 kratosx 配置驱动的指标不冲突（cron_ 前缀独立命名空间）
	for _, c := range []*prometheus.CounterVec{executionTotal, misfireTotal, leaderSwitchTotal} {
		_ = prometheus.Register(c)
	}
	_ = prometheus.Register(executionDuration)
}

func taskIDLabel(taskID uint32) string {
	return strconv.FormatUint(uint64(taskID), 10)
}