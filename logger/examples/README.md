# Logger 模块示例

本目录包含 logger 模块的完整使用示例，从基础用法到生产环境配置。

## 📚 示例列表

### [01_basic](./01_basic/) - 基础用法

**适合人群**: 新手，快速入门

**包含内容**:
- 使用默认 logger
- 初始化自定义配置
- 不同日志级别使用
- 结构化字段 vs Map 字段
- 子 logger (With 方法)
- 创建独立实例

**运行**:
```bash
cd 01_basic && go run main.go
```

---

### [02_file_output](./02_file_output/) - 文件输出和日志切割

**适合人群**: 需要将日志持久化到文件的开发者

**包含内容**:
- 基本文件输出
- 日志切割配置（大小、时间、数量）
- 文件压缩
- 多输出（stdout + file）
- 级别分离（info.log + error.log）

**运行**:
```bash
cd 02_file_output && go run main.go
```

---

### [03_with_trace](./03_with_trace/) - Trace 集成

**适合人群**: 需要实现分布式追踪的开发者

**包含内容**:
- OpenTelemetry trace 集成
- 自动提取 trace_id 和 span_id
- Error 日志自动标记 span
- 日志作为 span event
- 嵌套 span 示例
- 完整业务流程示例

**运行**:
```bash
cd 03_with_trace && go run main.go
```

---

### [04_remote_signoz](./04_remote_signoz/) - SigNoz 远程日志

**适合人群**: 需要集中式日志管理和可观测性的团队

**包含内容**:
- OTLP 协议发送日志到 SigNoz
- Trace 和日志关联
- 同时输出到本地和远程
- 完整的可观测性示例

**前置条件**:
- Docker 或 Docker Compose
- 运行中的 SigNoz 实例

**运行**:
```bash
# 1. 启动 SigNoz
cd /path/to/signoz/deploy && ./install.sh

# 2. 运行示例
cd 04_remote_signoz && go run main.go

# 3. 访问 SigNoz UI
open http://localhost:3301
```

---

### [05_production](./05_production/) - 生产环境配置

**适合人群**: 准备部署到生产环境的团队

**包含内容**:
- 完整的生产级配置
- 环境感知（dev/staging/prod）
- 多输出策略
- 配置管理（环境变量 + 配置文件）
- 优雅关闭
- Docker 和 Kubernetes 部署示例

**运行**:
```bash
# 开发环境
ENV=development go run main.go

# 生产环境
ENV=production \
LOG_LEVEL=info \
LOG_FORMAT=json \
LOG_FILE_PATH=/var/log/myapp/app.log \
go run main.go
```

---

## 🎯 学习路径

### 初学者

1. **[01_basic](./01_basic/)** - 了解基本 API
2. **[02_file_output](./02_file_output/)** - 学习文件输出
3. **[05_production](./05_production/)** - 了解生产配置

### 需要分布式追踪

1. **[03_with_trace](./03_with_trace/)** - 学习 trace 集成
2. **[04_remote_signoz](./04_remote_signoz/)** - 完整可观测性

### 生产环境部署

1. **[05_production](./05_production/)** - 完整生产配置
2. **[04_remote_signoz](./04_remote_signoz/)** - 远程日志收集

---

## 📊 功能对比

| 功能 | 01_basic | 02_file | 03_trace | 04_signoz | 05_prod |
|------|----------|---------|----------|-----------|---------|
| Stdout 输出 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 文件输出 | ❌ | ✅ | ❌ | ❌ | ✅ |
| 日志切割 | ❌ | ✅ | ❌ | ❌ | ✅ |
| Trace 集成 | ❌ | ❌ | ✅ | ✅ | ✅ |
| 远程输出 | ❌ | ❌ | ❌ | ✅ | ✅ |
| 环境感知 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 优雅关闭 | ❌ | ❌ | ❌ | ❌ | ✅ |

---

## 🚀 快速开始

### 最简单的例子

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

### 生产环境例子

```go
package main

import (
    "context"
    "github.com/Si40Code/kit/logger"
)

func main() {
    // 初始化 logger
    err := logger.Init(
        logger.WithLevel(logger.InfoLevel),
        logger.WithFormat(logger.JSONFormat),
        logger.WithFile("/var/log/app.log",
            logger.WithFileMaxSize(100),
            logger.WithFileMaxAge(30),
            logger.WithFileMaxBackups(10),
        ),
        logger.WithOTLP("signoz.example.com:4317"),
        logger.WithTrace("my-service"),
    )
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
    
    ctx := context.Background()
    logger.Info(ctx, "Application started")
}
```

---

## 💡 常见问题

### Q: 应该使用全局 logger 还是实例化？

**A**: 两种方式都支持，根据场景选择：

```go
// 全局 logger（简单场景）
logger.Info(ctx, "message")

// 实例化（需要多个配置）
l1, _ := logger.New(logger.WithLevel(logger.InfoLevel))
l2, _ := logger.New(logger.WithLevel(logger.DebugLevel))
```

### Q: 如何在生产环境降低日志开销？

**A**: 
1. 使用合适的日志级别（生产用 Info，开发用 Debug）
2. 使用结构化字段而不是字符串拼接
3. 避免记录大对象
4. 使用异步输出（OTLP 自动批量发送）

### Q: 日志文件切割不生效？

**A**: 检查：
1. 文件路径和权限
2. MaxSize 配置（默认 100MB）
3. 是否有足够的磁盘空间

### Q: 如何与现有的 zap logger 迁移？

**A**: kit/logger 基于 zap 实现，API 设计相似：

```go
// 旧代码（zap）
zap.L().Info("message", zap.String("key", "value"))

// 新代码（kit/logger）
logger.Info(ctx, "message", "key", "value")
```

---

## 📖 更多资源

- [Logger 模块 README](../) - 完整 API 文档
- [Kit 项目主页](../../) - 了解更多模块
- [OpenTelemetry 文档](https://opentelemetry.io/docs/) - Trace 集成
- [SigNoz 文档](https://signoz.io/docs/) - 可观测性平台

---

## 🤝 贡献

发现问题或有改进建议？欢迎提交 Issue 或 PR！

