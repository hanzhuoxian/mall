package server

import (
	"net"

	"github.com/hanzhuoxian/mall/pkg/logger"
	"google.golang.org/grpc"
)

// GRPCServer 封装了标准 gRPC 服务器，持有监听地址并提供统一的启停接口。
type GRPCServer struct {
	*grpc.Server
	address string // 监听地址，格式为 host:port
}

// NewGRPCServer 创建一个新的 GRPCServer 实例，可通过 opts 传入自定义 grpc.ServerOption。
func NewGRPCServer(addr string, opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{Server: grpc.NewServer(opts...), address: addr}
}

// Run 在独立 goroutine 中启动 gRPC 服务并开始监听连接，返回监听错误（不包括 Serve 内部错误）。
func (g *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		return err
	}

	go func() {
		if err := g.Serve(listen); err != nil {
			logger.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	return nil
}

// Close 优雅停止 gRPC 服务，等待所有进行中的 RPC 完成后再关闭。
func (g *GRPCServer) Close() error {
	g.GracefulStop()
	return nil
}
