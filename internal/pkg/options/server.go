package options

// ServerRunOptions
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"sj" mapstructure:"middlewares"`
}

func NewServerOptions() *ServerRunOptions {
	return &ServerRunOptions{}
}
