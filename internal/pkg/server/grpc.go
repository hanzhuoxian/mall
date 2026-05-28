package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	address string
}

func NewGRPCServer(addr string, opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{Server: grpc.NewServer(opts...), address: addr}
}

func (g *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		return err
	}

	go func() {
		if err := g.Serve(listen); err != nil {
			log.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	return nil
}

func (g *GRPCServer) Close() error {
	g.GracefulStop()
	return nil
}
