package logger

import (
	"context"
	"fmt"
	"sync"
)

var (
	defaultLogger Logger
	mu            sync.RWMutex
)

func init() {
	// 初始化一个默认的 logger（stdout，info 级别）
	l, err := New(WithLevel(InfoLevel), WithStdout())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize default logger: %v", err))
	}
	defaultLogger = l
}

// Init 初始化全局默认 logger
func Init(opts ...Option) error {
	l, err := New(opts...)
	if err != nil {
		return err
	}

	mu.Lock()
	defaultLogger = l
	mu.Unlock()

	return nil
}

// Default 返回全局默认 logger
func Default() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLogger
}

// SetDefault 设置全局默认 logger
func SetDefault(l Logger) {
	mu.Lock()
	defaultLogger = l
	mu.Unlock()
}

// 包级便捷函数 - 使用默认 logger

// Debug 记录 debug 级别日志
func Debug(ctx context.Context, msg string, fields ...any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.Debug(ctx, msg, fields...)
}

// Info 记录 info 级别日志
func Info(ctx context.Context, msg string, fields ...any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.Info(ctx, msg, fields...)
}

// Warn 记录 warn 级别日志
func Warn(ctx context.Context, msg string, fields ...any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.Warn(ctx, msg, fields...)
}

// Error 记录 error 级别日志
func Error(ctx context.Context, msg string, fields ...any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.Error(ctx, msg, fields...)
}

// Fatal 记录 fatal 级别日志
func Fatal(ctx context.Context, msg string, fields ...any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.Fatal(ctx, msg, fields...)
}

// DebugMap 记录 debug 级别日志（map 字段）
func DebugMap(ctx context.Context, msg string, fields map[string]any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.DebugMap(ctx, msg, fields)
}

// InfoMap 记录 info 级别日志（map 字段）
func InfoMap(ctx context.Context, msg string, fields map[string]any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.InfoMap(ctx, msg, fields)
}

// WarnMap 记录 warn 级别日志（map 字段）
func WarnMap(ctx context.Context, msg string, fields map[string]any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.WarnMap(ctx, msg, fields)
}

// ErrorMap 记录 error 级别日志（map 字段）
func ErrorMap(ctx context.Context, msg string, fields map[string]any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.ErrorMap(ctx, msg, fields)
}

// FatalMap 记录 fatal 级别日志（map 字段）
func FatalMap(ctx context.Context, msg string, fields map[string]any) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	l.FatalMap(ctx, msg, fields)
}

// With 创建带预设字段的子 logger
func With(fields ...any) Logger {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	return l.With(fields...)
}

// Sync 刷新默认 logger 的缓冲区
func Sync() error {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	return l.Sync()
}
