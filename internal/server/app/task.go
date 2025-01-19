package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"
	ktypes "github.com/limes-cloud/kratosx/types"

	"github.com/limes-cloud/cron/api/cron/errors"
	pb "github.com/limes-cloud/cron/api/cron/server/task/v1"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/service"
	"github.com/limes-cloud/cron/internal/server/infra/dbs"
	"github.com/limes-cloud/cron/internal/server/infra/rpc"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Task struct {
	pb.UnimplementedTaskServer
	srv *service.Task
}

func NewTask(conf *conf.Config) *Task {
	return &Task{
		srv: service.NewTask(
			conf,
			dbs.NewTask(),
			dbs.NewWorker(),
			dbs.NewLog(),
			rpc.NewTaskClient(),
		),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewTask(c)
		pb.RegisterTaskServer(gs, srv)
		pb.RegisterTaskHTTPServer(hs, srv)
	})
}

// GetTaskGroup 获取指定的任务分组
func (s *Task) GetTaskGroup(c context.Context, req *pb.GetTaskGroupRequest) (*pb.GetTaskGroupReply, error) {
	result, err := s.srv.GetTaskGroup(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetTaskGroupReply{
		Id:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		CreatedAt:   uint32(result.CreatedAt),
		UpdatedAt:   uint32(result.UpdatedAt),
	}, nil
}

// ListTaskGroup 获取任务分组列表
func (s *Task) ListTaskGroup(c context.Context, req *pb.ListTaskGroupRequest) (*pb.ListTaskGroupReply, error) {
	list, total, err := s.srv.ListTaskGroup(kratosx.MustContext(c), &types.ListTaskGroupRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Order:    req.Order,
		OrderBy:  req.OrderBy,
		Name:     req.Name,
	})
	if err != nil {
		return nil, err
	}
	reply := pb.ListTaskGroupReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &pb.ListTaskGroupReply_TaskGroup{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   uint32(item.CreatedAt),
			UpdatedAt:   uint32(item.UpdatedAt),
		})
	}
	return &reply, nil
}

// CreateTaskGroup 创建任务分组
func (s *Task) CreateTaskGroup(c context.Context, req *pb.CreateTaskGroupRequest) (*pb.CreateTaskGroupReply, error) {
	id, err := s.srv.CreateTaskGroup(kratosx.MustContext(c), &entity.TaskGroup{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskGroupReply{Id: id}, nil
}

// UpdateTaskGroup 更新任务分组
func (s *Task) UpdateTaskGroup(c context.Context, req *pb.UpdateTaskGroupRequest) (*pb.UpdateTaskGroupReply, error) {
	if err := s.srv.UpdateTaskGroup(kratosx.MustContext(c), &entity.TaskGroup{
		BaseModel:   ktypes.BaseModel{Id: req.Id},
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		return nil, err
	}
	return &pb.UpdateTaskGroupReply{}, nil
}

// DeleteTaskGroup 删除任务分组
func (s *Task) DeleteTaskGroup(c context.Context, req *pb.DeleteTaskGroupRequest) (*pb.DeleteTaskGroupReply, error) {
	err := s.srv.DeleteTaskGroup(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskGroupReply{}, nil
}

// GetTask 获取指定的任务信息
func (s *Task) GetTask(c context.Context, req *pb.GetTaskRequest) (*pb.GetTaskReply, error) {
	result, err := s.srv.GetTask(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	reply := pb.GetTaskReply{}
	if err := valx.Transform(result, &reply); err != nil {
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// ListTask 获取任务信息列表
func (s *Task) ListTask(c context.Context, req *pb.ListTaskRequest) (*pb.ListTaskReply, error) {
	result, total, err := s.srv.ListTask(kratosx.MustContext(c), &types.ListTaskRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		GroupId:  req.GroupId,
		Name:     req.Name,
		Tag:      req.Tag,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	reply := pb.ListTaskReply{Total: total}
	if err := valx.Transform(result, &reply.List); err != nil {
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// CreateTask 创建任务信息
func (s *Task) CreateTask(c context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskReply, error) {
	if (req.Start != nil || req.End != nil) && (req.Start == nil || req.End == nil) {
		return nil, errors.ParamsError()
	}
	id, err := s.srv.CreateTask(kratosx.MustContext(c), &entity.Task{
		GroupId:       req.GroupId,
		Name:          req.Name,
		Tag:           req.Tag,
		Spec:          req.Spec,
		Status:        req.Status,
		WorkerType:    req.WorkerType,
		WorkerGroupId: req.WorkerGroupId,
		WorkerId:      req.WorkerId,
		ExecType:      req.ExecType,
		ExecValue:     req.ExecValue,
		ExpectCode:    req.ExpectCode,
		RetryCount:    req.RetryCount,
		RetryWaitTime: req.RetryWaitTime,
		MaxExecTime:   req.MaxExecTime,
		Description:   req.Description,
		Start:         req.Start,
		End:           req.End,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskReply{Id: id}, nil
}

// UpdateTask 更新任务信息
func (s *Task) UpdateTask(c context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskReply, error) {
	if (req.Start != nil || req.End != nil) && (req.Start == nil || req.End == nil) {
		return nil, errors.ParamsError()
	}
	if err := s.srv.UpdateTask(kratosx.MustContext(c), &entity.Task{
		BaseModel:     ktypes.BaseModel{Id: req.Id},
		GroupId:       req.GroupId,
		Name:          req.Name,
		Tag:           req.Tag,
		Spec:          req.Spec,
		WorkerType:    req.WorkerType,
		WorkerGroupId: req.WorkerGroupId,
		WorkerId:      req.WorkerId,
		ExecType:      req.ExecType,
		ExecValue:     req.ExecValue,
		ExpectCode:    req.ExpectCode,
		RetryCount:    req.RetryCount,
		RetryWaitTime: req.RetryWaitTime,
		MaxExecTime:   req.MaxExecTime,
		Description:   req.Description,
		Start:         req.Start,
		End:           req.End,
	}); err != nil {
		return nil, err
	}

	return &pb.UpdateTaskReply{}, nil
}

// UpdateTaskStatus 更新任务信息状态
func (s *Task) UpdateTaskStatus(c context.Context, req *pb.UpdateTaskStatusRequest) (*pb.UpdateTaskStatusReply, error) {
	return &pb.UpdateTaskStatusReply{}, s.srv.UpdateTaskStatus(kratosx.MustContext(c), req.Id, req.Status)
}

// DeleteTask 删除任务信息
func (s *Task) DeleteTask(c context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskReply, error) {
	err := s.srv.DeleteTask(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskReply{}, nil
}

// ExecTask 执行任务信息
func (s *Task) ExecTask(c context.Context, req *pb.ExecTaskRequest) (*pb.ExecTaskReply, error) {
	err := s.srv.ExecTask(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ExecTaskReply{}, nil
}

// CancelExecTask 取消任务信息
func (s *Task) CancelExecTask(c context.Context, req *pb.CancelExecTaskRequest) (*pb.CancelExecTaskReply, error) {
	err := s.srv.CancelExecTask(kratosx.MustContext(c), req.Uuid)
	if err != nil {
		return nil, err
	}
	return &pb.CancelExecTaskReply{}, nil
}
