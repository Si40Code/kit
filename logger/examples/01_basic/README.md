# 示例 1: 基础用法

演示 logger 模块的基本使用方法。

## 功能展示

1. **使用默认 logger** - 开箱即用
2. **初始化自定义 logger** - 配置日志级别和格式
3. **不同日志级别** - Debug, Info, Warn, Error, Fatal
4. **结构化字段** - key-value 对方式
5. **Map 字段** - 使用 map 传递字段
6. **子 logger** - 使用 With 创建带预设字段的 logger
7. **独立实例** - 创建多个 logger 实例
8. **性能测试** - 测试日志记录性能

## 运行示例

```bash
cd logger/examples/01_basic
go run main.go
```

## 输出示例

```
2024-01-15T10:30:00.123+0800    INFO    使用默认 logger
2024-01-15T10:30:00.124+0800    DEBUG   这是 debug 日志
2024-01-15T10:30:00.125+0800    INFO    这是 info 日志
2024-01-15T10:30:00.126+0800    WARN    这是 warn 日志
2024-01-15T10:30:00.127+0800    ERROR   这是 error 日志
2024-01-15T10:30:00.128+0800    INFO    用户登录    {"user_id": 12345, "username": "alice", "ip": "192.168.1.1"}
2024-01-15T10:30:00.129+0800    INFO    订单创建    {"order_id": "ORD-2024-001", "amount": 99.99, "currency": "USD", "items": 3}
```

## 关键代码

### 1. 使用默认 logger

```go
ctx := context.Background()
logger.Info(ctx, "使用默认 logger")
```

### 2. 初始化自定义配置

```go
err := logger.Init(
    logger.WithLevel(logger.DebugLevel),
    logger.WithFormat(logger.ConsoleFormat),
    logger.WithStdout(),
)
```

### 3. 结构化字段

```go
logger.Info(ctx, "用户登录",
    "user_id", 12345,
    "username", "alice",
    "ip", "192.168.1.1",
)
```

### 4. Map 字段

```go
logger.InfoMap(ctx, "订单创建", map[string]any{
    "order_id": "ORD-2024-001",
    "amount":   99.99,
    "currency": "USD",
})
```

### 5. 子 logger

```go
userLogger := logger.With("user_id", 12345, "session", "abc123")
userLogger.Info(ctx, "用户执行操作", "action", "update_profile")
```

## 最佳实践

1. **使用 context** - 始终传递 context，方便 trace 集成
2. **结构化日志** - 使用 key-value 字段而不是字符串拼接
3. **合理的日志级别** - 根据重要性选择合适的级别
4. **刷新缓冲区** - 程序退出前调用 `logger.Sync()`

## 下一步

- [示例 2: 文件输出](../02_file_output/) - 学习如何将日志输出到文件
- [示例 3: Trace 集成](../03_with_trace/) - 学习如何与 OpenTelemetry 集成

