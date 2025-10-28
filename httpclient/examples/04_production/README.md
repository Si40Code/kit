# 示例 4: 生产环境配置

这个示例演示了在生产环境中如何配置和使用 httpclient 包。

## 功能特性

- 完整的 Trace + Log + Metric 集成
- 重试机制
- 超时配置
- 连接池优化
- 并发请求处理
- 性能统计

## 运行示例

```bash
cd examples/04_production
go run main.go
```

## 生产级配置

### 1. 完整的客户端配置

```go
client := httpclient.New(
    // 日志配置
    httpclient.WithLogger(logger.L()),
    
    // Trace 配置
    httpclient.WithTrace("production-service"),
    
    // Metric 配置
    httpclient.WithMetric(recorder),
    
    // 超时配置
    httpclient.WithTimeout(10*time.Second),
    
    // 重试配置（重试 3 次，初始等待 100ms，最大等待 2s）
    httpclient.WithRetry(3, 100*time.Millisecond, 2*time.Second),
    
    // 日志大小限制
    httpclient.WithMaxBodyLogSize(5*1024), // 5KB
    
    // 连接池配置
    httpclient.WithMaxIdleConns(100),
    httpclient.WithIdleConnTimeout(90*time.Second),
    httpclient.WithKeepAlive(30*time.Second),
)
```

### 2. 重试机制

客户端会自动重试失败的请求：
- 支持配置重试次数
- 支持配置重试等待时间
- 自动使用指数退避策略

### 3. 连接复用

通过配置连接池参数，实现高效的连接复用：
- `MaxIdleConns`: 最大空闲连接数
- `IdleConnTimeout`: 空闲连接超时时间
- `KeepAlive`: keep-alive 时间

## 使用场景

### 1. 单个请求

```go
resp, err := client.R(ctx).
    SetHeader("X-Request-ID", "req-123").
    SetBody(request).
    Post("https://api.example.com/endpoint")
```

### 2. 并发请求

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        reqCtx, span := tracer.Start(ctx, "request")
        defer span.End()
        
        resp, err := client.R(reqCtx).Get(url)
        // 处理响应
    }(i)
}
wg.Wait()
```

### 3. 业务场景

每个请求都会：
1. 自动创建 span 并关联到父 span
2. 记录完整的请求/响应日志
3. 收集详细的性能指标
4. 失败时自动重试

## 监控指标

### 关键指标

- **总请求数**: 总的 HTTP 请求次数
- **成功率**: 成功请求占比
- **平均响应时间**: 请求的平均总耗时
- **DNS 查询时间**: DNS 解析的平均耗时
- **TCP 连接时间**: TCP 连接建立的平均耗时
- **TLS 握手时间**: TLS 握手的平均耗时
- **服务器处理时间**: 服务器处理请求的平均耗时
- **连接复用率**: 连接复用的比例

### 输出示例

```
=== 性能统计 ===
total_requests: 16
success_requests: 15
failed_requests: 1
success_rate: 93.75%
avg_total_time_ms: 245
avg_dns_lookup_ms: 2
avg_tcp_conn_ms: 3
avg_tls_handshake_ms: 1
status_codes: map[200:15 500:1]
```

## 最佳实践

### 1. 超时配置

根据服务的 SLA 设置合理的超时时间：
```go
httpclient.WithTimeout(10*time.Second)
```

### 2. 重试配置

对于幂等操作启用重试：
```go
httpclient.WithRetry(3, 100*time.Millisecond, 2*time.Second)
```

### 3. 日志大小限制

避免记录过大的请求/响应体：
```go
httpclient.WithMaxBodyLogSize(5*1024) // 5KB
```

### 4. 连接池优化

根据并发量调整连接池大小：
```go
httpclient.WithMaxIdleConns(100)
httpclient.WithIdleConnTimeout(90*time.Second)
```

### 5. Context 传递

始终传递 context 以支持超时和取消：
```go
resp, err := client.R(ctx).Get(url)
```

### 6. Trace 集成

在业务逻辑中创建 span，HTTP 请求会自动作为子 span：
```go
ctx, span := tracer.Start(ctx, "business-operation")
defer span.End()

// HTTP 请求会自动创建 child span
resp, err := client.R(ctx).Get(url)
```

## 生产环境检查清单

- [ ] 配置合理的超时时间
- [ ] 启用重试机制（针对幂等操作）
- [ ] 配置日志大小限制
- [ ] 优化连接池参数
- [ ] 集成 Trace 系统（SigNoz/Jaeger）
- [ ] 集成 Metric 系统（Prometheus）
- [ ] 添加请求 ID 追踪
- [ ] 配置 TLS 证书验证
- [ ] 监控连接复用率
- [ ] 设置告警规则

