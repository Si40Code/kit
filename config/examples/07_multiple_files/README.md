# 示例 7: 多配置文件

展示如何使用多个配置文件，实现分层配置管理。

## 运行示例

```bash
cd config/examples/07_multiple_files
go run main.go
```

## 学习内容

1. **WithFile 多次调用** - 加载多个配置文件
2. **WithFiles** - 一次性加载多个配置文件
3. **配置分层** - base + env + local
4. **配置优先级** - 后加载的覆盖先加载的
5. **实际应用场景** - 多环境部署、模块化配置

## 配置文件说明

### 文件结构

```
07_multiple_files/
├── config-base.yaml          # 基础配置（通用）
├── config-dev.yaml           # 开发环境配置
├── config-prod.yaml          # 生产环境配置
├── config-local.yaml.example # 本地配置示例
└── main.go
```

### 配置分层

| 文件 | 用途 | 是否提交 Git |
|-----|------|-------------|
| config-base.yaml | 所有环境通用配置 | ✅ 是 |
| config-dev.yaml | 开发环境配置 | ✅ 是 |
| config-prod.yaml | 生产环境配置 | ✅ 是 |
| config-local.yaml | 本地个人配置 | ❌ 否 (.gitignore) |

## 使用方式

### 方式 1: 多次调用 WithFile

```go
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-dev.yaml"),
)
```

### 方式 2: 使用 WithFiles

```go
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-dev.yaml",
    ),
)
```

### 方式 3: 根据环境动态加载

```go
env := os.Getenv("ENV")
if env == "" {
    env = "dev"
}

config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile(fmt.Sprintf("config-%s.yaml", env)),
)
```

## 配置优先级

配置文件按加载顺序，后加载的会覆盖先加载的：

```
config-base.yaml
    ↓ (覆盖)
config-dev.yaml
    ↓ (覆盖)
config-local.yaml
    ↓ (覆盖)
环境变量
    ↓ (覆盖)
远程配置
```

**示例：**

```yaml
# config-base.yaml
server:
  port: 8080

# config-dev.yaml
server:
  port: 9090

# 最终值: 9090 (dev 覆盖了 base)
```

## 实际应用场景

### 场景 1: 多环境部署

```go
// 生产环境
// ENV=prod go run main.go

// 开发环境
// ENV=dev go run main.go

env := os.Getenv("ENV")
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile(fmt.Sprintf("config-%s.yaml", env)),
)
```

### 场景 2: 功能模块化配置

将不同功能的配置拆分到不同文件：

```go
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-database.yaml",   // 数据库配置
        "config-redis.yaml",      // Redis 配置
        "config-rabbitmq.yaml",   // 消息队列配置
        "config-elasticsearch.yaml", // ES 配置
    ),
)
```

**优点：**
- 配置文件更小，易于维护
- 可以按需加载配置
- 团队协作时减少冲突

### 场景 3: 团队协作

```
项目/
├── config-base.yaml          # 提交到 Git
├── config-dev.yaml           # 提交到 Git
├── config-prod.yaml          # 提交到 Git
├── config-local.yaml.example # 提交到 Git
└── config-local.yaml         # .gitignore (不提交)
```

**.gitignore**
```
config-local.yaml
```

**工作流：**
1. 开发者 clone 项目
2. 复制 `config-local.yaml.example` 为 `config-local.yaml`
3. 根据个人需要修改 `config-local.yaml`
4. 本地配置不会影响其他人

### 场景 4: 可选配置文件

```go
configFiles := []string{"config-base.yaml"}

// 检查环境配置是否存在
envConfig := fmt.Sprintf("config-%s.yaml", env)
if _, err := os.Stat(envConfig); err == nil {
    configFiles = append(configFiles, envConfig)
}

// 检查本地配置是否存在
if _, err := os.Stat("config-local.yaml"); err == nil {
    configFiles = append(configFiles, "config-local.yaml")
}

config.Init(config.WithFiles(configFiles...))
```

## 配置合并示例

**config-base.yaml:**
```yaml
app:
  name: myapp
  version: 1.0.0
server:
  host: 0.0.0.0
  port: 8080
database:
  host: localhost
```

**config-dev.yaml:**
```yaml
app:
  debug: true
server:
  port: 9090
database:
  host: dev-db.local
  password: dev-pass
```

**合并后的最终配置：**
```yaml
app:
  name: myapp        # 来自 base
  version: 1.0.0     # 来自 base
  debug: true        # 来自 dev
server:
  host: 0.0.0.0      # 来自 base
  port: 9090         # dev 覆盖了 base 的 8080
database:
  host: dev-db.local # dev 覆盖了 base 的 localhost
  password: dev-pass # 来自 dev
```

## 最佳实践

### 1. 配置分层清晰

```
✅ 好的分层:
- config-base.yaml    (通用配置)
- config-dev.yaml     (开发环境)
- config-prod.yaml    (生产环境)
- config-local.yaml   (个人配置)

❌ 不好的分层:
- config1.yaml
- config2.yaml
- final-config.yaml
```

### 2. 命名规范统一

```
✅ 推荐:
- config-{environment}.yaml
- config-{module}.yaml

❌ 不推荐:
- dev.yaml
- production-settings.yaml
- my_config.yml
```

### 3. 基础配置完整

config-base.yaml 应该包含所有必需的配置项，环境配置只覆盖需要变化的部分。

### 4. 提供示例文件

```bash
# 为敏感配置提供示例
cp config-local.yaml config-local.yaml.example
# 清除敏感信息
# 提交 .example 文件到 Git
```

### 5. 文档化配置

在 README 中说明：
- 每个配置文件的用途
- 配置加载顺序
- 如何创建本地配置
- 哪些文件不应提交

## 注意事项

1. **配置文件顺序很重要** - 后加载的会覆盖先加载的
2. **检查文件是否存在** - 对于可选配置文件，先检查再加载
3. **敏感信息管理** - 不要在配置文件中存储密码等敏感信息
4. **.gitignore 设置** - 确保本地配置不会被提交
5. **合并策略** - koanf 会深度合并嵌套的配置

## 扩展阅读

- [示例 1: 基础用法](../01_basic_usage/)
- [示例 2: 环境变量覆盖](../02_env_override/)
- [示例 6: 默认值功能](../06_default_values/)
- [Config 模块完整文档](../../README.md)



