# 示例 3: 文件监控（热更新）

展示如何监控配置文件变更，实现配置热更新，无需重启应用。

## 运行示例

```bash
cd config/examples/03_file_watch
go run main.go
```

## 学习内容

1. **启用文件监控** - `WithFileWatcher()` 选项
2. **注册变更回调** - `OnChange()` 方法
3. **多组件监听** - 多个组件同时监听配置变更
4. **实时生效** - 配置变更后立即生效

## 测试步骤

1. 运行程序：`go run main.go`
2. 修改 `config.yaml` 文件（例如：将 `server.port` 从 8080 改为 9090）
3. 保存文件
4. 观察控制台输出的配置变更通知
5. 按 `Ctrl+C` 退出

## 实际应用场景

### 场景 1: 动态调整日志级别
```go
config.OnChange(func() {
    newLevel := config.GetString("log.level")
    logger.SetLevel(newLevel) // 无需重启即可调整日志级别
})
```

### 场景 2: 动态调整连接池大小
```go
config.OnChange(func() {
    maxConns := config.GetInt("database.max_connections")
    db.SetMaxOpenConns(maxConns) // 动态调整连接池
})
```

### 场景 3: 动态开关功能
```go
config.OnChange(func() {
    enabled := config.GetBool("feature.new_feature_enabled")
    featureFlag.Set("new_feature", enabled) // 动态开关功能
})
```

## 注意事项

- ⚠️ 并非所有配置都适合热更新（如服务器端口，需要重启）
- ✅ 适合热更新的配置：日志级别、超时时间、功能开关、限流参数等
- ✅ 在回调函数中进行必要的验证和错误处理

