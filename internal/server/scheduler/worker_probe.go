package scheduler

import (
	"context"
	"time"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/server/types"
)

// workerProbeInterval 节点探测周期
const workerProbeInterval = 30 * time.Second

// workerProbeOnce 探测一次所有启用节点（gRPC 直连）：
// 更新 last_check_at（恒刷新）；成功刷新 last_heartbeat 并清空 last_error，失败记录 last_error。
// 故障时间段回溯：last_check_at 与 last_heartbeat 的差距即连续不可用时长。
func (s *Scheduler) workerProbeOnce() {
	ctx := kratosx.MustContext(context.Background())
	workers, err := s.cfg.WorkerRepo.ListEnabledWorkers(ctx)
	if err != nil {
		s.cfg.Logger.Warnw("msg", "probe: list workers fail", "err", err.Error())
		return
	}
	for _, wk := range workers {
		now := time.Now().Unix()
		probeErr := s.cfg.TaskClient.CheckHealthy(ctx, &types.CheckWorkerRequest{IP: wk.Ip, Ak: wk.Ak, Sk: wk.Sk})
		if probeErr != nil {
			// 失败：仅刷新 last_check_at + last_error（last_heartbeat 保持上次成功时间，即离线起点）
			if err := s.cfg.WorkerRepo.UpdateProbeResult(ctx, wk.Id, now, nil, probeErr.Error()); err != nil {
				s.cfg.Logger.Warnw("msg", "probe: update fail", "id", wk.Id, "err", err.Error())
			}
			s.cfg.Logger.Warnw("msg", "probe: worker unreachable", "id", wk.Id, "ip", wk.Ip, "err", probeErr.Error())
			continue
		}
		// 成功：刷新 last_heartbeat + 清空 last_error
		heartbeat := now
		if err := s.cfg.WorkerRepo.UpdateProbeResult(ctx, wk.Id, now, &heartbeat, ""); err != nil {
			s.cfg.Logger.Warnw("msg", "probe: update fail", "id", wk.Id, "err", err.Error())
		}
		s.cfg.Logger.Debugw("msg", "probe: worker ok", "id", wk.Id, "ip", wk.Ip)
	}
}

// workerProbeLoop 节点探活循环（随 s.ctx 退出；仅 leader 运行——Scheduler 只有 leader 启动）
func (s *Scheduler) workerProbeLoop() {
	ticker := time.NewTicker(workerProbeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.workerProbeOnce()
		}
	}
}