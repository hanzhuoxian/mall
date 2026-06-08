package userserver

import (
	"net"

	"google.golang.org/grpc"

	"github.com/hanzhuoxian/mall/pkg/log"
)

// GRPCServer 封装了 gRPC 服务器，持有监听地址并提供统一的启停接口。
type GRPCServer struct {
	*grpc.Server
	address string // 监听地址，格式为 host:port
}

// Run 在独立 goroutine 中启动 gRPC 服务并开始监听连接。
func (g *GRPCServer) Run() {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		log.Errorf("grpc run err: %w", err)
	}

	go func() {
		if err := g.Serve(listen); err != nil {
			log.Errorf("grpc run err: %w", err)
		}
	}()
}

// Close 优雅停止 gRPC 服务，等待所有进行中的 RPC 完成后再关闭。
func (g *GRPCServer) Close() {
	g.GracefulStop()
}
