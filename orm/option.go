package orm

import (
	"time"

	"github.com/Si40Code/kit/logger"
)

// Option 配置选项函数
type Option func(*options)

// options 配置选项结构体
type options struct {
	// 日志配置
	logger               logger.Logger
	enableLog            bool
	slowThreshold        time.Duration // 慢查询阈值
	ignoreRecordNotFound bool          // 全局配置：忽略 RecordNotFound 错误

	// Trace 配置
	enableTrace bool
	serviceName string

	// Metric 配置
	enableMetric   bool
	metricRecorder MetricRecorder

	// 数据库连接池配置
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

// newOptions 创建默认配置
func newOptions(opts ...Option) *options {
	o := &options{
		enableLog:            true,
		slowThreshold:        200 * time.Millisecond, // 默认 200ms 为慢查询
		ignoreRecordNotFound: false,
		enableTrace:          false,
		serviceName:          "orm-client",
		enableMetric:         false,
		maxIdleConns:         10,
		maxOpenConns:         100,
		connMaxLifetime:      time.Hour,
		connMaxIdleTime:      10 * time.Minute,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithLogger 设置日志记录器
func WithLogger(l logger.Logger) Option {
	return func(o *options) {
		o.logger = l
		o.enableLog = true
	}
}

// WithDisableLog 禁用日志记录
func WithDisableLog() Option {
	return func(o *options) {
		o.enableLog = false
	}
}

// WithSlowThreshold 设置慢查询阈值
func WithSlowThreshold(threshold time.Duration) Option {
	return func(o *options) {
		o.slowThreshold = threshold
	}
}

// WithIgnoreRecordNotFoundError 全局配置：查询无数据时不返回错误
func WithIgnoreRecordNotFoundError() Option {
	return func(o *options) {
		o.ignoreRecordNotFound = true
	}
}

// WithTrace 启用链路追踪
func WithTrace(serviceName string) Option {
	return func(o *options) {
		o.enableTrace = true
		o.serviceName = serviceName
	}
}

// WithMetric 启用指标监控
func WithMetric(recorder MetricRecorder) Option {
	return func(o *options) {
		o.enableMetric = true
		o.metricRecorder = recorder
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(count int) Option {
	return func(o *options) {
		o.maxIdleConns = count
	}
}

// WithMaxOpenConns 设置最大打开连接数
func WithMaxOpenConns(count int) Option {
	return func(o *options) {
		o.maxOpenConns = count
	}
}

// WithConnMaxLifetime 设置连接最大生命周期
func WithConnMaxLifetime(lifetime time.Duration) Option {
	return func(o *options) {
		o.connMaxLifetime = lifetime
	}
}

// WithConnMaxIdleTime 设置连接最大空闲时间
func WithConnMaxIdleTime(idleTime time.Duration) Option {
	return func(o *options) {
		o.connMaxIdleTime = idleTime
	}
}
