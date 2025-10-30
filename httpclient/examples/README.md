# HTTP Client 示例

本目录包含了 httpclient 包的各种使用示例，从基本用法到生产环境配置。

## 示例列表

### [01_basic](01_basic/) - 基本用法

演示 httpclient 的基本功能：
- 创建客户端
- GET/POST 请求
- 设置查询参数
- 自定义配置

**运行**:
```bash
cd 01_basic && go run main.go
```

### [02_with_trace](02_with_trace/) - Trace 集成

演示如何集成 OpenTelemetry trace：
- 自动创建 HTTP span
- Span 层级结构
- Trace context 传播
- 业务场景示例

**运行**:
```bash
cd 02_with_trace && go run main.go
```

### [03_with_metric](03_with_metric/) - Metric 集成

演示如何收集 HTTP 请求指标：
- 实现 MetricRecorder 接口
- 收集详细的网络性能数据
- 批量请求统计
- 聚合分析

**运行**:
```bash
cd 03_with_metric && go run main.go
```

### [04_production](04_production/) - 生产环境配置

演示生产环境的完整配置：
- Trace + Log + Metric 完整集成
- 超时和重试配置
- 连接池优化
- 并发请求处理
- 性能统计

**运行**:
```bash
cd 04_production && go run main.go
```

### [05_file_operations](05_file_operations/) - 文件操作和日志优化

演示文件上传/下载的智能日志记录：
- 文件上传的元信息记录
- 文件下载的智能检测
- 敏感头信息自动过滤
- 大文件响应的日志截断

**运行**:
```bash
cd 05_file_operations && go run main.go
```

### [06_insecure_skip_verify](06_insecure_skip_verify/) - 忽略证书验证

演示如何跳过 SSL/TLS 证书验证：
- 使用 `WithInsecureSkipVerify()` 选项
- 访问使用自签名证书的服务
- ⚠️ 仅限开发/测试环境使用

**运行**:
```bash
cd 06_insecure_skip_verify && go run main.go
```

## 学习路径

1. **初学者**: 从 `01_basic` 开始，了解基本的 HTTP 请求
2. **中级**: 学习 `02_with_trace` 和 `03_with_metric`，掌握可观测性集成
3. **高级**: 参考 `04_production`，了解生产环境的最佳实践
4. **特殊场景**: 学习 `05_file_operations` 和 `06_insecure_skip_verify`，处理文件操作和证书问题

## 快速对比

| 特性 | 01_basic | 02_with_trace | 03_with_metric | 04_production | 05_file_operations | 06_insecure_skip_verify |
|------|----------|---------------|----------------|---------------|--------------------|------------------------|
| 基本请求 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 日志记录 | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Trace 集成 | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ |
| Metric 收集 | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ |
| 重试配置 | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| 并发请求 | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| 连接池优化 | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| 文件上传 | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| 文件下载 | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| 敏感信息过滤 | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| 跳过证书验证 | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

## 常见模式

### 创建客户端

```go
client := httpclient.New(
    httpclient.WithLogger(logger.Default()),
    httpclient.WithTrace("my-service"),
    httpclient.WithMetric(recorder),
)
```

### 发起请求

```go
resp, err := client.R(ctx).
    SetHeader("Authorization", "Bearer token").
    SetBody(data).
    Post("https://api.example.com/endpoint")
```

### 处理响应

```go
if err != nil {
    logger.Error(ctx, "请求失败", map[string]interface{}{"error": err})
    return
}

if resp.StatusCode() != 200 {
    logger.Warn(ctx, "非 200 响应", map[string]interface{}{"status": resp.StatusCode()})
}
```

## 依赖说明

所有示例都依赖以下包：
- `github.com/Si40Code/kit/httpclient` - HTTP 客户端
- `github.com/Si40Code/kit/logger` - 日志包

部分示例的额外依赖：
- `go.opentelemetry.io/otel` - OpenTelemetry (02, 04)
- `go.opentelemetry.io/otel/sdk` - OTel SDK (02, 04)
- `go.opentelemetry.io/otel/exporters/stdout/stdouttrace` - Stdout exporter (02, 04)

## 测试 API

所有示例使用 [httpbin.org](https://httpbin.org) 作为测试 API：
- `GET /get` - 测试 GET 请求
- `POST /post` - 测试 POST 请求
- `GET /status/:code` - 返回指定状态码
- `GET /delay/:n` - 延迟 n 秒响应

## 下一步

- 查看 [主 README](../README.md) 了解更多配置选项
- 集成到你的项目中
- 根据需要实现自定义的 MetricRecorder
- 连接到 SigNoz/Jaeger 进行分布式追踪

