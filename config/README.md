# Config 模块

强大且易用的配置管理模块，基于 [koanf](https://github.com/knadh/koanf) 封装。

## ✨ 特性

- ✅ **多种配置源**: 支持文件（**YAML/JSON/TOML**）、环境变量、远程配置
- ✅ **多配置文件**: 支持加载多个配置文件，实现分层配置管理
- ✅ **自动格式识别**: 根据文件扩展名自动选择解析器
- ✅ **默认值支持**: 两种方式设置默认值，应用开箱即用
- ✅ **配置热更新**: 文件监控、远程配置推送自动生效
- ✅ **配置变更日志**: 自动记录配置变更历史，便于审计
- ✅ **敏感信息脱敏**: 自动隐藏密码、token 等敏感信息
- ✅ **线程安全**: 支持并发读取
- ✅ **类型安全**: 提供类型化的读取方法
- ✅ **简单易用**: 清晰的 API 设计，丰富的使用示例

## 🚀 快速开始

### 安装

```bash
go get github.com/silin/go-pkg-sdk/config
```

### 基础使用

```go
package main

import (
    "fmt"
    "github.com/silin/go-pkg-sdk/config"
)

func main() {
    // 初始化配置（支持 YAML、JSON、TOML 格式）
    if err := config.Init(config.WithFile("config.yaml")); err != nil {
        panic(err)
    }

    // 读取配置
    appName := config.GetString("app.name")
    port := config.GetInt("server.port")
    debug := config.GetBool("app.debug")

    fmt.Printf("App: %s, Port: %d, Debug: %v\n", appName, port, debug)
}
```

## 📖 使用示例

我们提供了 8 个详细的使用示例，涵盖各种实际场景：

### 1️⃣ [基础用法](./examples/01_basic_usage/) - 从文件读取配置

展示如何从 YAML 文件读取配置，包括：
- 读取各种类型的配置值
- 结构化读取（Unmarshal）
- 配置文件示例

### 2️⃣ [环境变量覆盖](./examples/02_env_override/) - 环境变量覆盖文件配置

展示如何使用环境变量覆盖文件配置，包括：
- 环境变量前缀设置
- 优先级说明
- 实际应用场景（开发/测试/生产环境）

### 3️⃣ [文件监控](./examples/03_file_watch/) - 配置文件热更新

展示如何监控配置文件变更，实现热更新，包括：
- 启用文件监控
- 注册变更回调
- 动态调整应用行为

### 4️⃣ [远程配置](./examples/04_remote_config/) - 接入 Apollo 配置中心

展示如何接入远程配置中心（以 Apollo 为例），包括：
- 实现 RemoteProvider 接口
- 远程配置热更新
- 本地配置作为兜底

### 5️⃣ [配置变更通知](./examples/05_change_notification/) - 配置变更日志和通知

展示配置变更的日志记录和通知机制，包括：
- 自动记录配置变更
- 敏感信息脱敏
- 多个组件监听配置变更

### 6️⃣ [默认值功能](./examples/06_default_values/) - 设置和使用默认值

展示如何使用默认值功能，包括：
- WithDefaults 设置全局默认值
- GetXxxOr 方法指定局部默认值
- Exists 检查配置是否存在
- 配置优先级说明

### 7️⃣ [多配置文件](./examples/07_multiple_files/) - 分层配置管理

展示如何使用多个配置文件实现分层配置，包括：
- WithFile 多次调用加载多个文件
- WithFiles 一次性加载多个文件
- 配置分层架构（base + env + local）
- 多环境部署和团队协作

### 8️⃣ [业务模块化配置](./examples/08_business_modules/) - 业务配置分离

展示如何将不同业务模块的配置拆分到独立文件，包括：
- 业务模块配置分离（sms.yaml、email.yaml、payment.yaml）
- 条件加载业务模块
- 配置验证和团队协作
- 微服务架构配置管理

## 🔧 API 参考

### 初始化

#### `Init(opts ...Option) error`

初始化配置模块。

**参数：**
- `opts`: 配置选项（可选）

**返回：**
- `error`: 初始化错误

**示例：**

```go
// 只从文件加载
config.Init(config.WithFile("config.yaml"))

// 文件 + 环境变量
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),
)

// 文件 + 环境变量 + 文件监控
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),
    config.WithFileWatcher(),
)

// 完整配置（文件 + 环境变量 + 远程配置 + 文件监控）
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),
    config.WithRemote(apolloProvider),
    config.WithFileWatcher(),
)
```

### 配置选项

#### `WithFile(path string) Option`

加载单个配置文件。**自动根据文件扩展名识别格式**，支持 YAML (`.yaml`, `.yml`)、JSON (`.json`)、TOML (`.toml`) 格式。可以多次调用加载多个文件。

**示例：**
```go
// YAML 格式
config.Init(config.WithFile("config.yaml"))

// JSON 格式
config.Init(config.WithFile("config.json"))

// TOML 格式
config.Init(config.WithFile("config.toml"))

// 多个文件（多次调用，可以混合不同格式）
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-dev.json"),
    config.WithFile("config-local.toml"),
)
```

#### `WithFiles(paths ...string) Option`

一次性加载多个配置文件（按顺序加载，后面的覆盖前面的）。

**示例：**
```go
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-dev.yaml",
        "config-local.yaml",
    ),
)

// 根据环境动态加载
env := os.Getenv("ENV")
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile(fmt.Sprintf("config-%s.yaml", env)),
)
```

#### `WithEnv(prefix string) Option`

从环境变量加载配置。环境变量名格式：`PREFIX_KEY_NAME`

**示例：**
```go
config.WithEnv("APP_")

// APP_SERVER_PORT=8080 -> server.port = 8080
// APP_DATABASE_HOST=localhost -> database.host = "localhost"
```

#### `WithFileWatcher() Option`

启用配置文件监控，文件变更时自动重新加载。

#### `WithRemote(provider RemoteProvider) Option`

从远程配置中心加载配置（如 Apollo、Nacos）。

#### `WithDefaults(defaults map[string]interface{}) Option`

设置默认配置值（优先级最低）。

**示例：**
```go
defaults := map[string]interface{}{
    "server.port": 8080,
    "app.timeout": 30,
    "log.level":   "info",
}

config.Init(
    config.WithDefaults(defaults),
    config.WithFile("config.yaml"),
)
```

#### `WithDefaultStruct(defaultStruct interface{}) Option`

从结构体设置默认配置值。

**示例：**
```go
type DefaultConfig struct {
    Server struct {
        Port int `koanf:"port"`
    } `koanf:"server"`
}

defaults := DefaultConfig{}
defaults.Server.Port = 8080

config.Init(
    config.WithDefaultStruct(defaults),
    config.WithFile("config.yaml"),
)
```

### 读取配置

#### `GetString(path string) string`

读取字符串配置。

**示例：**
```go
host := config.GetString("database.host")
```

#### `GetInt(path string) int`

读取整数配置。

**示例：**
```go
port := config.GetInt("server.port")
```

#### `GetBool(path string) bool`

读取布尔配置。

**示例：**
```go
debug := config.GetBool("app.debug")
```

#### `GetFloat64(path string) float64`

读取浮点数配置。

**示例：**
```go
ratio := config.GetFloat64("cache.ratio")
```

#### `GetStringSlice(path string) []string`

读取字符串数组配置。

**示例：**
```go
hosts := config.GetStringSlice("redis.hosts")
```

#### `GetStringOr(path, defaultValue string) string`

读取字符串配置，如果不存在则返回默认值。

**示例：**
```go
logFile := config.GetStringOr("log.file", "/var/log/app.log")
```

#### `GetIntOr(path string, defaultValue int) int`

读取整数配置，如果不存在则返回默认值。

**示例：**
```go
maxRetry := config.GetIntOr("http.max_retry", 3)
```

#### `GetBoolOr(path string, defaultValue bool) bool`

读取布尔配置，如果不存在则返回默认值。

**示例：**
```go
enableCache := config.GetBoolOr("cache.enabled", true)
```

#### `GetFloat64Or(path string, defaultValue float64) float64`

读取浮点数配置，如果不存在则返回默认值。

**示例：**
```go
ratio := config.GetFloat64Or("cache.ratio", 0.75)
```

#### `GetStringSliceOr(path string, defaultValue []string) []string`

读取字符串数组配置，如果不存在则返回默认值。

**示例：**
```go
allowedIPs := config.GetStringSliceOr("security.allowed_ips", []string{"127.0.0.1"})
```

#### `Exists(path string) bool`

检查配置键是否存在。

**示例：**
```go
if config.Exists("app.name") {
    name := config.GetString("app.name")
} else {
    name := "default-name"
}
```

#### `Unmarshal(path string, out interface{}) error`

将配置反序列化到结构体。

**示例：**
```go
type DatabaseConfig struct {
    Host     string `koanf:"host"`
    Port     int    `koanf:"port"`
    Username string `koanf:"username"`
    Password string `koanf:"password"`
}

var dbConfig DatabaseConfig
if err := config.Unmarshal("database", &dbConfig); err != nil {
    log.Fatal(err)
}
```

### 配置变更通知

#### `OnChange(callback func())`

注册配置变更回调函数。当配置发生变更时（文件变更或远程配置推送），会调用所有注册的回调函数。

**示例：**
```go
config.OnChange(func() {
    fmt.Println("Config changed!")
    // 重新读取配置
    newValue := config.GetString("some.key")
    // 更新应用行为
})
```

## 📋 配置优先级

配置的加载顺序和优先级（从低到高）：

1. **默认值** - 最低优先级（WithDefaults）
2. **文件配置** - 基础配置
3. **环境变量** - 覆盖文件配置
4. **远程配置** - 最高优先级

**示例：**

```yaml
# config.yaml
server:
  port: 8080
```

```bash
# 环境变量
export APP_SERVER_PORT=9090
```

```go
config.Init(
    config.WithFile("config.yaml"),
    config.WithEnv("APP_"),
)

// 实际值为 9090（环境变量覆盖了文件配置）
port := config.GetInt("server.port")
```

## 🔒 敏感信息脱敏

配置变更日志会自动脱敏敏感信息，关键词包括：

- `password`
- `secret`
- `token`
- `key`

**示例日志：**

```json
{
  "type": "config_change",
  "source": "file",
  "key": "database.password",
  "old": "******",
  "new": "******",
  "change": "UPDATE",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 🏗️ 远程配置接入

要接入远程配置中心（如 Apollo、Nacos），需要实现 `RemoteProvider` 接口：

```go
type RemoteProvider interface {
    Load(ctx context.Context, k *koanf.Koanf) error
    Watch(ctx context.Context, onChange func(map[string]interface{})) error
}
```

**Apollo 示例：**

参见 [examples/04_remote_config](./examples/04_remote_config/)

## 📝 配置文件格式

Config 模块支持三种配置文件格式，**自动根据文件扩展名选择解析器**：

| 格式 | 扩展名 | 特点 |
|------|--------|------|
| YAML | `.yaml`, `.yml` | 可读性最好，支持注释，适合人工编辑 |
| JSON | `.json` | 最通用，易于程序生成和解析 |
| TOML | `.toml` | 结构清晰，配置明确，适合配置文件 |

### YAML 格式示例

```yaml
app:
  name: my-application
  version: 1.0.0
  environment: production
  debug: false

server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 30
  write_timeout: 30

database:
  host: localhost
  port: 3306
  username: root
  password: secret123
  database: mydb
  max_connections: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

log:
  level: info
  format: json
  output: stdout
```

### JSON 格式示例

```json
{
  "app": {
    "name": "my-application",
    "version": "1.0.0",
    "environment": "production",
    "debug": false
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": 30,
    "write_timeout": 30
  },
  "database": {
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "secret123",
    "database": "mydb",
    "max_connections": 100
  }
}
```

### TOML 格式示例

```toml
[app]
name = "my-application"
version = "1.0.0"
environment = "production"
debug = false

[server]
host = "0.0.0.0"
port = 8080
read_timeout = 30
write_timeout = 30

[database]
host = "localhost"
port = 3306
username = "root"
password = "secret123"
database = "mydb"
max_connections = 100
```

### 混合使用多种格式

你可以同时使用不同格式的配置文件，后面加载的会覆盖前面的：

```go
config.Init(
    config.WithFile("config-base.yaml"),    // 基础配置用 YAML
    config.WithFile("config-env.json"),     // 环境配置用 JSON
    config.WithFile("config-local.toml"),   // 本地配置用 TOML
)
```

## 💡 最佳实践

### 1. 配置分层

```yaml
# config-base.yaml - 基础配置
app:
  name: my-app

# config-dev.yaml - 开发环境
app:
  debug: true

# config-prod.yaml - 生产环境
app:
  debug: false
```

### 2. 环境变量用于敏感信息

```bash
# 不要在配置文件中存储敏感信息
export APP_DATABASE_PASSWORD=secret123
export APP_API_TOKEN=xyz789
```

### 3. 配置验证

```go
config.Init(config.WithFile("config.yaml"))

// 验证必需的配置
if config.GetString("database.host") == "" {
    log.Fatal("database.host is required")
}

if config.GetInt("server.port") == 0 {
    log.Fatal("server.port is required")
}
```

### 4. 使用结构体

```go
// 定义配置结构体
type AppConfig struct {
    Server   ServerConfig   `koanf:"server"`
    Database DatabaseConfig `koanf:"database"`
    Redis    RedisConfig    `koanf:"redis"`
}

// 一次性读取所有配置
var cfg AppConfig
if err := config.Unmarshal("", &cfg); err != nil {
    log.Fatal(err)
}
```

## 🐛 常见问题

### Q: 配置文件路径找不到？

A: 使用绝对路径或相对于工作目录的路径：

```go
// 绝对路径
config.Init(config.WithFile("/etc/myapp/config.yaml"))

// 相对于工作目录
config.Init(config.WithFile("./configs/config.yaml"))
```

### Q: 环境变量不生效？

A: 确保环境变量名格式正确：

```go
config.Init(config.WithEnv("APP_"))

// ✅ 正确: APP_SERVER_PORT -> server.port
// ❌ 错误: SERVER_PORT -> 不会生效
```

### Q: 配置变更后应用没有响应？

A: 确保启用了文件监控并注册了回调：

```go
config.Init(
    config.WithFile("config.yaml"),
    config.WithFileWatcher(), // 必须启用
)

config.OnChange(func() {
    // 必须注册回调
    fmt.Println("Config changed!")
})
```

## 📚 更多资源

- [完整示例代码](./examples/)
- [项目架构文档](../ARCHITECTURE.md)
- [API 文档](https://pkg.go.dev/github.com/silin/go-pkg-sdk/config)
- [问题反馈](https://github.com/silin/go-pkg-sdk/issues)

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE)

