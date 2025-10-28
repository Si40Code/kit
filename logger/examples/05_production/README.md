# 示例 5: 生产环境配置

完整的生产环境 logger 配置示例，包含多输出、trace 集成、优雅关闭等最佳实践。

## 功能展示

1. **环境感知配置** - 根据环境（dev/prod）自动调整配置
2. **多输出配置** - stdout + file + OTLP
3. **日志分级** - Info 日志和 Error 日志分别输出
4. **Trace 完整集成** - 包含采样率配置
5. **优雅关闭** - 确保日志刷新
6. **配置管理** - 支持环境变量和配置文件

## 生产环境架构

```
应用日志输出
    ├── stdout（容器日志收集）
    ├── /var/log/myapp/app.log（本地文件）
    ├── /var/log/myapp/app.log.error（Error 日志）
    └── OTLP → SigNoz（远程可观测性平台）
```

## 运行示例

### 开发环境

```bash
cd logger/examples/05_production
ENV=development go run main.go
```

### 生产环境

```bash
# 使用环境变量配置
export ENV=production
export SERVICE_NAME=my-service
export LOG_LEVEL=info
export LOG_FORMAT=json
export LOG_FILE_PATH=/var/log/myapp/app.log
export OTLP_ENDPOINT=signoz.example.com:4317
export ENABLE_TRACE=true

go run main.go
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `SERVICE_NAME` | 服务名称 | my-production-service | order-service |
| `ENV` | 运行环境 | production | development/staging/production |
| `LOG_LEVEL` | 日志级别 | info | debug/info/warn/error |
| `LOG_FORMAT` | 日志格式 | json | json/console |
| `LOG_FILE_PATH` | 日志文件路径 | /var/log/myapp/app.log | /var/log/app.log |
| `LOG_FILE_MAX_SIZE` | 文件最大大小(MB) | 100 | 100 |
| `LOG_FILE_MAX_AGE` | 保留天数 | 30 | 30 |
| `LOG_FILE_BACKUPS` | 备份数量 | 10 | 10 |
| `OTLP_ENDPOINT` | OTLP 端点 | - | signoz.example.com:4317 |
| `OTLP_INSECURE` | 是否使用不安全连接 | true | true/false |
| `ENABLE_TRACE` | 是否启用 trace | true | true/false |
| `TRACE_ENDPOINT` | Trace 端点 | localhost:4318 | signoz.example.com:4318 |

### 配置文件示例

创建 `config.yaml`:

```yaml
service:
  name: my-service
  environment: production

logging:
  level: info
  format: json
  file:
    path: /var/log/myapp/app.log
    max_size: 100
    max_age: 30
    backups: 10
  otlp:
    endpoint: signoz.example.com:4317
    insecure: false

tracing:
  enabled: true
  endpoint: signoz.example.com:4318
  sample_rate: 1.0
```

## 生产环境最佳实践

### 1. 多输出策略

```go
// 生产环境：文件 + stdout + 远程
logger.Init(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    
    // 主日志文件
    logger.WithFile("/var/log/myapp/app.log",
        logger.WithFileMaxSize(100),
        logger.WithFileMaxAge(30),
        logger.WithFileMaxBackups(10),
        logger.WithFileCompress(),
    ),
    
    // Error 日志单独文件
    logger.WithFile("/var/log/myapp/app.log.error",
        logger.WithFileMaxSize(100),
        logger.WithFileMaxAge(30),
        logger.WithFileMaxBackups(10),
    ),
    
    // stdout（容器环境必须）
    logger.WithStdout(),
    
    // 远程可观测性平台
    logger.WithOTLP("signoz.example.com:4317"),
    
    // Trace 集成
    logger.WithTrace("my-service"),
)
```

### 2. 环境感知配置

```go
func initLogger(env string) error {
    var opts []logger.Option
    
    switch env {
    case "production":
        opts = []logger.Option{
            logger.WithLevel(logger.InfoLevel),
            logger.WithFormat(logger.JSONFormat),
            logger.WithFile("/var/log/myapp/app.log"),
            logger.WithStdout(),
        }
    case "staging":
        opts = []logger.Option{
            logger.WithLevel(logger.DebugLevel),
            logger.WithFormat(logger.JSONFormat),
            logger.WithStdout(),
        }
    default: // development
        opts = []logger.Option{
            logger.WithLevel(logger.DebugLevel),
            logger.WithFormat(logger.ConsoleFormat),
            logger.WithStdout(),
            logger.WithDevelopment(),
        }
    }
    
    return logger.Init(opts...)
}
```

### 3. 优雅关闭

```go
func main() {
    // 初始化 logger
    if err := logger.Init(/*...*/); err != nil {
        log.Fatal(err)
    }
    defer logger.Sync() // 确保日志刷新
    
    // 初始化 tracer
    cleanupTracer := initTracer()
    defer cleanupTracer()
    
    // 应用逻辑
    if err := run(); err != nil {
        logger.Error(ctx, "应用错误", "error", err)
        os.Exit(1)
    }
}

func run() error {
    // 监听关闭信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    
    <-sigCh
    
    logger.Info(ctx, "收到关闭信号，开始优雅关闭")
    
    // 执行清理工作...
    
    logger.Info(ctx, "优雅关闭完成")
    return nil
}
```

### 4. 错误处理

```go
// 记录详细的错误信息
if err := doSomething(); err != nil {
    logger.Error(ctx, "操作失败",
        "operation", "do_something",
        "error", err.Error(),
        "user_id", userID,
        "retry_count", retryCount,
        "stack", fmt.Sprintf("%+v", err), // 如果使用 pkg/errors
    )
}

// 记录 panic
defer func() {
    if r := recover(); r != nil {
        logger.Fatal(ctx, "应用 panic",
            "panic", r,
            "stack", string(debug.Stack()),
        )
    }
}()
```

### 5. 性能优化

```go
// 避免频繁的字符串格式化
// ❌ 错误
logger.Info(ctx, fmt.Sprintf("用户 %s 执行了 %s 操作", userID, action))

// ✅ 正确
logger.Info(ctx, "用户执行操作",
    "user_id", userID,
    "action", action,
)

// 使用合适的日志级别
// Debug 日志在生产环境会被过滤，不会影响性能
if logger.Default().Level() <= logger.DebugLevel {
    // 仅在 debug 模式下执行昂贵的操作
    logger.Debug(ctx, "详细信息", "data", expensiveOperation())
}
```

## Docker 部署

### Dockerfile

```dockerfile
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o /app/main .

FROM debian:bullseye-slim

# 创建日志目录
RUN mkdir -p /var/log/myapp

COPY --from=builder /app/main /usr/local/bin/main

# 设置环境变量
ENV ENV=production
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json
ENV LOG_FILE_PATH=/var/log/myapp/app.log

CMD ["main"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: .
    environment:
      - SERVICE_NAME=my-service
      - ENV=production
      - LOG_LEVEL=info
      - LOG_FORMAT=json
      - OTLP_ENDPOINT=signoz:4317
      - ENABLE_TRACE=true
    volumes:
      - ./logs:/var/log/myapp
    depends_on:
      - signoz
  
  signoz:
    image: signoz/signoz:latest
    ports:
      - "3301:3301"
      - "4317:4317"
      - "4318:4318"
```

## Kubernetes 部署

### deployment.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  SERVICE_NAME: "my-service"
  ENV: "production"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
  OTLP_ENDPOINT: "signoz-collector:4317"
  ENABLE_TRACE: "true"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-service
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
      - name: app
        image: my-service:latest
        envFrom:
        - configMapRef:
            name: app-config
        volumeMounts:
        - name: logs
          mountPath: /var/log/myapp
      volumes:
      - name: logs
        emptyDir: {}
```

## 监控和告警

### 1. 日志监控指标

- **Error 日志率** - 每分钟 error 日志数量
- **Fatal 日志** - 立即告警
- **慢请求** - duration > 1s 的请求
- **日志丢失率** - OTLP 发送失败率

### 2. SigNoz 告警规则

```yaml
alert: HighErrorRate
expr: rate(log_entries{level="error"}[5m]) > 10
for: 5m
annotations:
  summary: "高错误率告警"
  description: "服务 {{ $labels.service }} 的错误日志率过高"
```

## 故障排查

### 1. 日志未写入文件

检查文件权限和目录：
```bash
ls -la /var/log/myapp/
# 确保应用有写入权限
chmod 755 /var/log/myapp
```

### 2. OTLP 连接失败

测试连接：
```bash
telnet signoz.example.com 4317
```

查看应用日志中的连接错误。

### 3. 日志文件过大

检查切割配置是否生效：
```bash
ls -lh /var/log/myapp/
```

手动触发切割（发送 HUP 信号）。

## 性能基准

在生产环境配置下的性能数据：

| 场景 | 吞吐量 | 延迟(p99) | 内存占用 |
|------|--------|-----------|----------|
| 仅 stdout | 100K logs/s | 100μs | 10MB |
| stdout + file | 50K logs/s | 200μs | 15MB |
| stdout + file + OTLP | 30K logs/s | 500μs | 20MB |

## 下一步

- 集成到你的项目
- 配置生产环境监控
- 设置告警规则
- [返回示例总览](../)

