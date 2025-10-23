# 示例 5: 配置变更通知和日志

展示配置变更的日志记录和通知机制，包括敏感信息脱敏。

## 运行示例

```bash
cd config/examples/05_change_notification
go run main.go
```

## 学习内容

1. **配置变更日志** - 自动记录所有配置变更
2. **敏感信息脱敏** - 自动隐藏密码、token 等
3. **多组件通知** - 多个组件同时响应配置变更
4. **变更类型** - ADD、UPDATE、DELETE

## 配置变更日志格式

```json
{
  "type": "config_change",
  "source": "file",
  "key": "server.port",
  "old": "8080",
  "new": "9090",
  "change": "UPDATE",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 字段说明

- `type`: 日志类型，固定为 "config_change"
- `source`: 配置源，可能是 "file"、"apollo"、"nacos" 等
- `key`: 配置键名
- `old`: 旧值
- `new`: 新值
- `change`: 变更类型（ADD、UPDATE、DELETE）
- `timestamp`: 变更时间

## 敏感信息脱敏

包含以下关键词的配置会自动脱敏：

- `password`
- `secret`
- `token`
- `key`

### 示例

```json
{
  "key": "database.password",
  "old": "******",
  "new": "******",
  "change": "UPDATE"
}
```

## 实际应用场景

### 场景 1: 审计日志

将配置变更日志输出到审计系统，满足合规要求：

```go
config.OnChange(func() {
    // 获取变更日志，发送到审计系统
    auditLog.RecordConfigChange(...)
})
```

### 场景 2: 告警通知

关键配置变更时发送告警：

```go
config.OnChange(func() {
    newValue := config.GetString("critical.config")
    if isImportantChange(newValue) {
        alerting.SendAlert("配置已变更", newValue)
    }
})
```

### 场景 3: 配置回滚

记录配置历史，支持快速回滚：

```go
var configHistory []ConfigSnapshot

config.OnChange(func() {
    snapshot := captureCurrentConfig()
    configHistory = append(configHistory, snapshot)
})
```

## 多组件监听示例

```go
// 数据库连接池
config.OnChange(func() {
    host := config.GetString("database.host")
    port := config.GetInt("database.port")
    db.Reconnect(host, port)
})

// HTTP 服务器
config.OnChange(func() {
    port := config.GetInt("server.port")
    // 注意：端口变更通常需要重启服务器
    log.Printf("Server port changed to %d, restart required", port)
})

// 缓存管理器
config.OnChange(func() {
    ttl := config.GetInt("cache.ttl")
    cache.UpdateTTL(ttl)
})

// 日志系统
config.OnChange(func() {
    level := config.GetString("log.level")
    logger.SetLevel(level)
})
```

## 最佳实践

- ✅ 将配置变更日志输出到日志系统或文件
- ✅ 监控敏感配置的变更
- ✅ 在回调函数中进行必要的错误处理
- ✅ 避免在回调函数中执行耗时操作（使用 goroutine）
- ⚠️ 注意回调函数的执行顺序是不确定的

## 日志集成

可以将配置变更日志集成到日志系统：

```go
import "github.com/silin/go-pkg-sdk/logger"

config.OnChange(func() {
    logger.Info("Config changed",
        logger.String("source", "file"),
        logger.Any("changes", getChanges()),
    )
})
```

