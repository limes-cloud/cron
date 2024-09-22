package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/limes-cloud/cron/api/cron/client/v1"
	"github.com/limes-cloud/cron/internal/client/conf"
	"github.com/limes-cloud/cron/internal/client/factory"
	"github.com/limes-cloud/cron/internal/client/service"
)

type Task struct {
	pb.UnimplementedTaskServer
	srv *service.Task
}

func NewTask(conf *conf.Config) *Task {
	return &Task{
		srv: service.NewTask(
			conf,
			factory.New(conf),
		),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewTask(c)
		pb.RegisterTaskServer(gs, srv)
	})
}

func (s *Task) Healthy(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Task) CancelExecTask(_ context.Context, req *pb.CancelExecTaskRequest) (*emptypb.Empty, error) {
	s.srv.CancelExecTask(req.Uuid)
	return &emptypb.Empty{}, nil
}

func (s *Task) ExecTask(req *pb.ExecTaskRequest, res pb.Task_ExecTaskServer) error {
	task := &service.ExecTaskRequest{
		Id:            req.Id,
		Type:          req.Type,
		Value:         req.Value,
		ExpectCode:    req.ExpectCode,
		RetryCount:    req.RetryCount,
		RetryWaitTime: req.RetryWaitTime,
		MaxExecTime:   req.MaxExecTime,
		Uuid:          req.Uuid,
	}

	return s.srv.ExecTask(
		kratosx.MustContext(res.Context()),
		task,
		func(reply *service.ExecTaskReply) error {
			return res.Send(&pb.ExecTaskReply{
				Type:    reply.Type,
				Content: reply.Content,
				Time:    reply.Time,
			})
		},
	)
}
