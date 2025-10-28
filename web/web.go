package web

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Server Gin 服务器封装
type Server struct {
	engine  *gin.Engine
	options *options
}

// New 创建一个新的 Gin 服务器
func New(opts ...Option) *Server {
	// 应用选项
	options := newOptions(opts...)

	// 根据模式设置 Gin
	if options.mode == ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else if options.mode == TestMode {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建 Gin Engine
	engine := gin.New()

	// 应用中间件
	server := &Server{
		engine:  engine,
		options: options,
	}

	// 1. Recovery 中间件（必须第一个）
	if options.enableRecover {
		engine.Use(server.recoveryMiddleware())
	}

	// 2. CORS 中间件
	if options.enableCORS {
		engine.Use(server.corsMiddleware())
	}

	// 3. Trace 中间件
	if options.enableTrace {
		engine.Use(otelgin.Middleware(options.serviceName))
	}

	// 4. 请求日志和指标中间件（核心）
	engine.Use(server.loggingMiddleware())

	// 5. 用户自定义中间件
	for _, mw := range options.middlewares {
		engine.Use(mw)
	}

	return server
}

// Engine 返回底层的 Gin Engine，用于注册路由
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Run 启动服务器
func (s *Server) Run(addr string) error {
	if s.options.logger != nil {
		s.options.logger.Info(context.Background(), "Starting server", map[string]interface{}{
			"address":      addr,
			"mode":         s.options.mode,
			"service_name": s.options.serviceName,
		})
	}
	return s.engine.Run(addr)
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	// Gin 本身不支持优雅关闭，需要配合 http.Server 使用
	// 这里只是一个占位符，实际使用时需要使用 http.Server
	if s.options.logger != nil {
		s.options.logger.Info(ctx, "Server shutdown", nil)
	}
	return nil
}

// RunWithGracefulShutdown 启动服务器并支持优雅关闭
func (s *Server) RunWithGracefulShutdown(addr string) error {
	return runWithGracefulShutdown(s.engine, addr, s.options)
}
