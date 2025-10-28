# 示例 02: Trace 集成

本示例展示如何集成 OpenTelemetry Trace，实现数据库操作的完整链路追踪。

## 功能展示

1. **OpenTelemetry 初始化**：设置 TracerProvider 和 Exporter
2. **自动 Span 创建**：每个数据库查询自动创建独立的 span
3. **Span 层级关系**：展示业务 span 和数据库 span 的父子关系
4. **详细的 Span 属性**：记录 SQL、表名、耗时、影响行数等
5. **Context 传播**：自动传播 trace context

## Span 层级结构

```
业务操作 (CreateAndQueryUser)
├─ DB INSERT (users)
├─ DB SELECT (users)
└─ DB UPDATE (users)

业务操作 (BatchQueryUsers)
├─ DB SELECT (users)
└─ DB DELETE (users)
```

## Span 属性

每个数据库 span 包含以下属性：

| 属性 | 说明 | 示例 |
|------|------|------|
| `db.system` | 数据库系统 | `gorm` |
| `db.operation` | 操作类型 | `SELECT`, `INSERT`, `UPDATE`, `DELETE` |
| `db.table` | 表名 | `users` |
| `db.statement` | SQL 语句 | `SELECT * FROM users WHERE id = ?` |
| `db.rows_affected` | 影响行数 | `1` |
| `db.duration_ms` | 执行耗时（毫秒） | `5` |

## 运行示例

```bash
cd examples/02_with_trace
go run main.go
```

## 预期输出

### 日志输出

```
2024-01-01T10:00:00.000+0800    INFO    OpenTelemetry initialized
2024-01-01T10:00:00.100+0800    INFO    Successfully connected to database
2024-01-01T10:00:00.200+0800    INFO    Starting business operation: CreateAndQueryUser
2024-01-01T10:00:00.210+0800    INFO    database query executed    {"duration_ms": 8, "rows_affected": 1, "sql": "INSERT INTO `users` ...", "trace_id": "abc123...", "span_id": "def456..."}
2024-01-01T10:00:00.220+0800    INFO    User created    {"id": 1}
...
```

### Trace 输出

你将看到详细的 trace 输出，包含所有 span 的层级关系和属性：

```json
{
  "Name": "CreateAndQueryUser",
  "SpanContext": {
    "TraceID": "abc123...",
    "SpanID": "def456..."
  },
  "Parent": {...},
  "SpanKind": "Internal",
  "StartTime": "2024-01-01T10:00:00.200Z",
  "EndTime": "2024-01-01T10:00:01.000Z",
  "Attributes": [...],
  "Events": null,
  "Links": null,
  "Status": {
    "Code": "Ok"
  },
  "ChildSpans": [
    {
      "Name": "DB INSERT",
      "Attributes": [
        {"Key": "db.system", "Value": "gorm"},
        {"Key": "db.operation", "Value": "INSERT"},
        {"Key": "db.table", "Value": "users"},
        {"Key": "db.statement", "Value": "INSERT INTO `users` ..."},
        {"Key": "db.rows_affected", "Value": 1},
        {"Key": "db.duration_ms", "Value": 8}
      ]
    },
    ...
  ]
}
```

## 代码说明

### 1. 初始化 OpenTelemetry

```go
exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
tp := trace.NewTracerProvider(
    trace.WithBatcher(exporter),
)
otel.SetTracerProvider(tp)
```

在生产环境中，你应该使用 OTLP exporter 发送到 SigNoz、Jaeger 等后端：

```go
import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

exporter, err := otlptracehttp.New(
    context.Background(),
    otlptracehttp.WithEndpoint("signoz:4318"),
    otlptracehttp.WithInsecure(),
)
```

### 2. 启用 ORM Trace

```go
client, err := orm.New(
    mysql.Open(dsn),
    orm.WithTrace("orm-client"), // 设置 service name
)
```

### 3. 创建业务 Span

```go
tracer := otel.Tracer("business-logic")
ctx, span := tracer.Start(ctx, "CreateAndQueryUser")
defer span.End()

// 在这个 context 下的所有数据库操作都会成为子 span
client.WithContext(ctx).Create(&user)
```

### 4. Context 传播

关键是使用 `WithContext(ctx)`：

```go
// ✅ 正确：会创建子 span
client.WithContext(ctx).First(&user, id)

// ❌ 错误：不会关联到父 span
client.First(&user, id)
```

## 集成到 SigNoz

### 1. 修改 Exporter

```go
import (
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/sdk/resource"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer() *trace.TracerProvider {
    res, err := resource.New(
        context.Background(),
        resource.WithAttributes(
            semconv.ServiceName("my-service"),
            semconv.ServiceVersion("1.0.0"),
        ),
    )
    if err != nil {
        panic(err)
    }

    exporter, err := otlptracehttp.New(
        context.Background(),
        otlptracehttp.WithEndpoint("signoz:4318"),
        otlptracehttp.WithInsecure(),
    )
    if err != nil {
        panic(err)
    }

    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(res),
    )
    
    otel.SetTracerProvider(tp)
    return tp
}
```

### 2. 在 SigNoz 中查看

访问 SigNoz UI (通常是 http://localhost:3301)，你将看到：

1. **服务列表**：显示 `my-service` 和 `orm-client`
2. **Trace 详情**：展开可以看到完整的调用链
3. **数据库性能**：按表、操作类型聚合的性能统计
4. **慢查询识别**：自动高亮慢查询

## 常见问题

### Q1: 为什么看不到 trace 数据？

确保：
1. 已正确初始化 TracerProvider
2. 使用了 `WithContext(ctx)` 传递 context
3. 调用了 `span.End()` 结束 span
4. Exporter 配置正确（端点、认证等）

### Q2: 如何减少 trace 开销？

```go
// 使用采样器
tp := trace.NewTracerProvider(
    trace.WithSampler(trace.TraceIDRatioBased(0.1)), // 10% 采样率
    ...
)
```

### Q3: 如何自定义 span 属性？

```go
import "go.opentelemetry.io/otel/attribute"

span.SetAttributes(
    attribute.String("user.id", "123"),
    attribute.Int("user.age", 25),
)
```

## 最佳实践

1. **始终传递 Context**：确保 trace 信息能正确传播
2. **合理命名 Span**：使用描述性的 span 名称
3. **设置服务名**：便于在监控系统中区分
4. **采样策略**：生产环境使用合理的采样率
5. **资源属性**：添加服务版本、环境等信息

## 下一步

查看其他示例：
- [03_with_metric](../03_with_metric/) - Metric 集成
- [04_production](../04_production/) - 生产环境完整配置

