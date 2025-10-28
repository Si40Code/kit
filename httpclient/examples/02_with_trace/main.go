package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Si40Code/kit/httpclient"
	"github.com/Si40Code/kit/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	// 初始化 OpenTelemetry tracer
	cleanup := initTracer("httpclient-trace-example")
	defer cleanup()

	// 初始化 logger（启用 trace 集成）
	err := logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
		logger.WithTrace("httpclient-trace-example"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	fmt.Println("=== 示例 1: 基本 Trace 集成 ===")
	example1BasicTrace()

	fmt.Println("\n=== 示例 2: 嵌套 Span ===")
	example2NestedSpans()

	fmt.Println("\n=== 示例 3: 业务场景 - 用户注册流程 ===")
	example3BusinessScenario()
}

// 示例 1: 基本 Trace 集成
func example1BasicTrace() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	// 创建父 span
	ctx, span := tracer.Start(ctx, "example1-operation")
	defer span.End()

	// 创建 HTTP 客户端（启用 trace）
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithTrace("httpclient-example"),
	)

	// 发起请求 - 会自动创建 child span
	resp, err := client.R(ctx).
		Get("https://httpbin.org/get")
	if err != nil {
		logger.Error(ctx, "请求失败", map[string]interface{}{"error": err})
		return
	}

	logger.Info(ctx, "请求成功",
		map[string]interface{}{
			"status_code": resp.StatusCode(),
		},
	)

	fmt.Println("✓ HTTP 请求的 span 已自动创建并关联到父 span")
}

// 示例 2: 嵌套 Span
func example2NestedSpans() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	// 父 span: 用户查询流程
	ctx, parentSpan := tracer.Start(ctx, "user-query-flow")
	defer parentSpan.End()

	logger.Info(ctx, "开始用户查询流程")

	// 创建客户端
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithTrace("httpclient-example"),
	)

	// 子操作 1: 获取用户信息
	{
		ctx, childSpan1 := tracer.Start(ctx, "get-user-info")
		resp, err := client.R(ctx).
			Get("https://httpbin.org/get?userId=123")
		if err != nil {
			logger.Error(ctx, "获取用户信息失败", map[string]interface{}{"error": err})
		} else {
			logger.Info(ctx, "获取用户信息成功", map[string]interface{}{"status": resp.StatusCode()})
		}
		childSpan1.End()
	}

	// 子操作 2: 获取用户订单
	{
		ctx, childSpan2 := tracer.Start(ctx, "get-user-orders")
		time.Sleep(100 * time.Millisecond)
		resp, err := client.R(ctx).
			Get("https://httpbin.org/get?userId=123&type=orders")
		if err != nil {
			logger.Error(ctx, "获取用户订单失败", map[string]interface{}{"error": err})
		} else {
			logger.Info(ctx, "获取用户订单成功", map[string]interface{}{"status": resp.StatusCode()})
		}
		childSpan2.End()
	}

	logger.Info(ctx, "用户查询流程完成")
	fmt.Println("✓ 嵌套 span 已正确关联")
}

// 示例 3: 业务场景 - 用户注册流程
func example3BusinessScenario() {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	// 模拟用户注册请求
	ctx, span := tracer.Start(ctx, "user-registration")
	defer span.End()

	userID := "user-12345"
	email := "bob@example.com"

	logger.Info(ctx, "收到注册请求",
		map[string]interface{}{
			"user_id": userID,
			"email":   email,
		},
	)

	// 创建客户端
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithTrace("httpclient-example"),
	)

	// 步骤 1: 验证邮箱
	{
		ctx, validateSpan := tracer.Start(ctx, "validate-email")
		logger.Info(ctx, "验证邮箱", map[string]interface{}{"email": email})

		resp, err := client.R(ctx).
			SetQueryParams(map[string]string{
				"email": email,
			}).
			Get("https://httpbin.org/get")

		if err != nil {
			logger.Error(ctx, "邮箱验证失败", map[string]interface{}{"error": err})
		} else {
			logger.Info(ctx, "邮箱验证成功", map[string]interface{}{"status": resp.StatusCode()})
		}
		validateSpan.End()
	}

	// 步骤 2: 创建用户
	{
		ctx, createSpan := tracer.Start(ctx, "create-user")
		logger.Info(ctx, "创建用户记录", map[string]interface{}{"user_id": userID})

		type User struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
		}

		resp, err := client.R(ctx).
			SetBody(User{UserID: userID, Email: email}).
			Post("https://httpbin.org/post")

		if err != nil {
			logger.Error(ctx, "创建用户失败", map[string]interface{}{"error": err})
		} else {
			logger.Info(ctx, "创建用户成功", map[string]interface{}{"status": resp.StatusCode()})
		}
		createSpan.End()
	}

	// 步骤 3: 发送欢迎邮件
	{
		ctx, emailSpan := tracer.Start(ctx, "send-welcome-email")
		logger.Info(ctx, "发送欢迎邮件", map[string]interface{}{"email": email})

		resp, err := client.R(ctx).
			SetBody(map[string]string{
				"to":      email,
				"subject": "欢迎注册",
			}).
			Post("https://httpbin.org/post")

		if err != nil {
			logger.Error(ctx, "邮件发送失败", map[string]interface{}{"error": err})
		} else {
			logger.Info(ctx, "邮件发送成功", map[string]interface{}{"status": resp.StatusCode()})
		}
		emailSpan.End()
	}

	logger.Info(ctx, "注册流程完成",
		map[string]interface{}{
			"user_id": userID,
		},
	)

	fmt.Println("✓ 完整业务流程的 trace 已正确关联")
}

// initTracer 初始化 OpenTelemetry tracer
func initTracer(serviceName string) func() {
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

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
