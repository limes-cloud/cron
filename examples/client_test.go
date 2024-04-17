package examples

import (
	"context"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"

	v1 "github.com/limes-cloud/cron/api/client/v1"
)

func readShellValue(path string) string {
	data, _ := os.ReadFile(path)
	return string(data)
}

func TestClient(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint("localhost:8121"))
	if err != nil {
		t.Error(err)
	}
	client := v1.NewServiceClient(conn)
	_, err = client.Healthy(ctx, &empty.Empty{})
	t.Error(err)
	// req := &v1.ExecTaskRequest{
	//	Id:            1,
	//	Uuid:          "1",
	//	Type:          "shell",
	//	Value:         readShellValue("script/for.sh"),
	//	MaxExecTime:   5,
	//	RetryCount:    0,
	//	RetryWaitTime: 2,
	// }
	//
	// stream, err := client.ExecTask(ctx, req)
	// if err != nil {
	//	t.Error(err)
	// }
	//
	// for {
	//	msg, err := stream.Recv()
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		if errors.IsExecTaskFail(err) {
	//			break
	//		}
	//		for {
	//			if stream, err = client.ExecTask(ctx, req); err == nil {
	//				break
	//			}
	//			t := rand.Intn(10)
	//			time.Sleep(time.Duration(t) * time.Second)
	//		}
	//	}
	//	fmt.Println(msg)
	// }
}