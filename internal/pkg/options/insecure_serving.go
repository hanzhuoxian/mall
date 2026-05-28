package options

import (
	"fmt"
	"net"
	"strconv"

	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/spf13/pflag"
)

type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

func (i *InsecureServingOptions) ApplyTo(c *server.Config) error {
	c.InsecureServing = &server.InsecureServingInfo{
		Address: net.JoinHostPort(i.BindAddress, strconv.Itoa(i.BindPort)),
	}
	return nil
}

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
