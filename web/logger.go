package web

import (
	"context"
	"fmt"
)

// Logger 日志接口
type Logger interface {
	// Info 信息日志
	Info(ctx context.Context, msg string, fields map[string]interface{})
	// Warn 警告日志
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	// Error 错误日志
	Error(ctx context.Context, msg string, fields map[string]interface{})
}

// defaultLogger 默认日志实现（输出到标准输出）
type defaultLogger struct{}

func (l *defaultLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	// 简单实现，实际使用时应该使用专业的日志库
	fmt.Printf("[INFO] %s %v\n", msg, fields)
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("[WARN] %s %v\n", msg, fields)
}

func (l *defaultLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s %v\n", msg, fields)
}
