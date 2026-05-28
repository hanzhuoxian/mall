package options

import (
	"github.com/hanzhuoxian/mall/internal/pkg/server"
	"github.com/spf13/pflag"
)

// ServerRunOptions
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"sj" mapstructure:"middlewares"`
}

func NewServerOptions() *ServerRunOptions {
	defaults := server.NewConfig()
	return &ServerRunOptions{
		Mode:        defaults.Mode,
		Healthz:     defaults.Healthz,
		Middlewares: defaults.Middlewares,
	}
}

func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares
	return nil
}

func (o *ServerRunOptions) Validate() []error {
	var errors []error
	return errors
}

func (s *ServerRunOptions) Flags() (fs pflag.FlagSet) {
	fs.StringVar(&s.Mode, "server.mode", s.Mode, "server mode, Supported server mode: debug, test, release")
	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, "enable healthz")
	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, "server middlewares")
	return
}
