package options

import (
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/pflag"
)

// GRPCOptions are for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
type GRPCOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
}

// NewGRPCOptions 返回带有合理默认值的 GRPCOptions 实例（地址 0.0.0.0:8081，消息大小 4MB）。
func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		MaxMsgSize:  4 * 1024 * 1024,
	}
}

// Address 返回 host:port 格式的地址。
func (g *GRPCOptions) Address() string {
	return net.JoinHostPort(g.BindAddress, strconv.Itoa(g.BindPort))
}

// Validate 校验 gRPC 选项合法性，确保端口在有效范围内。
func (g *GRPCOptions) Validate() []error {
	errors := []error{}
	if g.BindPort < 0 || g.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--insecure-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				g.BindPort,
			),
		)
	}
	return errors
}

// AddFlags 向指定 FlagSet 注册 gRPC 服务器参数的命令行 flag。
func (g *GRPCOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&g.BindAddress, "grpc.bind-address", g.BindAddress, ""+
		"The IP address on which to serve the --grpc.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&g.BindPort, "grpc.bind-port", g.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated grpc accesg. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")

	fs.IntVar(&g.MaxMsgSize, "grpc.max-msg-size", g.MaxMsgSize, "gRPC max message size.")
}
