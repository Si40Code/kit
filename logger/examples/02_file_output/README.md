# 示例 2: 文件输出和日志切割

演示如何将日志输出到文件，以及配置日志切割和清理规则。

## 功能展示

1. **基本文件输出** - 将日志写入文件
2. **文件切割配置** - 基于大小、时间和数量的切割
3. **多输出** - 同时输出到控制台和文件
4. **级别分离** - 不同级别日志输出到不同文件

## 运行示例

```bash
cd logger/examples/02_file_output
go run main.go
```

运行后会在当前目录创建 `./logs` 文件夹，包含以下文件：
- `app.log` - 基本日志文件
- `app-rotated.log` - 配置了切割的日志文件
- `both.log` - 多输出示例的日志
- `info.log` - Info 级别日志
- `error.log` - Error 级别日志

## 日志切割配置

### 配置项说明

| 配置项 | 说明 | 默认值 |
|-------|------|-------|
| `WithFileMaxSize(mb)` | 单个文件最大大小（MB） | 100MB |
| `WithFileMaxAge(days)` | 文件最大保留天数 | 7 天 |
| `WithFileMaxBackups(count)` | 最大备份文件数量 | 3 个 |
| `WithFileCompress()` | 是否压缩旧文件 | false |

### 示例代码

```go
l, err := logger.New(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithFile("./logs/app.log",
        logger.WithFileMaxSize(100),      // 100MB 切割
        logger.WithFileMaxAge(7),         // 保留 7 天
        logger.WithFileMaxBackups(3),     // 保留 3 个备份
        logger.WithFileCompress(),        // 压缩旧文件
    ),
)
```

## 切割规则

日志文件会在以下情况下自动切割：

1. **文件大小超限** - 当文件大小超过 `MaxSize` 时切割
2. **时间超限** - 当文件存在时间超过 `MaxAge` 天时删除
3. **数量超限** - 当备份文件数量超过 `MaxBackups` 时删除最旧的

### 切割后的文件名格式

```
app.log                    # 当前日志文件
app-2024-01-15T10-30-00.log  # 备份文件（带时间戳）
app-2024-01-15T10-30-00.log.gz # 压缩的备份文件
```

## 多输出配置

同时输出到多个目标：

```go
l, err := logger.New(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithStdout(),                // 输出到控制台
    logger.WithFile("./logs/app.log"),  // 输出到文件
    logger.WithOTLP("localhost:4317"),  // 输出到远程
)
```

## 级别分离

为不同级别的日志创建不同的 logger：

```go
// Info 日志
infoLogger, _ := logger.New(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFile("./logs/info.log"),
)

// Error 日志
errorLogger, _ := logger.New(
    logger.WithLevel(logger.ErrorLevel),
    logger.WithFile("./logs/error.log"),
)
```

## 生产环境建议

### 1. 文件路径

```go
// 使用标准路径
logger.WithFile("/var/log/myapp/app.log")

// 或使用环境变量
logPath := os.Getenv("LOG_PATH")
logger.WithFile(logPath)
```

### 2. 切割配置

```go
// 生产环境推荐配置
logger.WithFile("/var/log/myapp/app.log",
    logger.WithFileMaxSize(100),      // 100MB
    logger.WithFileMaxAge(30),        // 保留 30 天
    logger.WithFileMaxBackups(10),    // 保留 10 个备份
    logger.WithFileCompress(),        // 压缩节省空间
)
```

### 3. 权限管理

```bash
# 创建日志目录并设置权限
sudo mkdir -p /var/log/myapp
sudo chown myapp:myapp /var/log/myapp
sudo chmod 755 /var/log/myapp
```

### 4. 日志清理脚本

虽然 lumberjack 会自动清理，但也可以配置额外的清理脚本：

```bash
#!/bin/bash
# cleanup-logs.sh

# 删除 30 天前的压缩日志
find /var/log/myapp -name "*.gz" -mtime +30 -delete

# 删除 90 天前的所有日志
find /var/log/myapp -name "*.log*" -mtime +90 -delete
```

## 注意事项

1. **文件权限** - 确保应用有权限写入日志目录
2. **磁盘空间** - 定期监控磁盘空间，防止日志占满磁盘
3. **性能影响** - 文件 I/O 会影响性能，考虑使用异步写入
4. **日志格式** - 生产环境推荐使用 JSON 格式便于解析

## 下一步

- [示例 3: Trace 集成](../03_with_trace/) - 学习如何与 OpenTelemetry 集成
- [示例 5: 生产环境配置](../05_production/) - 完整的生产环境配置

