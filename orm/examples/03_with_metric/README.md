# 示例 03: Metric 集成

本示例展示如何集成 Metric 监控，收集和分析数据库查询性能指标。

## 功能展示

1. **实现 MetricRecorder 接口**：自定义 metric 记录器
2. **实时 Metric 收集**：记录每个查询的详细指标
3. **性能统计分析**：按操作类型、表名聚合统计
4. **错误监控**：追踪查询错误
5. **性能洞察**：识别慢查询和热点表

## Metric 数据结构

每个数据库查询都会产生以下指标：

```go
type MetricData struct {
    Operation    string        // 操作类型：SELECT, INSERT, UPDATE, DELETE
    Table        string        // 表名
    SQL          string        // 完整 SQL 语句
    Duration     time.Duration // 查询耗时
    RowsAffected int64        // 影响/返回的行数
    Error        error        // 错误（如果有）
}
```

## 运行示例

```bash
cd examples/03_with_metric
go run main.go
```

## 预期输出

### 实时 Metric 输出

```
🔄 Starting database operations...

📊 Metric: INSERT on users, duration=8ms, rows=1
📊 Metric: INSERT on users, duration=7ms, rows=1
📊 Metric: INSERT on users, duration=6ms, rows=1
📊 Metric: SELECT on users, duration=3ms, rows=3
📊 Metric: SELECT on users, duration=2ms, rows=2
📊 Metric: SELECT on users, duration=1ms, rows=1
📊 Metric: UPDATE on users, duration=5ms, rows=1
📊 Metric: UPDATE on users, duration=4ms, rows=1
📊 Metric: DELETE on users, duration=3ms, rows=1
📊 Metric: SELECT on users, duration=2ms, rows=0, error=record not found
```

### 统计摘要

```
============================================================
📈 Metrics Summary
============================================================

Total Queries: 10
Total Duration: 41ms
Average Duration: 4ms
Error Count: 1

By Operation Type:
--------------------------------------------------
Operation  |    Count | Total(ms) | Avg(ms)
--------------------------------------------------
INSERT     |        3 |        21 |        7
SELECT     |        4 |        8  |        2 (errors: 1)
UPDATE     |        2 |        9  |        4
DELETE     |        1 |        3  |        3

By Table:
--------------------------------------------------
Table      |    Count | Total(ms) | Avg(ms)
--------------------------------------------------
users      |       10 |        41 |        4
============================================================
```

## 代码说明

### 1. 实现 MetricRecorder 接口

```go
type SimpleMetricRecorder struct {
    mu      sync.Mutex
    metrics []orm.MetricData
}

func (r *SimpleMetricRecorder) RecordQuery(data orm.MetricData) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // 存储 metric
    r.metrics = append(r.metrics, data)
    
    // 可以在这里发送到监控系统
    // sendToPrometheus(data)
    // sendToSigNoz(data)
}
```

### 2. 启用 Metric

```go
metricRecorder := NewSimpleMetricRecorder()

client, err := orm.New(
    mysql.Open(dsn),
    orm.WithMetric(metricRecorder), // 启用 metric
)
```

### 3. 自动收集

一旦启用，所有数据库操作都会自动记录 metric，无需额外代码。

## 集成到 Prometheus

### 1. 创建 Prometheus Metric Recorder

```go
package main

import (
    "github.com/Si40Code/kit/orm"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetricRecorder struct {
    queryDuration *prometheus.HistogramVec
    queryTotal    *prometheus.CounterVec
    queryErrors   *prometheus.CounterVec
}

func NewPrometheusMetricRecorder() *PrometheusMetricRecorder {
    return &PrometheusMetricRecorder{
        queryDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "db_query_duration_seconds",
                Help:    "Database query duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
            []string{"operation", "table"},
        ),
        queryTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "db_query_total",
                Help: "Total number of database queries",
            },
            []string{"operation", "table", "status"},
        ),
        queryErrors: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "db_query_errors_total",
                Help: "Total number of database query errors",
            },
            []string{"operation", "table"},
        ),
    }
}

func (r *PrometheusMetricRecorder) RecordQuery(data orm.MetricData) {
    labels := prometheus.Labels{
        "operation": data.Operation,
        "table":     data.Table,
    }

    // 记录耗时
    r.queryDuration.With(labels).Observe(data.Duration.Seconds())

    // 记录查询总数
    status := "success"
    if data.Error != nil {
        status = "error"
        r.queryErrors.With(labels).Inc()
    }
    
    statusLabels := prometheus.Labels{
        "operation": data.Operation,
        "table":     data.Table,
        "status":    status,
    }
    r.queryTotal.With(statusLabels).Inc()
}
```

### 2. 暴露 Metrics 端点

```go
import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // ... 初始化 ORM client ...

    // 启动 metrics 服务器
    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":2112", nil)
    
    // ... 你的业务逻辑 ...
}
```

### 3. Prometheus 配置

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'my-app'
    static_configs:
      - targets: ['localhost:2112']
```

### 4. 查询示例

```promql
# 平均查询耗时（按操作类型）
rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m])

# QPS（每秒查询数）
rate(db_query_total[1m])

# 错误率
rate(db_query_errors_total[5m]) / rate(db_query_total[5m])

# P95 延迟
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))
```

## 集成到 SigNoz

### 1. 使用 OTLP Metric Exporter

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initMeter() *metric.MeterProvider {
    res, _ := resource.New(
        context.Background(),
        resource.WithAttributes(
            semconv.ServiceName("my-service"),
        ),
    )

    exporter, _ := otlpmetrichttp.New(
        context.Background(),
        otlpmetrichttp.WithEndpoint("signoz:4318"),
        otlpmetrichttp.WithInsecure(),
    )

    mp := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(exporter)),
        metric.WithResource(res),
    )

    otel.SetMeterProvider(mp)
    return mp
}
```

### 2. 创建 OTLP Metric Recorder

```go
type OTLPMetricRecorder struct {
    meter          metric.Meter
    queryDuration  metric.Float64Histogram
    queryCounter   metric.Int64Counter
    errorCounter   metric.Int64Counter
}

func NewOTLPMetricRecorder() *OTLPMetricRecorder {
    meter := otel.Meter("orm-client")
    
    queryDuration, _ := meter.Float64Histogram(
        "db.query.duration",
        metric.WithDescription("Database query duration in milliseconds"),
        metric.WithUnit("ms"),
    )
    
    queryCounter, _ := meter.Int64Counter(
        "db.query.count",
        metric.WithDescription("Total number of database queries"),
    )
    
    errorCounter, _ := meter.Int64Counter(
        "db.query.errors",
        metric.WithDescription("Total number of database query errors"),
    )

    return &OTLPMetricRecorder{
        meter:         meter,
        queryDuration: queryDuration,
        queryCounter:  queryCounter,
        errorCounter:  errorCounter,
    }
}

func (r *OTLPMetricRecorder) RecordQuery(data orm.MetricData) {
    ctx := context.Background()
    
    attrs := metric.WithAttributes(
        attribute.String("db.operation", data.Operation),
        attribute.String("db.table", data.Table),
    )

    r.queryDuration.Record(ctx, float64(data.Duration.Milliseconds()), attrs)
    r.queryCounter.Add(ctx, 1, attrs)
    
    if data.Error != nil {
        r.errorCounter.Add(ctx, 1, attrs)
    }
}
```

## 监控面板示例

### 关键指标

1. **QPS (Queries Per Second)**：数据库查询频率
2. **平均延迟**：查询平均耗时
3. **P95/P99 延迟**：查询延迟分位数
4. **错误率**：查询失败比例
5. **慢查询数**：超过阈值的查询数
6. **热点表**：查询最频繁的表

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "Database Metrics",
    "panels": [
      {
        "title": "QPS",
        "targets": [
          {"expr": "rate(db_query_total[1m])"}
        ]
      },
      {
        "title": "Average Latency",
        "targets": [
          {"expr": "rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m])"}
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {"expr": "rate(db_query_errors_total[5m]) / rate(db_query_total[5m])"}
        ]
      }
    ]
  }
}
```

## 最佳实践

1. **选择合适的 Bucket**：根据查询特点设置合理的延迟分桶
2. **添加业务标签**：可以在 metric 中添加业务相关的标签
3. **设置告警规则**：对慢查询、错误率等设置告警
4. **定期分析**：定期查看 metric 数据，优化性能
5. **关联 Trace**：结合 trace 数据深入分析慢查询

## 常见问题

### Q1: Metric 会影响性能吗？

影响很小。Metric 收集是异步的，不会阻塞数据库操作。建议：
- 避免在 RecordQuery 中执行耗时操作
- 使用缓冲区批量发送 metric

### Q2: 如何减少 Metric 数据量？

```go
// 只记录慢查询
func (r *MyRecorder) RecordQuery(data orm.MetricData) {
    if data.Duration > 100*time.Millisecond {
        // 只记录慢查询
        r.recordSlow(data)
    }
}
```

### Q3: 如何添加自定义标签？

你可以扩展 MetricData 或在 RecordQuery 中添加额外的标签信息。

## 下一步

查看其他示例：
- [04_production](../04_production/) - 生产环境完整配置（包含 log + trace + metric）

