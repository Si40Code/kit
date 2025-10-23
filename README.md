# go-pkg-sdk

> 一个模块化、易用、生产就绪的 Go SDK，帮助你快速替换老项目中的常用模块。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 🎯 项目目标

`go-pkg-sdk` 旨在提供一套开箱即用的 Go 模块，帮助开发者：

- 🚀 快速替换老项目中的遗留模块
- 📦 独立使用每个模块，无强制依赖
- 📚 通过丰富的示例快速上手
- 🏗️ 统一的 API 设计，降低学习成本
- ⚡ 生产级性能和稳定性

## 📦 模块列表

### ✅ Config - 配置管理模块

强大且易用的配置管理，基于 [koanf](https://github.com/knadh/koanf) 封装。

**特性：**
- 支持多种配置源：文件（YAML/JSON/TOML）、环境变量、远程配置
- 配置热更新：文件监控、远程配置推送
- 配置变更日志：自动记录配置变更历史
- 敏感信息脱敏：自动隐藏密码、token 等
- 线程安全：并发读取无问题

**快速开始：**

```go
import "github.com/silin/go-pkg-sdk/config"

// 初始化配置
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),
    config.WithFileWatcher(),
)

// 读取配置
dbHost := config.GetString("database.host")
dbPort := config.GetInt("database.port")

// 结构化读取
var dbConfig DatabaseConfig
config.Unmarshal("database", &dbConfig)

// 监听配置变更
config.OnChange(func() {
    log.Println("Config changed!")
})
```

**更多示例：** [config/examples](./config/examples)

---

### 🔄 Logger - 日志模块（即将推出）

高性能的日志模块，基于 [zap](https://github.com/uber-go/zap) 封装。

---

### 🌐 HTTPClient - HTTP 客户端模块（即将推出）

简单易用的 HTTP 客户端，基于标准库 `net/http`。

---

## 🚀 快速开始

### 安装

```bash
go get github.com/silin/go-pkg-sdk
```

### 单独使用某个模块

```go
// 只需要配置模块
import "github.com/silin/go-pkg-sdk/config"

func main() {
    config.Init(config.WithFile("config.yaml"))
    
    appName := config.GetString("app.name")
    fmt.Println("App:", appName)
}
```

### 组合使用多个模块

```go
import (
    "github.com/silin/go-pkg-sdk/config"
    "github.com/silin/go-pkg-sdk/logger"
    "github.com/silin/go-pkg-sdk/httpclient"
)

func main() {
    // 1. 初始化配置
    config.Init(
        config.WithFile("config.yaml"),
        config.WithEnv("APP_"),
    )
    
    // 2. 基于配置初始化日志
    logger.Init(
        logger.WithLevel(config.GetString("log.level")),
        logger.WithFormat(config.GetString("log.format")),
    )
    
    // 3. 创建 HTTP 客户端
    client := httpclient.New(
        httpclient.WithTimeout(config.GetInt("http.timeout")),
        httpclient.WithLogger(logger.Default()),
    )
    
    logger.Info("App started", 
        logger.String("name", config.GetString("app.name")),
    )
}
```

## 📖 文档

- [架构设计](./ARCHITECTURE.md) - 了解项目的设计理念和目录结构
- [快速开始](./docs/getting_started.md) - 从零开始的入门指南
- [API 参考](./docs/api_reference.md) - 详细的 API 文档
- [最佳实践](./examples/best_practices/) - 生产环境使用建议
- [迁移指南](./examples/migration_guide/) - 从其他库迁移到本 SDK

## 🌟 设计原则

### 1. 模块独立性

每个模块都可以**独立使用**，不强制依赖其他模块：

```go
// ✅ 只使用 config，不需要其他模块
import "github.com/silin/go-pkg-sdk/config"

// ✅ 只使用 logger，不需要其他模块
import "github.com/silin/go-pkg-sdk/logger"
```

### 2. 易用性优先

提供简洁的 API 和丰富的示例：

```go
// 简单直观的 API
config.Init(config.WithFile("config.yaml"))
value := config.GetString("key")

// 每个模块都有 5+ 个实际用例
// 参见 config/examples/
```

### 3. 统一的 API 风格

所有模块采用一致的 Options 模式：

```go
// Config 模块
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),
)

// Logger 模块（即将推出）
logger.Init(
    logger.WithLevel("info"),
    logger.WithFormat("json"),
)

// HTTPClient 模块（即将推出）
client := httpclient.New(
    httpclient.WithTimeout(30),
    httpclient.WithRetry(3),
)
```

### 4. 生产就绪

- ✅ 完善的错误处理
- ✅ 线程安全
- ✅ 性能优化
- ✅ 完整的单元测试
- ✅ 实战场景验证

## 📂 项目结构

```
go-pkg-sdk/
├── config/              # 配置管理模块
│   ├── examples/        # 5+ 使用示例
│   └── README.md        # 模块文档
├── logger/              # 日志模块（即将推出）
├── httpclient/          # HTTP 客户端模块（即将推出）
├── examples/            # 综合示例和最佳实践
└── docs/                # 项目文档
```

详细结构请参考 [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔧 配置模块详细说明

### 支持的配置源

| 配置源 | 说明 | 优先级 |
|--------|------|--------|
| 文件 | 支持 YAML、JSON、TOML | 低 |
| 环境变量 | 覆盖文件配置 | 中 |
| 远程配置 | Apollo、Nacos 等 | 高 |

### 使用示例

#### 1. 基础用法

```go
// examples/config/01_basic_usage/main.go
config.Init(config.WithFile("config.yaml"))

host := config.GetString("server.host")
port := config.GetInt("server.port")
debug := config.GetBool("app.debug")
```

#### 2. 环境变量覆盖

```go
// examples/config/02_env_override/main.go
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),  // APP_SERVER_PORT=8080
)

// 环境变量 APP_SERVER_PORT 会覆盖配置文件中的 server.port
port := config.GetInt("server.port")
```

#### 3. 文件监控（热更新）

```go
// examples/config/03_file_watch/main.go
config.Init(
    config.WithFile("config.yaml"),
    config.WithFileWatcher(),  // 启用文件监控
)

config.OnChange(func() {
    fmt.Println("配置已更新！")
    newValue := config.GetString("some.value")
})
```

#### 4. 远程配置（Apollo）

```go
// examples/config/04_remote_config/main.go
apolloProvider := NewApolloProvider(apolloConfig)

config.Init(
    config.WithFile("config.yaml"),  // 本地兜底配置
    config.WithRemote(apolloProvider), // 远程配置
)
```

#### 5. 配置变更通知

```go
// examples/config/05_change_notification/main.go
// 配置变更会自动输出 JSON 格式的日志
// {"type":"config_change","source":"file","key":"server.port","old":"8080","new":"9090","change":"UPDATE","timestamp":"2024-01-01T12:00:00Z"}
```

## 📝 配置文件示例

```yaml
# config.yaml
app:
  name: my-app
  version: 1.0.0
  debug: false

server:
  host: 0.0.0.0
  port: 8080
  timeout: 30

database:
  host: localhost
  port: 3306
  username: root
  password: secret123  # 变更日志中会自动脱敏为 ******
  database: mydb

log:
  level: info
  format: json
  output: stdout
```

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/amazing`)
3. 提交变更 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing`)
5. 提交 Pull Request

## 📜 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🗺️ 路线图

- [x] **v0.1** - Config 模块
- [ ] **v0.2** - Logger 模块
- [ ] **v0.3** - HTTPClient 模块
- [ ] **v0.4** - Cache 模块
- [ ] **v0.5** - Database 模块
- [ ] **v1.0** - 正式版本发布

## ❓ 常见问题

### Q: 为什么不直接使用 viper、zap？

A: 本 SDK 是对优秀开源库的封装，提供了：
- 更简洁的 API
- 统一的使用风格
- 开箱即用的最佳实践
- 丰富的使用示例
- 生产环境验证的配置

### Q: 可以只使用某一个模块吗？

A: 完全可以！所有模块都是独立的，按需导入即可：

```go
import "github.com/silin/go-pkg-sdk/config"
```

### Q: 如何从旧项目迁移？

A: 我们提供了详细的[迁移指南](./examples/migration_guide/)，涵盖：
- 从 viper 迁移到 config
- 从 logrus 迁移到 logger
- 从 resty 迁移到 httpclient

## 📮 联系方式

- 问题反馈：[GitHub Issues](https://github.com/silin/go-pkg-sdk/issues)
- 功能建议：[GitHub Discussions](https://github.com/silin/go-pkg-sdk/discussions)

---

⭐ 如果这个项目对你有帮助，请给个 Star！

