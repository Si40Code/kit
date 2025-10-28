package main

import (
	"context"
	"fmt"
	"log"
	"sync"
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
	// 1. 初始化 Tracer
	cleanup := initTracer("production-service")
	defer cleanup()

	// 2. 初始化 Logger
	err := logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.JSONFormat),
		logger.WithStdout(),
		logger.WithTrace("production-service"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// 3. 初始化 Metric Recorder
	recorder := NewProductionMetricRecorder()

	// 4. 创建生产级 HTTP 客户端
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithTrace("production-service"),
		httpclient.WithMetric(recorder),
		httpclient.WithTimeout(10*time.Second),
		httpclient.WithRetry(3, 100*time.Millisecond, 2*time.Second),
		httpclient.WithMaxBodyLogSize(5*1024), // 5KB
		httpclient.WithMaxIdleConns(100),
		httpclient.WithIdleConnTimeout(90*time.Second),
		httpclient.WithKeepAlive(30*time.Second),
	)

	fmt.Println("=== 生产环境配置示例 ===")

	// 示例 1: 单个请求
	example1SingleRequest(client)

	// 示例 2: 并发请求
	example2ConcurrentRequests(client)

	// 示例 3: 重试机制
	example3RetryMechanism(client)

	// 示例 4: 查看性能统计
	example4ViewStats(recorder)
}

// ProductionMetricRecorder 生产级 metric 记录器
type ProductionMetricRecorder struct {
	mu      sync.RWMutex
	metrics []httpclient.MetricData
}

func NewProductionMetricRecorder() *ProductionMetricRecorder {
	return &ProductionMetricRecorder{
		metrics: make([]httpclient.MetricData, 0),
	}
}

func (r *ProductionMetricRecorder) RecordRequest(data httpclient.MetricData) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = append(r.metrics, data)
}

func (r *ProductionMetricRecorder) GetStats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.metrics) == 0 {
		return map[string]interface{}{
			"total_requests": 0,
		}
	}

	var totalTime, totalDNS, totalTCP, totalTLS int64
	var success, failed int
	statusCodes := make(map[int]int)

	for _, m := range r.metrics {
		totalTime += m.TotalTime.Milliseconds()
		totalDNS += m.DNSLookup.Milliseconds()
		totalTCP += m.TCPConn.Milliseconds()
		totalTLS += m.TLSHandshake.Milliseconds()

		statusCodes[m.StatusCode]++

		if m.StatusCode >= 200 && m.StatusCode < 400 {
			success++
		} else {
			failed++
		}
	}

	count := int64(len(r.metrics))

	return map[string]interface{}{
		"total_requests":       count,
		"success_requests":     success,
		"failed_requests":      failed,
		"success_rate":         fmt.Sprintf("%.2f%%", float64(success)/float64(count)*100),
		"avg_total_time_ms":    totalTime / count,
		"avg_dns_lookup_ms":    totalDNS / count,
		"avg_tcp_conn_ms":      totalTCP / count,
		"avg_tls_handshake_ms": totalTLS / count,
		"status_codes":         statusCodes,
	}
}

// 示例 1: 单个请求
func example1SingleRequest(client *httpclient.Client) {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	ctx, span := tracer.Start(ctx, "single-request-demo")
	defer span.End()

	logger.Info(ctx, "开始单个请求示例", nil)

	type APIRequest struct {
		UserID string `json:"user_id"`
		Action string `json:"action"`
	}

	resp, err := client.R(ctx).
		SetHeader("X-Request-ID", "req-123").
		SetBody(APIRequest{
			UserID: "user-456",
			Action: "query",
		}).
		Post("https://httpbin.org/post")
	if err != nil {
		logger.Error(ctx, "请求失败", map[string]interface{}{"error": err})
		return
	}

	logger.Info(ctx, "请求成功",
		map[string]interface{}{
			"status_code": resp.StatusCode(),
		},
	)

	fmt.Println("✓ 单个请求完成")
}

// 示例 2: 并发请求
func example2ConcurrentRequests(client *httpclient.Client) {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	ctx, span := tracer.Start(ctx, "concurrent-requests-demo")
	defer span.End()

	logger.Info(ctx, "开始并发请求示例", nil)

	var wg sync.WaitGroup
	concurrency := 5

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// 为每个 goroutine 创建子 span
			reqCtx, reqSpan := tracer.Start(ctx, fmt.Sprintf("concurrent-request-%d", index))
			defer reqSpan.End()

			resp, err := client.R(reqCtx).
				SetQueryParam("id", fmt.Sprintf("%d", index)).
				Get("https://httpbin.org/get")
			if err != nil {
				logger.Error(reqCtx, "并发请求失败",
					map[string]interface{}{
						"index": index,
						"error": err,
					},
				)
				return
			}

			logger.Info(reqCtx, "并发请求成功",
				map[string]interface{}{
					"index":       index,
					"status_code": resp.StatusCode(),
				},
			)
		}(i)
	}

	wg.Wait()
	logger.Info(ctx, "并发请求完成", nil)
	fmt.Println("✓ 并发请求完成")
}

// 示例 3: 重试机制
func example3RetryMechanism(client *httpclient.Client) {
	ctx := context.Background()
	tracer := otel.Tracer("example")

	ctx, span := tracer.Start(ctx, "retry-mechanism-demo")
	defer span.End()

	logger.Info(ctx, "开始重试机制示例", nil)

	// 请求一个会返回 500 的端点（模拟错误）
	resp, err := client.R(ctx).
		Get("https://httpbin.org/status/500")

	if err != nil {
		logger.Error(ctx, "请求失败（重试后）", map[string]interface{}{"error": err})
	} else {
		logger.Warn(ctx, "请求返回错误状态码",
			map[string]interface{}{
				"status_code": resp.StatusCode(),
			},
		)
	}

	fmt.Println("✓ 重试机制示例完成")
}

// 示例 4: 查看性能统计
func example4ViewStats(recorder *ProductionMetricRecorder) {
	fmt.Println("\n=== 性能统计 ===")

	stats := recorder.GetStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	fmt.Println("\n✓ 统计信息已生成")
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
