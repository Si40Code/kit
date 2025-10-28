package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Si40Code/kit/logger"
)

func main() {
	// 创建日志目录
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	fmt.Println("=== 示例 1: 基本文件输出 ===")
	example1BasicFileOutput()

	fmt.Println("\n=== 示例 2: 文件切割配置 ===")
	example2FileRotation()

	fmt.Println("\n=== 示例 3: 同时输出到 stdout 和文件 ===")
	example3MultipleOutputs()

	fmt.Println("\n=== 示例 4: 不同级别输出到不同文件 ===")
	example4LevelBasedFiles()

	fmt.Println("✅ 所有示例完成！请查看 ./logs 目录下的日志文件")
}

// 示例 1: 基本文件输出
func example1BasicFileOutput() {
	l, err := logger.New(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithFile("./logs/app.log"),
	)
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	ctx := context.Background()
	l.Info(ctx, "应用启动", "version", "1.0.0")
	l.Info(ctx, "配置加载完成", "config_file", "config.yaml")
	l.Warn(ctx, "磁盘空间不足", "available", "10GB")

	fmt.Println("✓ 日志已写入 ./logs/app.log")
}

// 示例 2: 文件切割配置
func example2FileRotation() {
	l, err := logger.New(
		logger.WithLevel(logger.DebugLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithFile("./logs/app-rotated.log",
			logger.WithFileMaxSize(1),        // 1MB
			logger.WithFileMaxAge(7),         // 保留 7 天
			logger.WithFileMaxBackups(3),     // 保留 3 个备份
			logger.WithFileCompress(),        // 压缩旧文件
		),
	)
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	ctx := context.Background()

	// 写入大量日志以触发切割（仅演示，实际需要更多数据）
	for i := 0; i < 100; i++ {
		l.Info(ctx, "处理请求",
			"request_id", fmt.Sprintf("req-%d", i),
			"method", "GET",
			"path", "/api/users",
			"status", 200,
			"duration", fmt.Sprintf("%dms", 50+i%100),
			"user_agent", "Mozilla/5.0 (compatible; example)",
		)
	}

	fmt.Println("✓ 日志已写入 ./logs/app-rotated.log（配置了切割规则）")
	fmt.Println("  - 文件大小超过 1MB 时自动切割")
	fmt.Println("  - 保留最近 3 个备份文件")
	fmt.Println("  - 超过 7 天的日志自动删除")
	fmt.Println("  - 旧文件自动压缩")
}

// 示例 3: 同时输出到 stdout 和文件
func example3MultipleOutputs() {
	l, err := logger.New(
		logger.WithLevel(logger.DebugLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),                // 输出到控制台
		logger.WithFile("./logs/both.log"), // 同时输出到文件
	)
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	ctx := context.Background()
	l.Info(ctx, "这条日志同时出现在控制台和文件中")
	l.Debug(ctx, "调试信息", "module", "database", "query_time", "15ms")
	l.Error(ctx, "遇到错误", "error", "connection timeout")

	fmt.Println("\n✓ 日志已输出到控制台并写入 ./logs/both.log")
}

// 示例 4: 不同级别输出到不同文件
func example4LevelBasedFiles() {
	// Info 及以上级别的日志
	infoLogger, err := logger.New(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithFile("./logs/info.log"),
	)
	if err != nil {
		panic(err)
	}
	defer infoLogger.Sync()

	// Error 及以上级别的日志
	errorLogger, err := logger.New(
		logger.WithLevel(logger.ErrorLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithFile("./logs/error.log"),
	)
	if err != nil {
		panic(err)
	}
	defer errorLogger.Sync()

	ctx := context.Background()

	// Info 日志只会出现在 info.log
	infoLogger.Info(ctx, "正常操作", "action", "create_user")
	infoLogger.Info(ctx, "正常操作", "action", "send_email")

	// Error 日志会出现在 error.log
	errorLogger.Error(ctx, "数据库连接失败", "error", "timeout")
	errorLogger.Error(ctx, "API 调用失败", "error", "rate limit exceeded")

	// 也可以同时记录到两个 logger
	infoLogger.Error(ctx, "这条 error 日志会出现在 info.log")
	errorLogger.Error(ctx, "这条 error 日志会出现在 error.log")

	fmt.Println("✓ 不同级别的日志已分别写入:")
	fmt.Println("  - ./logs/info.log  (Info 及以上)")
	fmt.Println("  - ./logs/error.log (Error 及以上)")
}

