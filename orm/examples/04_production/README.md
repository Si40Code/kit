# 示例 04: 生产环境配置

本示例展示完整的生产环境配置，包含所有最佳实践。

## 功能展示

1. **完整的可观测性**：Log + Trace + Metric 三件套
2. **生产级配置**：连接池、超时、重试等
3. **健康检查**：数据库连接状态监控
4. **优雅关闭**：正确清理资源
5. **错误处理**：重试机制和错误恢复
6. **事务支持**：ACID 保证

## 配置清单

### Logger 配置

```go
logger.Init(
    logger.WithLevel(logger.InfoLevel),      // Info 级别
    logger.WithFormat(logger.JSONFormat),     // JSON 格式
    logger.WithStdout(),                      // 标准输出
    logger.WithFile(                          // 文件输出
        "/var/log/app/app.log",
        logger.WithFileMaxSize(100),          // 100MB 轮转
        logger.WithFileMaxAge(7),             // 保留 7 天
        logger.WithFileMaxBackups(3),         // 3 个备份
    ),
    logger.WithTrace("my-service"),           // Trace 集成
    logger.WithCaller(true),                  // 记录调用者
)
```

### ORM 配置

```go
client, err := orm.New(
    mysql.Open(dsn),
    // 日志配置
    orm.WithLogger(logger.L()),
    orm.WithSlowThreshold(200*time.Millisecond),
    orm.WithIgnoreRecordNotFoundError(),
    
    // Trace 配置
    orm.WithTrace("orm-client"),
    
    // Metric 配置
    orm.WithMetric(metricRecorder),
    
    // 连接池配置（关键！）
    orm.WithMaxIdleConns(10),
    orm.WithMaxOpenConns(100),
    orm.WithConnMaxLifetime(time.Hour),
    orm.WithConnMaxIdleTime(10*time.Minute),
)
```

### Trace 配置

```go
import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

exporter, err := otlptracehttp.New(
    context.Background(),
    otlptracehttp.WithEndpoint("signoz:4318"),
    otlptracehttp.WithInsecure(),
)

tp := trace.NewTracerProvider(
    trace.WithBatcher(exporter),
    trace.WithSampler(trace.TraceIDRatioBased(0.1)), // 10% 采样
)

otel.SetTracerProvider(tp)
```

## 连接池配置详解

### 参数说明

| 参数 | 默认值 | 推荐值 | 说明 |
|------|--------|--------|------|
| `MaxIdleConns` | 10 | 10-20 | 最大空闲连接数 |
| `MaxOpenConns` | 100 | 50-200 | 最大打开连接数 |
| `ConnMaxLifetime` | 1h | 1h | 连接最大生命周期 |
| `ConnMaxIdleTime` | 10m | 5-10m | 连接最大空闲时间 |

### 如何确定合适的值？

#### 1. MaxOpenConns

```
MaxOpenConns = (并发请求数 × 平均查询时间) / 目标响应时间
```

示例：
- 并发请求：1000 QPS
- 平均查询时间：10ms
- 目标响应时间：50ms

```
MaxOpenConns = (1000 × 0.01) / 0.05 = 200
```

#### 2. MaxIdleConns

通常设置为 `MaxOpenConns` 的 10-20%：

```
MaxIdleConns = MaxOpenConns × 0.15 = 200 × 0.15 = 30
```

#### 3. 监控和调整

通过健康检查观察：

```go
stats := sqlDB.Stats()
fmt.Printf("Open: %d, InUse: %d, Idle: %d, WaitCount: %d\n",
    stats.OpenConnections,
    stats.InUse,
    stats.Idle,
    stats.WaitCount,
)
```

调整策略：
- `WaitCount` 持续增长 → 增加 `MaxOpenConns`
- `Idle` 总是很高 → 减少 `MaxIdleConns`
- `WaitDuration` 过长 → 增加 `MaxOpenConns`

## 错误处理最佳实践

### 1. 重试机制

```go
maxRetries := 3
for i := 0; i < maxRetries; i++ {
    if err := client.WithContext(ctx).Create(&user).Error; err != nil {
        logger.Warn(ctx, "Operation failed, retrying",
            "attempt", i+1,
            "error", err,
        )
        time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
        continue
    }
    break
}
```

### 2. 事务处理

```go
err := client.WithContext(ctx).Transaction(func(tx *orm.Client) error {
    // 操作 1
    if err := tx.Create(&user).Error; err != nil {
        return err // 自动回滚
    }
    
    // 操作 2
    if err := tx.Create(&profile).Error; err != nil {
        return err // 自动回滚
    }
    
    return nil // 自动提交
})

if err != nil {
    logger.Error(ctx, "Transaction failed", "error", err)
}
```

### 3. 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 查询会在 5 秒后超时
err := client.WithContext(ctx).Find(&users).Error
```

## 健康检查

### 1. HTTP 端点

```go
import "net/http"

http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    if err := checkDatabaseHealth(client, r.Context()); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "unhealthy",
            "error":  err.Error(),
        })
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
})
```

### 2. Kubernetes Liveness/Readiness Probe

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    livenessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 30
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
```

## 监控和告警

### 关键指标

1. **连接池指标**
   - `db_pool_open_connections`
   - `db_pool_in_use`
   - `db_pool_idle`
   - `db_pool_wait_count`

2. **查询性能指标**
   - `db_query_duration_seconds`
   - `db_slow_query_count`
   - `db_query_errors_total`

3. **业务指标**
   - `db_query_total` (按操作类型)
   - `db_rows_affected_total`

### 告警规则示例

```yaml
# Prometheus Alert Rules
groups:
  - name: database
    rules:
      # 慢查询告警
      - alert: HighSlowQueryRate
        expr: rate(db_slow_query_count[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High slow query rate detected"
      
      # 错误率告警
      - alert: HighDatabaseErrorRate
        expr: rate(db_query_errors_total[5m]) / rate(db_query_total[5m]) > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High database error rate (> 1%)"
      
      # 连接池耗尽告警
      - alert: DatabaseConnectionPoolExhausted
        expr: db_pool_wait_count > 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Database connection pool wait count is high"
```

## 性能优化建议

### 1. 索引优化

```go
type User struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"size:100;index"`              // 单列索引
    Email string `gorm:"size:100;uniqueIndex"`        // 唯一索引
    Age   int    `gorm:"index:idx_age_name,priority:1"` // 复合索引
}
```

### 2. 查询优化

```go
// ✅ 好：只查询需要的字段
client.Select("id", "name").Find(&users)

// ❌ 坏：查询所有字段
client.Find(&users)

// ✅ 好：使用预加载
client.Preload("Profile").Find(&users)

// ❌ 坏：N+1 查询
for _, user := range users {
    client.First(&user.Profile, "user_id = ?", user.ID)
}
```

### 3. 批量操作

```go
// ✅ 好：批量插入
client.CreateInBatches(users, 100)

// ❌ 坏：逐条插入
for _, user := range users {
    client.Create(&user)
}
```

### 4. 连接复用

```go
// ✅ 好：复用连接
db := client.WithContext(ctx)
db.First(&user1)
db.First(&user2)

// ❌ 坏：每次创建新连接
client.WithContext(ctx).First(&user1)
client.WithContext(ctx).First(&user2)
```

## 部署检查清单

- [ ] 配置从环境变量或配置文件读取
- [ ] 日志输出到文件和 stdout
- [ ] 启用 JSON 格式日志
- [ ] 配置合适的连接池参数
- [ ] 启用 Trace 并连接到后端
- [ ] 启用 Metric 并暴露端点
- [ ] 实现健康检查端点
- [ ] 配置慢查询阈值
- [ ] 设置监控告警规则
- [ ] 测试故障恢复机制
- [ ] 准备数据库迁移脚本
- [ ] 配置数据库备份策略

## 运行示例

```bash
cd examples/04_production
go run main.go
```

## 环境变量配置示例

```bash
# .env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=testdb

DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=3600
DB_SLOW_THRESHOLD=200

OTEL_EXPORTER_OTLP_ENDPOINT=signoz:4318
OTEL_SERVICE_NAME=my-service
OTEL_TRACE_SAMPLER=traceidratio
OTEL_TRACE_SAMPLER_ARG=0.1

LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=/var/log/app/app.log
```

## 生产环境清单

### 数据库配置

- ✅ 使用连接池
- ✅ 设置超时
- ✅ 启用慢查询日志
- ✅ 配置读写分离（如需要）
- ✅ 使用事务
- ✅ 定期备份

### 可观测性

- ✅ 结构化日志
- ✅ 分布式追踪
- ✅ 性能指标
- ✅ 健康检查
- ✅ 告警规则

### 安全性

- ✅ 参数化查询（防 SQL 注入）
- ✅ 最小权限原则
- ✅ 加密敏感信息
- ✅ 审计日志

### 高可用

- ✅ 自动重试
- ✅ 故障转移
- ✅ 优雅降级
- ✅ 熔断器

## 相关资源

- [GORM 官方文档](https://gorm.io/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
- [SigNoz](https://signoz.io/)
- [Prometheus](https://prometheus.io/)

