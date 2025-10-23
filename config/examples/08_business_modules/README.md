# 示例 8: 业务模块化配置

展示如何将不同业务模块的配置拆分到独立的配置文件中，实现业务配置的模块化管理。

## 运行示例

```bash
cd config/examples/08_business_modules
go run main.go
```

## 学习内容

1. **业务模块拆分** - 将不同业务配置分离到独立文件
2. **条件加载** - 根据环境或功能开关加载不同模块
3. **配置验证** - 检查业务模块配置的完整性
4. **模块化架构** - 便于团队协作和维护

## 配置文件说明

### 文件结构

```
08_business_modules/
├── config-base.yaml          # 基础配置（通用）
├── config-sms.yaml           # 短信服务配置
├── config-email.yaml         # 邮件服务配置
├── config-payment.yaml       # 支付服务配置
├── config-monitoring.yaml    # 监控服务配置（生产环境）
└── main.go
```

### 配置模块说明

| 文件 | 业务模块 | 说明 |
|-----|---------|------|
| config-base.yaml | 基础配置 | 应用通用配置、数据库、日志等 |
| config-sms.yaml | 短信服务 | SMS 提供商配置、模板、策略 |
| config-email.yaml | 邮件服务 | SMTP 配置、模板、队列 |
| config-payment.yaml | 支付服务 | 支付网关、风控、策略 |
| config-monitoring.yaml | 监控服务 | Prometheus、Grafana、告警 |

## 使用方式

### 方式 1: 加载所有业务模块

```go
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-sms.yaml"),
    config.WithFile("config-email.yaml"),
    config.WithFile("config-payment.yaml"),
)
```

### 方式 2: 使用 WithFiles 批量加载

```go
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-sms.yaml",
        "config-email.yaml",
        "config-payment.yaml",
    ),
)
```

### 方式 3: 根据环境条件加载

```go
configFiles := []string{"config-base.yaml"}

env := os.Getenv("APP_ENV")
switch env {
case "dev":
    configFiles = append(configFiles,
        "config-sms.yaml",
        "config-email.yaml",
        "config-payment.yaml",
    )
case "prod":
    configFiles = append(configFiles,
        "config-sms.yaml",
        "config-email.yaml",
        "config-payment.yaml",
        "config-monitoring.yaml",  // 生产环境额外监控
    )
case "test":
    configFiles = append(configFiles,
        "config-sms.yaml",
        "config-email.yaml",
        // 测试环境不加载支付配置
    )
}

config.Init(config.WithFiles(configFiles...))
```

### 方式 4: 根据功能开关加载

```go
configFiles := []string{"config-base.yaml"}

if os.Getenv("ENABLE_SMS") == "true" {
    configFiles = append(configFiles, "config-sms.yaml")
}
if os.Getenv("ENABLE_EMAIL") == "true" {
    configFiles = append(configFiles, "config-email.yaml")
}
if os.Getenv("ENABLE_PAYMENT") == "true" {
    configFiles = append(configFiles, "config-payment.yaml")
}

config.Init(config.WithFiles(configFiles...))
```

## 配置读取示例

### SMS 配置

```go
// 基础配置
provider := config.GetString("sms.provider")
apiKey := config.GetString("sms.api_key")
rateLimit := config.GetInt("sms.rate_limit")

// 模板配置
verifyTemplate := config.GetString("sms.templates.verification.code")
notifyTemplate := config.GetString("sms.templates.notification.code")

// 策略配置
batchSize := config.GetInt("sms.strategy.batch_size")
maxRetry := config.GetInt("sms.strategy.max_retry")
```

### Email 配置

```go
// SMTP 配置
smtpHost := config.GetString("email.smtp.host")
smtpPort := config.GetInt("email.smtp.port")
smtpUser := config.GetString("email.smtp.username")

// 发送者配置
fromName := config.GetString("email.from.name")
fromAddress := config.GetString("email.from.address")

// 模板配置
welcomeSubject := config.GetString("email.templates.welcome.subject")
welcomeTemplate := config.GetString("email.templates.welcome.template")
```

### Payment 配置

```go
// 基础配置
gateway := config.GetString("payment.gateway")
currency := config.GetString("payment.currency")
timeout := config.GetInt("payment.timeout")

// 支付宝配置
alipayAppId := config.GetString("payment.alipay.app_id")
alipayPrivateKey := config.GetString("payment.alipay.private_key")

// 风控配置
maxAmount := config.GetInt("payment.risk_control.max_amount")
dailyLimit := config.GetInt("payment.risk_control.daily_limit")
```

## 实际应用场景

### 场景 1: 微服务架构

每个微服务有自己的业务配置：

```go
// 用户服务
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-user.yaml"),
    config.WithFile("config-auth.yaml"),
)

// 订单服务
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-order.yaml"),
    config.WithFile("config-inventory.yaml"),
)

// 支付服务
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-payment.yaml"),
    config.WithFile("config-notification.yaml"),
)
```

### 场景 2: 第三方服务集成

```go
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-sms.yaml"),       // 短信服务
    config.WithFile("config-email.yaml"),     // 邮件服务
    config.WithFile("config-wechat.yaml"),    // 微信服务
    config.WithFile("config-alipay.yaml"),    // 支付宝
    config.WithFile("config-stripe.yaml"),    // Stripe
)
```

### 场景 3: 功能模块化

```go
config.Init(
    config.WithFile("config-base.yaml"),
    config.WithFile("config-auth.yaml"),      // 认证模块
    config.WithFile("config-cache.yaml"),     // 缓存模块
    config.WithFile("config-search.yaml"),    // 搜索模块
    config.WithFile("config-analytics.yaml"), // 分析模块
    config.WithFile("config-reporting.yaml"), // 报表模块
)
```

### 场景 4: 环境差异化

```go
// 开发环境：加载所有模块
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-sms.yaml",
        "config-email.yaml",
        "config-payment.yaml",
    ),
)

// 生产环境：额外加载监控
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-sms.yaml",
        "config-email.yaml",
        "config-payment.yaml",
        "config-monitoring.yaml",  // 生产环境专用
    ),
)

// 测试环境：只加载必要模块
config.Init(
    config.WithFiles(
        "config-base.yaml",
        "config-sms.yaml",
        "config-email.yaml",
        // 不加载支付和监控配置
    ),
)
```

## 配置验证

### 检查业务模块配置完整性

```go
func validateBusinessConfig() error {
    // SMS 配置验证
    if !config.Exists("sms.provider") {
        return fmt.Errorf("SMS provider not configured")
    }
    
    // Email 配置验证
    if !config.Exists("email.smtp.host") {
        return fmt.Errorf("Email SMTP host not configured")
    }
    
    // Payment 配置验证
    if !config.Exists("payment.gateway") {
        return fmt.Errorf("Payment gateway not configured")
    }
    
    return nil
}
```

### 条件性配置验证

```go
func validateEnabledModules() error {
    if config.Exists("sms.provider") {
        // 验证 SMS 配置
        if config.GetString("sms.api_key") == "" {
            return fmt.Errorf("SMS API key is required")
        }
    }
    
    if config.Exists("email.smtp.host") {
        // 验证 Email 配置
        if config.GetString("email.smtp.username") == "" {
            return fmt.Errorf("Email SMTP username is required")
        }
    }
    
    return nil
}
```

## 最佳实践

### 1. 配置文件命名规范

```
✅ 推荐:
- config-base.yaml      (基础配置)
- config-sms.yaml       (短信服务)
- config-email.yaml      (邮件服务)
- config-payment.yaml    (支付服务)
- config-auth.yaml       (认证服务)

❌ 不推荐:
- config1.yaml
- sms_config.yaml
- email-settings.yml
```

### 2. 配置结构设计

```yaml
# 每个业务模块有独立的配置段
sms:
  provider: aliyun
  api_key: "xxx"
  templates:
    verification:
      code: "SMS_001"
      
email:
  smtp:
    host: smtp.example.com
  templates:
    welcome:
      subject: "欢迎"
```

### 3. 环境差异化

```go
// 根据环境加载不同配置
env := os.Getenv("ENV")
configFiles := []string{"config-base.yaml"}

switch env {
case "dev":
    configFiles = append(configFiles, "config-sms.yaml", "config-email.yaml")
case "prod":
    configFiles = append(configFiles, "config-sms.yaml", "config-email.yaml", "config-monitoring.yaml")
}
```

### 4. 功能开关

```go
// 使用环境变量控制功能模块
if os.Getenv("ENABLE_SMS") == "true" {
    configFiles = append(configFiles, "config-sms.yaml")
}
```

### 5. 团队协作

```
项目/
├── config-base.yaml          # 基础配置（所有团队）
├── config-sms.yaml           # SMS 团队负责
├── config-email.yaml         # Email 团队负责
├── config-payment.yaml       # Payment 团队负责
└── config-monitoring.yaml    # DevOps 团队负责
```

## 注意事项

1. **配置键冲突** - 不同业务模块使用不同的配置前缀
2. **依赖关系** - 某些业务模块可能依赖其他模块的配置
3. **敏感信息** - 业务配置中的敏感信息应使用环境变量
4. **配置验证** - 加载后验证各业务模块配置的完整性
5. **文档化** - 每个业务模块的配置应有清晰的文档说明

## 扩展阅读

- [示例 7: 多配置文件](../07_multiple_files/)
- [示例 6: 默认值功能](../06_default_values/)
- [Config 模块完整文档](../../README.md)

