package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Si40Code/kit/httpclient"
	"github.com/Si40Code/kit/logger"
)

func main() {
	// 初始化 logger
	err := logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	fmt.Println("=== 示例 1: 基本 Metric 记录 ===")
	example1BasicMetric()

	fmt.Println("\n=== 示例 2: 批量请求 Metric ===")
	example2BatchMetric()

	fmt.Println("\n=== 示例 3: 查看聚合统计 ===")
	example3ViewStats()
}

// SimpleMetricRecorder 简单的 metric 记录器实现
type SimpleMetricRecorder struct {
	mu      sync.Mutex
	metrics []httpclient.MetricData
}

func NewSimpleMetricRecorder() *SimpleMetricRecorder {
	return &SimpleMetricRecorder{
		metrics: make([]httpclient.MetricData, 0),
	}
}

func (r *SimpleMetricRecorder) RecordRequest(data httpclient.MetricData) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = append(r.metrics, data)

	// 打印详细的 metric 信息
	fmt.Printf("\n[Metric] HTTP Request:\n")
	fmt.Printf("  Method:        %s\n", data.Method)
	fmt.Printf("  Path:          %s\n", data.Path)
	fmt.Printf("  Host:          %s\n", data.Host)
	fmt.Printf("  Status Code:   %d\n", data.StatusCode)
	fmt.Printf("  Total Time:    %v\n", data.TotalTime)
	fmt.Printf("  DNS Lookup:    %v\n", data.DNSLookup)
	fmt.Printf("  TCP Conn:      %v\n", data.TCPConn)
	fmt.Printf("  TLS Handshake: %v\n", data.TLSHandshake)
	fmt.Printf("  Server Time:   %v\n", data.ServerTime)
	fmt.Printf("  Response Time: %v\n", data.ResponseTime)
	fmt.Printf("  Conn Reused:   %v\n", data.IsConnReused)
	fmt.Printf("  Conn Was Idle: %v\n", data.IsConnWasIdle)
	fmt.Printf("  Remote Addr:   %s\n", data.RemoteAddr)
}

func (r *SimpleMetricRecorder) GetMetrics() []httpclient.MetricData {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]httpclient.MetricData{}, r.metrics...)
}

func (r *SimpleMetricRecorder) GetStats() map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.metrics) == 0 {
		return map[string]interface{}{
			"total_requests": 0,
		}
	}

	var totalTime, totalDNS, totalTCP, totalTLS, totalServer int64
	var connReused, connIdle int
	statusCodes := make(map[int]int)

	for _, m := range r.metrics {
		totalTime += m.TotalTime.Milliseconds()
		totalDNS += m.DNSLookup.Milliseconds()
		totalTCP += m.TCPConn.Milliseconds()
		totalTLS += m.TLSHandshake.Milliseconds()
		totalServer += m.ServerTime.Milliseconds()

		if m.IsConnReused {
			connReused++
		}
		if m.IsConnWasIdle {
			connIdle++
		}

		statusCodes[m.StatusCode]++
	}

	count := int64(len(r.metrics))

	return map[string]interface{}{
		"total_requests":       count,
		"avg_total_time_ms":    totalTime / count,
		"avg_dns_lookup_ms":    totalDNS / count,
		"avg_tcp_conn_ms":      totalTCP / count,
		"avg_tls_handshake_ms": totalTLS / count,
		"avg_server_time_ms":   totalServer / count,
		"conn_reuse_rate":      fmt.Sprintf("%.2f%%", float64(connReused)/float64(count)*100),
		"conn_idle_rate":       fmt.Sprintf("%.2f%%", float64(connIdle)/float64(count)*100),
		"status_codes":         statusCodes,
	}
}

// 示例 1: 基本 Metric 记录
func example1BasicMetric() {
	// 创建 metric 记录器
	recorder := NewSimpleMetricRecorder()

	// 创建客户端
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithMetric(recorder),
	)

	// 发起请求
	resp, err := client.R(context.Background()).
		Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("\n响应状态码: %d\n", resp.StatusCode())
	fmt.Println("✓ Metric 已记录")
}

// 示例 2: 批量请求 Metric
func example2BatchMetric() {
	recorder := NewSimpleMetricRecorder()

	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithMetric(recorder),
		httpclient.WithDisableLog(), // 禁用日志以便查看 metric
	)

	endpoints := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/status/200",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/json",
	}

	fmt.Println("\n开始批量请求...")

	for i, endpoint := range endpoints {
		resp, err := client.R(context.Background()).
			Get(endpoint)
		if err != nil {
			log.Printf("请求 %d 失败: %v", i+1, err)
			continue
		}

		fmt.Printf("请求 %d 完成，状态码: %d\n", i+1, resp.StatusCode())
	}

	fmt.Println("\n✓ 所有请求的 metric 已记录")
}

// 示例 3: 查看聚合统计
func example3ViewStats() {
	recorder := NewSimpleMetricRecorder()

	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithMetric(recorder),
		httpclient.WithDisableLog(),
	)

	// 发起多个请求
	fmt.Println("发起 10 个请求...")
	for i := 0; i < 10; i++ {
		_, _ = client.R(context.Background()).
			Get("https://httpbin.org/get")
	}

	// 获取统计信息
	stats := recorder.GetStats()

	fmt.Println("\n=== 聚合统计信息 ===")
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	fmt.Println("\n✓ 统计信息已生成")
}
