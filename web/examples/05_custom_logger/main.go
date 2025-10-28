package main

import (
	"context"
	"fmt"

	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
)

// CustomLogger 自定义日志实现
type CustomLogger struct {
	prefix string
}

func NewCustomLogger(prefix string) *CustomLogger {
	return &CustomLogger{prefix: prefix}
}

func (l *CustomLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("[%s][INFO] %s: %v\n", l.prefix, msg, fields)
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("[%s][WARN] %s: %v\n", l.prefix, msg, fields)
}

func (l *CustomLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("[%s][ERROR] %s: %v\n", l.prefix, msg, fields)
}

func main() {
	// 创建自定义日志记录器
	customLogger := NewCustomLogger("MY-APP")

	// 创建服务器，使用自定义日志
	server := web.New(
		web.WithMode(web.DebugMode),
		web.WithServiceName("custom-logger-example"),
		web.WithLogger(customLogger),
		web.WithMaxBodyLogSize(2048), // 限制日志大小
	)

	engine := server.Engine()

	engine.GET("/hello", func(c *gin.Context) {
		web.Success(c, gin.H{
			"message": "Hello with custom logger!",
		})
	})

	engine.POST("/test", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, 400, err.Error())
			return
		}
		web.Success(c, gin.H{"received": req})
	})

	server.RunWithGracefulShutdown(":8080")
}
