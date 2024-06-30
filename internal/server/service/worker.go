package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"

	"github.com/limes-cloud/cron/api/cron/errors"
	pb "github.com/limes-cloud/cron/api/cron/server/worker/v1"
	"github.com/limes-cloud/cron/internal/server/biz/worker"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/data"
	fw "github.com/limes-cloud/cron/internal/server/factory/worker"
)

type WorkerService struct {
	pb.UnimplementedWorkerServer
	uc *worker.UseCase
}

func NewWorkerService(conf *conf.Config) *WorkerService {
	return &WorkerService{
		uc: worker.NewUseCase(conf, data.NewWorkerRepo(), fw.New()),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewWorkerService(c)
		pb.RegisterWorkerHTTPServer(hs, srv)
		pb.RegisterWorkerServer(gs, srv)
	})
}

// GetWorkerGroup 获取指定的节点分组
func (s *WorkerService) GetWorkerGroup(c context.Context, req *pb.GetWorkerGroupRequest) (*pb.GetWorkerGroupReply, error) {
	var (
		in  = worker.GetWorkerGroupRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, err := s.uc.GetWorkerGroup(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.GetWorkerGroupReply{}
	if err := valx.Transform(result, &reply); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// ListWorkerGroup 获取节点分组列表
func (s *WorkerService) ListWorkerGroup(c context.Context, req *pb.ListWorkerGroupRequest) (*pb.ListWorkerGroupReply, error) {
	var (
		in  = worker.ListWorkerGroupRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, total, err := s.uc.ListWorkerGroup(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.ListWorkerGroupReply{Total: total}
	if err := valx.Transform(result, &reply.List); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	return &reply, nil
}

// CreateWorkerGroup 创建节点分组
func (s *WorkerService) CreateWorkerGroup(c context.Context, req *pb.CreateWorkerGroupRequest) (*pb.CreateWorkerGroupReply, error) {
	var (
		in  = worker.WorkerGroup{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	id, err := s.uc.CreateWorkerGroup(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &pb.CreateWorkerGroupReply{Id: id}, nil
}

// UpdateWorkerGroup 更新节点分组
func (s *WorkerService) UpdateWorkerGroup(c context.Context, req *pb.UpdateWorkerGroupRequest) (*pb.UpdateWorkerGroupReply, error) {
	var (
		in  = worker.WorkerGroup{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	if err := s.uc.UpdateWorkerGroup(ctx, &in); err != nil {
		return nil, err
	}

	return &pb.UpdateWorkerGroupReply{}, nil
}

// DeleteWorkerGroup 删除节点分组
func (s *WorkerService) DeleteWorkerGroup(c context.Context, req *pb.DeleteWorkerGroupRequest) (*pb.DeleteWorkerGroupReply, error) {
	err := s.uc.DeleteWorkerGroup(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteWorkerGroupReply{}, nil
}

// GetWorker 获取指定的节点信息
func (s *WorkerService) GetWorker(c context.Context, req *pb.GetWorkerRequest) (*pb.GetWorkerReply, error) {
	var (
		in  = worker.GetWorkerRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, err := s.uc.GetWorker(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.GetWorkerReply{}
	if err := valx.Transform(result, &reply); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// ListWorker 获取节点信息列表
func (s *WorkerService) ListWorker(c context.Context, req *pb.ListWorkerRequest) (*pb.ListWorkerReply, error) {
	var (
		in  = worker.ListWorkerRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, total, err := s.uc.ListWorker(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.ListWorkerReply{Total: total}
	if err := valx.Transform(result, &reply.List); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	return &reply, nil
}

// CreateWorker 创建节点信息
func (s *WorkerService) CreateWorker(c context.Context, req *pb.CreateWorkerRequest) (*pb.CreateWorkerReply, error) {
	var (
		in  = worker.Worker{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	id, err := s.uc.CreateWorker(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &pb.CreateWorkerReply{Id: id}, nil
}

// UpdateWorker 更新节点信息
func (s *WorkerService) UpdateWorker(c context.Context, req *pb.UpdateWorkerRequest) (*pb.UpdateWorkerReply, error) {
	var (
		in  = worker.Worker{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	if err := s.uc.UpdateWorker(ctx, &in); err != nil {
		return nil, err
	}

	return &pb.UpdateWorkerReply{}, nil
}

// UpdateWorkerStatus 更新节点信息状态
func (s *WorkerService) UpdateWorkerStatus(c context.Context, req *pb.UpdateWorkerStatusRequest) (*pb.UpdateWorkerStatusReply, error) {
	return &pb.UpdateWorkerStatusReply{}, s.uc.UpdateWorkerStatus(kratosx.MustContext(c), req.Id, req.Status)
}

// DeleteWorker 删除节点信息
func (s *WorkerService) DeleteWorker(c context.Context, req *pb.DeleteWorkerRequest) (*pb.DeleteWorkerReply, error) {
	err := s.uc.DeleteWorker(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteWorkerReply{}, nil
}
