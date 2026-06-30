package grpcclient

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// UserClient wraps the generated gRPC client and the underlying connection
// so the connection can be closed cleanly on shutdown.
type UserClient struct {
	userv1.UserServiceClient
	conn *grpc.ClientConn
}

// NewUserClient 建立到 userserver 的 gRPC 连接并返回 UserClient。
func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}
	return &UserClient{
		UserServiceClient: userv1.NewUserServiceClient(conn),
		conn:              conn,
	}, nil
}

// Close 关闭底层 gRPC 连接。
func (c *UserClient) Close() error {
	return c.conn.Close()
}
