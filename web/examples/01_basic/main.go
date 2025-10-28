package main

import (
	"net/http"

	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建服务器
	server := web.New(
		web.WithMode(web.DebugMode),
		web.WithServiceName("basic-example"),
	)

	// 注册路由
	engine := server.Engine()
	
	engine.GET("/ping", func(c *gin.Context) {
		web.Success(c, gin.H{
			"message": "pong",
		})
	})

	engine.POST("/echo", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		web.Success(c, req)
	})

	// 启动服务器
	server.RunWithGracefulShutdown(":8080")
}
