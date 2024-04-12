package service

import (
	"context"

	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/limes-cloud/cron/api/errors"
	v1 "github.com/limes-cloud/cron/api/server/v1"
	"github.com/limes-cloud/cron/internal/server/biz"
)

func (s *Service) PageWorkerGroup(ctx context.Context, in *v1.PageWorkerGroupRequest) (*v1.PageWorkerGroupReply, error) {
	var req biz.PageWorkerGroupRequest
	if err := util.Transform(in, &req); err != nil {
		return nil, errors.Transform()
	}

	list, total, err := s.worker.PageWorkerGroup(kratosx.MustContext(ctx), &req)
	if err != nil {
		return nil, err
	}

	reply := v1.PageWorkerGroupReply{Total: total}
	if err := util.Transform(list, &reply.List); err != nil {
		return nil, errors.Transform()
	}

	return &reply, nil
}

func (s *Service) AddWorkerGroup(ctx context.Context, in *v1.AddWorkerGroupRequest) (*v1.AddWorkerGroupReply, error) {
	wk := biz.WorkerGroup{}
	if err := util.Transform(in, &wk); err != nil {
		return nil, errors.TransformFormat(err.Error())
	}

	id, err := s.worker.AddWorkerGroup(kratosx.MustContext(ctx), &wk)
	if err != nil {
		return nil, err
	}

	return &v1.AddWorkerGroupReply{Id: id}, nil
}

func (s *Service) UpdateWorkerGroup(ctx context.Context, in *v1.UpdateWorkerGroupRequest) (*emptypb.Empty, error) {
	wk := biz.WorkerGroup{}
	if err := util.Transform(in, &wk); err != nil {
		return nil, errors.TransformFormat(err.Error())
	}

	return nil, s.worker.UpdateWorkerGroup(kratosx.MustContext(ctx), &wk)
}

func (s *Service) DeleteWorkerGroup(ctx context.Context, in *v1.DeleteWorkerGroupRequest) (*emptypb.Empty, error) {
	return nil, s.worker.DeleteWorkerGroup(kratosx.MustContext(ctx), in.Id)
}

func (s *Service) PageWorker(ctx context.Context, in *v1.PageWorkerRequest) (*v1.PageWorkerReply, error) {
	var req biz.PageWorkerRequest
	if err := util.Transform(in, &req); err != nil {
		return nil, errors.Transform()
	}

	list, total, err := s.worker.PageWorker(kratosx.MustContext(ctx), &req)
	if err != nil {
		return nil, err
	}

	reply := v1.PageWorkerReply{Total: total}
	if err := util.Transform(list, &reply.List); err != nil {
		return nil, errors.Transform()
	}

	return &reply, nil
}

func (s *Service) AddWorker(ctx context.Context, in *v1.AddWorkerRequest) (*v1.AddWorkerReply, error) {
	wk := biz.Worker{}
	if err := util.Transform(in, &wk); err != nil {
		return nil, errors.TransformFormat(err.Error())
	}

	id, err := s.worker.AddWorker(kratosx.MustContext(ctx), &wk)
	if err != nil {
		return nil, err
	}

	return &v1.AddWorkerReply{Id: id}, nil
}

func (s *Service) UpdateWorker(ctx context.Context, in *v1.UpdateWorkerRequest) (*emptypb.Empty, error) {
	wk := biz.Worker{}
	if err := util.Transform(in, &wk); err != nil {
		return nil, errors.TransformFormat(err.Error())
	}

	return nil, s.worker.UpdateWorker(kratosx.MustContext(ctx), &wk)
}

func (s *Service) EnableWorker(ctx context.Context, in *v1.EnableWorkerRequest) (*emptypb.Empty, error) {
	return nil, s.worker.EnableWorker(kratosx.MustContext(ctx), in.Id)
}

func (s *Service) DisableWorker(ctx context.Context, in *v1.DisableWorkerRequest) (*emptypb.Empty, error) {
	return nil, s.worker.DisableWorker(kratosx.MustContext(ctx), in.Id)
}

func (s *Service) DeleteWorker(ctx context.Context, in *v1.DeleteWorkerRequest) (*emptypb.Empty, error) {
	return nil, s.worker.DeleteWorker(kratosx.MustContext(ctx), in.Id)
}
