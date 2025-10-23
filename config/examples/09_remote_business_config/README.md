# 示例 9: 远程业务配置加载

展示如何通过远程配置中心（如 Apollo）加载业务模块配置（如 sms.yaml）。

## 运行示例

```bash
cd config/examples/09_remote_business_config
go run main.go
```

## 学习内容

1. **远程加载业务配置** - 从 Apollo 加载 SMS、Email 等业务配置
2. **本地兜底策略** - 本地配置文件作为兜底
3. **命名空间组织** - 使用 Apollo 命名空间组织业务配置
4. **远程配置热更新** - 动态调整业务参数
5. **配置分层策略** - 本地 + 远程的最佳实践

## 配置架构

### 方式 1: 完全使用远程配置

```go
// 所有业务配置都从 Apollo 读取
apolloProvider := NewApolloProvider(&ApolloConfig{
    AppID:      "your-app-id",
    Cluster:    "default",
    Namespaces: []string{
        "application",  // 基础配置
        "sms",          // SMS 业务配置
        "email",        // Email 业务配置
        "payment",      // Payment 业务配置
    },
    ServerURL: "http://apollo.example.com",
})

config.Init(
    config.WithFile("config-base.yaml"),  // 最基础的兜底
    config.WithRemote(apolloProvider),    // 远程业务配置
)
```

### 方式 2: 本地兜底 + 远程覆盖

```go
// 本地有完整配置，远程只覆盖需要动态调整的部分
config.Init(
    config.WithFile("config-base.yaml"),     // 基础配置
    config.WithFile("config-sms.yaml"),      // 本地 SMS 配置（兜底）
    config.WithFile("config-email.yaml"),    // 本地 Email 配置（兜底）
    config.WithRemote(apolloProvider),       // 远程覆盖动态参数
)
```

### 方式 3: 分层配置策略（推荐）

```go
defaults := map[string]interface{}{
    "sms.rate_limit": 50,
    "sms.timeout":    10,
}

config.Init(
    config.WithDefaults(defaults),           // 1. 代码默认值
    config.WithFile("config-base.yaml"),     // 2. 本地基础配置
    config.WithFile("config-sms.yaml"),      // 3. 本地业务配置
    config.WithEnv("APP_"),                  // 4. 环境变量
    config.WithRemote(apolloProvider),       // 5. 远程动态配置
)
```

## Apollo 配置示例

### Apollo 控制台配置

**命名空间: sms**

```properties
# Apollo 配置格式
sms.provider = aliyun
sms.api_key = apollo-api-key-123
sms.api_secret = apollo-api-secret-456
sms.sign_name = Apollo应用
sms.rate_limit = 200
sms.templates.verification.code = SMS_APOLLO_001
sms.templates.verification.content = 您的验证码是：{code}，5分钟内有效
```

**命名空间: email**

```properties
email.smtp.host = smtp.apollo.com
email.smtp.port = 587
email.smtp.username = apollo@example.com
email.from.name = Apollo通知
email.from.address = apollo@example.com
```

## RemoteProvider 实现示例

### Apollo Provider 实现

```go
package config

import (
    "context"
    "github.com/apolloconfig/agollo/v4"
    "github.com/apolloconfig/agollo/v4/env/config"
)

type ApolloProvider struct {
    client agollo.Client
}

type ApolloConfig struct {
    AppID      string
    Cluster    string
    Namespaces []string
    ServerURL  string
}

func NewApolloProvider(cfg *ApolloConfig) (*ApolloProvider, error) {
    c := &config.AppConfig{
        AppID:         cfg.AppID,
        Cluster:       cfg.Cluster,
        IP:            cfg.ServerURL,
        NamespaceName: strings.Join(cfg.Namespaces, ","),
    }
    
    client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
        return c, nil
    })
    if err != nil {
        return nil, err
    }
    
    return &ApolloProvider{client: client}, nil
}

func (p *ApolloProvider) Load(ctx context.Context, k *koanf.Koanf) error {
    // 从 Apollo 读取所有配置
    for _, namespace := range p.client.GetConfigCache().GetNamespaces() {
        cache := p.client.GetConfigCache(namespace)
        for key, value := range cache.GetContent() {
            k.Set(key, value)
        }
    }
    return nil
}

func (p *ApolloProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
    // 注册 Apollo 监听器
    p.client.AddChangeListener(&apolloChangeListener{
        onChange: onChange,
    })
    return nil
}

type apolloChangeListener struct {
    onChange func(map[string]interface{})
}

func (l *apolloChangeListener) OnChange(event *storage.ChangeEvent) {
    changes := make(map[string]interface{})
    for key, change := range event.Changes {
        changes[key] = change.NewValue
    }
    l.onChange(changes)
}

func (l *apolloChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
    // 处理全量变更
}
```

## 实际应用场景

### 场景 1: 动态限流调整

```go
// 运维人员在 Apollo 控制台修改
// sms.rate_limit: 100 → 200

config.OnChange(func() {
    newLimit := config.GetInt("sms.rate_limit")
    smsService.UpdateRateLimit(newLimit)
    log.Printf("SMS rate limit updated to %d", newLimit)
})
```

### 场景 2: 功能开关

```go
// 运维人员在 Apollo 控制台修改
// sms.enabled: true → false

config.OnChange(func() {
    enabled := config.GetBool("sms.enabled")
    if !enabled {
        smsService.Stop()
        log.Println("SMS service disabled")
    } else {
        smsService.Start()
        log.Println("SMS service enabled")
    }
})
```

### 场景 3: A/B 测试

```go
// 运维人员在 Apollo 为不同实例配置不同的模板
// 实例 A: sms.templates.verification.code = SMS_001
// 实例 B: sms.templates.verification.code = SMS_002

templateCode := config.GetString("sms.templates.verification.code")
smsService.SetTemplate(templateCode)
```

### 场景 4: 灰度发布

```go
// 先在 Apollo 为部分实例配置新功能
// new_feature.enabled = true (只对部分实例可见)

if config.GetBool("new_feature.enabled") {
    // 使用新功能
    newFeatureService.Start()
} else {
    // 使用旧功能
    oldFeatureService.Start()
}
```

## 配置优先级

从低到高：

```
1. 代码默认值 (WithDefaults)
2. 本地基础配置 (config-base.yaml)
3. 本地业务配置 (config-sms.yaml)
4. 环境变量 (WithEnv)
5. 远程配置 (WithRemote)          ← 最高优先级
```

## 配置读取示例

### SMS 配置读取

```go
// 基础配置（可能来自本地或远程）
provider := config.GetString("sms.provider")
apiKey := config.GetString("sms.api_key")

// 动态参数（通常来自远程）
rateLimit := config.GetInt("sms.rate_limit")
enabled := config.GetBool("sms.enabled")

// 模板配置
templateCode := config.GetString("sms.templates.verification.code")
```

### 配置验证

```go
func validateSMSConfig() error {
    if !config.Exists("sms.provider") {
        return fmt.Errorf("SMS provider not configured")
    }
    
    if config.GetString("sms.api_key") == "" {
        return fmt.Errorf("SMS API key is required")
    }
    
    if config.GetInt("sms.rate_limit") <= 0 {
        return fmt.Errorf("SMS rate limit must be positive")
    }
    
    return nil
}
```

## 最佳实践

### 1. 配置分层

```
层次                  | 内容                    | 变更频率
---------------------|------------------------|----------
代码默认值            | 最基础的默认配置         | 从不
本地基础配置          | 应用通用配置            | 很少
本地业务配置          | 业务模块默认配置         | 偶尔
环境变量              | 敏感信息、环境差异       | 部署时
远程配置              | 动态参数、功能开关       | 频繁
```

### 2. 远程配置内容

**✅ 适合放远程配置:**
- 限流参数、超时时间
- 功能开关
- A/B 测试配置
- 灰度发布配置
- 业务规则参数

**❌ 不适合放远程配置:**
- API Key、密码等敏感信息（用环境变量）
- 很少变化的配置
- 必需的基础配置

### 3. 本地兜底策略

```go
// ✅ 推荐：本地有完整配置
config.Init(
    config.WithFile("config-sms.yaml"),   // 本地完整配置
    config.WithRemote(apolloProvider),    // 远程覆盖动态参数
)

// ❌ 不推荐：完全依赖远程配置
config.Init(
    config.WithRemote(apolloProvider),    // 远程配置不可用时无法启动
)
```

### 4. 错误处理

```go
apolloProvider, err := NewApolloProvider(apolloConfig)
if err != nil {
    log.Printf("Failed to connect to Apollo: %v", err)
    log.Println("Using local config as fallback")
    // 降级到本地配置
    config.Init(config.WithFile("config-sms.yaml"))
} else {
    config.Init(
        config.WithFile("config-sms.yaml"),
        config.WithRemote(apolloProvider),
    )
}
```

### 5. 监控和告警

```go
config.OnChange(func() {
    // 记录配置变更
    log.Printf("Config changed at %s", time.Now())
    
    // 关键配置变更告警
    if config.GetInt("sms.rate_limit") < 10 {
        alerting.Send("SMS rate limit too low!")
    }
})
```

## 注意事项

1. **远程配置可用性** - 远程配置不可用时要能降级到本地配置
2. **配置验证** - 远程配置变更后要验证配置的合法性
3. **变更审计** - 记录所有远程配置的变更历史
4. **敏感信息** - 不要在远程配置中存储敏感信息
5. **变更通知** - 重要配置变更要通知相关人员

## 扩展阅读

- [示例 4: 远程配置](../04_remote_config/)
- [示例 8: 业务模块化配置](../08_business_modules/)
- [Config 模块完整文档](../../README.md)

## 常见问题

### Q: Apollo 如何配置多个命名空间？

```go
apolloConfig := &ApolloConfig{
    Namespaces: []string{
        "application",  // 基础配置
        "sms",          // SMS 配置
        "email",        // Email 配置
        "payment",      // Payment 配置
    },
}
```

### Q: 如何只加载某些业务模块的配置？

```go
// 根据功能开关决定加载哪些命名空间
namespaces := []string{"application"}

if os.Getenv("ENABLE_SMS") == "true" {
    namespaces = append(namespaces, "sms")
}
if os.Getenv("ENABLE_EMAIL") == "true" {
    namespaces = append(namespaces, "email")
}

apolloConfig.Namespaces = namespaces
```

### Q: 远程配置失败怎么办？

```go
// 方式 1: 降级到本地配置
apolloProvider, err := NewApolloProvider(apolloConfig)
if err != nil {
    log.Println("Using local config")
    config.Init(config.WithFile("config-sms.yaml"))
    return
}

// 方式 2: 本地配置兜底
config.Init(
    config.WithFile("config-sms.yaml"),   // 兜底
    config.WithRemote(apolloProvider),    // 尽力而为
)
```

