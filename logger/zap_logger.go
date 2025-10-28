package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger 基于 zap 的 Logger 实现
type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	opts   *options
}

// New 创建新的 logger 实例
func New(opts ...Option) (Logger, error) {
	options := newOptions(opts...)

	// 创建所有输出的 cores
	cores, err := createCores(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create cores: %w", err)
	}

	// 组合所有 cores
	core := zapcore.NewTee(cores...)

	// 创建 zap logger 选项
	zapOpts := []zap.Option{}

	if options.caller {
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if options.stacktrace {
		zapOpts = append(zapOpts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if options.development {
		zapOpts = append(zapOpts, zap.Development())
	}

	// 创建 zap.Logger
	zlog := zap.New(core, zapOpts...)

	return &zapLogger{
		logger: zlog,
		sugar:  zlog.Sugar(),
		opts:   options,
	}, nil
}

// Debug 记录 debug 级别日志
func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...any) {
	l.log(ctx, DebugLevel, msg, fields...)
}

// Info 记录 info 级别日志
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...any) {
	l.log(ctx, InfoLevel, msg, fields...)
}

// Warn 记录 warn 级别日志
func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...any) {
	l.log(ctx, WarnLevel, msg, fields...)
}

// Error 记录 error 级别日志
func (l *zapLogger) Error(ctx context.Context, msg string, fields ...any) {
	// Error 日志需要标记 span 为 error
	if l.opts.enableTrace {
		markSpanError(ctx, msg)
	}
	l.log(ctx, ErrorLevel, msg, fields...)
}

// Fatal 记录 fatal 级别日志
func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...any) {
	// Fatal 日志也需要标记 span 为 error
	if l.opts.enableTrace {
		markSpanError(ctx, msg)
	}
	l.log(ctx, FatalLevel, msg, fields...)
}

// DebugMap 记录 debug 级别日志（map 字段）
func (l *zapLogger) DebugMap(ctx context.Context, msg string, fields map[string]any) {
	l.logMap(ctx, DebugLevel, msg, fields)
}

// InfoMap 记录 info 级别日志（map 字段）
func (l *zapLogger) InfoMap(ctx context.Context, msg string, fields map[string]any) {
	l.logMap(ctx, InfoLevel, msg, fields)
}

// WarnMap 记录 warn 级别日志（map 字段）
func (l *zapLogger) WarnMap(ctx context.Context, msg string, fields map[string]any) {
	l.logMap(ctx, WarnLevel, msg, fields)
}

// ErrorMap 记录 error 级别日志（map 字段）
func (l *zapLogger) ErrorMap(ctx context.Context, msg string, fields map[string]any) {
	if l.opts.enableTrace {
		markSpanError(ctx, msg)
	}
	l.logMap(ctx, ErrorLevel, msg, fields)
}

// FatalMap 记录 fatal 级别日志（map 字段）
func (l *zapLogger) FatalMap(ctx context.Context, msg string, fields map[string]any) {
	if l.opts.enableTrace {
		markSpanError(ctx, msg)
	}
	l.logMap(ctx, FatalLevel, msg, fields)
}

// With 创建带预设字段的子 logger
func (l *zapLogger) With(fields ...any) Logger {
	return &zapLogger{
		logger: l.logger.With(l.convertToZapFields(fields...)...),
		sugar:  l.sugar.With(fields...),
		opts:   l.opts,
	}
}

// Sync 刷新缓冲区
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// log 内部日志记录方法（结构化字段）
func (l *zapLogger) log(ctx context.Context, level Level, msg string, fields ...any) {
	// 添加 trace 信息
	var zapFields []zap.Field
	if l.opts.enableTrace {
		zapFields = extractTraceFields(ctx)
		addSpanEvent(ctx, level, msg)
	}

	// 转换字段为 zap.Field
	zapFields = append(zapFields, l.convertToZapFields(fields...)...)

	// 记录日志
	switch level {
	case DebugLevel:
		l.logger.Debug(msg, zapFields...)
	case InfoLevel:
		l.logger.Info(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case FatalLevel:
		l.logger.Fatal(msg, zapFields...)
	}
}

// logMap 内部日志记录方法（map 字段）
func (l *zapLogger) logMap(ctx context.Context, level Level, msg string, fields map[string]any) {
	// 添加 trace 信息
	var zapFields []zap.Field
	if l.opts.enableTrace {
		zapFields = extractTraceFields(ctx)
		addSpanEvent(ctx, level, msg)
	}

	// 转换 map 为 zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	// 记录日志
	switch level {
	case DebugLevel:
		l.logger.Debug(msg, zapFields...)
	case InfoLevel:
		l.logger.Info(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case FatalLevel:
		l.logger.Fatal(msg, zapFields...)
	}
}

// convertToZapFields 转换字段为 zap.Field
func (l *zapLogger) convertToZapFields(fields ...any) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	var zapFields []zap.Field

	// 如果字段是偶数，尝试解析为 key-value 对
	if len(fields)%2 == 0 {
		for i := 0; i < len(fields); i += 2 {
			if key, ok := fields[i].(string); ok {
				zapFields = append(zapFields, zap.Any(key, fields[i+1]))
			}
		}
	}

	return zapFields
}

