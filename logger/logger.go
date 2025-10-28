package logger

import "context"

// Logger 日志接口
type Logger interface {
	// 结构化字段方式（key-value 对）
	Debug(ctx context.Context, msg string, fields ...any)
	Info(ctx context.Context, msg string, fields ...any)
	Warn(ctx context.Context, msg string, fields ...any)
	Error(ctx context.Context, msg string, fields ...any)
	Fatal(ctx context.Context, msg string, fields ...any)

	// Map 字段方式
	DebugMap(ctx context.Context, msg string, fields map[string]any)
	InfoMap(ctx context.Context, msg string, fields map[string]any)
	WarnMap(ctx context.Context, msg string, fields map[string]any)
	ErrorMap(ctx context.Context, msg string, fields map[string]any)
	FatalMap(ctx context.Context, msg string, fields map[string]any)

	// With 方法创建带预设字段的子 logger
	With(fields ...any) Logger

	// Sync 刷新缓冲区
	Sync() error
}

// Level 日志级别
type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

// ParseLevel 从字符串解析日志级别
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Format 日志格式
type Format string

const (
	JSONFormat    Format = "json"
	ConsoleFormat Format = "console"
)

