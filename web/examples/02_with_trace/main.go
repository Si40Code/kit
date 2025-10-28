package main

import (
	"context"
	"log"

	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() func() {
	// 创建 stdout exporter（输出到终端，生产环境建议使用 Jaeger/Otlp）
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(), // 格式化输出
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("trace-example"),
		)),
	)
	
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// 初始化 Tracer
	cleanup := initTracer()
	defer cleanup()

	// 创建服务器（启用 Trace）
	server := web.New(
		web.WithMode(web.ReleaseMode),
		web.WithServiceName("trace-example"),
		web.WithTrace(),
	)

	engine := server.Engine()
	
	engine.GET("/user/:id", func(c *gin.Context) {
		// 从 context 中获取 tracer 并创建子 span
		ctx := c.Request.Context()
		tracer := otel.Tracer("user-service")
		_, span := tracer.Start(ctx, "get-user-from-db")
		defer span.End()

		// 模拟数据库查询
		userID := c.Param("id")
		web.Success(c, gin.H{
			"id":   userID,
			"name": "John Doe",
		})
	})

	server.RunWithGracefulShutdown(":8080")
}
