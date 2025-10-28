# 示例 3: Metric 集成

这个示例演示了如何使用 httpclient 包的 metric 记录功能。

## 功能特性

- 记录详细的 HTTP 请求指标
- 包含网络层面的性能数据（DNS、TCP、TLS 等）
- 支持自定义 MetricRecorder 实现
- 可集成到 Prometheus、SigNoz 等监控系统

## 运行示例

```bash
cd examples/03_with_metric
go run main.go
```

## 代码说明

### 1. 实现 MetricRecorder 接口

```go
type SimpleMetricRecorder struct {
    mu      sync.Mutex
    metrics []httpclient.MetricData
}

func (r *SimpleMetricRecorder) RecordRequest(data httpclient.MetricData) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.metrics = append(r.metrics, data)
}
```

### 2. 创建启用 Metric 的客户端

```go
recorder := NewSimpleMetricRecorder()

client := httpclient.New(
    httpclient.WithLogger(logger.L()),
    httpclient.WithMetric(recorder),
)
```

### 3. 发起请求（自动记录 metric）

```go
resp, err := client.R(context.Background()).
    Get("https://httpbin.org/get")
```

## MetricData 包含的字段

每个 HTTP 请求会记录以下指标：

| 字段 | 类型 | 说明 |
|------|------|------|
| Method | string | HTTP 方法 |
| Host | string | 主机名 |
| Path | string | 请求路径 |
| StatusCode | int | HTTP 状态码 |
| TotalTime | Duration | 总耗时 |
| DNSLookup | Duration | DNS 查询时间 |
| TCPConn | Duration | TCP 连接时间 |
| TLSHandshake | Duration | TLS 握手时间 |
| ServerTime | Duration | 服务器处理时间 |
| ResponseTime | Duration | 响应传输时间 |
| IsConnReused | bool | 连接是否复用 |
| IsConnWasIdle | bool | 连接是否空闲 |
| ConnIdleTime | Duration | 连接空闲时间 |
| RequestAttempt | int | 请求尝试次数 |
| RemoteAddr | string | 远程地址 |

## 与 Prometheus 集成示例

```go
type PrometheusRecorder struct {
    requestDuration *prometheus.HistogramVec
    requestTotal    *prometheus.CounterVec
}

func (r *PrometheusRecorder) RecordRequest(data httpclient.MetricData) {
    labels := prometheus.Labels{
        "method": data.Method,
        "status": strconv.Itoa(data.StatusCode),
    }
    
    r.requestDuration.With(labels).Observe(data.TotalTime.Seconds())
    r.requestTotal.With(labels).Inc()
}
```

## 输出示例

```
[Metric] HTTP Request:
  Method:        GET
  Path:          /get
  Host:          https://httpbin.org/get
  Status Code:   200
  Total Time:    523ms
  DNS Lookup:    45ms
  TCP Conn:      89ms
  TLS Handshake: 124ms
  Server Time:   265ms
  Response Time: 15ms
  Conn Reused:   false
  Remote Addr:   54.221.198.78:443

=== 聚合统计信息 ===
total_requests: 10
avg_total_time_ms: 234
avg_dns_lookup_ms: 3
avg_tcp_conn_ms: 2
avg_tls_handshake_ms: 1
avg_server_time_ms: 228
conn_reuse_rate: 90.00%
conn_idle_rate: 80.00%
status_codes: map[200:10]
```

## 应用场景

1. **性能监控**: 实时监控 HTTP 请求性能
2. **问题诊断**: 快速定位网络问题（DNS、连接、TLS 等）
3. **容量规划**: 分析请求模式和性能趋势
4. **SLA 监控**: 跟踪响应时间和成功率

