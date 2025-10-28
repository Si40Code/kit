package main

import (
	"context"
	"time"

	"github.com/Si40Code/kit/logger"
)

func main() {
	// 示例 1: 使用默认 logger（已自动初始化）
	ctx := context.Background()
	logger.Info(ctx, "使用默认 logger")

	// 示例 2: 初始化自定义 logger
	err := logger.Init(
		logger.WithLevel(logger.DebugLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
	)
	if err != nil {
		panic(err)
	}

	// 示例 3: 不同日志级别
	logger.Debug(ctx, "这是 debug 日志")
	logger.Info(ctx, "这是 info 日志")
	logger.Warn(ctx, "这是 warn 日志")
	logger.Error(ctx, "这是 error 日志")

	// 示例 4: 结构化字段（key-value 对）
	logger.Info(ctx, "用户登录",
		"user_id", 12345,
		"username", "alice",
		"ip", "192.168.1.1",
	)

	// 示例 5: Map 字段方式
	logger.InfoMap(ctx, "订单创建", map[string]any{
		"order_id": "ORD-2024-001",
		"amount":   99.99,
		"currency": "USD",
		"items":    3,
	})

	// 示例 6: 使用 With 创建子 logger
	userLogger := logger.With("user_id", 12345, "session", "abc123")
	userLogger.Info(ctx, "用户执行操作", "action", "update_profile")
	userLogger.Info(ctx, "用户执行操作", "action", "change_password")

	// 示例 7: 创建独立 logger 实例
	customLogger, err := logger.New(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithStdout(),
	)
	if err != nil {
		panic(err)
	}

	customLogger.Info(ctx, "这是独立的 logger 实例")

	// 示例 8: 性能测试
	start := time.Now()
	for i := 0; i < 1000; i++ {
		logger.Info(ctx, "性能测试", "iteration", i)
	}
	elapsed := time.Since(start)
	logger.Info(ctx, "性能测试完成",
		"iterations", 1000,
		"elapsed", elapsed.String(),
		"avg_per_log", (elapsed / 1000).String(),
	)

	// 确保所有日志都被刷新
	logger.Sync()
}

