# 示例 8: 使用 Kit Logger

演示如何在 web 服务中集成 kit logger 模块。

## 功能展示

1. **Kit Logger 集成** - 使用 kit/logger 替代默认 logger
2. **适配器模式** - 无缝集成到 web 模块
3. **统一日志管理** - 应用和 web 框架使用同一个 logger

## 运行示例

```bash
cd web/examples/08_with_kit_logger
go run main.go
```

然后访问：
- http://localhost:8080/api/hello - 正常请求
- http://localhost:8080/api/error - 错误请求

## 关键代码

### 1. 初始化 Kit Logger

```go
err := logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithStdout(),
)
```

### 2. 创建适配器并注入到 Web 模块

```go
server := web.New(
    web.WithLogger(web.NewLoggerAdapter(logger.Default())),
)
```

### 3. 在处理器中使用 Logger

```go
engine.GET("/api/hello", func(c *gin.Context) {
    ctx := c.Request.Context()
    
    // 直接使用 kit logger
    logger.Info(ctx, "处理请求",
        "path", "/api/hello",
        "method", "GET",
    )
    
    web.Success(c, gin.H{"message": "Hello!"})
})
```

## 优势

### 1. 统一的日志管理

所有日志通过 kit logger 统一管理，便于：
- 统一配置（级别、格式、输出）
- 集中收集（OTLP、文件）
- Trace 集成

### 2. 灵活的配置

```go
// 开发环境
logger.Init(
    logger.WithLevel(logger.DebugLevel),
    logger.WithFormat(logger.ConsoleFormat),
    logger.WithStdout(),
)

// 生产环境
logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithFile("/var/log/app.log"),
    logger.WithOTLP("signoz:4317"),
)
```

### 3. 与 Trace 集成

```go
// 启用 trace
logger.Init(
    logger.WithTrace("my-service"),
    logger.WithOTLP("signoz:4317"),
)

// Web 中间件和应用代码的日志都会包含 trace 信息
```

## 完整生产环境示例

```go
package main

import (
    "github.com/Si40Code/kit/logger"
    "github.com/Si40Code/kit/web"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化 logger
    logger.Init(
        logger.WithLevel(logger.InfoLevel),
        logger.WithFormat(logger.JSONFormat),
        logger.WithStdout(),
        logger.WithFile("/var/log/app.log",
            logger.WithFileMaxSize(100),
            logger.WithFileMaxAge(30),
        ),
        logger.WithOTLP("signoz:4317"),
        logger.WithTrace("my-service"),
    )
    defer logger.Sync()
    
    // 创建 web server
    server := web.New(
        web.WithMode(web.ReleaseMode),
        web.WithServiceName("my-service"),
        web.WithLogger(web.NewLoggerAdapter(logger.Default())),
        web.WithTrace(),  // 启用 web trace
    )
    
    // 注册路由
    engine := server.Engine()
    engine.GET("/api/users", handleUsers)
    
    // 启动服务器
    server.RunWithGracefulShutdown(":8080")
}

func handleUsers(c *gin.Context) {
    ctx := c.Request.Context()
    
    // 应用日志会自动包含 trace_id
    logger.Info(ctx, "查询用户列表")
    
    // 业务逻辑...
    users := getUsers()
    
    logger.Info(ctx, "返回用户列表", "count", len(users))
    web.Success(c, users)
}
```

## 对比

### 使用默认 logger

```go
server := web.New(
    web.WithMode(web.ReleaseMode),
    // 使用默认的简单 logger
)
```

- ✅ 简单快速
- ❌ 功能有限
- ❌ 无法统一管理
- ❌ 无 trace 集成

### 使用 kit logger

```go
logger.Init(logger.WithOTLP("signoz:4317"))

server := web.New(
    web.WithLogger(web.NewLoggerAdapter(logger.Default())),
)
```

- ✅ 功能完整
- ✅ 统一管理
- ✅ Trace 集成
- ✅ 多种输出

## 下一步

- [Logger 模块文档](../../../logger/) - 了解更多 logger 功能
- [示例 7: SigNoz 集成](../07_with_signoz/) - 完整可观测性
- [Logger 示例 5: 生产配置](../../../logger/examples/05_production/) - 生产环境最佳实践

