# ORM 示例

本目录包含 ORM 模块的使用示例，从基础到生产环境的完整配置。

## 示例列表

### [01_basic](01_basic/) - 基础用法 ⭐ 推荐新手

展示 ORM 的基本使用方法：
- 连接数据库
- 基本 CRUD 操作
- 日志集成
- 慢查询检测
- RecordNotFound 错误处理

**适合场景**：快速上手，理解基本概念

### [02_with_trace](02_with_trace/) - Trace 集成

展示如何集成 OpenTelemetry Trace：
- 初始化 TracerProvider
- 自动创建数据库 span
- Span 层级关系
- Span 属性详解

**适合场景**：需要分布式追踪的应用

### [03_with_metric](03_with_metric/) - Metric 集成

展示如何收集性能指标：
- 实现 MetricRecorder 接口
- 收集查询性能数据
- 统计分析
- 集成 Prometheus/SigNoz

**适合场景**：需要性能监控的应用

### [04_production](04_production/) - 生产环境配置 ⭐ 推荐生产

展示完整的生产环境配置：
- Log + Trace + Metric 三件套
- 连接池配置
- 健康检查
- 错误处理和重试
- 事务支持
- 优雅关闭

**适合场景**：生产环境部署

## 快速开始

### 1. 准备数据库

所有示例默认使用 MySQL，请确保已安装 MySQL 并创建测试数据库：

```sql
CREATE DATABASE testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 修改连接信息

编辑示例代码，修改 DSN：

```go
dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
//       ^^^^  ^^^^^^^^      ^^^^^^^^^^^      ^^^^^^
//       用户  密码          主机:端口         数据库名
```

### 3. 运行示例

```bash
cd 01_basic
go run main.go
```

## 学习路径

```
01_basic          → 学习基础用法和日志
    ↓
02_with_trace     → 添加 trace 支持
    ↓
03_with_metric    → 添加 metric 支持
    ↓
04_production     → 生产环境完整配置
```

## 示例对比

| 特性 | 01_basic | 02_with_trace | 03_with_metric | 04_production |
|------|----------|---------------|----------------|---------------|
| 日志 | ✅ | ✅ | ✅ | ✅ |
| Trace | ❌ | ✅ | ❌ | ✅ |
| Metric | ❌ | ❌ | ✅ | ✅ |
| 连接池配置 | ❌ | ❌ | ❌ | ✅ |
| 健康检查 | ❌ | ❌ | ❌ | ✅ |
| 重试机制 | ❌ | ❌ | ❌ | ✅ |
| 事务 | ❌ | ❌ | ❌ | ✅ |

## 常见问题

### Q1: 需要安装什么依赖？

所有示例都自动包含必要的依赖。运行 `go run main.go` 会自动下载。

### Q2: 支持哪些数据库？

支持所有 GORM 支持的数据库：
- MySQL
- PostgreSQL
- SQLite
- SQL Server
- ClickHouse
- TiDB

只需更换驱动即可：

```go
import "gorm.io/driver/postgres"

client, err := orm.New(
    postgres.Open(dsn),
    ...
)
```

### Q3: 如何在生产环境使用？

参考 [04_production](04_production/) 示例，包含所有生产环境最佳实践。

### Q4: 如何迁移现有 GORM 代码？

ORM 客户端完全兼容 GORM。只需：

```go
// 原来的代码
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
db.First(&user, 1)

// 迁移到 kit/orm
client, err := orm.New(mysql.Open(dsn), orm.WithLogger(...))
client.First(&user, 1)
```

### Q5: 性能开销有多大？

- **日志**: 异步记录，几乎无影响
- **Trace**: 使用采样器，可控制开销（建议生产环境 10% 采样）
- **Metric**: 异步收集，可忽略

## 更多资源

- [ORM 模块 README](../README.md) - 完整文档
- [GORM 官方文档](https://gorm.io/) - GORM 使用指南
- [Kit 主仓库](https://github.com/Si40Code/kit) - 查看其他模块

