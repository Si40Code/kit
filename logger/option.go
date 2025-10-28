package logger

import "time"

// Option 配置选项函数
type Option func(*options)

// options 配置选项结构体
type options struct {
	level       Level
	format      Format
	outputs     []Output
	enableTrace bool
	serviceName string
	development bool
	caller      bool
	stacktrace  bool
}

// Output 输出配置
type Output struct {
	Type   OutputType
	Config OutputConfig
}

// OutputType 输出类型
type OutputType string

const (
	StdoutOutput OutputType = "stdout"
	FileOutput   OutputType = "file"
	OTLPOutput   OutputType = "otlp"
)

// OutputConfig 输出配置详情
type OutputConfig struct {
	// File 配置
	FilePath   string
	MaxSize    int  // MB
	MaxAge     int  // days
	MaxBackups int  // 保留文件数量
	Compress   bool // 是否压缩

	// OTLP 配置
	Endpoint string
	Insecure bool
	Headers  map[string]string
	Timeout  time.Duration
}

// newOptions 创建默认配置
func newOptions(opts ...Option) *options {
	o := &options{
		level:       InfoLevel,
		format:      ConsoleFormat,
		outputs:     []Output{{Type: StdoutOutput}},
		enableTrace: false,
		development: false,
		caller:      true,
		stacktrace:  false,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithLevel 设置日志级别
func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithFormat 设置日志格式
func WithFormat(format Format) Option {
	return func(o *options) {
		o.format = format
	}
}

// WithStdout 添加标准输出
func WithStdout() Option {
	return func(o *options) {
		o.outputs = append(o.outputs, Output{
			Type: StdoutOutput,
		})
	}
}

// FileOption 文件选项
type FileOption func(*OutputConfig)

// WithFileMaxSize 设置文件最大大小（MB）
func WithFileMaxSize(size int) FileOption {
	return func(c *OutputConfig) {
		c.MaxSize = size
	}
}

// WithFileMaxAge 设置文件最大保留时间（天）
func WithFileMaxAge(days int) FileOption {
	return func(c *OutputConfig) {
		c.MaxAge = days
	}
}

// WithFileMaxBackups 设置最大备份文件数量
func WithFileMaxBackups(count int) FileOption {
	return func(c *OutputConfig) {
		c.MaxBackups = count
	}
}

// WithFileCompress 启用文件压缩
func WithFileCompress() FileOption {
	return func(c *OutputConfig) {
		c.Compress = true
	}
}

// WithFile 添加文件输出
func WithFile(path string, opts ...FileOption) Option {
	return func(o *options) {
		cfg := OutputConfig{
			FilePath:   path,
			MaxSize:    100,  // 默认 100MB
			MaxAge:     7,    // 默认保留 7 天
			MaxBackups: 3,    // 默认保留 3 个备份
			Compress:   false,
		}

		for _, opt := range opts {
			opt(&cfg)
		}

		o.outputs = append(o.outputs, Output{
			Type:   FileOutput,
			Config: cfg,
		})
	}
}

// OTLPOption OTLP 选项
type OTLPOption func(*OutputConfig)

// WithOTLPInsecure 使用不安全连接
func WithOTLPInsecure() OTLPOption {
	return func(c *OutputConfig) {
		c.Insecure = true
	}
}

// WithOTLPHeaders 设置自定义 headers
func WithOTLPHeaders(headers map[string]string) OTLPOption {
	return func(c *OutputConfig) {
		c.Headers = headers
	}
}

// WithOTLPTimeout 设置连接超时
func WithOTLPTimeout(timeout time.Duration) OTLPOption {
	return func(c *OutputConfig) {
		c.Timeout = timeout
	}
}

// WithOTLP 添加 OTLP 输出（用于 SigNoz 等）
func WithOTLP(endpoint string, opts ...OTLPOption) Option {
	return func(o *options) {
		cfg := OutputConfig{
			Endpoint: endpoint,
			Insecure: true,
			Timeout:  5 * time.Second,
		}

		for _, opt := range opts {
			opt(&cfg)
		}

		o.outputs = append(o.outputs, Output{
			Type:   OTLPOutput,
			Config: cfg,
		})
	}
}

// WithTrace 启用 trace 集成
func WithTrace(serviceName string) Option {
	return func(o *options) {
		o.enableTrace = true
		o.serviceName = serviceName
	}
}

// WithDevelopment 启用开发模式
func WithDevelopment() Option {
	return func(o *options) {
		o.development = true
		o.level = DebugLevel
		o.format = ConsoleFormat
		o.stacktrace = true
	}
}

// WithCaller 启用调用者信息
func WithCaller(enabled bool) Option {
	return func(o *options) {
		o.caller = enabled
	}
}

// WithStacktrace 启用堆栈跟踪（Error 级别及以上）
func WithStacktrace(enabled bool) Option {
	return func(o *options) {
		o.stacktrace = enabled
	}
}

