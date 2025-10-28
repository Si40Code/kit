package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Si40Code/kit/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Config 应用配置（生产环境应从配置文件或环境变量读取）
type Config struct {
	ServiceName string
	Environment string
	LogLevel    string
	LogFormat   string

	// 文件输出配置
	LogFilePath    string
	LogFileMaxSize int
	LogFileMaxAge  int
	LogFileBackups int

	// OTLP 配置
	OTLPEndpoint string
	OTLPInsecure bool

	// Trace 配置
	EnableTrace       bool
	TraceEndpoint     string
	TraceSampleRate   float64
}

func main() {
	// 加载配置
	cfg := loadConfig()

	// 初始化 logger
	if err := initLogger(cfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 初始化 tracer（如果启用）
	var cleanupTracer func()
	if cfg.EnableTrace {
		cleanupTracer = initTracer(cfg)
		defer cleanupTracer()
	}

	ctx := context.Background()

	// 应用启动日志
	logger.Info(ctx, "应用启动",
		"service", cfg.ServiceName,
		"environment", cfg.Environment,
		"version", "1.0.0",
	)

	// 运行应用
	if err := runApplication(ctx, cfg); err != nil {
		logger.Error(ctx, "应用运行失败", "error", err.Error())
		os.Exit(1)
	}

	logger.Info(ctx, "应用正常退出")
}

// loadConfig 加载配置（生产环境应从配置文件或环境变量读取）
func loadConfig() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "my-production-service"),
		Environment: getEnv("ENV", "production"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		LogFormat:   getEnv("LOG_FORMAT", "json"),

		// 文件输出
		LogFilePath:    getEnv("LOG_FILE_PATH", "/var/log/myapp/app.log"),
		LogFileMaxSize: getEnvInt("LOG_FILE_MAX_SIZE", 100),
		LogFileMaxAge:  getEnvInt("LOG_FILE_MAX_AGE", 30),
		LogFileBackups: getEnvInt("LOG_FILE_BACKUPS", 10),

		// OTLP（SigNoz/Jaeger）
		OTLPEndpoint: getEnv("OTLP_ENDPOINT", ""),
		OTLPInsecure: getEnvBool("OTLP_INSECURE", true),

		// Trace
		EnableTrace:     getEnvBool("ENABLE_TRACE", true),
		TraceEndpoint:   getEnv("TRACE_ENDPOINT", "localhost:4318"),
		TraceSampleRate: 1.0,
	}
}

// initLogger 初始化 logger
func initLogger(cfg *Config) error {
	opts := []logger.Option{
		logger.WithLevel(logger.ParseLevel(cfg.LogLevel)),
		logger.WithFormat(logger.Format(cfg.LogFormat)),
	}

	// 生产环境配置
	if cfg.Environment == "production" {
		// 1. 输出到文件（主要日志）
		opts = append(opts, logger.WithFile(cfg.LogFilePath,
			logger.WithFileMaxSize(cfg.LogFileMaxSize),
			logger.WithFileMaxAge(cfg.LogFileMaxAge),
			logger.WithFileMaxBackups(cfg.LogFileBackups),
			logger.WithFileCompress(),
		))

		// 2. Error 日志单独输出（可选）
		errorLogPath := cfg.LogFilePath + ".error"
		opts = append(opts, logger.WithFile(errorLogPath,
			logger.WithFileMaxSize(cfg.LogFileMaxSize),
			logger.WithFileMaxAge(cfg.LogFileMaxAge),
			logger.WithFileMaxBackups(cfg.LogFileBackups),
			logger.WithFileCompress(),
		))

		// 3. 同时输出到 stdout（便于容器日志收集）
		opts = append(opts, logger.WithStdout())

	} else {
		// 开发/测试环境：只输出到 stdout
		opts = append(opts,
			logger.WithStdout(),
			logger.WithDevelopment(),
		)
	}

	// 4. OTLP 输出（如果配置）
	if cfg.OTLPEndpoint != "" {
		otlpOpts := []logger.OTLPOption{}
		if cfg.OTLPInsecure {
			otlpOpts = append(otlpOpts, logger.WithOTLPInsecure())
		}
		opts = append(opts, logger.WithOTLP(cfg.OTLPEndpoint, otlpOpts...))
	}

	// 5. Trace 集成
	if cfg.EnableTrace {
		opts = append(opts, logger.WithTrace(cfg.ServiceName))
	}

	return logger.Init(opts...)
}

// initTracer 初始化 tracer
func initTracer(cfg *Config) func() {
	ctx := context.Background()

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(cfg.TraceEndpoint),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.TraceSampleRate)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error(context.Background(), "Tracer shutdown error", "error", err.Error())
		}
	}
}

// runApplication 运行应用主逻辑
func runApplication(ctx context.Context, cfg *Config) error {
	tracer := otel.Tracer(cfg.ServiceName)

	// 创建根 span
	ctx, span := tracer.Start(ctx, "application-lifecycle")
	defer span.End()

	// 启动任务
	logger.Info(ctx, "启动后台任务")

	// 模拟后台任务
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				taskCtx, taskSpan := tracer.Start(context.Background(), "background-task")
				logger.Debug(taskCtx, "执行后台任务")

				// 模拟任务处理
				processBackgroundTask(taskCtx, tracer)

				taskSpan.End()
			case <-ctx.Done():
				logger.Info(context.Background(), "后台任务停止")
				return
			}
		}
	}()

	// 模拟 HTTP 请求处理
	for i := 0; i < 10; i++ {
		handleHTTPRequest(ctx, tracer, i)
		time.Sleep(2 * time.Second)
	}

	// 等待优雅关闭信号
	return waitForShutdown(ctx)
}

// handleHTTPRequest 模拟 HTTP 请求处理
func handleHTTPRequest(ctx context.Context, tracer trace.Tracer, requestID int) {
	ctx, span := tracer.Start(ctx, "http-request")
	defer span.End()

	logger.Info(ctx, "收到 HTTP 请求",
		"request_id", requestID,
		"method", "GET",
		"path", "/api/users",
	)

	// 模拟数据库查询
	{
		ctx, dbSpan := tracer.Start(ctx, "database-query")
		logger.Debug(ctx, "执行数据库查询", "query", "SELECT * FROM users")
		time.Sleep(50 * time.Millisecond)
		dbSpan.End()
	}

	// 模拟缓存查询
	{
		ctx, cacheSpan := tracer.Start(ctx, "cache-lookup")
		logger.Debug(ctx, "查询缓存", "key", "user:12345")
		time.Sleep(10 * time.Millisecond)
		cacheSpan.End()
	}

	// 模拟错误（10% 概率）
	if requestID%10 == 0 {
		logger.Error(ctx, "请求处理失败",
			"request_id", requestID,
			"error", "database connection timeout",
		)
	} else {
		logger.Info(ctx, "请求处理完成",
			"request_id", requestID,
			"status", 200,
			"duration_ms", 60,
		)
	}
}

// processBackgroundTask 处理后台任务
func processBackgroundTask(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "process-batch")
	defer span.End()

	logger.Debug(ctx, "开始处理批量任务")

	// 模拟批量处理
	for i := 0; i < 10; i++ {
		logger.Debug(ctx, "处理项目", "index", i)
		time.Sleep(10 * time.Millisecond)
	}

	logger.Debug(ctx, "批量任务完成", "count", 10)
}

// waitForShutdown 等待优雅关闭信号
func waitForShutdown(ctx context.Context) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	logger.Info(ctx, "应用运行中，等待关闭信号...")

	sig := <-sigCh
	logger.Info(ctx, "收到关闭信号", "signal", sig.String())

	// 优雅关闭
	logger.Info(ctx, "开始优雅关闭...")
	time.Sleep(2 * time.Second) // 模拟清理工作
	logger.Info(ctx, "优雅关闭完成")

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// getEnvBool 获取布尔环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

