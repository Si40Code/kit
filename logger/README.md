# Logger 模块

高性能、易用的日志模块，基于 [zap](https://github.com/uber-go/zap) 封装。

## ✨ 特性

- ✅ **五级日志** - Debug、Info、Warn、Error、Fatal
- ✅ **双 API 风格** - 支持结构化字段和 Map 字段两种方式
- ✅ **多种输出** - stdout、文件、远程（OTLP 协议）
- ✅ **日志切割** - 基于大小、时间、数量的自动切割
- ✅ **Trace 集成** - 自动关联 OpenTelemetry trace
- ✅ **Error 标记** - Error/Fatal 日志自动标记 span
- ✅ **全局 + 实例** - 同时支持全局 logger 和独立实例
- ✅ **零依赖配置** - 模块完全独立，不强制依赖其他 kit 模块
- ✅ **生产就绪** - 高性能、线程安全、优雅关闭

## 🚀 快速开始

### 安装

```bash
go get github.com/Si40Code/kit/logger
```

### 基础使用

```go
package main

import (
    "context"
    "github.com/Si40Code/kit/logger"
)

func main() {
    ctx := context.Background()
    
    // 使用默认 logger（开箱即用）
    logger.Info(ctx, "应用启动成功")
    
    // 结构化字段
    logger.Info(ctx, "用户登录",
        "user_id", 12345,
        "username", "alice",
        "ip", "192.168.1.1",
    )
    
    // Map 字段
    logger.InfoMap(ctx, "订单创建", map[string]any{
        "order_id": "ORD-2024-001",
        "amount": 99.99,
    })
}
```

### 自定义配置

```go
err := logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithStdout(),
    logger.WithFile("/var/log/app.log",
        logger.WithFileMaxSize(100),    // 100MB
        logger.WithFileMaxAge(7),       // 保留 7 天
        logger.WithFileMaxBackups(3),   // 保留 3 个备份
    ),
)
if err != nil {
    panic(err)
}
defer logger.Sync()
```

## 📖 使用示例

我们提供了 5 个详细的使用示例：

| 示例 | 说明 | 适合场景 |
|------|------|---------|
| [01_basic](./examples/01_basic/) | 基础用法 | 快速入门 |
| [02_file_output](./examples/02_file_output/) | 文件输出和切割 | 日志持久化 |
| [03_with_trace](./examples/03_with_trace/) | Trace 集成 | 分布式追踪 |
| [04_remote_signoz](./examples/04_remote_signoz/) | SigNoz 远程日志 | 集中式管理 |
| [05_production](./examples/05_production/) | 生产环境配置 | 生产部署 |

[查看所有示例 →](./examples/)

## 🔧 API 参考

### 日志级别

```go
type Level int8

const (
    DebugLevel Level = iota - 1  // 调试信息
    InfoLevel                     // 一般信息
    WarnLevel                     // 警告信息
    ErrorLevel                    // 错误信息
    FatalLevel                    // 致命错误（程序退出）
)
```

### 初始化

#### `Init(opts ...Option) error`

初始化全局默认 logger。

```go
err := logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithStdout(),
)
```

#### `New(opts ...Option) (Logger, error)`

创建新的 logger 实例。

```go
l, err := logger.New(
    logger.WithLevel(logger.DebugLevel),
    logger.WithStdout(),
)
```

### 配置选项

#### 基础配置

##### `WithLevel(level Level) Option`

设置日志级别。

```go
logger.WithLevel(logger.DebugLevel)  // debug, info, warn, error, fatal
logger.WithLevel(logger.InfoLevel)   // info, warn, error, fatal
logger.WithLevel(logger.ErrorLevel)  // error, fatal
```

##### `WithFormat(format Format) Option`

设置日志格式。

```go
logger.WithFormat(logger.JSONFormat)     // JSON 格式（生产推荐）
logger.WithFormat(logger.ConsoleFormat)  // 控制台格式（开发推荐）
```

##### `WithDevelopment() Option`

启用开发模式（Debug 级别 + 控制台格式 + 堆栈跟踪）。

```go
logger.WithDevelopment()
```

#### 输出配置

##### `WithStdout() Option`

添加标准输出。

```go
logger.WithStdout()
```

##### `WithFile(path string, opts ...FileOption) Option`

添加文件输出。

```go
logger.WithFile("/var/log/app.log",
    logger.WithFileMaxSize(100),      // 最大 100MB
    logger.WithFileMaxAge(7),         // 保留 7 天
    logger.WithFileMaxBackups(3),     // 保留 3 个备份
    logger.WithFileCompress(),        // 压缩旧文件
)
```

**文件选项**:
- `WithFileMaxSize(mb int)` - 文件最大大小（MB），默认 100
- `WithFileMaxAge(days int)` - 文件最大保留天数，默认 7
- `WithFileMaxBackups(count int)` - 最大备份文件数，默认 3
- `WithFileCompress()` - 启用文件压缩

##### `WithOTLP(endpoint string, opts ...OTLPOption) Option`

添加 OTLP 输出（用于 SigNoz、Jaeger 等）。

```go
logger.WithOTLP("localhost:4317",
    logger.WithOTLPInsecure(),  // 使用不安全连接（开发环境）
    logger.WithOTLPTimeout(5*time.Second),
)
```

**OTLP 选项**:
- `WithOTLPInsecure()` - 使用不安全连接（HTTP）
- `WithOTLPHeaders(headers map[string]string)` - 自定义 headers
- `WithOTLPTimeout(timeout time.Duration)` - 连接超时

#### Trace 配置

##### `WithTrace(serviceName string) Option`

启用 OpenTelemetry trace 集成。

```go
logger.WithTrace("my-service")
```

启用后，日志会自动：
- 从 context 提取 trace_id 和 span_id
- Error/Fatal 日志标记 span 为 error
- 将日志作为 span event 记录

### 日志记录

#### 结构化字段方式

```go
logger.Debug(ctx, "调试信息", "key1", "value1", "key2", value2)
logger.Info(ctx, "一般信息", "key", "value")
logger.Warn(ctx, "警告信息", "key", "value")
logger.Error(ctx, "错误信息", "key", "value")
logger.Fatal(ctx, "致命错误", "key", "value")  // 会终止程序
```

#### Map 字段方式

```go
logger.DebugMap(ctx, "调试信息", map[string]any{"key": "value"})
logger.InfoMap(ctx, "一般信息", map[string]any{"key": "value"})
logger.WarnMap(ctx, "警告信息", map[string]any{"key": "value"})
logger.ErrorMap(ctx, "错误信息", map[string]any{"key": "value"})
logger.FatalMap(ctx, "致命错误", map[string]any{"key": "value"})
```

#### 子 Logger

```go
// 创建带预设字段的子 logger
userLogger := logger.With("user_id", 12345, "session", "abc123")

// 后续日志会自动包含这些字段
userLogger.Info(ctx, "用户执行操作", "action", "login")
userLogger.Info(ctx, "用户执行操作", "action", "logout")
```

#### 独立实例

```go
// 创建独立的 logger 实例
appLogger, _ := logger.New(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFile("/var/log/app.log"),
)

debugLogger, _ := logger.New(
    logger.WithLevel(logger.DebugLevel),
    logger.WithFile("/var/log/debug.log"),
)

appLogger.Info(ctx, "应用日志")
debugLogger.Debug(ctx, "调试日志")
```

### 刷新和同步

#### `Sync() error`

刷新缓冲区，确保所有日志都被写入。

```go
defer logger.Sync()  // 程序退出前调用
```

## 🎯 使用场景

### 场景 1: Web 应用

```go
func main() {
    // 初始化 logger
    logger.Init(
        logger.WithLevel(logger.InfoLevel),
        logger.WithFormat(logger.JSONFormat),
        logger.WithStdout(),
        logger.WithFile("/var/log/app.log"),
    )
    defer logger.Sync()
    
    // 启动 web 服务器
    http.HandleFunc("/api/users", handleUsers)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
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

### 场景 2: 微服务 + Trace

```go
func main() {
    // 初始化 tracer
    initTracer()
    
    // 初始化 logger（启用 trace）
    logger.Init(
        logger.WithLevel(logger.InfoLevel),
        logger.WithTrace("order-service"),
        logger.WithOTLP("signoz:4317"),
    )
    defer logger.Sync()
    
    // 业务逻辑
    processOrder()
}

func processOrder() {
    ctx := context.Background()
    tracer := otel.Tracer("order-service")
    
    // 创建 span
    ctx, span := tracer.Start(ctx, "process-order")
    defer span.End()
    
    // 日志会自动包含 trace_id 和 span_id
    logger.Info(ctx, "开始处理订单", "order_id", "ORD-001")
    
    // 模拟错误
    if err := validateOrder(); err != nil {
        // Error 日志会自动标记 span 为 error
        logger.Error(ctx, "订单验证失败", "error", err)
        return
    }
    
    logger.Info(ctx, "订单处理完成")
}
```

### 场景 3: 批处理任务

```go
func main() {
    // 批处理任务配置
    logger.Init(
        logger.WithLevel(logger.DebugLevel),
        logger.WithFile("/var/log/batch.log",
            logger.WithFileMaxSize(500),   // 500MB
            logger.WithFileMaxAge(90),     // 保留 90 天
            logger.WithFileCompress(),     // 压缩
        ),
    )
    defer logger.Sync()
    
    ctx := context.Background()
    logger.Info(ctx, "批处理任务开始")
    
    for i := 0; i < 10000; i++ {
        processItem(ctx, i)
    }
    
    logger.Info(ctx, "批处理任务完成")
}

func processItem(ctx context.Context, id int) {
    logger.Debug(ctx, "处理项目", "id", id)
    // 处理逻辑...
}
```

## 🔒 Trace 集成

### 自动提取 Trace 信息

当启用 trace 集成后，日志会自动包含 trace 信息：

```json
{
  "timestamp": "2024-01-15T10:30:00.123Z",
  "level": "info",
  "message": "用户登录",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "user_id": 12345
}
```

### Error 日志标记 Span

Error 和 Fatal 日志会自动：
1. 将 span 状态设置为 `codes.Error`
2. 调用 `span.RecordError()`
3. 添加日志到 span events

```go
logger.Error(ctx, "数据库连接失败", "error", "timeout")

// 等价于：
// span.SetStatus(codes.Error, "数据库连接失败")
// span.RecordError(errors.New("数据库连接失败"))
// span.AddEvent("数据库连接失败")
```

### 完整示例

```go
// 初始化
logger.Init(
    logger.WithTrace("my-service"),
    logger.WithOTLP("signoz:4317"),
)

// 使用
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

logger.Info(ctx, "操作开始")        // 包含 trace_id
logger.Error(ctx, "操作失败")       // 标记 span 为 error
```

## 📊 输出格式

### JSON 格式（生产推荐）

```json
{
  "timestamp": "2024-01-15T10:30:00.123Z",
  "level": "info",
  "caller": "main.go:25",
  "message": "用户登录",
  "user_id": 12345,
  "username": "alice",
  "ip": "192.168.1.1"
}
```

### Console 格式（开发推荐）

```
2024-01-15T10:30:00.123+0800    INFO    main.go:25    用户登录    {"user_id": 12345, "username": "alice", "ip": "192.168.1.1"}
```

## 💡 最佳实践

### 1. 始终传递 Context

```go
// ✅ 正确 - 传递 context
func processOrder(ctx context.Context, orderID string) {
    logger.Info(ctx, "处理订单", "order_id", orderID)
}

// ❌ 错误 - 丢失 trace 信息
func processOrder(orderID string) {
    ctx := context.Background()  // 新的 context，没有 trace 信息
    logger.Info(ctx, "处理订单", "order_id", orderID)
}
```

### 2. 使用结构化字段

```go
// ✅ 正确 - 结构化字段
logger.Info(ctx, "用户登录", "user_id", userID, "ip", ip)

// ❌ 错误 - 字符串拼接
logger.Info(ctx, fmt.Sprintf("用户 %d 从 %s 登录", userID, ip))
```

### 3. 选择合适的日志级别

```go
logger.Debug(ctx, "详细调试信息")                    // 开发环境
logger.Info(ctx, "重要业务事件")                      // 正常流程
logger.Warn(ctx, "可能的问题", "disk_usage", "90%")  // 需要关注
logger.Error(ctx, "错误", "error", err)              // 需要处理
logger.Fatal(ctx, "致命错误", "error", err)          // 程序无法继续
```

### 4. 避免敏感信息

```go
// ❌ 错误 - 记录密码
logger.Info(ctx, "用户登录", "password", password)

// ✅ 正确 - 脱敏
logger.Info(ctx, "用户登录", "password", "***")

// ✅ 正确 - 只记录必要信息
logger.Info(ctx, "用户登录", "user_id", userID)
```

### 5. 优雅关闭

```go
func main() {
    logger.Init(/*...*/)
    defer logger.Sync()  // 确保日志刷新
    
    // 应用逻辑...
}
```

## 🐛 常见问题

### Q: 日志级别如何工作？

A: 设置的级别及以上的日志会被记录：

```go
logger.Init(logger.WithLevel(logger.InfoLevel))

logger.Debug(ctx, "debug")  // ❌ 不会记录
logger.Info(ctx, "info")    // ✅ 会记录
logger.Warn(ctx, "warn")    // ✅ 会记录
logger.Error(ctx, "error")  // ✅ 会记录
```

### Q: 如何同时输出到多个目标？

A: 使用多个 WithXxx 选项：

```go
logger.Init(
    logger.WithStdout(),                    // 控制台
    logger.WithFile("/var/log/app.log"),    // 文件
    logger.WithOTLP("signoz:4317"),         // 远程
)
```

### Q: Fatal 日志会终止程序吗？

A: 是的，Fatal 会调用 `os.Exit(1)` 终止程序。使用前请确保：
- 已经刷新所有日志（调用 Sync）
- 不需要执行 defer 语句
- 真的需要终止程序

### Q: 性能如何？

A: 基于 zap 实现，性能优异：
- 每条日志 < 1μs（微秒）
- 零内存分配（大多数情况）
- 异步 OTLP 发送不阻塞主流程

### Q: 如何与 web 模块集成？

A: 参见 [web 模块适配器](../web/)。

## 📚 更多资源

- [使用示例](./examples/) - 5 个完整示例
- [Kit 项目主页](../) - 了解更多模块
- [Zap 文档](https://github.com/uber-go/zap) - 底层日志库
- [OpenTelemetry](https://opentelemetry.io/) - Trace 标准
- [SigNoz](https://signoz.io/) - 可观测性平台

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE)

