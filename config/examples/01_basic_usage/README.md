# 示例 1: 基础用法

展示如何从配置文件读取各种类型的配置。支持 **YAML**、**JSON**、**TOML** 三种格式。

## 运行示例

```bash
cd config/examples/01_basic_usage

# 使用 YAML 格式（默认）
go run main.go

# 使用 JSON 格式
go run main.go -format json

# 使用 TOML 格式
go run main.go -format toml

# 测试配置校验失败的情况
go run main-invalid.go
```

## 学习内容

1. **初始化配置** - 从文件加载配置
2. **读取字符串** - `GetString()`
3. **读取整数** - `GetInt()`
4. **读取布尔值** - `GetBool()`
5. **读取浮点数** - `GetFloat64()`
6. **读取数组** - `GetStringSlice()`
7. **结构化读取** - `Unmarshal()` 到结构体
8. **读取嵌套配置** - 使用点号访问嵌套字段
9. **配置校验** - 验证配置的完整性和格式正确性

## 配置校验功能

### 校验器使用方式

```go
// 创建校验器
validator := NewConfigValidator()

// 链式调用校验方法
validator.
    Required("app.name", "应用名称").
    Required("database.host", "数据库主机").
    Port("server.port", "服务器端口").
    Email("contact.email", "联系邮箱").
    URL("api.base_url", "API基础URL").
    In("log.level", "日志级别", []string{"debug", "info", "warn", "error"}).
    MinLength("app.name", "应用名称", 3).
    MaxLength("app.name", "应用名称", 50)

// 执行校验
if err := validator.Validate(); err != nil {
    log.Fatal(err)
}
```

### 支持的校验规则

| 校验方法 | 说明 | 示例 |
|---------|------|------|
| `Required()` | 检查必填字段 | `Required("app.name", "应用名称")` |
| `RequiredInt()` | 检查必填整数字段 | `RequiredInt("server.port", "服务器端口")` |
| `Email()` | 检查邮箱格式 | `Email("contact.email", "联系邮箱")` |
| `URL()` | 检查URL格式 | `URL("api.base_url", "API基础URL")` |
| `Port()` | 检查端口号范围(1-65535) | `Port("server.port", "服务器端口")` |
| `Host()` | 检查主机地址格式 | `Host("server.host", "服务器主机")` |
| `In()` | 检查值是否在指定范围内 | `In("log.level", "日志级别", []string{"debug", "info"})` |
| `MinLength()` | 检查最小长度 | `MinLength("app.name", "应用名称", 3)` |
| `MaxLength()` | 检查最大长度 | `MaxLength("app.name", "应用名称", 50)` |

### 校验示例

**正常配置校验通过：**
```
✅ 所有配置校验通过！
```

**配置校验失败：**
```
❌ 配置校验失败:
配置校验失败:
配置校验失败 [应用名称]: 不能为空
配置校验失败 [服务器端口]: 端口号必须在 1-65535 范围内
配置校验失败 [联系邮箱]: 邮箱格式不正确
配置校验失败 [API基础URL]: URL格式不正确
```

## 预期输出

```
=== Config 基础用法示例 ===

✅ 配置初始化成功

📖 示例 1: 读取字符串配置
  应用名称: example-app
  应用版本: 1.0.0

📖 示例 2: 读取整数配置
  服务器端口: 8080
  数据库端口: 3306

📖 示例 3: 读取布尔配置
  调试模式: true

📖 示例 9: 配置校验
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ 所有配置校验通过！

📖 示例 10: 演示校验失败情况
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❌ 预期的校验失败:
配置校验失败:
配置校验失败 [不存在的字段]: 不能为空
配置校验失败 [日志级别]: 值必须是以下之一: only_debug
配置校验失败 [应用名称]: 邮箱格式不正确

💡 配置校验提示:
   - 使用 NewConfigValidator() 创建校验器
   - 链式调用各种校验方法
   - 最后调用 Validate() 执行所有校验
   - 支持必填字段、格式校验、范围校验等
```

## 配置文件说明

示例提供了多种配置文件：

- `config.yaml` - 正常的YAML配置文件
- `config-invalid.yaml` - 包含无效值的配置文件（用于测试校验失败）
- `config.json` - JSON 格式配置文件  
- `config.toml` - TOML 格式配置文件

所有配置文件都包含：

- 应用信息（app）
- 服务器配置（server）
- 数据库配置（database）
- 日志配置（log）
- 联系信息（contact）
- API配置（api）
- 安全配置（security）

## 支持的配置格式

Config 模块会自动根据文件扩展名选择合适的解析器：

| 格式 | 扩展名 | 示例 |
|------|--------|------|
| YAML | `.yaml`, `.yml` | `config.yaml` |
| JSON | `.json` | `config.json` |
| TOML | `.toml` | `config.toml` |

**优势对比：**

- **YAML**: 可读性最好，支持注释，适合人工编辑
- **JSON**: 最通用，易于程序生成和解析
- **TOML**: 结构清晰，配置明确，适合配置文件

**示例：**

```go
// YAML 格式
config.Init(config.WithFile("config.yaml"))

// JSON 格式
config.Init(config.WithFile("config.json"))

// TOML 格式
config.Init(config.WithFile("config.toml"))
```

