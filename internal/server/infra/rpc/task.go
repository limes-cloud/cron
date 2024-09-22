package rpc

import (
	"sync"

	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	midmetadata "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/limes-cloud/kratosx"
	"github.com/limes-cloud/kratosx/config"
	"github.com/limes-cloud/kratosx/library/logger"
	"github.com/limes-cloud/kratosx/library/signature"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/limes-cloud/cron/api/cron/client/v1"
	"github.com/limes-cloud/cron/internal/server/types"
)

const (
	TaskClientServer = "TaskClient"
)

var (
	taskIns  *TaskClient
	taskOnce sync.Once
)

type TaskClient struct {
}

type connect struct {
	ip string
	ak string
	sk string
}

func NewTaskClient() *TaskClient {
	taskOnce.Do(func() {
		taskIns = &TaskClient{}
	})
	return taskIns
}

func (i TaskClient) client(ctx kratosx.Context, req *connect) (pb.TaskClient, error) {
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(req.ip),
		grpc.WithMiddleware(
			midmetadata.Client(),
			logging.Client(logger.Instance()),
			circuitbreaker.Client(),
			tracing.Client(),
			signature.Instance().Client(&config.Signature{
				Enable: true,
				Ak:     req.ak,
				Sk:     req.sk,
			}),
		),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewTaskClient(conn), nil
}

func (i TaskClient) ExecTask(ctx kratosx.Context, req *types.ExecTaskRequest, recv func(msg *types.ExecTaskLog)) error {
	client, err := i.client(ctx, &connect{
		ip: req.IP,
		ak: req.Ak,
		sk: req.Sk,
	})
	if err != nil {
		return err
	}
	stream, err := client.ExecTask(ctx, &pb.ExecTaskRequest{
		Id:            req.Id,
		Type:          req.ExecType,
		Value:         req.ExecValue,
		ExpectCode:    req.ExpectCode,
		RetryCount:    req.RetryCount,
		RetryWaitTime: req.RetryWaitTime,
		MaxExecTime:   req.MaxExecTime,
		Uuid:          req.Uuid,
	})
	if err != nil {
		return err
	}

	for {
		msg, rer := stream.Recv()
		if rer != nil {
			return err
		}
		recv(&types.ExecTaskLog{
			Type:    msg.Type,
			Content: msg.Content,
			Time:    int64(msg.Time),
		})
	}
}

func (i TaskClient) CancelExec(ctx kratosx.Context, req *types.CancelExecRequest) error {
	client, err := i.client(ctx, &connect{
		ip: req.IP,
		ak: req.Ak,
		sk: req.Sk,
	})
	if err != nil {
		return err
	}
	_, err = client.CancelExecTask(ctx, &pb.CancelExecTaskRequest{Uuid: req.Uuid})
	return err
}

func (i TaskClient) CheckHealthy(ctx kratosx.Context, req *types.CheckWorkerRequest) error {
	client, err := i.client(ctx, &connect{
		ip: req.IP,
		ak: req.Ak,
		sk: req.Sk,
	})
	if err != nil {
		return err
	}
	_, err = client.Healthy(ctx, &emptypb.Empty{})
	return err
}
