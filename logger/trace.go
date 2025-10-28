package logger

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// extractTraceFields 从 context 提取 trace 信息
func extractTraceFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return nil
	}

	spanCtx := span.SpanContext()
	if !spanCtx.IsValid() {
		return nil
	}

	return []zap.Field{
		zap.String("trace_id", spanCtx.TraceID().String()),
		zap.String("span_id", spanCtx.SpanID().String()),
	}
}

// markSpanError 标记 span 为 error 状态
func markSpanError(ctx context.Context, msg string) {
	if ctx == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return
	}

	span.SetStatus(codes.Error, msg)
	span.RecordError(errors.New(msg))
}

// addSpanEvent 添加日志到 span event
func addSpanEvent(ctx context.Context, level Level, msg string) {
	if ctx == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return
	}

	// 根据日志级别设置事件属性
	eventName := level.String() + ": " + msg
	span.AddEvent(eventName)
}

