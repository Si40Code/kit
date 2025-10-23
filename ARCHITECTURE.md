# go-pkg-sdk 架构设计

## 项目概述

`go-pkg-sdk` 是一个模块化的 Go SDK，旨在提供开箱即用的常用功能模块，帮助快速替换老项目中的相关组件。

## 设计原则

1. **模块独立性**: 每个模块都可以独立使用，不强制依赖其他模块
2. **易用性优先**: 提供简洁的 API 和丰富的示例代码
3. **生产就绪**: 所有模块都经过实战验证，可直接用于生产环境
4. **统一风格**: 所有模块采用一致的 API 设计风格（Options 模式）

## 目录结构

```
go-pkg-sdk/
├── go.mod                          # 主模块定义
├── README.md                       # 项目总览
├── ARCHITECTURE.md                 # 架构文档（本文件）
├── LICENSE                         # 开源许可证
│
├── config/                         # 配置管理模块
│   ├── config.go                   # 配置核心实现
│   ├── option.go                   # 配置选项
│   ├── watcher.go                  # 文件监控
│   ├── change_logger.go            # 配置变更日志
│   ├── remote.go                   # 远程配置接口定义
│   ├── provider/                   # 配置提供者实现
│   │   ├── file.go                 # 文件配置提供者
│   │   ├── env.go                  # 环境变量提供者
│   │   └── apollo.go               # Apollo 远程配置（示例）
│   ├── examples/                   # 配置模块示例
│   │   ├── 01_basic_usage/         # 基础用法
│   │   ├── 02_env_override/        # 环境变量覆盖
│   │   ├── 03_file_watch/          # 文件监控
│   │   ├── 04_remote_config/       # 远程配置
│   │   └── 05_change_notification/ # 变更通知
│   └── README.md                   # 配置模块文档
│
├── logger/                         # 日志模块
│   ├── logger.go                   # 日志核心实现
│   ├── option.go                   # 日志选项
│   ├── level.go                    # 日志级别
│   ├── formatter.go                # 日志格式化
│   ├── writer.go                   # 日志输出
│   ├── rotation.go                 # 日志轮转
│   ├── examples/                   # 日志模块示例
│   │   ├── 01_basic_logging/       # 基础日志
│   │   ├── 02_structured_logging/  # 结构化日志
│   │   ├── 03_log_rotation/        # 日志轮转
│   │   ├── 04_context_logging/     # 上下文日志
│   │   └── 05_performance/         # 性能优化
│   └── README.md                   # 日志模块文档
│
├── httpclient/                     # HTTP 客户端模块
│   ├── client.go                   # HTTP 客户端核心
│   ├── option.go                   # 客户端选项
│   ├── request.go                  # 请求构造器
│   ├── response.go                 # 响应处理器
│   ├── retry.go                    # 重试策略
│   ├── circuit_breaker.go          # 熔断器
│   ├── middleware/                 # 中间件
│   │   ├── logging.go              # 日志中间件
│   │   ├── tracing.go              # 链路追踪
│   │   ├── metrics.go              # 指标采集
│   │   └── auth.go                 # 认证中间件
│   ├── examples/                   # HTTP 客户端示例
│   │   ├── 01_basic_request/       # 基础请求
│   │   ├── 02_rest_api/            # RESTful API
│   │   ├── 03_retry_backoff/       # 重试和退避
│   │   ├── 04_circuit_breaker/     # 熔断器
│   │   ├── 05_middleware/          # 中间件使用
│   │   └── 06_performance/         # 性能优化
│   └── README.md                   # HTTP 客户端文档
│
├── examples/                       # 综合示例
│   ├── full_app/                   # 完整应用示例
│   │   ├── main.go                 # 使用所有模块
│   │   ├── config.yaml             # 配置文件
│   │   └── README.md               # 说明文档
│   ├── migration_guide/            # 迁移指南
│   │   ├── from_logrus.md          # 从 logrus 迁移
│   │   ├── from_viper.md           # 从 viper 迁移
│   │   └── from_resty.md           # 从 resty 迁移
│   └── best_practices/             # 最佳实践
│       ├── error_handling.md       # 错误处理
│       ├── testing.md              # 测试指南
│       └── production.md           # 生产环境部署
│
└── docs/                           # 文档目录
    ├── getting_started.md          # 快速开始
    ├── configuration.md            # 配置指南
    ├── api_reference.md            # API 参考
    └── faq.md                      # 常见问题
```

## 模块导入方式

### 单独使用某个模块

```go
// 只使用配置模块
import "github.com/silin/go-pkg-sdk/config"

// 只使用日志模块
import "github.com/silin/go-pkg-sdk/logger"

// 只使用 HTTP 客户端模块
import "github.com/silin/go-pkg-sdk/httpclient"
```

### 组合使用多个模块

```go
import (
    "github.com/silin/go-pkg-sdk/config"
    "github.com/silin/go-pkg-sdk/logger"
    "github.com/silin/go-pkg-sdk/httpclient"
)

func main() {
    // 初始化配置
    config.Init(config.WithFile("config.yaml"))
    
    // 使用配置初始化日志
    logger.Init(
        logger.WithLevel(config.GetString("log.level")),
        logger.WithFormat(config.GetString("log.format")),
    )
    
    // 创建 HTTP 客户端
    client := httpclient.New(
        httpclient.WithTimeout(config.GetInt("http.timeout")),
        httpclient.WithLogger(logger.Default()),
    )
}
```

## 各模块特性

### Config 模块

- ✅ 支持多种配置源：文件、环境变量、远程配置
- ✅ 文件格式：YAML、JSON、TOML
- ✅ 配置热更新：文件监控、远程配置推送
- ✅ 配置变更日志：自动记录配置变更
- ✅ 敏感信息脱敏：自动隐藏密码、token 等

### Logger 模块

- ⏳ 基于 zap 封装，高性能
- ⏳ 支持多种输出：控制台、文件、网络
- ⏳ 日志轮转：按大小、时间轮转
- ⏳ 结构化日志：JSON 格式
- ⏳ 上下文日志：自动注入 trace_id、user_id 等
- ⏳ 多级别：Debug、Info、Warn、Error、Fatal

### HTTPClient 模块

- ⏳ 基于标准库 net/http 封装
- ⏳ 链式 API：简洁易用
- ⏳ 自动重试：支持自定义重试策略
- ⏳ 熔断器：保护下游服务
- ⏳ 中间件：日志、追踪、指标、认证
- ⏳ 连接池：高性能连接管理

## API 设计规范

所有模块都遵循统一的 API 设计模式：

```go
// 1. Options 模式配置
type Option func(*options)

// 2. 简单初始化
func Init(opts ...Option) error

// 3. 链式调用
func New(opts ...Option) *Client

// 4. 上下文支持
func (c *Client) Do(ctx context.Context, ...) error

// 5. 错误处理
// - 返回具体的错误类型
// - 支持 errors.Is/As
```

## 依赖管理

每个模块的依赖都明确声明在 `go.mod` 中：

```go
require (
    // config 模块依赖
    github.com/knadh/koanf/v2 v2.1.1
    github.com/fsnotify/fsnotify v1.7.0
    
    // logger 模块依赖
    go.uber.org/zap v1.26.0
    
    // httpclient 模块依赖
    // 使用标准库，无额外依赖
)
```

## 版本管理

- 使用语义化版本 (Semantic Versioning)
- 主版本号：不兼容的 API 变更
- 次版本号：向下兼容的功能新增
- 修订号：向下兼容的问题修复

## 测试策略

- 单元测试：覆盖率 > 80%
- 集成测试：真实场景测试
- 性能测试：Benchmark
- 示例测试：确保示例代码可运行

## 文档要求

每个模块必须提供：

1. README.md - 模块概览和快速开始
2. examples/ - 至少 5 个实际用例
3. API 文档 - godoc 注释
4. 最佳实践 - 生产环境使用建议

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交变更
4. 通过所有测试
5. 提交 Pull Request

## 路线图

- [x] Config 模块 v1.0
- [ ] Logger 模块 v1.0
- [ ] HTTPClient 模块 v1.0
- [ ] Cache 模块 v2.0
- [ ] Database 模块 v2.0
- [ ] Queue 模块 v2.0

