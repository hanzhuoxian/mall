package userserver

import (
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/hanzhuoxian/mall/internal/userserver/controller"
	"github.com/hanzhuoxian/mall/internal/userserver/service"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// installRoutes 注册所有 gRPC 服务到服务器。
func installRoutes(s *server.GRPCServer, svc service.Service) {
	userv1.RegisterUserServiceServer(s.Server, controller.NewUserServiceServer(svc))
}
