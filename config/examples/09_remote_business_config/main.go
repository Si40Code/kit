package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/Si40Code/kit/config"
)

func main() {
	fmt.Println("=== 远程业务配置加载示例 ===\n")

	// ========================================
	// 示例 1: 本地配置 + 远程业务配置
	// ========================================
	fmt.Println("📝 示例 1: 本地基础配置 + 远程业务配置")
	fmt.Println("  架构:")
	fmt.Println("  • config-base.yaml (本地) - 基础配置")
	fmt.Println("  • SMS 配置 (远程 Apollo) - 业务配置")
	fmt.Println("  • Email 配置 (远程 Apollo) - 业务配置")
	fmt.Println()

	// 创建远程配置提供者
	apolloProvider := NewMockRemoteProvider(map[string]interface{}{
		// SMS 配置（从远程配置中心读取）
		"sms.provider":                       "aliyun",
		"sms.api_key":                        "remote-api-key-123",
		"sms.api_secret":                     "remote-api-secret-456",
		"sms.sign_name":                      "远程配置应用",
		"sms.rate_limit":                     200,
		"sms.templates.verification.code":    "SMS_REMOTE_001",
		"sms.templates.verification.content": "您的验证码是：{code}，5分钟内有效",

		// Email 配置（从远程配置中心读取）
		"email.smtp.host":     "smtp.remote.com",
		"email.smtp.port":     587,
		"email.smtp.username": "remote@example.com",
		"email.from.name":     "远程通知",
		"email.from.address":  "remote@example.com",
	})

	if err := config.Init(
		config.WithFile("config-base.yaml"), // 本地基础配置
		config.WithRemote(apolloProvider),   // 远程业务配置
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	// 读取远程 SMS 配置
	fmt.Println("📖 远程 SMS 配置:")
	fmt.Printf("  Provider: %s\n", config.GetString("sms.provider"))
	fmt.Printf("  API Key: %s\n", config.GetString("sms.api_key"))
	fmt.Printf("  Sign Name: %s\n", config.GetString("sms.sign_name"))
	fmt.Printf("  Rate Limit: %d/分钟\n", config.GetInt("sms.rate_limit"))
	fmt.Println()

	// 读取远程 Email 配置
	fmt.Println("📖 远程 Email 配置:")
	fmt.Printf("  SMTP Host: %s\n", config.GetString("email.smtp.host"))
	fmt.Printf("  SMTP Port: %d\n", config.GetInt("email.smtp.port"))
	fmt.Printf("  From Name: %s\n", config.GetString("email.from.name"))
	fmt.Println()

	// ========================================
	// 示例 2: 本地兜底 + 远程覆盖
	// ========================================
	fmt.Println("📝 示例 2: 本地配置兜底 + 远程配置覆盖")
	fmt.Println("  策略:")
	fmt.Println("  1. 本地 config-sms.yaml (兜底配置)")
	fmt.Println("  2. 远程 Apollo SMS 配置 (覆盖)")
	fmt.Println("  优点: 远程配置不可用时，仍可使用本地配置")
	fmt.Println()

	apolloProvider2 := NewMockRemoteProvider(map[string]interface{}{
		// 远程只配置需要动态调整的部分
		"sms.rate_limit":         300, // 动态调整限流
		"sms.strategy.max_retry": 5,   // 动态调整重试
	})

	if err := config.Init(
		config.WithFile("config-base.yaml"),
		config.WithFile("config-sms.yaml"), // 本地 SMS 配置（兜底）
		config.WithRemote(apolloProvider2), // 远程覆盖部分配置
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	fmt.Println("📖 合并后的配置:")
	fmt.Printf("  Provider: %s (来自本地)\n", config.GetString("sms.provider"))
	fmt.Printf("  API Key: %s (来自本地)\n", config.GetString("sms.api_key"))
	fmt.Printf("  Rate Limit: %d (来自远程，覆盖了本地)\n", config.GetInt("sms.rate_limit"))
	fmt.Println()

	// ========================================
	// 示例 3: 按命名空间加载业务配置
	// ========================================
	fmt.Println("📝 示例 3: 按命名空间加载业务配置")
	fmt.Println("  Apollo 命名空间:")
	fmt.Println("  • application (基础配置)")
	fmt.Println("  • sms (短信业务配置)")
	fmt.Println("  • email (邮件业务配置)")
	fmt.Println("  • payment (支付业务配置)")
	fmt.Println()

	// 模拟从不同命名空间加载配置
	// 注意：实际使用时，可以创建多个 RemoteProvider
	// 这里简化演示，使用一个 provider 包含所有命名空间的配置
	allNamespaces := NewMockRemoteProvider(map[string]interface{}{
		// application 命名空间
		"app.name": "remote-app",

		// sms 命名空间
		"sms.provider":   "aliyun",
		"sms.api_key":    "namespace-sms-key",
		"sms.rate_limit": 150,

		// email 命名空间
		"email.smtp.host": "smtp.namespace.com",
		"email.smtp.port": 587,

		// payment 命名空间
		"payment.gateway": "stripe",
		"payment.timeout": 60,
	})

	if err := config.Init(
		config.WithFile("config-base.yaml"),
		config.WithRemote(allNamespaces),
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")
	fmt.Println("📖 各命名空间配置:")
	fmt.Printf("  [application] app.name: %s\n", config.GetString("app.name"))
	fmt.Printf("  [sms] provider: %s\n", config.GetString("sms.provider"))
	fmt.Printf("  [email] smtp.host: %s\n", config.GetString("email.smtp.host"))
	fmt.Printf("  [payment] gateway: %s\n", config.GetString("payment.gateway"))
	fmt.Println()

	// ========================================
	// 示例 4: 远程配置热更新
	// ========================================
	fmt.Println("📝 示例 4: 远程配置热更新")
	fmt.Println("  模拟 Apollo 推送配置变更")
	fmt.Println()

	dynamicProvider := NewDynamicRemoteProvider()
	dynamicProvider.SetConfig(map[string]interface{}{
		"sms.rate_limit": 100,
		"sms.enabled":    true,
	})

	if err := config.Init(
		config.WithFile("config-base.yaml"),
		config.WithRemote(dynamicProvider),
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 注册变更回调
	config.OnChange(func() {
		fmt.Println("  🔔 检测到配置变更")
		fmt.Printf("    SMS Rate Limit: %d\n", config.GetInt("sms.rate_limit"))
		fmt.Printf("    SMS Enabled: %v\n", config.GetBool("sms.enabled"))
	})

	fmt.Println("✅ 配置加载成功")
	fmt.Printf("  初始 Rate Limit: %d\n\n", config.GetInt("sms.rate_limit"))

	// 模拟远程配置推送
	fmt.Println("  ⏳ 3秒后模拟配置推送...\n")
	time.Sleep(3 * time.Second)

	fmt.Println("  📡 Apollo 推送新配置:")
	fmt.Println("    sms.rate_limit: 100 → 200")
	fmt.Println("    sms.enabled: true → false")
	fmt.Println()

	dynamicProvider.UpdateConfig(map[string]interface{}{
		"sms.rate_limit": 200,
		"sms.enabled":    false,
	})

	time.Sleep(1 * time.Second)

	fmt.Printf("  ✅ 配置已更新: Rate Limit = %d\n\n", config.GetInt("sms.rate_limit"))

	// ========================================
	// 实际应用场景
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 实际应用场景:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Println("场景 1: Apollo 多命名空间")
	fmt.Println("```go")
	fmt.Println("// Apollo 配置")
	fmt.Println("apolloConfig := &ApolloConfig{")
	fmt.Println("    AppID:     \"your-app-id\",")
	fmt.Println("    Cluster:   \"default\",")
	fmt.Println("    Namespaces: []string{")
	fmt.Println("        \"application\",  // 基础配置")
	fmt.Println("        \"sms\",          // SMS 业务配置")
	fmt.Println("        \"email\",        // Email 业务配置")
	fmt.Println("        \"payment\",      // Payment 业务配置")
	fmt.Println("    },")
	fmt.Println("    ServerURL: \"http://apollo.example.com\",")
	fmt.Println("}")
	fmt.Println()
	fmt.Println("apolloProvider := NewApolloProvider(apolloConfig)")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithFile(\"config-base.yaml\"),  // 本地兜底")
	fmt.Println("    config.WithRemote(apolloProvider),      // 远程业务配置")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("场景 2: 分层配置策略")
	fmt.Println("```go")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithDefaults(defaults),          // 1. 代码默认值")
	fmt.Println("    config.WithFile(\"config-base.yaml\"),    // 2. 本地基础配置")
	fmt.Println("    config.WithFile(\"config-sms.yaml\"),     // 3. 本地 SMS 配置")
	fmt.Println("    config.WithEnv(\"APP_\"),                 // 4. 环境变量")
	fmt.Println("    config.WithRemote(apolloProvider),      // 5. 远程动态配置")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("场景 3: 动态功能开关")
	fmt.Println("  运维人员在 Apollo 控制台修改:")
	fmt.Println("  • sms.enabled: true → false  (关闭 SMS 功能)")
	fmt.Println("  • sms.rate_limit: 100 → 50   (降低发送频率)")
	fmt.Println("  所有应用实例实时生效，无需重启")
	fmt.Println()

	// ========================================
	// 最佳实践
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ 最佳实践:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("1. 本地配置作为兜底")
	fmt.Println("   • 本地文件包含所有必需配置")
	fmt.Println("   • 远程配置不可用时仍能正常启动")
	fmt.Println()
	fmt.Println("2. 远程配置负责动态部分")
	fmt.Println("   • 限流参数、超时时间等运行时参数")
	fmt.Println("   • 功能开关")
	fmt.Println("   • A/B 测试配置")
	fmt.Println()
	fmt.Println("3. 敏感信息不放远程配置")
	fmt.Println("   • API Key、密码等通过环境变量")
	fmt.Println("   • 或使用加密的配置中心")
	fmt.Println()
	fmt.Println("4. 使用命名空间组织配置")
	fmt.Println("   • application: 基础配置")
	fmt.Println("   • {business}: 业务配置 (sms, email, payment)")
	fmt.Println("   • {env}: 环境配置 (dev, prod)")
	fmt.Println()

	fmt.Println("✨ 所有示例执行完成！")
}

// ============ Mock Remote Provider ============

type MockRemoteProvider struct {
	config map[string]interface{}
}

func NewMockRemoteProvider(cfg map[string]interface{}) *MockRemoteProvider {
	return &MockRemoteProvider{config: cfg}
}

func (p *MockRemoteProvider) Load(ctx context.Context, k *koanf.Koanf) error {
	fmt.Println("  [Remote] 从配置中心加载配置...")
	time.Sleep(100 * time.Millisecond)

	// 将远程配置加载到 koanf
	for key, value := range p.config {
		k.Set(key, value)
	}

	fmt.Printf("  [Remote] 加载了 %d 个配置项\n", len(p.config))
	return nil
}

func (p *MockRemoteProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
	// 模拟监听，不推送
	return nil
}

// ============ Dynamic Remote Provider (支持热更新) ============

type DynamicRemoteProvider struct {
	config     map[string]interface{}
	onChangeFn func(map[string]interface{})
}

func NewDynamicRemoteProvider() *DynamicRemoteProvider {
	return &DynamicRemoteProvider{
		config: make(map[string]interface{}),
	}
}

func (p *DynamicRemoteProvider) SetConfig(cfg map[string]interface{}) {
	p.config = cfg
}

func (p *DynamicRemoteProvider) UpdateConfig(cfg map[string]interface{}) {
	p.config = cfg
	if p.onChangeFn != nil {
		p.onChangeFn(cfg)
	}
}

func (p *DynamicRemoteProvider) Load(ctx context.Context, k *koanf.Koanf) error {
	// 将配置加载到 koanf
	for key, value := range p.config {
		k.Set(key, value)
	}
	return nil
}

func (p *DynamicRemoteProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
	p.onChangeFn = onChange
	return nil
}
