package options

import (
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/pflag"

	"github.com/hanzhuoxian/mall/internal/pkg/server"
)

// InsecureServingOptions 包含 HTTP（非 TLS）服务的监听地址和端口配置。
type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

// NewInsecureServingOptions 返回默认监听 127.0.0.1:8080 的 InsecureServingOptions 实例。
func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

// ApplyTo 将 HTTP 监听地址写入 server.Config 的 InsecureServing 字段。
func (i *InsecureServingOptions) ApplyTo(c *server.Config) error {
	c.InsecureServing = &server.InsecureServingInfo{
		Address: net.JoinHostPort(i.BindAddress, strconv.Itoa(i.BindPort)),
	}
	return nil
}

// Validate 校验 HTTP 服务选项合法性，确保端口在 0-65535 范围内。
func (i *InsecureServingOptions) Validate() []error {
	var errors []error

	if i.BindPort < 0 || i.BindPort > 65535 {
		errors = append(errors,
			fmt.Errorf(
				"--insecure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				i.BindPort,
			))
	}
	return errors
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (i *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&i.BindAddress, "insecure.bind-address", i.BindAddress, ""+
		"The IP address on which to serve the --insecure.bind-port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	fs.IntVar(&i.BindPort, "insecure.bind-port", i.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")
}
