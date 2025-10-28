# Logger 快速开始

5 分钟学会使用 kit/logger 模块。

## 1. 安装

```bash
go get github.com/Si40Code/kit/logger
```

## 2. 最简单的例子

```go
package main

import (
    "context"
    "github.com/Si40Code/kit/logger"
)

func main() {
    ctx := context.Background()
    
    // 使用默认 logger（开箱即用）
    logger.Info(ctx, "Hello, Logger!")
}
```

运行：
```bash
go run main.go
```

输出：
```
2024-01-15T10:30:00.123+0800    INFO    main.go:11    Hello, Logger!
```

## 3. 自定义配置

```go
package main

import (
    "context"
    "github.com/Si40Code/kit/logger"
)

func main() {
    // 初始化配置
    logger.Init(
        logger.WithLevel(logger.InfoLevel),
        logger.WithFormat(logger.JSONFormat),
        logger.WithStdout(),
    )
    defer logger.Sync()
    
    ctx := context.Background()
    
    // 记录结构化日志
    logger.Info(ctx, "应用启动",
        "version", "1.0.0",
        "env", "production",
    )
}
```

输出：
```json
{"level":"info","timestamp":"2024-01-15T10:30:00.123Z","msg":"应用启动","version":"1.0.0","env":"production"}
```

## 4. 文件输出

```go
logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithFile("/var/log/app.log",
        logger.WithFileMaxSize(100),    // 100MB 切割
        logger.WithFileMaxAge(7),       // 保留 7 天
        logger.WithFileMaxBackups(3),   // 保留 3 个备份
    ),
)
```

## 5. 多种日志级别

```go
ctx := context.Background()

logger.Debug(ctx, "调试信息", "detail", "...")
logger.Info(ctx, "一般信息", "user_id", 12345)
logger.Warn(ctx, "警告信息", "disk_usage", "90%")
logger.Error(ctx, "错误信息", "error", err)
logger.Fatal(ctx, "致命错误", "error", err)  // 会终止程序
```

## 6. Map 字段方式

```go
logger.InfoMap(ctx, "用户登录", map[string]any{
    "user_id": 12345,
    "username": "alice",
    "ip": "192.168.1.1",
    "timestamp": time.Now().Unix(),
})
```

## 7. 子 Logger

```go
// 创建带预设字段的子 logger
userLogger := logger.With("user_id", 12345, "session", "abc123")

// 后续日志自动包含这些字段
userLogger.Info(ctx, "操作", "action", "login")
userLogger.Info(ctx, "操作", "action", "logout")
```

## 8. Trace 集成

```go
// 初始化（启用 trace）
logger.Init(
    logger.WithTrace("my-service"),
    logger.WithStdout(),
)

// 使用
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

logger.Info(ctx, "操作开始")        // 自动包含 trace_id
logger.Error(ctx, "操作失败")       // 自动标记 span 为 error
```

## 9. 远程日志（SigNoz）

```go
logger.Init(
    logger.WithOTLP("signoz-host:4317"),
    logger.WithTrace("my-service"),
)
```

## 10. 生产环境配置

```go
logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    
    // 输出到文件
    logger.WithFile("/var/log/app.log",
        logger.WithFileMaxSize(100),
        logger.WithFileMaxAge(30),
        logger.WithFileMaxBackups(10),
        logger.WithFileCompress(),
    ),
    
    // 同时输出到 stdout（容器环境）
    logger.WithStdout(),
    
    // 发送到 SigNoz
    logger.WithOTLP("signoz:4317"),
    
    // 启用 trace
    logger.WithTrace("my-service"),
)
```

## 常用场景

### Web 应用

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    logger.Info(ctx, "收到请求",
        "method", r.Method,
        "path", r.URL.Path,
        "ip", r.RemoteAddr,
    )
    
    // 处理请求...
    
    logger.Info(ctx, "请求完成", "status", 200)
}
```

### 错误处理

```go
if err := doSomething(); err != nil {
    logger.Error(ctx, "操作失败",
        "operation", "do_something",
        "error", err.Error(),
        "user_id", userID,
    )
    return err
}
```

### 性能分析

```go
start := time.Now()

// 执行操作...

logger.Info(ctx, "操作完成",
    "operation", "process_batch",
    "duration_ms", time.Since(start).Milliseconds(),
    "count", 1000,
)
```

## 下一步

- [完整文档](./README.md) - 详细的 API 参考
- [示例代码](./examples/) - 5 个完整示例
- [最佳实践](./examples/05_production/) - 生产环境配置

## 帮助

遇到问题？查看：
- [常见问题](./README.md#常见问题)
- [GitHub Issues](https://github.com/Si40Code/kit/issues)

