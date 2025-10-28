# 示例 2: Trace 集成

这个示例演示了如何使用 httpclient 包的 OpenTelemetry trace 集成功能。

## 功能特性

- 自动为每个 HTTP 请求创建 span
- 自动传播 trace context
- 记录详细的请求/响应信息到 span
- 错误请求自动标记 span 为 error
- 与业务 span 完美集成

## 运行示例

```bash
cd examples/02_with_trace
go run main.go
```

## 代码说明

### 1. 初始化 Tracer

```go
cleanup := initTracer("httpclient-trace-example")
defer cleanup()
```

### 2. 创建启用 Trace 的客户端

```go
client := httpclient.New(
    httpclient.WithLogger(logger.L()),
    httpclient.WithTrace("httpclient-example"),
)
```

### 3. 在业务 Span 中发起请求

```go
ctx, span := tracer.Start(ctx, "my-operation")
defer span.End()

// HTTP 请求会自动创建 child span
resp, err := client.R(ctx).Get("https://httpbin.org/get")
```

## Trace 信息包含

每个 HTTP 请求的 span 会自动记录：

- `http.method`: HTTP 方法
- `http.url`: 请求 URL
- `http.status_code`: 响应状态码
- `http.response.time_ms`: 总响应时间
- `http.dns_lookup_ms`: DNS 查询时间
- `http.tcp_conn_ms`: TCP 连接时间
- `http.tls_handshake_ms`: TLS 握手时间
- `http.server_time_ms`: 服务器处理时间
- `http.conn_reused`: 连接是否复用

## Span 层级结构

```
user-registration (业务 span)
├── validate-email (业务子 span)
│   └── HTTP GET (HTTP 请求 span)
├── create-user (业务子 span)
│   └── HTTP POST (HTTP 请求 span)
└── send-welcome-email (业务子 span)
    └── HTTP POST (HTTP 请求 span)
```

## 输出示例

程序会输出：
1. 详细的日志（包含 trace_id 和 span_id）
2. Trace 信息（stdout exporter 格式）
3. 完整的 span 层级结构

## 与 SigNoz/Jaeger 集成

将 stdout exporter 替换为 OTLP exporter 即可：

```go
exporter, _ := otlptrace.New(
    context.Background(),
    otlptracehttp.NewClient(
        otlptracehttp.WithEndpoint("localhost:4318"),
        otlptracehttp.WithInsecure(),
    ),
)
```

