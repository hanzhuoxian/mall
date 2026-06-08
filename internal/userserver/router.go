package userserver

import (
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/controller"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// initGRPCServer 注册所有 gRPC 服务到服务器。
func installRoutes(s *server.GRPCServer) {
	userv1.RegisterUserServiceServer(s.Server, controller.NewUserServiceServer())
}
