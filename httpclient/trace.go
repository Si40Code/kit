package httpclient

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// createSpan 为 HTTP 请求创建一个新的 span
func createSpan(ctx context.Context, serviceName, method, url string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)
	spanName := fmt.Sprintf("HTTP %s", method)

	ctx, span := tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("http.method", method),
			attribute.String("http.url", url),
		),
	)

	return ctx, span
}

// setSpanAttributes 设置 span 属性
func setSpanAttributes(span trace.Span, attrs map[string]interface{}) {
	if !span.IsRecording() {
		return
	}

	for key, value := range attrs {
		switch v := value.(type) {
		case string:
			span.SetAttributes(attribute.String(key, v))
		case int:
			span.SetAttributes(attribute.Int(key, v))
		case int64:
			span.SetAttributes(attribute.Int64(key, v))
		case float64:
			span.SetAttributes(attribute.Float64(key, v))
		case bool:
			span.SetAttributes(attribute.Bool(key, v))
		default:
			span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
}

// markSpanError 标记 span 为错误状态
func markSpanError(span trace.Span, err error, msg string) {
	if !span.IsRecording() {
		return
	}

	span.SetStatus(codes.Error, msg)
	if err != nil {
		span.RecordError(err)
	}
}

// markSpanSuccess 标记 span 为成功状态
func markSpanSuccess(span trace.Span) {
	if !span.IsRecording() {
		return
	}

	span.SetStatus(codes.Ok, "")
}

// getSpanFromContext 从 context 中获取 span
func getSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
