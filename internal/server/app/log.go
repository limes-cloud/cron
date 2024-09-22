package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"

	pb "github.com/limes-cloud/cron/api/cron/server/log/v1"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/domain/service"
	"github.com/limes-cloud/cron/internal/server/infra/dbs"
	"github.com/limes-cloud/cron/internal/server/types"
)

type Log struct {
	pb.UnimplementedLogServer
	srv *service.Log
}

func NewLog(conf *conf.Config) *Log {
	return &Log{
		srv: service.NewLog(
			conf,
			dbs.NewLog(),
		),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewLog(c)
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
	list, total, err := s.srv.ListLog(kratosx.MustContext(c), &types.ListLogRequest{
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
