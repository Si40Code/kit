package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Si40Code/kit/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	// 初始化 OpenTelemetry tracer
	cleanup := initTracer("trace-example")
	defer cleanup()

	// 初始化 logger（启用 trace 集成）
	err := logger.Init(
		logger.WithLevel(logger.DebugLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
		logger.WithTrace("trace-example"), // 启用 trace 集成
	)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	fmt.Println("=== 示例 1: 基本 Trace 集成 ===")
	example1BasicTrace()

	fmt.Println("\n=== 示例 2: Error 日志标记 Span ===")
	example2ErrorSpan()

	fmt.Println("\n=== 示例 3: 嵌套 Span ===")
	example3NestedSpans()

	fmt.Println("\n=== 示例 4: 业务场景示例 ===")
	example4BusinessScenario()
}

// 示例 1: 基本 Trace 集成
func example1BasicTrace() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	// 创建 span
	ctx, span := tracer.Start(ctx, "example1-operation")
	defer span.End()

	// 记录日志会自动包含 trace_id 和 span_id
	logger.Info(ctx, "操作开始")
	logger.Info(ctx, "处理数据", "count", 100)
	logger.Info(ctx, "操作完成")

	fmt.Println("✓ 日志已包含 trace_id 和 span_id")
}

// 示例 2: Error 日志自动标记 Span
func example2ErrorSpan() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	ctx, span := tracer.Start(ctx, "example2-error-operation")
	defer span.End()

	logger.Info(ctx, "开始执行操作")

	// 模拟错误
	time.Sleep(100 * time.Millisecond)

	// Error 日志会自动将 span 状态设置为 error
	logger.Error(ctx, "操作失败", "error", "database connection timeout")

	fmt.Println("✓ Span 已被标记为 error 状态")
}

// 示例 3: 嵌套 Span
func example3NestedSpans() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	// 父 span
	ctx, parentSpan := tracer.Start(ctx, "parent-operation")
	defer parentSpan.End()

	logger.Info(ctx, "父操作开始")

	// 子 span 1
	{
		ctx, childSpan1 := tracer.Start(ctx, "child-operation-1")
		logger.Info(ctx, "子操作 1 执行")
		time.Sleep(50 * time.Millisecond)
		childSpan1.End()
	}

	// 子 span 2
	{
		ctx, childSpan2 := tracer.Start(ctx, "child-operation-2")
		logger.Info(ctx, "子操作 2 执行")
		time.Sleep(50 * time.Millisecond)
		childSpan2.End()
	}

	logger.Info(ctx, "父操作完成")

	fmt.Println("✓ 嵌套 span 的日志已正确关联")
}

// 示例 4: 业务场景示例 - 用户注册流程
func example4BusinessScenario() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	// 模拟用户注册请求
	ctx, span := tracer.Start(ctx, "user-registration")
	defer span.End()

	userID := "user-12345"
	email := "alice@example.com"

	logger.Info(ctx, "收到注册请求",
		"user_id", userID,
		"email", email,
	)

	// 步骤 1: 验证邮箱
	{
		ctx, validateSpan := tracer.Start(ctx, "validate-email")
		logger.Info(ctx, "验证邮箱", "email", email)
		time.Sleep(100 * time.Millisecond)
		validateSpan.End()
	}

	// 步骤 2: 创建用户
	{
		ctx, createSpan := tracer.Start(ctx, "create-user")
		logger.Info(ctx, "创建用户记录", "user_id", userID)
		time.Sleep(150 * time.Millisecond)
		createSpan.End()
	}

	// 步骤 3: 发送欢迎邮件（模拟失败）
	{
		ctx, emailSpan := tracer.Start(ctx, "send-welcome-email")
		logger.Info(ctx, "发送欢迎邮件", "email", email)
		time.Sleep(50 * time.Millisecond)

		// 模拟邮件发送失败
		logger.Error(ctx, "邮件发送失败",
			"error", "SMTP connection timeout",
			"email", email,
		)
		emailSpan.End()
	}

	// 步骤 4: 记录注册事件
	{
		ctx, eventSpan := tracer.Start(ctx, "log-registration-event")
		logger.Info(ctx, "记录注册事件",
			"user_id", userID,
			"timestamp", time.Now().Unix(),
		)
		time.Sleep(30 * time.Millisecond)
		eventSpan.End()
	}

	logger.Warn(ctx, "注册流程完成但邮件发送失败",
		"user_id", userID,
		"email", email,
	)

	fmt.Println("✓ 完整业务流程的 trace 和日志已关联")
}

// initTracer 初始化 OpenTelemetry tracer（使用 stdout exporter 用于演示）
func initTracer(serviceName string) func() {
	// 创建 stdout exporter（用于演示）
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
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

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}
}

