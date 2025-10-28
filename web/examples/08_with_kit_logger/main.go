package main

import (
	"log"

	"github.com/Si40Code/kit/logger"
	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化 kit logger
	err := logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithStdout(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// 2. 创建 web server（使用 kit logger 适配器）
	server := web.New(
		web.WithMode(web.ReleaseMode),
		web.WithServiceName("kit-logger-example"),
		web.WithLogger(web.NewLoggerAdapter(logger.Default())),
	)

	engine := server.Engine()

	// 3. 路由示例
	engine.GET("/api/hello", func(c *gin.Context) {
		ctx := c.Request.Context()

		// 使用 kit logger 记录日志
		logger.Info(ctx, "处理请求",
			"path", "/api/hello",
			"method", "GET",
		)

		web.Success(c, gin.H{
			"message": "Hello from kit logger!",
		})
	})

	engine.GET("/api/error", func(c *gin.Context) {
		ctx := c.Request.Context()

		// 模拟错误
		logger.Error(ctx, "处理请求失败",
			"path", "/api/error",
			"error", "something went wrong",
		)

		web.Error(c, 500, "Internal Server Error")
	})

	// 4. 启动服务器
	log.Printf("🚀 Server starting on :8080")
	log.Printf("📖 使用 kit logger 模块")
	log.Printf("GET http://localhost:8080/api/hello")
	log.Printf("GET http://localhost:8080/api/error")

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

