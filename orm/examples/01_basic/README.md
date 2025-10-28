# 示例 01: 基础用法

本示例展示了 ORM 客户端的基本使用方法，包括日志集成。

## 功能展示

1. **连接数据库**：使用 MySQL 驱动连接数据库
2. **日志集成**：自动记录所有 SQL 查询日志
3. **基本 CRUD 操作**：
   - Create（创建）
   - Read（查询）
   - Update（更新）
   - Delete（删除）
4. **慢查询检测**：自动识别并警告慢查询
5. **错误处理**：展示 RecordNotFound 错误处理

## 运行前准备

### 1. 安装 MySQL

确保你已经安装了 MySQL 数据库。

### 2. 创建数据库

```sql
CREATE DATABASE testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 修改连接信息

编辑 `main.go`，修改数据库连接字符串：

```go
dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
//       ^^^^  ^^^^^^^^      ^^^^^^^^^^^      ^^^^^^
//       用户  密码          主机:端口         数据库名
```

## 运行示例

```bash
cd examples/01_basic
go run main.go
```

## 预期输出

你将看到类似以下的日志输出：

```
2024-01-01T10:00:00.000+0800    INFO    Successfully connected to database
2024-01-01T10:00:00.100+0800    INFO    database query executed    {"duration_ms": 5, "rows_affected": 0, "sql": "CREATE TABLE `users` ..."}
2024-01-01T10:00:00.200+0800    INFO    Creating user...
2024-01-01T10:00:00.210+0800    INFO    database query executed    {"duration_ms": 8, "rows_affected": 1, "sql": "INSERT INTO `users` ..."}
2024-01-01T10:00:00.220+0800    INFO    User created successfully    {"id": 1}
2024-01-01T10:00:00.230+0800    INFO    Querying user by ID...
2024-01-01T10:00:00.235+0800    INFO    database query executed    {"duration_ms": 3, "rows_affected": 1, "sql": "SELECT * FROM `users` WHERE `id` = 1"}
2024-01-01T10:00:00.240+0800    INFO    User found    {"name": "Alice", "email": "alice@example.com"}
...
```

## 代码说明

### 1. 初始化 Logger

```go
logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.ConsoleFormat),
    logger.WithStdout(),
)
```

### 2. 创建 ORM 客户端

```go
client, err := orm.New(
    mysql.Open(dsn),
    orm.WithLogger(logger.L()),
    orm.WithSlowThreshold(100*time.Millisecond),
)
```

- `mysql.Open(dsn)`：使用 MySQL 驱动
- `WithLogger()`：集成日志记录
- `WithSlowThreshold()`：设置慢查询阈值为 100ms

### 3. 使用 Context

所有数据库操作都应该使用 `WithContext(ctx)`：

```go
client.WithContext(ctx).Create(&user)
```

这样可以：
- 支持超时控制
- 传递 trace 信息
- 记录上下文日志

### 4. 忽略 RecordNotFound 错误

```go
// 方法 1: 全局配置
client, err := orm.New(
    mysql.Open(dsn),
    orm.WithIgnoreRecordNotFoundError(),
)

// 方法 2: 单次查询
client.WithIgnoreRecordNotFound().First(&user, id)
```

## 常见问题

### Q1: 连接失败怎么办？

确保：
1. MySQL 服务已启动
2. 用户名和密码正确
3. 数据库已创建
4. 防火墙允许连接

### Q2: 如何查看 SQL 语句？

所有 SQL 都会自动记录到日志中，查看 `sql` 字段即可。

### Q3: 如何调整日志级别？

```go
logger.Init(
    logger.WithLevel(logger.DebugLevel), // 更详细的日志
)
```

## 下一步

查看其他示例：
- [02_with_trace](../02_with_trace/) - Trace 集成
- [03_with_metric](../03_with_metric/) - Metric 集成
- [04_production](../04_production/) - 生产环境配置

