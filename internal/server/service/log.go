package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/valx"

	"github.com/limes-cloud/cron/api/cron/errors"
	pb "github.com/limes-cloud/cron/api/cron/server/log/v1"
	"github.com/limes-cloud/cron/internal/server/biz/log"
	"github.com/limes-cloud/cron/internal/server/conf"
	"github.com/limes-cloud/cron/internal/server/data"
)

type LogService struct {
	pb.UnimplementedLogServer
	uc *log.UseCase
}

func NewLogService(conf *conf.Config) *LogService {
	return &LogService{
		uc: log.NewUseCase(conf, data.NewLogRepo()),
	}
}

func init() {
	register(func(c *conf.Config, hs *http.Server, gs *grpc.Server) {
		srv := NewLogService(c)
		pb.RegisterLogHTTPServer(hs, srv)
		pb.RegisterLogServer(gs, srv)
	})
}

// GetLog 获取指定的日志信息
func (s *LogService) GetLog(c context.Context, req *pb.GetLogRequest) (*pb.GetLogReply, error) {
	var (
		in  = log.GetLogRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, err := s.uc.GetLog(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.GetLogReply{}
	if err := valx.Transform(result, &reply); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// ListLog 获取日志信息列表
func (s *LogService) ListLog(c context.Context, req *pb.ListLogRequest) (*pb.ListLogReply, error) {
	var (
		in  = log.ListLogRequest{}
		ctx = kratosx.MustContext(c)
	)

	if err := valx.Transform(req, &in); err != nil {
		ctx.Logger().Warnw("msg", "req transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	result, total, err := s.uc.ListLog(ctx, &in)
	if err != nil {
		return nil, err
	}

	reply := pb.ListLogReply{Total: total}
	if err := valx.Transform(result, &reply.List); err != nil {
		ctx.Logger().Warnw("msg", "reply transform err", "err", err.Error())
		return nil, errors.TransformError()
	}

	return &reply, nil
}

// DeleteLog 删除日志信息
func (s *LogService) DeleteLog(c context.Context, req *pb.DeleteLogRequest) (*pb.DeleteLogReply, error) {
	total, err := s.uc.DeleteLog(kratosx.MustContext(c), req.Ids)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteLogReply{Total: total}, nil
}
