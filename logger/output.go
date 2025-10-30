package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	zapotlpencoder "github.com/SigNoz/zap_otlp/zap_otlp_encoder"
	zapotlpsync "github.com/SigNoz/zap_otlp/zap_otlp_sync"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/natefinch/lumberjack.v2"
)

// createCores 创建所有输出的 cores
func createCores(opts *options) ([]zapcore.Core, error) {
	var cores []zapcore.Core
	level := zapLevel(opts.level)

	for _, output := range opts.outputs {
		var core zapcore.Core
		var err error

		switch output.Type {
		case StdoutOutput:
			core, err = createStdoutCore(opts.format, opts.development, level)
		case FileOutput:
			core, err = createFileCore(output.Config, opts.format, opts.development, level)
		case OTLPOutput:
			core, err = createOTLPCore(output.Config, opts.serviceName, opts.resourceAttributes, level)
		default:
			return nil, fmt.Errorf("unknown output type: %s", output.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create %s core: %w", output.Type, err)
		}

		cores = append(cores, core)
	}

	if len(cores) == 0 {
		// 如果没有配置任何输出，默认使用 stdout
		core, err := createStdoutCore(opts.format, opts.development, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, core)
	}

	return cores, nil
}

// createStdoutCore 创建标准输出 core
func createStdoutCore(format Format, development bool, level zapcore.Level) (zapcore.Core, error) {
	encoder := createEncoder(format, development)
	writer := zapcore.Lock(os.Stdout)
	return zapcore.NewCore(encoder, writer, level), nil
}

// createFileCore 创建文件输出 core
func createFileCore(cfg OutputConfig, format Format, development bool, level zapcore.Level) (zapcore.Core, error) {
	if cfg.FilePath == "" {
		return nil, fmt.Errorf("file path is required for file output")
	}

	// 使用 lumberjack 实现日志切割
	writer := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,    // MB
		MaxAge:     cfg.MaxAge,     // days
		MaxBackups: cfg.MaxBackups, // 保留文件数
		Compress:   cfg.Compress,   // 是否压缩
		LocalTime:  true,           // 使用本地时间
	}

	encoder := createEncoder(format, development)
	writeSyncer := zapcore.AddSync(writer)

	return zapcore.NewCore(encoder, writeSyncer, level), nil
}

// createOTLPCore 创建 OTLP 输出 core（用于 SigNoz）
func createOTLPCore(cfg OutputConfig, serviceName string, resourceAttrs []attribute.KeyValue, level zapcore.Level) (zapcore.Core, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OTLP output")
	}

	if serviceName == "" {
		serviceName = "unknown-service"
	}

	// 创建 gRPC 连接选项
	var dialOpts []grpc.DialOption

	if cfg.Insecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	// 设置超时
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	dialOpts = append(dialOpts, grpc.WithTimeout(timeout), grpc.WithBlock())

	// 建立连接
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Endpoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OTLP endpoint: %w", err)
	}

	// 创建 OTLP encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	otlpEncoder := zapotlpencoder.NewOTLPEncoder(encoderConfig)

	// 构建 resource attributes
	resourceAttributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
	}
	// 添加自定义 resource attributes
	resourceAttributes = append(resourceAttributes, resourceAttrs...)

	// 创建 OTLP Syncer
	otlpSyncer := zapotlpsync.NewOtlpSyncer(conn, zapotlpsync.Options{
		BatchSize:      100,
		BatchInterval:  5 * time.Second,
		ResourceSchema: semconv.SchemaURL,
		Resource: resource.NewWithAttributes(
			semconv.SchemaURL,
			resourceAttributes...,
		),
	})

	return zapcore.NewCore(otlpEncoder, otlpSyncer, level), nil
}
