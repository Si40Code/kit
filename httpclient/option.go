package httpclient

import (
	"time"

	"github.com/Si40Code/kit/logger"
)

// Option 配置选项函数
type Option func(*options)

// options 配置选项结构体
type options struct {
	// 日志配置
	logger         logger.Logger
	enableLog      bool
	maxBodyLogSize int64 // 最大 body 日志大小（字节）

	// Trace 配置
	enableTrace bool
	serviceName string

	// Metric 配置
	enableMetric   bool
	metricRecorder MetricRecorder

	// HTTP 客户端配置
	timeout            time.Duration
	retryCount         int
	retryWaitTime      time.Duration
	retryMaxWaitTime   time.Duration
	insecureSkipVerify bool

	// 高级配置
	maxIdleConns        int
	idleConnTimeout     time.Duration
	tlsHandshakeTimeout time.Duration
	keepAlive           time.Duration
}

// newOptions 创建默认配置
func newOptions(opts ...Option) *options {
	o := &options{
		enableLog:           true,
		maxBodyLogSize:      10 * 1024, // 10KB
		enableTrace:         false,
		serviceName:         "http-client",
		enableMetric:        false,
		timeout:             30 * time.Second,
		retryCount:          0,
		retryWaitTime:       100 * time.Millisecond,
		retryMaxWaitTime:    2 * time.Second,
		insecureSkipVerify:  false,
		maxIdleConns:        100,
		idleConnTimeout:     90 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
		keepAlive:           30 * time.Second,
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

// WithMaxBodyLogSize 设置最大 body 日志大小
func WithMaxBodyLogSize(size int64) Option {
	return func(o *options) {
		o.maxBodyLogSize = size
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

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithRetry 设置重试配置
func WithRetry(count int, waitTime, maxWaitTime time.Duration) Option {
	return func(o *options) {
		o.retryCount = count
		o.retryWaitTime = waitTime
		o.retryMaxWaitTime = maxWaitTime
	}
}

// WithInsecureSkipVerify 跳过 TLS 证书验证（不推荐在生产环境使用）
func WithInsecureSkipVerify() Option {
	return func(o *options) {
		o.insecureSkipVerify = true
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(count int) Option {
	return func(o *options) {
		o.maxIdleConns = count
	}
}

// WithIdleConnTimeout 设置空闲连接超时时间
func WithIdleConnTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.idleConnTimeout = timeout
	}
}

// WithTLSHandshakeTimeout 设置 TLS 握手超时时间
func WithTLSHandshakeTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.tlsHandshakeTimeout = timeout
	}
}

// WithKeepAlive 设置 keep-alive 时间
func WithKeepAlive(duration time.Duration) Option {
	return func(o *options) {
		o.keepAlive = duration
	}
}
