# kit 架构设计

## 🎯 项目概述

`kit` 是一个模块化的 Go 工具包，旨在提供开箱即用的常用功能模块，帮助快速构建可观测的应用。

## 🌟 设计原则

1. **模块独立性**: 每个模块都可以独立使用，不强制依赖其他模块
2. **可观测性优先**: 所有模块内置 log/trace/metric 支持
3. **接口与实现分离**: 清晰的抽象层，易于扩展和测试
4. **生产就绪**: 所有模块都经过实战验证

## 📂 目录结构

```
kit/
├── config/              # 配置管理模块
│   ├── config.go        # 核心API
│   ├── option.go         # Options模式
│   ├── watcher.go        # 文件监控
│   ├── provider/         # 配置提供者
│   └── examples/         # 使用示例
├── log/                  # 日志模块（计划中）
├── trace/                # 追踪模块（计划中）
├── metric/               # 指标模块（计划中）
├── http/                 # HTTP客户端（计划中）
├── db/                   # 数据库封装（计划中）
├── examples/             # 综合示例
└── docs/                 # 项目文档
```

## 🔧 模块设计

### Config 模块

**核心文件：**
- `config.go` - 全局API，提供类型安全的配置读取
- `option.go` - Options模式配置
- `watcher.go` - 文件监控实现
- `provider/` - 配置源（文件/环境变量/远程）

**使用方式：**
```go
config.Init(config.WithFile("app.yaml"))
value := config.GetString("key")
```

### Logger 模块（计划中）

**核心接口：**
```go
type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    WithContext(ctx context.Context) Logger
}
```

### HTTP 客户端（计划中）

**核心功能：**
- 自动记录请求/响应日志
- 自动生成trace span
- 自动上报metrics
- 支持重试、熔断等

### DB 封装（计划中）

**核心功能：**
- GORM Logger适配器
- 自动记录SQL日志
- 自动生成trace
- 慢查询警告

## 🎨 API设计规范

所有模块遵循统一的API设计：

### 1. Options模式

```go
// Config模块
config.Init(
    config.WithFile("app.yaml"),
    config.WithEnv("APP_"),
)

// Logger模块
logger := log.New(
    log.WithLevel("info"),
    log.WithFormat("json"),
)
```

### 2. 接口与实现分离

```
log/
├── logger.go      # 接口定义
└── zap/           # zap实现
    └── logger.go
```

### 3. 接受接口，返回结构体

```go
// 构造函数接受接口
func NewDB(logger Logger, config Config) *DB {
    // ...
}

// 返回具体类型
return &DB{...}
```

## 📦 依赖管理

每个模块的依赖都声明在根目录的 `go.mod` 中：

```go
require (
    // config模块
    github.com/knadh/koanf/v2 v2.1.1
    github.com/fsnotify/fsnotify v1.7.0
    
    // log模块
    go.uber.org/zap v1.26.0
)
```

## 🧪 测试策略

- **单元测试**: 覆盖率 > 80%
- **集成测试**: 真实场景测试
- **性能测试**: Benchmark
- **示例测试**: 确保示例代码可运行

## 📝 文档要求

每个模块必须提供：
1. README.md - 模块概览和快速开始
2. examples/ - 至少 5 个实际用例
3. 代码注释 - godoc 格式
4. 最佳实践 - 生产环境使用建议

## 🚀 路线图

- [x] **v0.1** - Config 模块
- [ ] **v0.2** - Logger 模块
- [ ] **v0.3** - Trace 模块
- [ ] **v0.4** - HTTP 客户端
- [ ] **v0.5** - DB 封装
- [ ] **v1.0** - 正式版本发布

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交变更
4. 通过所有测试
5. 提交 Pull Request

