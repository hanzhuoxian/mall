package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// GRPCOptions are for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
type GRPCOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
}

func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		MaxMsgSize:  4 * 1024 * 1024,
	}
}

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
