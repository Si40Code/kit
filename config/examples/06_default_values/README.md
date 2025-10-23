# 示例 6: 默认值功能

展示如何使用默认值功能，提供两种设置默认值的方式。

## 运行示例

```bash
cd config/examples/06_default_values
go run main.go
```

## 学习内容

1. **WithDefaults** - 使用 Map 设置全局默认值
2. **WithDefaultStruct** - 使用结构体设置默认值
3. **GetXxxOr** - 读取时指定默认值
4. **Exists** - 检查配置是否存在
5. **配置优先级** - 默认值的优先级最低

## 两种设置默认值的方式

### 方式 1: 使用 WithDefaults (推荐用于全局默认值)

```go
defaults := map[string]interface{}{
    "app.name":    "default-app",
    "server.port": 8080,
    "app.timeout": 30,
}

config.Init(
    config.WithDefaults(defaults),
    config.WithFile("config.yaml"),
)

// 如果配置文件中没有 app.timeout，将使用默认值 30
timeout := config.GetInt("app.timeout")
```

### 方式 2: 使用 GetXxxOr (推荐用于局部默认值)

```go
config.Init(config.WithFile("config.yaml"))

// 读取时指定默认值
logFile := config.GetStringOr("log.file", "/var/log/app.log")
maxRetry := config.GetIntOr("http.max_retry", 3)
enableCache := config.GetBoolOr("cache.enabled", true)
```

## 配置优先级

从低到高：
1. **默认值** (WithDefaults) - 最低优先级
2. **配置文件** (WithFile)
3. **环境变量** (WithEnv)
4. **远程配置** (WithRemote) - 最高优先级

## 实际应用场景

### 场景 1: 应用启动时提供合理的默认值

```go
defaults := map[string]interface{}{
    "server.host":        "0.0.0.0",
    "server.port":        8080,
    "database.max_conns": 100,
    "database.min_conns": 10,
    "log.level":          "info",
}

config.Init(
    config.WithDefaults(defaults),
    config.WithFile("config.yaml"),  // 可选的配置文件
)

// 即使没有配置文件，应用也能正常启动
```

### 场景 2: 平滑升级，新增配置有默认值

```go
// v1.0 版本的配置
defaults := map[string]interface{}{
    "server.port": 8080,
}

// v2.0 新增功能，增加新配置
defaults["cache.enabled"] = true
defaults["cache.ttl"] = 300

config.Init(
    config.WithDefaults(defaults),
    config.WithFile("config.yaml"),  // 老版本配置文件
)

// 即使老配置文件中没有 cache 相关配置，新功能也能正常工作
```

### 场景 3: 简化配置文件

```yaml
# 只需配置非默认值
app:
  name: my-app

server:
  port: 9090  # 覆盖默认的 8080

# 其他配置都使用默认值
```

## GetXxxOr 方法列表

- `GetStringOr(path, defaultValue)` - 读取字符串，不存在返回默认值
- `GetIntOr(path, defaultValue)` - 读取整数，不存在返回默认值
- `GetBoolOr(path, defaultValue)` - 读取布尔值，不存在返回默认值
- `GetFloat64Or(path, defaultValue)` - 读取浮点数，不存在返回默认值
- `GetStringSliceOr(path, defaultValue)` - 读取字符串数组，不存在返回默认值

## Exists 方法

检查配置键是否存在：

```go
if config.Exists("app.name") {
    name := config.GetString("app.name")
} else {
    name := "default-name"
}

// 或者直接使用 GetStringOr
name := config.GetStringOr("app.name", "default-name")
```

## 最佳实践

### ✅ 推荐做法

**1. 全局默认值使用 WithDefaults**
```go
defaults := map[string]interface{}{
    "server.port": 8080,
    "log.level":   "info",
}
config.Init(config.WithDefaults(defaults))
```

**2. 局部默认值使用 GetXxxOr**
```go
// 某些可选的配置项
maxRetry := config.GetIntOr("http.max_retry", 3)
timeout := config.GetIntOr("http.timeout", 30)
```

**3. 结合使用**
```go
// 基础默认值
defaults := map[string]interface{}{
    "server.port": 8080,
    "log.level":   "info",
}

config.Init(
    config.WithDefaults(defaults),
    config.WithFile("config.yaml"),
)

// 特殊情况使用 GetXxxOr
customValue := config.GetStringOr("custom.value", "special-default")
```

### ❌ 避免的做法

**不要在多处重复设置相同的默认值**
```go
// ❌ 不好
port1 := config.GetIntOr("server.port", 8080)
port2 := config.GetIntOr("server.port", 8080)  // 重复

// ✅ 好
defaults := map[string]interface{}{
    "server.port": 8080,  // 统一管理
}
```

## 注意事项

1. **默认值优先级最低** - 任何其他配置源都会覆盖默认值
2. **类型要匹配** - 默认值的类型应该与实际使用时一致
3. **文档化默认值** - 在文档中说明各配置项的默认值
4. **合理的默认值** - 默认值应该是安全的、适用于大多数场景的

## 完整示例

```go
package main

import (
    "github.com/silin/go-pkg-sdk/config"
)

func main() {
    // 设置默认值
    defaults := map[string]interface{}{
        "app.name":           "my-app",
        "server.host":        "0.0.0.0",
        "server.port":        8080,
        "database.max_conns": 100,
        "log.level":          "info",
    }

    // 初始化配置
    config.Init(
        config.WithDefaults(defaults),
        config.WithFile("config.yaml"),
        config.WithEnv("APP_"),
    )

    // 读取配置（会按优先级使用默认值）
    port := config.GetInt("server.port")
    
    // 读取可选配置，指定默认值
    timeout := config.GetIntOr("server.timeout", 30)
    
    // 检查配置是否存在
    if config.Exists("feature.new_feature") {
        enabled := config.GetBool("feature.new_feature")
    }
}
```

