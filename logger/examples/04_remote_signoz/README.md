# 示例 4: SigNoz 远程日志

演示如何将日志发送到 SigNoz 可观测性平台，实现日志、trace 和 metrics 的统一管理。

## 功能展示

1. **OTLP 协议发送日志** - 通过 gRPC 发送日志到 SigNoz
2. **Trace 集成** - 日志与 trace 自动关联
3. **多输出** - 同时输出到 stdout 和 SigNoz
4. **完整业务流程** - 展示实际业务场景的可观测性

## 前置条件

### 安装 SigNoz

#### 方式 1: Docker Compose（推荐）

```bash
git clone https://github.com/SigNoz/signoz.git
cd signoz/deploy
./install.sh
```

#### 方式 2: Docker（快速启动）

```bash
docker run -d --name signoz \
  -p 3301:3301 \
  -p 4317:4317 \
  -p 4318:4318 \
  signoz/signoz
```

### 验证 SigNoz 运行

访问 http://localhost:3301 确认 SigNoz UI 可以访问。

## 运行示例

```bash
cd logger/examples/04_remote_signoz
go run main.go
```

## SigNoz 端口说明

| 端口 | 协议 | 用途 |
|------|------|------|
| 3301 | HTTP | SigNoz Web UI |
| 4317 | gRPC | OTLP 接收端点（日志/trace/metrics） |
| 4318 | HTTP | OTLP HTTP 接收端点 |

## 配置说明

### 1. Logger 配置

```go
err := logger.Init(
    logger.WithLevel(logger.DebugLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithStdout(), // 本地也输出
    logger.WithOTLP("localhost:4317", // SigNoz gRPC 端点
        logger.WithOTLPInsecure(), // 本地环境使用不安全连接
    ),
    logger.WithTrace(serviceName), // 启用 trace 集成
)
```

### 2. Trace 配置

```go
exporter, _ := otlptrace.New(
    ctx,
    otlptracehttp.NewClient(
        otlptracehttp.WithEndpoint("localhost:4318"), // SigNoz HTTP 端点
        otlptracehttp.WithInsecure(),
    ),
)

tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resource.NewWithAttributes(
        semconv.ServiceNameKey.String(serviceName),
    )),
)
```

## 在 SigNoz 中查看

### 1. 查看 Traces

1. 访问 http://localhost:3301
2. 点击左侧菜单 "Traces"
3. 可以看到所有 trace 记录
4. 点击具体的 trace 查看详细信息

### 2. 查看 Logs

1. 点击左侧菜单 "Logs"
2. 可以看到所有日志记录
3. 支持按时间、级别、服务名等过滤

### 3. 关联查看

1. 在 Traces 页面点击某个 trace
2. 在 trace 详情页可以看到关联的日志
3. 在 Logs 页面可以通过 trace_id 搜索相关日志

### 4. 错误分析

1. 在 Traces 页面筛选 error 状态的 trace
2. 查看 error span 的详细信息
3. 查看关联的 error 日志

## 业务场景示例

### HTTP 请求处理

```go
ctx, span := tracer.Start(ctx, "http-request-handler")
defer span.End()

logger.Info(ctx, "收到 HTTP 请求",
    "method", "GET",
    "path", "/api/users",
    "client_ip", "192.168.1.100",
)

// 处理请求...

logger.Info(ctx, "请求处理完成",
    "status", 200,
    "duration_ms", 50,
)
```

在 SigNoz 中可以看到：
- 完整的请求 trace
- 每个日志的详细信息
- 请求的耗时分析

### 订单履行流程

```go
ctx, span := tracer.Start(ctx, "order-fulfillment")
defer span.End()

logger.Info(ctx, "开始订单履行流程", "order_id", orderID)

// 步骤 1: 检查库存
ctx, checkSpan := tracer.Start(ctx, "check-inventory")
logger.Info(ctx, "检查库存", "order_id", orderID)
checkSpan.End()

// 步骤 2: 扣减库存
ctx, deductSpan := tracer.Start(ctx, "deduct-inventory")
logger.Info(ctx, "扣减库存", "order_id", orderID)
deductSpan.End()

// 步骤 3: 创建发货单
ctx, shippingSpan := tracer.Start(ctx, "create-shipping")
logger.Info(ctx, "创建发货单", "order_id", orderID)
shippingSpan.End()
```

在 SigNoz 中可以看到：
- 订单履行的完整流程
- 每个步骤的耗时
- 步骤间的依赖关系
- 每个步骤的详细日志

## 生产环境配置

### 1. 使用环境变量

```go
signozEndpoint := os.Getenv("SIGNOZ_ENDPOINT")
if signozEndpoint == "" {
    signozEndpoint = "localhost:4317"
}

logger.Init(
    logger.WithOTLP(signozEndpoint,
        logger.WithOTLPInsecure(),
    ),
)
```

### 2. 使用 TLS

```go
logger.Init(
    logger.WithOTLP("signoz.example.com:4317",
        // 不使用 WithOTLPInsecure()，默认使用 TLS
    ),
)
```

### 3. 设置超时

```go
logger.Init(
    logger.WithOTLP("signoz.example.com:4317",
        logger.WithOTLPTimeout(10*time.Second),
    ),
)
```

### 4. 添加自定义 Headers

```go
logger.Init(
    logger.WithOTLP("signoz.example.com:4317",
        logger.WithOTLPHeaders(map[string]string{
            "Authorization": "Bearer token",
            "X-Custom-Header": "value",
        }),
    ),
)
```

## Docker Compose 示例

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  app:
    build: .
    environment:
      - SIGNOZ_ENDPOINT=signoz:4317
      - ENV=production
    depends_on:
      - signoz
  
  signoz:
    image: signoz/signoz:latest
    ports:
      - "3301:3301"
      - "4317:4317"
      - "4318:4318"
```

## 性能优化

### 1. 批量发送

OTLP exporter 会自动批量发送日志，默认配置：
- 批量大小: 100 条
- 批量间隔: 5 秒

### 2. 异步发送

日志发送是异步的，不会阻塞主流程：

```go
// 日志记录立即返回
logger.Info(ctx, "处理请求")

// 日志在后台批量发送到 SigNoz
```

### 3. 降级处理

如果 SigNoz 不可用，日志仍会输出到 stdout：

```go
logger.Init(
    logger.WithStdout(),                // 本地输出（降级）
    logger.WithOTLP("signoz:4317"),     // 远程输出（可能失败）
)
```

## 常见问题

### 1. 连接失败

**错误**: `failed to connect to OTLP endpoint: context deadline exceeded`

**解决**:
- 检查 SigNoz 是否运行: `docker ps | grep signoz`
- 检查端口是否正确: 4317 (gRPC) 或 4318 (HTTP)
- 检查网络连接: `telnet localhost 4317`

### 2. 日志未显示

**可能原因**:
- 日志级别过滤
- SigNoz 索引延迟（等待几秒）
- 服务名不匹配

**解决**:
- 在 SigNoz UI 调整日志级别过滤器
- 刷新页面或等待片刻
- 检查 `serviceName` 配置

### 3. Trace 和日志未关联

**原因**: 未启用 trace 集成

**解决**:
```go
logger.Init(
    logger.WithTrace(serviceName), // 必须启用
    logger.WithOTLP("signoz:4317"),
)
```

## 下一步

- [示例 5: 生产环境配置](../05_production/) - 完整的生产级配置
- [SigNoz 文档](https://signoz.io/docs/) - 了解更多 SigNoz 功能

