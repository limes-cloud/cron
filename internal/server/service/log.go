package service

import (
	"context"

	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/pkg/util"

	"github.com/limes-cloud/cron/api/errors"
	v1 "github.com/limes-cloud/cron/api/server/v1"
	"github.com/limes-cloud/cron/internal/server/biz"
)

func (s *Service) PageLog(ctx context.Context, in *v1.PageLogRequest) (*v1.PageLogReply, error) {
	var req biz.PageLogRequest
	if err := util.Transform(in, &req); err != nil {
		return nil, errors.Transform()
	}

	list, total, err := s.log.PageLog(kratosx.MustContext(ctx), &req)
	if err != nil {
		return nil, err
	}

	reply := v1.PageLogReply{Total: total}
	if err := util.Transform(list, &reply.List); err != nil {
		return nil, errors.Transform()
	}

	return &reply, nil
}

func (s *Service) GetLog(ctx context.Context, in *v1.GetLogRequest) (*v1.Log, error) {
	log, err := s.log.GetLog(kratosx.MustContext(ctx), in.Id)
	if err != nil {
		return nil, err
	}
	reply := v1.Log{}
	if err := util.Transform(log, &reply); err != nil {
		return nil, errors.Transform()
	}

	return &reply, nil
}
