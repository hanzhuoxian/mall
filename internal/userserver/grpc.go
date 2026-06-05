package userserver

import (
	"net"

	"github.com/hanzhuoxian/mall/pkg/log"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	address string
}

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

func (g *GRPCServer) Close() {
	g.GracefulStop()
}
