package web

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Mode 运行模式
type Mode string

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
	TestMode    Mode = "test"
)

// Option 配置选项
type Option func(*options)

// options 配置选项结构体
type options struct {
	// 基础配置
	mode        Mode
	serviceName string

	// 日志配置
	logger            Logger
	skipPaths         []string // 跳过日志记录的路径
	maxBodyLogSize    int64    // 最大 body 日志大小（字节）
	slowRequestThresh time.Duration

	// Trace 配置
	enableTrace bool

	// Metric 配置
	enableMetric bool
	metricRecorder MetricRecorder

	// 中间件配置
	enableRecover bool
	enableCORS    bool
	middlewares   []gin.HandlerFunc

	// 响应配置
	enablePrettyJSON bool

	// 文件上传配置
	maxMultipartMemory int64 // 最大文件内存（字节）
}

// newOptions 创建默认配置
func newOptions(opts ...Option) *options {
	o := &options{
		mode:               DebugMode,
		serviceName:        "web-service",
		skipPaths:          []string{"/health", "/metrics"},
		maxBodyLogSize:     10 * 1024,  // 10KB
		slowRequestThresh:  1 * time.Second,
		enableTrace:        false,
		enableMetric:       false,
		enableRecover:      true,
		enableCORS:         false,
		enablePrettyJSON:   false,
		maxMultipartMemory: 32 << 20, // 32MB
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithMode 设置运行模式
func WithMode(mode Mode) Option {
	return func(o *options) {
		o.mode = mode
	}
}

// WithServiceName 设置服务名称
func WithServiceName(name string) Option {
	return func(o *options) {
		o.serviceName = name
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithSkipPaths 设置跳过日志记录的路径
func WithSkipPaths(paths ...string) Option {
	return func(o *options) {
		o.skipPaths = append(o.skipPaths, paths...)
	}
}

// WithMaxBodyLogSize 设置最大 body 日志大小
func WithMaxBodyLogSize(size int64) Option {
	return func(o *options) {
		o.maxBodyLogSize = size
	}
}

// WithSlowRequestThreshold 设置慢请求阈值
func WithSlowRequestThreshold(duration time.Duration) Option {
	return func(o *options) {
		o.slowRequestThresh = duration
	}
}

// WithTrace 启用链路追踪
func WithTrace() Option {
	return func(o *options) {
		o.enableTrace = true
	}
}

// WithMetric 启用指标监控
func WithMetric(recorder MetricRecorder) Option {
	return func(o *options) {
		o.enableMetric = true
		o.metricRecorder = recorder
	}
}

// WithRecover 启用 panic 恢复
func WithRecover() Option {
	return func(o *options) {
		o.enableRecover = true
	}
}

// WithCORS 启用 CORS
func WithCORS() Option {
	return func(o *options) {
		o.enableCORS = true
	}
}

// WithMiddleware 添加自定义中间件
func WithMiddleware(middlewares ...gin.HandlerFunc) Option {
	return func(o *options) {
		o.middlewares = append(o.middlewares, middlewares...)
	}
}

// WithPrettyJSON 启用格式化 JSON 响应
func WithPrettyJSON() Option {
	return func(o *options) {
		o.enablePrettyJSON = true
	}
}

// WithMaxMultipartMemory 设置最大文件上传内存
func WithMaxMultipartMemory(size int64) Option {
	return func(o *options) {
		o.maxMultipartMemory = size
	}
}
