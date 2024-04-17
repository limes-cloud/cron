package service

import (
	"context"

	"github.com/limes-cloud/kratosx"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/limes-cloud/cron/api/client/v1"
	"github.com/limes-cloud/cron/internal/client/biz"
	"github.com/limes-cloud/cron/internal/client/conf"
	"github.com/limes-cloud/cron/internal/client/factory"
)

type Service struct {
	v1.UnimplementedServiceServer
	conf *conf.Config
	task *biz.TaskUseCase
}

func New(conf *conf.Config) *Service {
	return &Service{
		conf: conf,
		task: biz.NewTaskUseCase(conf, factory.New(conf)),
	}
}

func (s *Service) Healthy(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Service) CancelExecTask(_ context.Context, req *v1.CancelExecTaskRequest) (*emptypb.Empty, error) {
	s.task.CancelExecTask(req.Uuid)
	return &emptypb.Empty{}, nil
}

func (s *Service) ExecTask(req *v1.ExecTaskRequest, res v1.Service_ExecTaskServer) error {
	task := &biz.Task{
		Id:            req.Id,
		Type:          req.Type,
		Value:         req.Value,
		ExpectCode:    req.ExpectCode,
		RetryCount:    req.RetryCount,
		RetryWaitTime: req.RetryWaitTime,
		MaxExecTime:   req.MaxExecTime,
		Uuid:          req.Uuid,
	}

	return s.task.ExecTask(
		kratosx.MustContext(res.Context()),
		task,
		func(reply *biz.ExecTaskReply) error {
			return res.Send(&v1.ExecTaskReply{
				Type:    reply.Type,
				Content: reply.Content,
				Time:    reply.Time,
			})
		},
	)
}
