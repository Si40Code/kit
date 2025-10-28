package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// runWithGracefulShutdown 启动服务器并支持优雅关闭
func runWithGracefulShutdown(engine *gin.Engine, addr string, opts *options) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 启动服务器
	go func() {
		if opts.logger != nil {
			opts.logger.Info(context.Background(), "Starting server with graceful shutdown", map[string]interface{}{
				"address": addr,
			})
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if opts.logger != nil {
				opts.logger.Error(context.Background(), "Server error", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if opts.logger != nil {
		opts.logger.Info(context.Background(), "Shutting down server...", nil)
	}

	// 给 5 秒时间完成现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		if opts.logger != nil {
			opts.logger.Error(context.Background(), "Server forced to shutdown", map[string]interface{}{
				"error": err.Error(),
			})
		}
		return err
	}

	if opts.logger != nil {
		opts.logger.Info(context.Background(), "Server exited gracefully", nil)
	}
	return nil
}
