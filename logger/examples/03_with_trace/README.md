# 示例 3: Trace 集成

演示如何将 logger 与 OpenTelemetry trace 集成，实现日志与链路追踪的关联。

## 功能展示

1. **自动提取 trace 信息** - 从 context 自动提取 trace_id 和 span_id
2. **Error 日志标记 Span** - Error/Fatal 日志自动将 span 状态设置为 error
3. **日志作为 Span Event** - 日志自动记录为 span 事件
4. **嵌套 Span 支持** - 正确处理父子 span 关系
5. **业务场景示例** - 完整的用户注册流程 trace

## 运行示例

```bash
cd logger/examples/03_with_trace
go run main.go
```

## Trace 集成效果

### 普通日志（无 trace）

```
2024-01-15T10:30:00.123+0800    INFO    操作开始
```

### 带 trace 的日志

```
2024-01-15T10:30:00.123+0800    INFO    操作开始    {"trace_id": "4bf92f3577b34da6a3ce929d0e0e4736", "span_id": "00f067aa0ba902b7"}
```

## 关键代码

### 1. 启用 Trace 集成

```go
err := logger.Init(
    logger.WithLevel(logger.DebugLevel),
    logger.WithTrace("my-service"), // 启用 trace 集成
    logger.WithStdout(),
)
```

### 2. 创建 Span 并记录日志

```go
ctx := context.Background()
tracer := otel.Tracer("my-service")

// 创建 span
ctx, span := tracer.Start(ctx, "operation-name")
defer span.End()

// 记录日志会自动包含 trace_id 和 span_id
logger.Info(ctx, "操作开始")
logger.Error(ctx, "操作失败") // 自动标记 span 为 error
```

### 3. 嵌套 Span

```go
// 父 span
ctx, parentSpan := tracer.Start(ctx, "parent-operation")
defer parentSpan.End()

logger.Info(ctx, "父操作开始")

// 子 span
ctx, childSpan := tracer.Start(ctx, "child-operation")
logger.Info(ctx, "子操作执行") // 包含子 span 的 trace 信息
childSpan.End()

logger.Info(ctx, "父操作完成")
```

## Trace 信息字段

日志中会自动添加以下字段：

| 字段 | 说明 | 示例 |
|------|------|------|
| `trace_id` | Trace ID（16 字节，32 字符十六进制） | `4bf92f3577b34da6a3ce929d0e0e4736` |
| `span_id` | Span ID（8 字节，16 字符十六进制） | `00f067aa0ba902b7` |

## Error 日志的特殊处理

当记录 Error 或 Fatal 级别的日志时，logger 会自动：

1. **设置 Span 状态** - 将 span 状态设置为 `codes.Error`
2. **记录错误** - 调用 `span.RecordError()`
3. **添加事件** - 将日志作为 span 事件记录

```go
logger.Error(ctx, "数据库连接失败", "error", "timeout")

// 等价于：
// span.SetStatus(codes.Error, "数据库连接失败")
// span.RecordError(errors.New("数据库连接失败"))
// span.AddEvent("数据库连接失败")
```

## 与可观测性平台集成

### SigNoz

```go
logger.Init(
    logger.WithTrace("my-service"),
    logger.WithOTLP("signoz-host:4317"), // 日志发送到 SigNoz
)
```

在 SigNoz 中可以：
- 查看 trace 详情
- 点击 span 查看关联的日志
- 通过 trace_id 搜索日志
- 分析错误 span 和对应日志

### Jaeger

```go
// 初始化 Jaeger exporter
exporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(
    jaeger.WithEndpoint("http://localhost:14268/api/traces"),
))

tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
)
otel.SetTracerProvider(tp)

// 初始化 logger
logger.Init(
    logger.WithTrace("my-service"),
    logger.WithStdout(),
)
```

## 业务场景示例

### 用户注册流程

```go
ctx, span := tracer.Start(ctx, "user-registration")
defer span.End()

logger.Info(ctx, "收到注册请求", "email", email)

// 步骤 1: 验证邮箱
ctx, validateSpan := tracer.Start(ctx, "validate-email")
logger.Info(ctx, "验证邮箱", "email", email)
validateSpan.End()

// 步骤 2: 创建用户
ctx, createSpan := tracer.Start(ctx, "create-user")
logger.Info(ctx, "创建用户记录", "user_id", userID)
createSpan.End()

// 步骤 3: 发送欢迎邮件
ctx, emailSpan := tracer.Start(ctx, "send-welcome-email")
if err := sendEmail(email); err != nil {
    logger.Error(ctx, "邮件发送失败", "error", err)
}
emailSpan.End()

logger.Info(ctx, "注册流程完成", "user_id", userID)
```

在 trace 视图中可以看到：
- 完整的注册流程耗时
- 每个步骤的耗时分布
- 邮件发送失败的 error span
- 每个步骤的详细日志

## 最佳实践

### 1. 始终传递 Context

```go
// ✅ 正确
func processOrder(ctx context.Context, orderID string) {
    logger.Info(ctx, "处理订单", "order_id", orderID)
}

// ❌ 错误
func processOrder(orderID string) {
    ctx := context.Background() // 会丢失 trace 信息
    logger.Info(ctx, "处理订单", "order_id", orderID)
}
```

### 2. 为关键操作创建 Span

```go
func handleRequest(ctx context.Context) {
    tracer := otel.Tracer("my-service")
    
    // 为关键操作创建 span
    ctx, span := tracer.Start(ctx, "database-query")
    defer span.End()
    
    logger.Info(ctx, "执行数据库查询")
    // ... 数据库操作
}
```

### 3. Error 日志包含详细信息

```go
// ✅ 正确 - 包含详细上下文
logger.Error(ctx, "订单处理失败",
    "order_id", orderID,
    "error", err.Error(),
    "retry_count", retryCount,
)

// ❌ 错误 - 信息不足
logger.Error(ctx, "失败")
```

### 4. 使用有意义的 Span 名称

```go
// ✅ 正确 - 清晰的操作名称
ctx, span := tracer.Start(ctx, "payment-process")
ctx, span := tracer.Start(ctx, "inventory-check")

// ❌ 错误 - 模糊的名称
ctx, span := tracer.Start(ctx, "process")
ctx, span := tracer.Start(ctx, "step1")
```

## 性能考虑

1. **Trace 提取开销** - 很小，可以忽略不计
2. **Span Event 记录** - 轻量级操作
3. **Error 标记** - 仅在 Error/Fatal 时执行

## 下一步

- [示例 4: SigNoz 远程日志](../04_remote_signoz/) - 完整的可观测性集成
- [示例 5: 生产环境配置](../05_production/) - 生产级配置

