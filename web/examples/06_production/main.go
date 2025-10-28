package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
)

// ProdLogger 生产环境日志记录器（示例）
type ProdLogger struct {
	env string
}

func NewProdLogger(env string) *ProdLogger {
	return &ProdLogger{env: env}
}

func (l *ProdLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	// 在生产环境中，应该使用专业的日志库（如 zap, logrus）
	fmt.Printf("%s | ENV=%s | %s | %v\n", time.Now().Format("2006-01-02 15:04:05"), l.env, msg, fields)
}

func (l *ProdLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("%s | ENV=%s | %s | %v\n", time.Now().Format("2006-01-02 15:04:05"), l.env, msg, fields)
}

func (l *ProdLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("%s | ENV=%s | %s | %v\n", time.Now().Format("2006-01-02 15:04:05"), l.env, msg, fields)
}

func main() {
	// 从环境变量读取配置
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// 根据环境选择日志级别
	var logger web.Logger
	if env == "production" {
		logger = NewProdLogger("prod")
	} else {
		logger = NewProdLogger("dev")
	}

	// 生产环境配置
	server := web.New(
		web.WithMode(web.ReleaseMode),
		web.WithServiceName("production-app"),
		web.WithLogger(logger),
		web.WithSkipPaths("/health", "/metrics", "/ready"), // 跳过健康检查的日志
		web.WithMaxBodyLogSize(4096),                      // 限制日志大小
		web.WithSlowRequestThreshold(3*time.Second),       // 慢请求阈值
		web.WithTrace(),                                   // 启用链路追踪
	)

	engine := server.Engine()

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		web.Success(c, gin.H{"status": "ok"})
	})

	// 就绪检查
	engine.GET("/ready", func(c *gin.Context) {
		web.Success(c, gin.H{"status": "ready"})
	})

	// 业务路由
	engine.GET("/api/users", func(c *gin.Context) {
		web.Success(c, gin.H{
			"users": []gin.H{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		})
	})

	engine.POST("/api/users", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, 400, err.Error())
			return
		}

		// 模拟创建用户
		web.Success(c, gin.H{
			"id":   3,
			"user": req,
		})
	})

	fmt.Printf("Starting production server on :8080 (env=%s)\n", env)
	server.RunWithGracefulShutdown(":8080")
}
