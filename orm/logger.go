package orm

import (
	"context"
	"errors"
	"time"

	"github.com/Si40Code/kit/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// gormLogger GORM 日志适配器，桥接到 kit/logger
type gormLogger struct {
	logger               logger.Logger
	slowThreshold        time.Duration
	ignoreRecordNotFound bool
}

// newGormLogger 创建新的 GORM 日志适配器
func newGormLogger(l logger.Logger, slowThreshold time.Duration, ignoreRecordNotFound bool) *gormLogger {
	return &gormLogger{
		logger:               l,
		slowThreshold:        slowThreshold,
		ignoreRecordNotFound: ignoreRecordNotFound,
	}
}

// LogMode 实现 GORM logger.Interface
func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info 记录 Info 级别日志
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logger != nil {
		l.logger.Info(ctx, msg, "data", data)
	}
}

// Warn 记录 Warn 级别日志
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logger != nil {
		l.logger.Warn(ctx, msg, "data", data)
	}
}

// Error 记录 Error 级别日志
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logger != nil {
		l.logger.Error(ctx, msg, "data", data)
	}
}

// Trace 记录 SQL 执行日志（核心方法）
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logger == nil {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := map[string]any{
		"duration_ms":   elapsed.Milliseconds(),
		"rows_affected": rows,
		"sql":           sql,
	}

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFound):
		// 错误日志（如果不是 RecordNotFound 或者没有配置忽略）
		fields["error"] = err.Error()
		l.logger.ErrorMap(ctx, "database query error", fields)

	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		// 慢查询日志
		fields["slow_threshold_ms"] = l.slowThreshold.Milliseconds()
		l.logger.WarnMap(ctx, "slow query detected", fields)

	default:
		// 正常查询日志
		l.logger.InfoMap(ctx, "database query executed", fields)
	}
}
