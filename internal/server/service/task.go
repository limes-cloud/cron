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

func (s *Service) PageTask(ctx context.Context, in *v1.PageTaskRequest) (*v1.PageTaskReply, error) {
	var req biz.PageTaskRequest
	if err := util.Transform(in, &req); err != nil {
		return nil, errors.Transform()
	}

	list, total, err := s.task.PageTask(kratosx.MustContext(ctx), &req)
	if err != nil {
		return nil, err
	}

	reply := v1.PageTaskReply{Total: total}
	if err := util.Transform(list, &reply.List); err != nil {
		return nil, errors.Transform()
	}

	return &reply, nil
}

func (s *Service) AddTask(ctx context.Context, in *v1.AddTaskRequest) (*v1.AddTaskReply, error) {
	wk := biz.Task{}
	if err := util.Transform(in, &wk); err != nil {
		return nil, errors.TransformFormat(err.Error())
	}

	id, err := s.task.AddTask(kratosx.MustContext(ctx), &wk)
	if err != nil {
		return nil, err
	}

	return &v1.AddTaskReply{Id: id}, nil
}

func (s *Service) UpdateTask(ctx context.Context, in *v1.UpdateTaskRequest) (*emptypb.Empty, error) {
	wk := biz.Task{}
	if err := util.Transform(in, &wk); err != nil {
		return nil, errors.TransformFormat(err.Error())
	}

	return nil, s.task.UpdateTask(kratosx.MustContext(ctx), &wk)
}

func (s *Service) EnableTask(ctx context.Context, in *v1.EnableTaskRequest) (*emptypb.Empty, error) {
	return nil, s.task.EnableTask(kratosx.MustContext(ctx), in.Id)
}

func (s *Service) DisableTask(ctx context.Context, in *v1.DisableTaskRequest) (*emptypb.Empty, error) {
	return nil, s.task.DisableTask(kratosx.MustContext(ctx), in.Id)
}

func (s *Service) CancelExecTask(ctx context.Context, in *v1.CancelExecTaskRequest) (*emptypb.Empty, error) {
	return nil, s.task.CancelTask(kratosx.MustContext(ctx), in.Uuid)
}

func (s *Service) DeleteTask(ctx context.Context, in *v1.DeleteTaskRequest) (*emptypb.Empty, error) {
	return nil, s.task.DeleteTask(kratosx.MustContext(ctx), in.Id)
}
