package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"

	"github.com/limes-cloud/cron/api/cron/errors"
	pb "github.com/limes-cloud/cron/api/cron/server/task/v1"
	"github.com/limes-cloud/cron/internal/server/biz/task"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/data"
	ft "github.com/limes-cloud/cron/internal/server/factory/task"
)

type TaskService struct {
	pb.UnimplementedTaskServer
	uc *task.UseCase
}

func NewTaskService(conf *conf.Config) *TaskService {
	return &TaskService{
		uc: task.NewUseCase(conf, data.NewTaskRepo(), ft.GlobalFactory()),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewTaskService(c)
		pb.RegisterTaskHTTPServer(hs, srv)
		pb.RegisterTaskServer(gs, srv)
	})
}

// GetTaskGroup 获取指定的任务分组
func (s *TaskService) GetTaskGroup(c context.Context, req *pb.GetTaskGroupRequest) (*pb.GetTaskGroupReply, error) {
	var (
		in  = task.GetTaskGroupRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, err := s.uc.GetTaskGroup(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.GetTaskGroupReply{}
	if err := valx.Transform(result, &reply); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// ListTaskGroup 获取任务分组列表
func (s *TaskService) ListTaskGroup(c context.Context, req *pb.ListTaskGroupRequest) (*pb.ListTaskGroupReply, error) {
	var (
		in  = task.ListTaskGroupRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, total, err := s.uc.ListTaskGroup(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.ListTaskGroupReply{Total: total}
	if err := valx.Transform(result, &reply.List); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	return &reply, nil
}

// CreateTaskGroup 创建任务分组
func (s *TaskService) CreateTaskGroup(c context.Context, req *pb.CreateTaskGroupRequest) (*pb.CreateTaskGroupReply, error) {
	var (
		in  = task.TaskGroup{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	id, err := s.uc.CreateTaskGroup(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskGroupReply{Id: id}, nil
}

// UpdateTaskGroup 更新任务分组
func (s *TaskService) UpdateTaskGroup(c context.Context, req *pb.UpdateTaskGroupRequest) (*pb.UpdateTaskGroupReply, error) {
	var (
		in  = task.TaskGroup{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	if err := s.uc.UpdateTaskGroup(ctx, &in); err != nil {
		return nil, err
	}

	return &pb.UpdateTaskGroupReply{}, nil
}

// DeleteTaskGroup 删除任务分组
func (s *TaskService) DeleteTaskGroup(c context.Context, req *pb.DeleteTaskGroupRequest) (*pb.DeleteTaskGroupReply, error) {
	err := s.uc.DeleteTaskGroup(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskGroupReply{}, nil
}

// GetTask 获取指定的任务信息
func (s *TaskService) GetTask(c context.Context, req *pb.GetTaskRequest) (*pb.GetTaskReply, error) {
	var (
		in  = task.GetTaskRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, err := s.uc.GetTask(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.GetTaskReply{}
	if err := valx.Transform(result, &reply); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// ListTask 获取任务信息列表
func (s *TaskService) ListTask(c context.Context, req *pb.ListTaskRequest) (*pb.ListTaskReply, error) {
	var (
		in  = task.ListTaskRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, total, err := s.uc.ListTask(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.ListTaskReply{Total: total}
	if err := valx.Transform(result, &reply.List); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	return &reply, nil
}

// CreateTask 创建任务信息
func (s *TaskService) CreateTask(c context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskReply, error) {
	var (
		in  = task.Task{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	id, err := s.uc.CreateTask(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskReply{Id: id}, nil
}

// UpdateTask 更新任务信息
func (s *TaskService) UpdateTask(c context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskReply, error) {
	var (
		in  = task.Task{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	if err := s.uc.UpdateTask(ctx, &in); err != nil {
		return nil, err
	}

	return &pb.UpdateTaskReply{}, nil
}

// UpdateTaskStatus 更新任务信息状态
func (s *TaskService) UpdateTaskStatus(c context.Context, req *pb.UpdateTaskStatusRequest) (*pb.UpdateTaskStatusReply, error) {
	return &pb.UpdateTaskStatusReply{}, s.uc.UpdateTaskStatus(kratosx.MustContext(c), req.Id, req.Status)
}

// DeleteTask 删除任务信息
func (s *TaskService) DeleteTask(c context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskReply, error) {
	err := s.uc.DeleteTask(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskReply{}, nil
}

// ExecTask 执行任务信息
func (s *TaskService) ExecTask(c context.Context, req *pb.ExecTaskRequest) (*pb.ExecTaskReply, error) {
	err := s.uc.ExecTask(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ExecTaskReply{}, nil
}

// CancelExecTask 取消任务信息
func (s *TaskService) CancelExecTask(c context.Context, req *pb.CancelExecTaskRequest) (*pb.CancelExecTaskReply, error) {
	err := s.uc.CancelExecTask(kratosx.MustContext(c), req.Uuid)
	if err != nil {
		return nil, err
	}
	return &pb.CancelExecTaskReply{}, nil
}
