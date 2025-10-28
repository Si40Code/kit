package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Si40Code/kit/web"
	zapotlp "github.com/SigNoz/zap_otlp"
	zapotlpencoder "github.com/SigNoz/zap_otlp/zap_otlp_encoder"
	zapotlpsync "github.com/SigNoz/zap_otlp/zap_otlp_sync"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// SigNozLogger 使用 Zap + OTLP 的日志记录器
type SigNozLogger struct {
	zapLogger *zap.Logger
}

func NewSigNozLogger(endpoint, serviceName string) (*SigNozLogger, error) {
	// 创建两个 encoder：一个用于控制台，一个用于 OTLP
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	otlpEncoder := zapotlpencoder.NewOTLPEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// 创建 gRPC 连接到 SigNoz
	ctx := context.Background()

	// 检查是否使用安全连接
	grpcInsecure := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")
	var secureOption grpc.DialOption
	if strings.ToLower(grpcInsecure) == "false" || grpcInsecure == "0" {
		secureOption = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.DialContext(
		ctx,
		endpoint,
		grpc.WithBlock(),
		secureOption,
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OTLP endpoint: %w", err)
	}

	// 创建 OTLP Syncer
	otlpSync := zapotlpsync.NewOtlpSyncer(conn, zapotlpsync.Options{
		BatchSize:      100,
		BatchInterval:  5 * time.Second,
		ResourceSchema: semconv.SchemaURL,
		Resource: resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	})

	// 创建 Zap Core，同时输出到控制台和 OTLP
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, os.Stdout, zapcore.InfoLevel),
		zapcore.NewCore(otlpEncoder, otlpSync, zapcore.InfoLevel),
	)

	// 创建 Logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &SigNozLogger{zapLogger: zapLogger}, nil
}

func (l *SigNozLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	// 添加 trace context
	zapFields = append(zapFields, zapotlp.SpanCtx(ctx))
	l.zapLogger.Info(msg, zapFields...)
}

func (l *SigNozLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	zapFields = append(zapFields, zapotlp.SpanCtx(ctx))
	l.zapLogger.Warn(msg, zapFields...)
}

func (l *SigNozLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	zapFields = append(zapFields, zapotlp.SpanCtx(ctx))
	l.zapLogger.Error(msg, zapFields...)
}

func (l *SigNozLogger) convertFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (l *SigNozLogger) Sync() error {
	return l.zapLogger.Sync()
}

// 使用 otelhttp 的标准 HTTP 指标，不再自定义业务指标

func initSigNozTracer(endpoint string, serviceName string) (func(), error) {
	// 创建 OTLP HTTP exporter
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}, nil
}

// initSigNozMetrics initializes OTLP HTTP metrics exporter and MeterProvider
func initSigNozMetrics(endpoint string, serviceName string) (func(), error) {
	exporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP metrics exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second))

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetMeterProvider(mp)

	return func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}, nil
}

func main() {
	signozEndpoint := "47.83.197.11"
	serviceName := "signoz-example"

	// 设置环境变量
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")

	// 初始化 SigNoz Tracer (HTTP endpoint for traces)
	traceCleanup, err := initSigNozTracer(signozEndpoint+":4318", serviceName)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer traceCleanup()

	// 初始化 SigNoz Logger (gRPC endpoint for logs)
	logger, err := NewSigNozLogger(signozEndpoint+":4317", serviceName)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 初始化 SigNoz Metrics (HTTP endpoint for metrics)
	metricCleanup, err := initSigNozMetrics(signozEndpoint+":4318", serviceName)
	if err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}
	defer metricCleanup()

	// 创建服务器（不启用 WithTrace，避免与 otelhttp 重复 Trace）
	server := web.New(
		web.WithMode(web.ReleaseMode),
		web.WithServiceName(serviceName),
		web.WithLogger(logger),
        // 依赖 otelhttp 的标准 HTTP 指标，不再注入业务指标
        web.WithSkipPaths("/health", "/metrics"),
	)

	engine := server.Engine()

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		web.Success(c, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 简单的用户列表
	engine.GET("/api/users", func(c *gin.Context) {
		web.Success(c, gin.H{
			"users": []gin.H{
				{"id": 1, "name": "Alice", "email": "alice@example.com"},
				{"id": 2, "name": "Bob", "email": "bob@example.com"},
				{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
			},
		})
	})

	// 获取单个用户
	engine.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		web.Success(c, gin.H{
			"id":    userID,
			"name":  "User " + userID,
			"email": fmt.Sprintf("user%s@example.com", userID),
		})
	})

	// 创建用户
	engine.POST("/api/users", func(c *gin.Context) {
		var req struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		web.Success(c, gin.H{
			"id":         4,
			"name":       req.Name,
			"email":      req.Email,
			"created_at": time.Now().Format(time.RFC3339),
		})
	})

	log.Printf("🚀 Starting server with SigNoz integration on :8080 (otelhttp)")
	log.Printf("📊 SigNoz UI: http://%s", signozEndpoint)
	log.Printf("✅ 日志通过 OTLP 直接发送到 SigNoz（无需额外脚本）")

	// 启动自动测试请求（5秒后开始）
	go func() {
		time.Sleep(5 * time.Second)
		log.Println("📡 Starting automatic test requests...")
		autoSendTestRequests()
	}()

	// 使用 otelhttp 包装 Gin，引入自动 Trace + http.server.* 指标
	srv := &http.Server{
		Addr:    ":8080",
		Handler: otelhttp.NewHandler(engine, serviceName),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}

// autoSendTestRequests 自动发送测试请求
func autoSendTestRequests() {
	client := &http.Client{Timeout: 5 * time.Second}
	baseURL := "http://localhost:8080"

	// 每10秒发送一轮测试请求
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🔄 Sending test requests...")

		// 1. 健康检查
		makeRequest(client, "GET", baseURL+"/health", nil)
		time.Sleep(100 * time.Millisecond)

		// 2. 获取用户列表
		makeRequest(client, "GET", baseURL+"/api/users", nil)
		time.Sleep(100 * time.Millisecond)

		// 3. 获取单个用户（循环几个 ID）
		for i := 1; i <= 3; i++ {
			makeRequest(client, "GET", fmt.Sprintf("%s/api/users/%d", baseURL, i), nil)
			time.Sleep(100 * time.Millisecond)
		}

		// 4. 创建用户
		userData := []byte(`{"name":"TestUser","email":"test@example.com"}`)
		makeRequest(client, "POST", baseURL+"/api/users", bytes.NewBuffer(userData))
		time.Sleep(100 * time.Millisecond)

		log.Println("✅ Test requests batch completed")
	}
}

// makeRequest 发送 HTTP 请求
func makeRequest(client *http.Client, method, url string, body io.Reader) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		return
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Request failed [%s %s]: %v", method, url, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("✓ [%s %s] → %d", method, url, resp.StatusCode)
}
