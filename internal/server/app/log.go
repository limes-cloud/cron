package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"

	pb "github.com/limes-cloud/cron/api/cron/server/log/v1"
	"github.com/limes-cloud/cron/internal/server/domain/service"
	"github.com/limes-cloud/cron/internal/server/infra/dbs"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Log struct {
	pb.UnimplementedLogServer
	srv *service.Log
}

func NewLog() *Log {
	return &Log{
		srv: service.NewLog(
			dbs.NewLog(),
			dbs.NewExecution(),
		),
	}
}

func init() {
	register(func(hs *http.Server, gs *grpc.Server) {
		srv := NewLog()
		pb.RegisterLogServer(gs, srv)
		pb.RegisterLogHTTPServer(hs, srv)
	})
}

// GetLog 获取指定的日志信息
func (s *Log) GetLog(c context.Context, req *pb.GetLogRequest) (*pb.GetLogReply, error) {
	log, err := s.srv.GetLog(kratosx.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetLogReply{
		Id:             log.Id,
		Uuid:           log.Uuid,
		WorkerId:       log.WorkerId,
		WorkerSnapshot: log.WorkerSnapshot,
		TaskId:         log.TaskId,
		TaskSnapshot:   log.TaskSnapshot,
		StartAt:        uint32(log.StartAt),
		EndAt:          uint32(log.EndAt),
		Content:        log.Content,
		Status:         log.Status,
	}, nil
}

// ListLog 获取日志信息列表
func (s *Log) ListLog(c context.Context, req *pb.ListLogRequest) (*pb.ListLogReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	list, total, err := s.srv.ListLog(kratosx.MustContext(c), &types.ListLogRequest{
		App:      app,
		Page:     req.Page,
		PageSize: req.PageSize,
		TaskId:   req.TaskId,
	})
	if err != nil {
		return nil, err
	}

	reply := pb.ListLogReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &pb.ListLogReply_Log{
			Id:       item.Id,
			Uuid:     item.Uuid,
			WorkerId: item.WorkerId,
			TaskId:   item.TaskId,
			StartAt:  uint32(item.StartAt),
			EndAt:    uint32(item.EndAt),
			Status:   item.Status,
		})
	}
	return &reply, nil
}

// ListExecution 获取执行历史列表
func (s *Log) ListExecution(c context.Context, req *pb.ListExecutionRequest) (*pb.ListExecutionReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	list, total, err := s.srv.ListExecution(kratosx.MustContext(c), &types.ListExecutionRequest{
		App:         app,
		Page:        req.Page,
		PageSize:    req.PageSize,
		TaskId:      req.TaskId,
		Status:      req.Status,
		TriggerType: req.TriggerType,
	})
	if err != nil {
		return nil, err
	}

	reply := pb.ListExecutionReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &pb.Execution{
			Id:          item.Id,
			TaskId:      item.TaskId,
			TriggerAt:   uint32(item.TriggerAt),
			TriggerType: item.TriggerType,
			WorkerId:    item.WorkerId,
			Status:      item.Status,
			RetryCount:  item.RetryCount,
			LogUuid:     valueStrOrNil(item.LogUuid),
			StartedAt:   uint32OrNil(item.StartedAt),
			FinishedAt:  uint32OrNil(item.FinishedAt),
			LastError:   valueStrOrNil(item.LastError),
			CreatedAt:   uint32(item.CreatedAt),
		})
	}
	return &reply, nil
}

// valueStrOrNil 空串返回 nil（proto optional string）
func valueStrOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// uint32OrNil 解引用 int64 指针（proto optional uint32）
func uint32OrNil(v *int64) *uint32 {
	if v == nil {
		return nil
	}
	r := uint32(*v)
	return &r
}
