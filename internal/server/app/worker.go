package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	ktypes "github.com/limes-cloud/kratosx/types"

	pb "github.com/limes-cloud/cron/api/cron/server/worker/v1"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/domain/entity"
	"github.com/limes-cloud/cron/internal/server/domain/service"
	"github.com/limes-cloud/cron/internal/server/infra/dbs"
	"github.com/limes-cloud/cron/internal/server/infra/rpc"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Worker struct {
	pb.UnimplementedWorkerServer
	srv *service.Worker
}

func NewWorker(conf *conf.Config) *Worker {
	return &Worker{
		srv: service.NewWorker(
			conf,
			dbs.NewWorker(),
			rpc.NewTaskClient(),
		),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewWorker(c)
		pb.RegisterWorkerServer(gs, srv)
		pb.RegisterWorkerHTTPServer(hs, srv)
	})
}

// GetWorkerGroup 获取指定的节点分组
func (s *Worker) GetWorkerGroup(c context.Context, req *pb.GetWorkerGroupRequest) (*pb.GetWorkerGroupReply, error) {
	result, err := s.srv.GetWorkerGroup(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetWorkerGroupReply{
		Id:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		CreatedAt:   uint32(result.CreatedAt),
		UpdatedAt:   uint32(result.UpdatedAt),
	}, nil
}

// ListWorkerGroup 获取节点分组列表
func (s *Worker) ListWorkerGroup(c context.Context, req *pb.ListWorkerGroupRequest) (*pb.ListWorkerGroupReply, error) {
	list, total, err := s.srv.ListWorkerGroup(kratosx.MustContext(c), &types.ListWorkerGroupRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
	})
	if err != nil {
		return nil, err
	}

	reply := pb.ListWorkerGroupReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &pb.ListWorkerGroupReply_WorkerGroup{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   uint32(item.CreatedAt),
			UpdatedAt:   uint32(item.UpdatedAt),
		})
	}
	return &reply, nil
}

// CreateWorkerGroup 创建节点分组
func (s *Worker) CreateWorkerGroup(c context.Context, req *pb.CreateWorkerGroupRequest) (*pb.CreateWorkerGroupReply, error) {
	id, err := s.srv.CreateWorkerGroup(kratosx.MustContext(c), &entity.WorkerGroup{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateWorkerGroupReply{Id: id}, nil
}

// UpdateWorkerGroup 更新节点分组
func (s *Worker) UpdateWorkerGroup(c context.Context, req *pb.UpdateWorkerGroupRequest) (*pb.UpdateWorkerGroupReply, error) {
	if err := s.srv.UpdateWorkerGroup(kratosx.MustContext(c), &entity.WorkerGroup{
		BaseModel:   ktypes.BaseModel{Id: req.Id},
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		return nil, err
	}

	return &pb.UpdateWorkerGroupReply{}, nil
}

// DeleteWorkerGroup 删除节点分组
func (s *Worker) DeleteWorkerGroup(c context.Context, req *pb.DeleteWorkerGroupRequest) (*pb.DeleteWorkerGroupReply, error) {
	err := s.srv.DeleteWorkerGroup(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteWorkerGroupReply{}, nil
}

// GetWorker 获取指定的节点信息
func (s *Worker) GetWorker(c context.Context, req *pb.GetWorkerRequest) (*pb.GetWorkerReply, error) {
	result, err := s.srv.GetWorker(kratosx.MustContext(c), &types.GetWorkerRequest{
		Id: req.Id,
		Ip: req.Ip,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetWorkerReply{
		Id:          result.Id,
		Name:        result.Name,
		Ip:          result.Ip,
		Ak:          result.Ak,
		Sk:          result.Sk,
		GroupId:     result.GroupId,
		Status:      result.Status,
		Description: result.Description,
		CreatedAt:   uint32(result.CreatedAt),
		UpdatedAt:   uint32(result.UpdatedAt),
	}, nil
}

// ListWorker 获取节点信息列表
func (s *Worker) ListWorker(c context.Context, req *pb.ListWorkerRequest) (*pb.ListWorkerReply, error) {
	list, total, err := s.srv.ListWorker(kratosx.MustContext(c), &types.ListWorkerRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Ip:       req.Ip,
		GroupId:  req.GroupId,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	reply := pb.ListWorkerReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &pb.ListWorkerReply_Worker{
			Id:          item.Id,
			Name:        item.Name,
			Ip:          item.Ip,
			Ak:          item.Ak,
			Sk:          item.Sk,
			GroupId:     item.GroupId,
			Status:      item.Status,
			Description: item.Description,
			CreatedAt:   uint32(item.CreatedAt),
			UpdatedAt:   uint32(item.UpdatedAt),
			Group: &pb.ListWorkerReply_Group{
				Name: item.Group.Name,
			},
		})
	}

	return &reply, nil
}

// CreateWorker 创建节点信息
func (s *Worker) CreateWorker(c context.Context, req *pb.CreateWorkerRequest) (*pb.CreateWorkerReply, error) {
	id, err := s.srv.CreateWorker(kratosx.MustContext(c), &entity.Worker{
		Name:        req.Name,
		Ip:          req.Ip,
		Ak:          req.Ak,
		Sk:          req.Sk,
		GroupId:     req.GroupId,
		Status:      req.Status,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateWorkerReply{Id: id}, nil
}

// UpdateWorker 更新节点信息
func (s *Worker) UpdateWorker(c context.Context, req *pb.UpdateWorkerRequest) (*pb.UpdateWorkerReply, error) {
	if err := s.srv.UpdateWorker(kratosx.MustContext(c), &entity.Worker{
		BaseModel:   ktypes.BaseModel{Id: req.Id},
		Name:        req.Name,
		Ip:          req.Ip,
		Ak:          req.Ak,
		Sk:          req.Sk,
		GroupId:     req.GroupId,
		Description: req.Description,
	}); err != nil {
		return nil, err
	}

	return &pb.UpdateWorkerReply{}, nil
}

// UpdateWorkerStatus 更新节点信息状态
func (s *Worker) UpdateWorkerStatus(c context.Context, req *pb.UpdateWorkerStatusRequest) (*pb.UpdateWorkerStatusReply, error) {
	return &pb.UpdateWorkerStatusReply{}, s.srv.UpdateWorkerStatus(kratosx.MustContext(c), req.Id, req.Status)
}

// DeleteWorker 删除节点信息
func (s *Worker) DeleteWorker(c context.Context, req *pb.DeleteWorkerRequest) (*pb.DeleteWorkerReply, error) {
	err := s.srv.DeleteWorker(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteWorkerReply{}, nil
}
