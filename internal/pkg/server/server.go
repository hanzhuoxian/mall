package server

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/sync/errgroup"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/pkg/logger"
	"github.com/hanzhuoxian/mall/pkg/version"
)

// APIServer 是基于 Gin 的 HTTP API 服务器，支持同时监听 HTTP 和 HTTPS，
// 并内置健康检查、Prometheus 指标和 pprof 性能分析接口。
type APIServer struct {
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo
	// ShutdownTimeout 是 Close 等待在途请求处理完成的上限，为 0 时回落到 DefaultShutdownTimeout。
	ShutdownTimeout time.Duration
	*gin.Engine
	insecureServer  *http.Server
	secureServer    *http.Server
	middlewares     []string
	healthz         bool
	enableMetrics   bool
	enableProfiling bool
}

// initAPIServer 依次执行 Setup、InstallMiddlewares 和 InstallAPIs，完成服务器初始化。
func initAPIServer(s *APIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// Setup 执行服务器自定义的额外初始化逻辑，当前为空占位，后续可在此扩展。
func (s *APIServer) Setup() {
}

// InstallMiddlewares 注册默认中间件（otelgin 链路追踪、recovery、requestid、context）及配置中指定的中间件。
func (s *APIServer) InstallMiddlewares() {
	// 传空字符串：otelgin 会从请求 Host 推断 server.address；服务标识（service.name）
	// 来自 OTel resource，与此参数无关。telemetry 未启用时全局 Tracer 为 no-op。
	s.Use(otelgin.Middleware(""))
	s.Use(gin.Recovery())
	s.Use(requestid.New())
	s.Use(middleware.Context())
	for _, m := range s.middlewares {
		if mv, ok := middleware.Middlewares[m]; ok {
			s.Use(mv)
		}
	}
}

// InstallAPIs 根据配置注册内置 API 路由：/healthz、/metrics、/version 及 pprof 路由。
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

// Run 启动 HTTP（和可选的 HTTPS）服务器，若启用健康检查则等待服务就绪后再返回。
func (s *APIServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	if s.SecureServingInfo != nil {
		s.secureServer = &http.Server{
			Addr:    s.SecureServingInfo.Address(),
			Handler: s,
		}
		eg.Go(func() error {
			if err := s.secureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		})
	}

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

// ping 循环请求 /healthz 接口，直到返回 200 或 ctx 超时，用于等待服务器启动完成。
func (s *APIServer) ping(ctx context.Context) error {
	addr := s.InsecureServingInfo.Address
	if strings.HasPrefix(addr, "0.0.0.0:") {
		addr = "127.0.0.1:" + addr[len("0.0.0.0:"):]
	}
	url := "http://" + addr + "/healthz"

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			ok := resp.StatusCode == http.StatusOK
			_ = resp.Body.Close()
			if ok {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}
}

// Close 优雅关闭 HTTP 和 HTTPS 服务器，等待在途请求处理完成，最多等待 ShutdownTimeout。
//
// 这里刻意使用 context.Background() 而非触发关闭的那个已取消的 ctx：
// Shutdown 需要一个仍然有效的截止时间来排空在途请求，传入已取消的 ctx 会让它立即返回。
func (s *APIServer) Close() {
	timeout := s.ShutdownTimeout
	if timeout <= 0 {
		timeout = DefaultShutdownTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if s.secureServer != nil {
		if err := s.secureServer.Shutdown(ctx); err != nil {
			logger.Warnf("Shutdown secure server failed: %s", err.Error())
		}
	}

	if s.insecureServer != nil {
		if err := s.insecureServer.Shutdown(ctx); err != nil {
			logger.Warnf("Shutdown insecure server failed: %s", err.Error())
		}
	}
}
