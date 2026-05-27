package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/pkg/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type APIServer struct {
	middlewares         []string
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo
	ShutdownTimeout     *time.Duration

	*gin.Engine
	healthz                      bool
	enableMetrics                bool
	enableProfiling              bool
	insecureServer, secureServer http.Server
}

func initAPIServer(s *APIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

func (s *APIServer) Setup() {
}

func (s *APIServer) InstallMiddlewares() {
	s.Use(gin.Recovery())
	s.Use(requestid.New())

	for _, m := range s.middlewares {
		if mv, ok := middleware.Middlewares[m]; ok {
			s.Use(mv)
		}
	}
}

func (s *APIServer) InstallAPIs() {
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	if s.enableMetrics {
		s.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// install pprof handler
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	s.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.Get())
	})
}
