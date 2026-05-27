package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RecommendedHomeDir = ".mall"

	RecommendedEnvPrefix = "MALL"
)

type Config struct {
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Jwt             *JwtInfo
}

type CertKey struct {
	CertFile string
	KeyFile  string
}

type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

type InsecureServingInfo struct {
	Address string
}

type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm:      "mall",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
	}
}

type CompletedConfig struct {
	*Config
}

func Complete(cfg *Config) CompletedConfig {
	return CompletedConfig{cfg}
}

func (c CompletedConfig) New() (*APIServer, error) {
	gin.SetMode(c.Mode)

	s := &APIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	initAPIServer(s)
	return s, nil
}
