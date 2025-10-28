# HTTP Client

生产级的 HTTP 客户端包，基于 [go-resty/resty](https://github.com/go-resty/resty)，提供完整的 Trace、日志和 Metric 集成。

## 特性

- ✅ **OpenTelemetry Trace 集成**: 自动创建和传播 span，无缝集成到分布式追踪系统
- ✅ **完整的日志记录**: 记录请求/响应的所有详情，包括网络性能指标
- ✅ **智能文件处理**: 自动检测文件上传/下载，只记录元信息而不记录文件内容
- ✅ **敏感信息保护**: 自动过滤 Authorization、API-Key 等敏感头信息
- ✅ **详细的 Metric**: 收集 DNS、TCP、TLS、服务器处理时间等性能数据
- ✅ **统一的配置模式**: 遵循 kit 的 Option 配置风格
- ✅ **自动重试**: 支持可配置的重试机制
- ✅ **连接池优化**: 高效的连接复用和管理
- ✅ **生产就绪**: 包含超时、重试、日志大小限制等生产级特性

## 快速开始

### 安装

```bash
go get github.com/Si40Code/kit/httpclient
```

### 基本用法

```go
package main

import (
    "context"
    "github.com/Si40Code/kit/httpclient"
    "github.com/Si40Code/kit/logger"
)

func main() {
    // 初始化 logger
    logger.Init(
        logger.WithLevel(logger.InfoLevel),
        logger.WithStdout(),
    )
    defer logger.Sync()

    // 创建 HTTP 客户端
    client := httpclient.New(
        httpclient.WithLogger(logger.L()),
    )

    // 发起请求
    resp, err := client.R(context.Background()).
        Get("https://api.example.com/users")

    if err != nil {
        panic(err)
    }

    println("Status:", resp.StatusCode())
}
```

## 配置选项

### 日志配置

```go
client := httpclient.New(
    // 设置 logger
    httpclient.WithLogger(logger.L()),
    
    // 禁用日志
    httpclient.WithDisableLog(),
    
    // 设置最大日志 body 大小（字节）
    httpclient.WithMaxBodyLogSize(10 * 1024),
)
```

### Trace 配置

```go
client := httpclient.New(
    // 启用 trace
    httpclient.WithTrace("my-service"),
)
```

### Metric 配置

```go
// 实现 MetricRecorder 接口
type MyMetricRecorder struct{}

func (r *MyMetricRecorder) RecordRequest(data httpclient.MetricData) {
    // 发送到 Prometheus、SigNoz 等
}

client := httpclient.New(
    // 启用 metric
    httpclient.WithMetric(&MyMetricRecorder{}),
)
```

### HTTP 客户端配置

```go
client := httpclient.New(
    // 设置超时
    httpclient.WithTimeout(10 * time.Second),
    
    // 设置重试（次数、初始等待、最大等待）
    httpclient.WithRetry(3, 100*time.Millisecond, 2*time.Second),
    
    // 跳过 TLS 证书验证（不推荐生产使用）
    httpclient.WithInsecureSkipVerify(),
)
```

### 连接池配置

```go
client := httpclient.New(
    // 最大空闲连接数
    httpclient.WithMaxIdleConns(100),
    
    // 空闲连接超时
    httpclient.WithIdleConnTimeout(90 * time.Second),
    
    // TLS 握手超时
    httpclient.WithTLSHandshakeTimeout(10 * time.Second),
    
    // Keep-Alive 时间
    httpclient.WithKeepAlive(30 * time.Second),
)
```

## Trace 集成

HTTP 客户端会自动集成到 OpenTelemetry trace 系统中：

```go
// 业务代码中创建 span
ctx, span := tracer.Start(ctx, "my-operation")
defer span.End()

// HTTP 请求会自动作为 child span
resp, err := client.R(ctx).Get("https://api.example.com/data")
```

每个 HTTP 请求的 span 会自动记录：
- HTTP 方法和 URL
- 请求/响应状态码
- 总响应时间
- DNS 查询时间
- TCP 连接时间
- TLS 握手时间
- 服务器处理时间
- 连接复用信息

## 日志记录

客户端会自动记录详细的请求日志：

```json
{
  "level": "info",
  "msg": "HTTP request completed successfully",
  "http.method": "GET",
  "http.url": "https://api.example.com/users",
  "http.status_code": 200,
  "http.total_time_ms": 234,
  "http.dns_lookup_ms": 3,
  "http.tcp_conn_ms": 45,
  "http.tls_handshake_ms": 89,
  "http.server_time_ms": 97,
  "http.conn_reused": true,
  "trace_id": "abc123...",
  "span_id": "def456..."
}
```

### 智能文件处理

对于文件上传/下载，日志只记录元信息：

```go
// 文件上传
resp, err := client.R(ctx).
    SetFile("file", "/path/to/file.pdf").
    Post("https://api.example.com/upload")

// 日志输出：
// "http.request.body": "[multipart/form-data, size: 1234567 bytes]"
```

```go
// 文件下载
resp, err := client.R(ctx).
    Get("https://api.example.com/download/file.pdf")

// 日志输出：
// "http.response.body": "[file download, size: 1234567 bytes]"
```

### 敏感信息保护

自动过滤敏感的请求头：

```go
resp, err := client.R(ctx).
    SetHeader("Authorization", "Bearer secret-token").
    SetHeader("X-API-Key", "my-api-key").
    Get(url)

// 日志输出：
// "http.request.headers": {
//   "Authorization": "******",
//   "X-API-Key": "******",
//   "User-Agent": "MyApp/1.0"
// }
```

## Metric 数据

通过实现 `MetricRecorder` 接口，可以收集以下指标：

```go
type MetricData struct {
    Method        string        // HTTP 方法
    Host          string        // 主机名
    Path          string        // 请求路径
    StatusCode    int           // 状态码
    TotalTime     time.Duration // 总耗时
    DNSLookup     time.Duration // DNS 查询时间
    TCPConn       time.Duration // TCP 连接时间
    TLSHandshake  time.Duration // TLS 握手时间
    ServerTime    time.Duration // 服务器处理时间
    ResponseTime  time.Duration // 响应传输时间
    IsConnReused  bool          // 连接是否复用
    IsConnWasIdle bool          // 连接是否空闲
    ConnIdleTime  time.Duration // 连接空闲时间
    RequestAttempt int          // 请求尝试次数
    RemoteAddr    string        // 远程地址
}
```

## 示例

查看 `examples/` 目录了解更多使用示例：

- [01_basic](examples/01_basic/README.md) - 基本用法
- [02_with_trace](examples/02_with_trace/README.md) - Trace 集成
- [03_with_metric](examples/03_with_metric/README.md) - Metric 集成
- [04_production](examples/04_production/README.md) - 生产环境配置
- [05_file_operations](examples/05_file_operations/README.md) - 文件操作和日志优化

## API 文档

### Client

```go
// New 创建一个新的 HTTP 客户端
func New(opts ...Option) *Client

// R 创建一个新的请求
func (c *Client) R(ctx context.Context) *resty.Request
```

### MetricRecorder 接口

```go
type MetricRecorder interface {
    RecordRequest(data MetricData)
}
```

### Option 函数

| 函数 | 说明 |
|------|------|
| `WithLogger(logger.Logger)` | 设置日志记录器 |
| `WithDisableLog()` | 禁用日志 |
| `WithMaxBodyLogSize(int64)` | 设置最大日志 body 大小 |
| `WithTrace(string)` | 启用 trace（服务名） |
| `WithMetric(MetricRecorder)` | 启用 metric |
| `WithTimeout(time.Duration)` | 设置超时时间 |
| `WithRetry(int, time.Duration, time.Duration)` | 设置重试配置 |
| `WithInsecureSkipVerify()` | 跳过 TLS 验证 |
| `WithMaxIdleConns(int)` | 设置最大空闲连接数 |
| `WithIdleConnTimeout(time.Duration)` | 设置空闲连接超时 |
| `WithTLSHandshakeTimeout(time.Duration)` | 设置 TLS 握手超时 |
| `WithKeepAlive(time.Duration)` | 设置 keep-alive 时间 |

## 最佳实践

1. **始终传递 context**: 使用 `client.R(ctx)` 而不是 `client.Client.R()`
2. **设置合理的超时**: 根据服务的 SLA 设置超时时间
3. **谨慎使用重试**: 只对幂等操作启用重试
4. **限制日志大小**: 避免记录过大的请求/响应体
5. **监控连接复用**: 通过 metric 监控连接复用率
6. **集成 trace**: 在生产环境中集成到 SigNoz/Jaeger

## 与 Resty 的关系

本包基于 [go-resty/resty](https://github.com/go-resty/resty)，并直接暴露 `*resty.Client`。这意味着：

- 可以使用所有 Resty 的原生功能
- 无需学习新的 API
- 通过钩子自动增强功能（trace、log、metric）

```go
// 可以直接使用 Resty 的所有方法
resp, err := client.R(ctx).
    SetHeader("Authorization", "Bearer token").
    SetQueryParam("page", "1").
    SetBody(data).
    Post("https://api.example.com/endpoint")
```

## 性能考虑

- **连接复用**: 默认启用连接池，自动复用连接
- **Keep-Alive**: 默认启用，减少连接建立开销
- **日志优化**: 可配置日志大小限制，避免性能影响
- **并发安全**: 客户端是并发安全的，可以在多个 goroutine 中共享

## License

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！

