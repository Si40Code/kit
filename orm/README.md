# ORM

生产级的 ORM 客户端包，基于 [GORM](https://gorm.io/)，提供完整的 Trace、日志和 Metric 集成。

## 特性

- ✅ **完整的日志记录**: 记录所有 SQL 查询，包括耗时、影响行数
- ✅ **OpenTelemetry Trace 集成**: 每个查询自动创建独立 span，无缝集成到分布式追踪系统
- ✅ **详细的 Metric**: 收集查询类型、表名、耗时、错误等性能数据
- ✅ **慢查询检测**: 自动识别并警告慢查询
- ✅ **灵活的错误处理**: 可配置查询无数据时不返回错误（全局+单次）
- ✅ **连接池管理**: 生产级连接池配置
- ✅ **统一的配置模式**: 遵循 kit 的 Option 配置风格
- ✅ **完全兼容 GORM**: 直接暴露 `*gorm.DB`，支持所有 GORM 功能

## 快速开始

### 安装

```bash
go get github.com/Si40Code/kit/orm
go get gorm.io/driver/mysql  # 或其他数据库驱动
```

### 基本用法

```go
package main

import (
    "context"
    "time"

    "github.com/Si40Code/kit/logger"
    "github.com/Si40Code/kit/orm"
    "gorm.io/driver/mysql"
)

type User struct {
    ID   uint   `gorm:"primarykey"`
    Name string
    Age  int
}

func main() {
    // 初始化 logger
    logger.Init(logger.WithStdout())
    defer logger.Sync()

    // 创建 ORM 客户端
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    client, err := orm.New(
        mysql.Open(dsn),
        orm.WithLogger(logger.L()),
        orm.WithSlowThreshold(100*time.Millisecond),
    )
    if err != nil {
        panic(err)
    }
    defer client.Close()

    ctx := context.Background()

    // 使用 GORM 的所有功能
    var user User
    client.WithContext(ctx).First(&user, 1)
}
```

## 配置选项

### 日志配置

```go
client, err := orm.New(
    mysql.Open(dsn),
    // 设置 logger
    orm.WithLogger(logger.L()),
    
    // 禁用日志
    orm.WithDisableLog(),
    
    // 设置慢查询阈值
    orm.WithSlowThreshold(200*time.Millisecond),
    
    // 全局忽略 RecordNotFound 错误
    orm.WithIgnoreRecordNotFoundError(),
)
```

### Trace 配置

```go
client, err := orm.New(
    mysql.Open(dsn),
    // 启用 trace（需要先初始化 OpenTelemetry）
    orm.WithTrace("my-service"),
)
```

每个数据库查询会自动创建一个 span，包含以下属性：

| 属性 | 说明 | 示例 |
|------|------|------|
| `db.system` | 数据库系统 | `gorm` |
| `db.operation` | 操作类型 | `SELECT`, `INSERT`, `UPDATE`, `DELETE` |
| `db.table` | 表名 | `users` |
| `db.statement` | SQL 语句 | `SELECT * FROM users WHERE id = ?` |
| `db.rows_affected` | 影响行数 | `1` |
| `db.duration_ms` | 执行耗时 | `5` |

### Metric 配置

```go
// 实现 MetricRecorder 接口
type MyMetricRecorder struct{}

func (r *MyMetricRecorder) RecordQuery(data orm.MetricData) {
    // 发送到 Prometheus、SigNoz 等
}

client, err := orm.New(
    mysql.Open(dsn),
    // 启用 metric
    orm.WithMetric(&MyMetricRecorder{}),
)
```

Metric 数据结构：

```go
type MetricData struct {
    Operation    string        // SELECT, INSERT, UPDATE, DELETE
    Table        string        // 表名
    SQL          string        // 完整 SQL 语句
    Duration     time.Duration // 查询耗时
    RowsAffected int64        // 影响/返回的行数
    Error        error        // 错误（如果有）
}
```

### 连接池配置

```go
client, err := orm.New(
    mysql.Open(dsn),
    // 最大空闲连接数
    orm.WithMaxIdleConns(10),
    
    // 最大打开连接数
    orm.WithMaxOpenConns(100),
    
    // 连接最大生命周期
    orm.WithConnMaxLifetime(time.Hour),
    
    // 连接最大空闲时间
    orm.WithConnMaxIdleTime(10*time.Minute),
)
```

## 核心功能

### 1. 日志记录

所有 SQL 查询都会自动记录到日志中：

```json
{
  "level": "info",
  "msg": "database query executed",
  "duration_ms": 5,
  "rows_affected": 1,
  "sql": "SELECT * FROM `users` WHERE `id` = 1"
}
```

慢查询会自动标记为 WARN 级别：

```json
{
  "level": "warn",
  "msg": "slow query detected",
  "duration_ms": 250,
  "slow_threshold_ms": 200,
  "sql": "SELECT * FROM `users` WHERE `age` > 18"
}
```

### 2. Trace 集成

```go
import "go.opentelemetry.io/otel"

// 创建业务 span
tracer := otel.Tracer("business-logic")
ctx, span := tracer.Start(ctx, "CreateUser")
defer span.End()

// 数据库操作会自动成为子 span
client.WithContext(ctx).Create(&user)
```

Span 层级结构：

```
CreateUser (业务 span)
└─ DB INSERT (数据库 span)
   ├─ db.operation: INSERT
   ├─ db.table: users
   ├─ db.statement: INSERT INTO `users` ...
   └─ db.duration_ms: 5
```

### 3. 错误处理

#### 全局配置

```go
client, err := orm.New(
    mysql.Open(dsn),
    orm.WithIgnoreRecordNotFoundError(), // 全局忽略
)

// 查询不存在的记录不会返回错误
var user User
err := client.First(&user, 99999).Error // err == nil
```

#### 单次查询覆盖

```go
// 默认会返回 RecordNotFound 错误
err := client.First(&user, 99999).Error // err != nil

// 使用 WithIgnoreRecordNotFound 单次忽略
err := client.WithIgnoreRecordNotFound().First(&user, 99999).Error // err == nil
```

### 4. Context 传播

**重要**: 始终使用 `WithContext(ctx)` 传递 context：

```go
// ✅ 正确：支持 trace、超时、取消
client.WithContext(ctx).First(&user, id)

// ❌ 错误：不支持 trace 传播
client.First(&user, id)
```

### 5. 事务支持

```go
err := client.WithContext(ctx).Transaction(func(tx *orm.Client) error {
    // 在事务中执行多个操作
    if err := tx.Create(&user).Error; err != nil {
        return err // 自动回滚
    }
    
    if err := tx.Create(&profile).Error; err != nil {
        return err // 自动回滚
    }
    
    return nil // 自动提交
})
```

## 示例

查看 `examples/` 目录了解更多使用示例：

- [01_basic](examples/01_basic/) - 基本用法和日志
- [02_with_trace](examples/02_with_trace/) - Trace 集成
- [03_with_metric](examples/03_with_metric/) - Metric 集成
- [04_production](examples/04_production/) - 生产环境完整配置

## API 文档

### Client

```go
// New 创建一个新的 ORM 客户端
func New(dialector gorm.Dialector, opts ...Option) (*Client, error)

// WithContext 返回一个新的带有指定 context 的客户端实例
func (c *Client) WithContext(ctx context.Context) *Client

// WithIgnoreRecordNotFound 返回一个新的客户端实例，该实例会忽略 RecordNotFound 错误
func (c *Client) WithIgnoreRecordNotFound() *Client

// Close 关闭数据库连接
func (c *Client) Close() error
```

`Client` 直接嵌入了 `*gorm.DB`，因此可以使用所有 GORM 的方法：

```go
client.Create(&user)
client.First(&user, id)
client.Where("age > ?", 18).Find(&users)
client.Model(&user).Update("age", 25)
client.Delete(&user)
// ... 所有 GORM 方法
```

### MetricRecorder 接口

```go
type MetricRecorder interface {
    RecordQuery(data MetricData)
}
```

### Option 函数

| 函数 | 说明 |
|------|------|
| `WithLogger(logger.Logger)` | 设置日志记录器 |
| `WithDisableLog()` | 禁用日志 |
| `WithSlowThreshold(time.Duration)` | 设置慢查询阈值 |
| `WithIgnoreRecordNotFoundError()` | 全局忽略 RecordNotFound 错误 |
| `WithTrace(string)` | 启用 trace（服务名） |
| `WithMetric(MetricRecorder)` | 启用 metric |
| `WithMaxIdleConns(int)` | 设置最大空闲连接数 |
| `WithMaxOpenConns(int)` | 设置最大打开连接数 |
| `WithConnMaxLifetime(time.Duration)` | 设置连接最大生命周期 |
| `WithConnMaxIdleTime(time.Duration)` | 设置连接最大空闲时间 |

## 支持的数据库

通过 GORM 驱动支持多种数据库：

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/driver/sqlserver"
)

// MySQL
orm.New(mysql.Open(dsn), ...)

// PostgreSQL
orm.New(postgres.Open(dsn), ...)

// SQLite
orm.New(sqlite.Open("test.db"), ...)

// SQL Server
orm.New(sqlserver.Open(dsn), ...)
```

## 最佳实践

1. **始终传递 Context**: 使用 `WithContext(ctx)` 确保 trace 信息正确传播
2. **设置合理的超时**: 使用 `context.WithTimeout` 控制查询超时
3. **配置连接池**: 根据负载调整 `MaxIdleConns` 和 `MaxOpenConns`
4. **监控慢查询**: 设置合理的 `SlowThreshold`，及时发现性能问题
5. **使用事务**: 对于需要原子性的操作使用 `Transaction`
6. **健康检查**: 定期检查数据库连接状态
7. **集成 Trace**: 在生产环境中集成到 SigNoz/Jaeger
8. **索引优化**: 为常用查询字段添加索引

## 性能考虑

- **连接复用**: 默认启用连接池，自动复用连接
- **日志优化**: 日志记录是异步的，不会阻塞查询
- **Trace 开销**: 使用采样器减少 trace 开销
- **Metric 开销**: Metric 收集是异步的，影响可忽略
- **并发安全**: 客户端是并发安全的，可以在多个 goroutine 中共享

## 与 GORM 的关系

本包基于 [GORM](https://gorm.io/)，并直接暴露 `*gorm.DB`。这意味着：

- 可以使用所有 GORM 的原生功能
- 无需学习新的 API
- 通过插件机制自动增强功能（trace、log、metric）
- 可以无缝迁移现有 GORM 代码

```go
// 可以直接使用 GORM 的所有方法
client.Model(&User{}).
    Where("age > ?", 18).
    Order("created_at DESC").
    Limit(10).
    Find(&users)
```

## 与其他 kit 模块集成

```go
import (
    "github.com/Si40Code/kit/config"
    "github.com/Si40Code/kit/logger"
    "github.com/Si40Code/kit/orm"
)

func main() {
    // 1. 初始化配置
    config.Init(config.WithFile("config.yaml"))
    
    // 2. 基于配置初始化日志
    logger.Init(
        logger.WithLevel(logger.ParseLevel(config.GetString("log.level"))),
        logger.WithTrace(config.GetString("service.name")),
    )
    
    // 3. 基于配置创建 ORM 客户端
    client, err := orm.New(
        mysql.Open(config.GetString("database.dsn")),
        orm.WithLogger(logger.L()),
        orm.WithTrace(config.GetString("service.name")),
        orm.WithSlowThreshold(config.GetDuration("database.slow_threshold")),
        orm.WithMaxOpenConns(config.GetInt("database.max_open_conns")),
    )
}
```

## License

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！

## 相关资源

- [GORM 官方文档](https://gorm.io/)
- [Kit 主仓库](https://github.com/Si40Code/kit)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)

