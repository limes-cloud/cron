package factory

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/client/conf"
	"github.com/limes-cloud/cron/internal/client/service"
)

func readShellValue(path string) string {
	data, _ := os.ReadFile(path)
	return string(data)
}

func TestFactory_ExecTask(t *testing.T) {
	factory := New(&conf.Config{Shell: "/bin/bash"})

	ctx := kratosx.MustContext(context.Background())
	err := factory.ExecTask(ctx, &service.ExecTaskRequest{
		Id:            1,
		Uuid:          "1",
		Type:          "shell",
		Value:         readShellValue("examples/for.sh"),
		MaxExecTime:   11,
		RetryCount:    3,
		RetryWaitTime: 2,
	}, func(reply *service.ExecTaskReply) error {
		fmt.Println(reply)
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}
