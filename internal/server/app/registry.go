package app

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type registryFunc func(hs *http.Server, gs *grpc.Server)

var registries []registryFunc

func register(fn registryFunc) {
	registries = append(registries, fn)
}

// Register 注册所有服务（gRPC + HTTP），由 kratosx.WithRegistrarServer 调用
func Register(hs *http.Server, gs *grpc.Server) {
	for _, registry := range registries {
		registry(hs, gs)
	}
}