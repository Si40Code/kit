# web - Gin 框架封装

> 基于 Gin 的生产级 HTTP 服务器封装，内置日志、链路追踪、指标监控支持

## ✨ 特性

- 🚀 **开箱即用** - 简单配置即可启动
- 📝 **自动日志** - 记录每个请求的详细信息
- 🔍 **链路追踪** - 集成 OpenTelemetry
- 📊 **指标监控** - 支持 Prometheus
- 🛡️ **智能处理** - 自动处理大文件和大响应
- 🔄 **优雅关闭** - 支持优雅关闭服务器
- ⚡ **高性能** - 基于 Gin，性能优异

## 📦 安装

```bash
go get github.com/Si40Code/kit/web
```

## 🚀 快速开始

```go
package main

import (
	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := web.New(
		web.WithMode(web.ReleaseMode),
		web.WithServiceName("my-service"),
	)

	engine := server.Engine()
	engine.GET("/ping", func(c *gin.Context) {
		web.Success(c, gin.H{"message": "pong"})
	})

	server.RunWithGracefulShutdown(":8080")
}
```

## 📚 示例

| 示例 | 说明 |
|------|------|
| [01_basic](./examples/01_basic) | 基础用法 |
| [02_with_trace](./examples/02_with_trace) | 集成链路追踪 |
| [03_with_metric](./examples/03_with_metric) | 集成指标监控 |
| [04_file_upload](./examples/04_file_upload) | 文件上传处理 |
| [05_custom_logger](./examples/05_custom_logger) | 自定义日志 |
| [06_production](./examples/06_production) | 生产环境配置 |

## 🔧 配置选项

### 基础配置

- `WithMode(mode)` - 设置运行模式（Debug/Release/Test）
- `WithServiceName(name)` - 设置服务名称

### 日志配置

- `WithLogger(logger)` - 自定义日志记录器
- `WithSkipPaths(paths...)` - 跳过特定路径的日志
- `WithMaxBodyLogSize(size)` - 设置最大 body 日志大小
- `WithSlowRequestThreshold(duration)` - 设置慢请求阈值

### 功能开关

- `WithTrace()` - 启用链路追踪
- `WithMetric(recorder)` - 启用指标监控
- `WithRecover()` - 启用 panic 恢复（默认开启）
- `WithCORS()` - 启用 CORS

### 文件上传

- `WithMaxMultipartMemory(size)` - 设置最大文件上传内存

## 🎯 最佳实践

### 1. 日志处理

- 敏感信息自动脱敏（密码、token 等）
- 大请求/响应自动截断
- 文件上传只记录元信息

### 2. 性能优化

- 跳过健康检查等高频端点的日志
- 合理设置 body 日志大小
- 使用异步指标上报

### 3. 生产环境

```go
server := web.New(
	web.WithMode(web.ReleaseMode),
	web.WithServiceName("prod-service"),
	web.WithLogger(yourLogger),
	web.WithTrace(),
	web.WithMetric(yourMetricRecorder),
	web.WithSkipPaths("/health", "/metrics"),
	web.WithMaxBodyLogSize(4096),
	web.WithSlowRequestThreshold(3*time.Second),
)
```

## 📝 日志格式

```json
{
  "client_ip": "127.0.0.1",
  "method": "POST",
  "path": "/api/user",
  "query": "id=123",
  "status": 200,
  "latency_ms": 45,
  "req_body": "{\"name\":\"John\"}",
  "resp_body": "{\"code\":0,\"data\":{...}}",
  "user_agent": "Mozilla/5.0..."
}
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
