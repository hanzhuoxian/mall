package userserver

import (
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	address string
}

func (g *GRPCServer) Run() {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
	}

	go func() {
		if err := g.Serve(listen); err != nil {
		}
	}()
}

func (g *GRPCServer) Close() {
	g.GracefulStop()
}
