package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Si40Code/kit/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	// 修改为你的 SigNoz 地址
	signozHost     = "47.83.197.11" // SigNoz 服务器地址（不需要 http:// 前缀）
	signozPort     = "4317"         // gRPC 端口（日志）
	signozHTTPPort = "4318"         // HTTP 端口（trace）
	serviceName    = "logger-signoz-example"
)

func main() {
	// 初始化 OpenTelemetry tracer（发送到 SigNoz）
	cleanup := initTracer()
	defer cleanup()

	// 初始化 logger（发送到 SigNoz + stdout）
	err := logger.Init(
		logger.WithLevel(logger.DebugLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithStdout(), // 同时输出到控制台
		logger.WithOTLP(signozHost+":"+signozPort,
			logger.WithOTLPInsecure(),
		),
		logger.WithTrace(serviceName), // 启用 trace 集成
	)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	fmt.Printf("🚀 应用启动，日志和 trace 将发送到 SigNoz\n")
	fmt.Printf("📊 SigNoz 地址: http://%s:3301\n\n", signozHost)

	// 运行示例
	runExamples()

	fmt.Println("\n✅ 所有示例完成！")
	fmt.Printf("📊 请访问 SigNoz UI 查看日志和 trace: http://%s:3301\n", signozHost)
	fmt.Println("   - Traces: 查看完整的调用链")
	fmt.Println("   - Logs: 查看所有日志记录")
	fmt.Println("   - 点击 trace 可以看到关联的日志")
}

func runExamples() {
	ctx := context.Background()
	tracer := otel.Tracer(serviceName)

	// 示例 1: 简单的 HTTP 请求处理
	{
		ctx, span := tracer.Start(ctx, "http-request-handler")
		defer span.End()

		logger.Info(ctx, "收到 HTTP 请求",
			"method", "GET",
			"path", "/api/users",
			"client_ip", "192.168.1.100",
		)

		// 模拟处理
		time.Sleep(50 * time.Millisecond)

		logger.Info(ctx, "请求处理完成",
			"status", 200,
			"duration_ms", 50,
		)
	}

	// 等待一下，让日志有时间发送
	time.Sleep(100 * time.Millisecond)

	// 示例 2: 数据库操作
	{
		ctx, span := tracer.Start(ctx, "database-operation")
		defer span.End()

		logger.Info(ctx, "开始数据库查询",
			"query", "SELECT * FROM users WHERE id = ?",
			"params", 12345,
		)

		// 模拟数据库查询
		time.Sleep(100 * time.Millisecond)

		logger.Info(ctx, "数据库查询完成",
			"rows", 1,
			"duration_ms", 100,
		)
	}

	time.Sleep(100 * time.Millisecond)

	// 示例 3: 业务逻辑错误
	{
		ctx, span := tracer.Start(ctx, "process-payment")
		defer span.End()

		logger.Info(ctx, "开始处理支付",
			"order_id", "ORD-2024-001",
			"amount", 99.99,
			"currency", "USD",
		)

		// 模拟支付处理
		time.Sleep(150 * time.Millisecond)

		// 模拟支付失败
		logger.Error(ctx, "支付处理失败",
			"order_id", "ORD-2024-001",
			"error", "insufficient funds",
			"balance", 50.00,
			"required", 99.99,
		)
	}

	time.Sleep(100 * time.Millisecond)

	// 示例 4: 复杂的业务流程
	{
		ctx, span := tracer.Start(ctx, "order-fulfillment")
		defer span.End()

		orderID := "ORD-2024-002"

		logger.Info(ctx, "开始订单履行流程", "order_id", orderID)

		// 步骤 1: 检查库存
		{
			ctx, checkSpan := tracer.Start(ctx, "check-inventory")
			logger.Info(ctx, "检查库存", "order_id", orderID)
			time.Sleep(50 * time.Millisecond)
			logger.Info(ctx, "库存充足", "available", 100)
			checkSpan.End()
		}

		// 步骤 2: 扣减库存
		{
			ctx, deductSpan := tracer.Start(ctx, "deduct-inventory")
			logger.Info(ctx, "扣减库存", "order_id", orderID, "quantity", 1)
			time.Sleep(50 * time.Millisecond)
			deductSpan.End()
		}

		// 步骤 3: 创建发货单
		{
			ctx, shippingSpan := tracer.Start(ctx, "create-shipping")
			logger.Info(ctx, "创建发货单", "order_id", orderID)
			time.Sleep(50 * time.Millisecond)
			logger.Info(ctx, "发货单创建成功", "shipping_id", "SHIP-2024-001")
			shippingSpan.End()
		}

		// 步骤 4: 通知用户（模拟失败）
		{
			ctx, notifySpan := tracer.Start(ctx, "notify-customer")
			logger.Info(ctx, "发送通知", "order_id", orderID, "channel", "email")
			time.Sleep(50 * time.Millisecond)

			// 模拟通知失败
			logger.Error(ctx, "通知发送失败",
				"order_id", orderID,
				"error", "email service unavailable",
			)
			notifySpan.End()
		}

		logger.Warn(ctx, "订单履行完成但通知失败",
			"order_id", orderID,
			"shipping_id", "SHIP-2024-001",
		)
	}

	time.Sleep(100 * time.Millisecond)

	// 示例 5: 高频日志（性能测试）
	{
		ctx, span := tracer.Start(ctx, "high-frequency-logging")
		defer span.End()

		logger.Info(ctx, "开始高频日志测试")

		start := time.Now()
		for i := 0; i < 100; i++ {
			logger.Debug(ctx, "处理项目",
				"index", i,
				"timestamp", time.Now().Unix(),
			)
		}
		elapsed := time.Since(start)

		logger.Info(ctx, "高频日志测试完成",
			"count", 100,
			"elapsed_ms", elapsed.Milliseconds(),
			"avg_per_log_us", elapsed.Microseconds()/100,
		)
	}
}

// initTracer 初始化 OpenTelemetry tracer（发送到 SigNoz）
func initTracer() func() {
	ctx := context.Background()

	// 创建 OTLP HTTP exporter
	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(signozHost+":"+signozHTTPPort),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(getEnv("ENV", "development")),
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
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
