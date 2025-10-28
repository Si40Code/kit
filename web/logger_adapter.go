package web

import (
	"context"

	"github.com/Si40Code/kit/logger"
)

// LoggerAdapter 将 kit/logger.Logger 适配到 web.Logger 接口
type LoggerAdapter struct {
	logger logger.Logger
}

// NewLoggerAdapter 创建新的 logger 适配器
func NewLoggerAdapter(l logger.Logger) Logger {
	return &LoggerAdapter{logger: l}
}

// Info 记录 info 级别日志
func (a *LoggerAdapter) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	a.logger.InfoMap(ctx, msg, convertFields(fields))
}

// Warn 记录 warn 级别日志
func (a *LoggerAdapter) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	a.logger.WarnMap(ctx, msg, convertFields(fields))
}

// Error 记录 error 级别日志
func (a *LoggerAdapter) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	a.logger.ErrorMap(ctx, msg, convertFields(fields))
}

// convertFields 转换字段类型
func convertFields(fields map[string]interface{}) map[string]any {
	if fields == nil {
		return nil
	}

	result := make(map[string]any, len(fields))
	for k, v := range fields {
		result[k] = v
	}
	return result
}

