package main

import (
	"fmt"
	"log"
	"os"

	"github.com/silin/go-pkg-sdk/config"
)

func main() {
	fmt.Println("=== 业务模块化配置示例 ===\n")

	// ========================================
	// 示例 1: 基础配置 + 业务模块配置
	// ========================================
	fmt.Println("📝 示例 1: 基础配置 + 业务模块配置")
	fmt.Println("  配置加载顺序:")
	fmt.Println("  1. config-base.yaml     (基础配置)")
	fmt.Println("  2. config-sms.yaml      (短信服务配置)")
	fmt.Println("  3. config-email.yaml    (邮件服务配置)")
	fmt.Println("  4. config-payment.yaml  (支付服务配置)")
	fmt.Println()

	if err := config.Init(
		config.WithFile("config-base.yaml"),
		config.WithFile("config-sms.yaml"),
		config.WithFile("config-email.yaml"),
		config.WithFile("config-payment.yaml"),
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	// 读取各业务模块的配置
	fmt.Println("📖 各业务模块配置:")
	fmt.Printf("  SMS 服务:\n")
	fmt.Printf("    Provider: %s\n", config.GetString("sms.provider"))
	fmt.Printf("    API Key: %s\n", config.GetString("sms.api_key"))
	fmt.Printf("    Rate Limit: %d/分钟\n", config.GetInt("sms.rate_limit"))

	fmt.Printf("  Email 服务:\n")
	fmt.Printf("    SMTP Host: %s\n", config.GetString("email.smtp.host"))
	fmt.Printf("    SMTP Port: %d\n", config.GetInt("email.smtp.port"))
	fmt.Printf("    From: %s\n", config.GetString("email.from"))

	fmt.Printf("  Payment 服务:\n")
	fmt.Printf("    Gateway: %s\n", config.GetString("payment.gateway"))
	fmt.Printf("    Currency: %s\n", config.GetString("payment.currency"))
	fmt.Printf("    Timeout: %d秒\n", config.GetInt("payment.timeout"))
	fmt.Println()

	// ========================================
	// 示例 2: 使用 WithFiles 批量加载业务模块
	// ========================================
	fmt.Println("📝 示例 2: 使用 WithFiles 批量加载业务模块")

	// 根据环境决定加载哪些业务模块
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 基础配置
	configFiles := []string{
		"config-base.yaml",
	}

	// 根据环境添加业务模块
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
			"config-monitoring.yaml", // 生产环境额外监控配置
		)
	case "test":
		configFiles = append(configFiles,
			"config-sms.yaml",
			"config-email.yaml",
			// 测试环境不加载支付配置，使用 mock
		)
	}

	fmt.Printf("  环境: %s\n", env)
	fmt.Printf("  配置文件: %v\n\n", configFiles)

	if err := config.Init(config.WithFiles(configFiles...)); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	// ========================================
	// 示例 3: 条件加载业务模块
	// ========================================
	fmt.Println("📝 示例 3: 条件加载业务模块")
	fmt.Println("  根据功能开关决定是否加载某些业务模块\n")

	// 检查功能开关
	enableSMS := os.Getenv("ENABLE_SMS") == "true"
	enableEmail := os.Getenv("ENABLE_EMAIL") == "true"
	enablePayment := os.Getenv("ENABLE_PAYMENT") == "true"

	fmt.Printf("  功能开关:\n")
	fmt.Printf("    ENABLE_SMS: %v\n", enableSMS)
	fmt.Printf("    ENABLE_EMAIL: %v\n", enableEmail)
	fmt.Printf("    ENABLE_PAYMENT: %v\n", enablePayment)
	fmt.Println()

	configFiles = []string{"config-base.yaml"}

	if enableSMS {
		configFiles = append(configFiles, "config-sms.yaml")
		fmt.Println("  ✅ 加载 SMS 配置")
	}
	if enableEmail {
		configFiles = append(configFiles, "config-email.yaml")
		fmt.Println("  ✅ 加载 Email 配置")
	}
	if enablePayment {
		configFiles = append(configFiles, "config-payment.yaml")
		fmt.Println("  ✅ 加载 Payment 配置")
	}

	fmt.Printf("\n  最终配置文件: %v\n\n", configFiles)

	if err := config.Init(config.WithFiles(configFiles...)); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	// ========================================
	// 示例 4: 业务模块配置验证
	// ========================================
	fmt.Println("📝 示例 4: 业务模块配置验证")
	fmt.Println("  检查各业务模块的必需配置是否存在\n")

	// SMS 配置验证
	if config.Exists("sms.provider") {
		fmt.Printf("  ✅ SMS 配置完整\n")
		fmt.Printf("    Provider: %s\n", config.GetString("sms.provider"))
	} else {
		fmt.Printf("  ❌ SMS 配置缺失\n")
	}

	// Email 配置验证
	if config.Exists("email.smtp.host") && config.Exists("email.smtp.port") {
		fmt.Printf("  ✅ Email 配置完整\n")
		fmt.Printf("    SMTP: %s:%d\n", config.GetString("email.smtp.host"), config.GetInt("email.smtp.port"))
	} else {
		fmt.Printf("  ❌ Email 配置缺失\n")
	}

	// Payment 配置验证
	if config.Exists("payment.gateway") {
		fmt.Printf("  ✅ Payment 配置完整\n")
		fmt.Printf("    Gateway: %s\n", config.GetString("payment.gateway"))
	} else {
		fmt.Printf("  ❌ Payment 配置缺失\n")
	}
	fmt.Println()

	// ========================================
	// 实际应用场景
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 实际应用场景:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Println("场景 1: 微服务架构")
	fmt.Println("```go")
	fmt.Println("// 每个微服务有自己的业务配置")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithFile(\"config-base.yaml\"),")
	fmt.Println("    config.WithFile(\"config-user.yaml\"),      // 用户服务")
	fmt.Println("    config.WithFile(\"config-order.yaml\"),     // 订单服务")
	fmt.Println("    config.WithFile(\"config-inventory.yaml\"), // 库存服务")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("场景 2: 第三方服务集成")
	fmt.Println("```go")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithFile(\"config-base.yaml\"),")
	fmt.Println("    config.WithFile(\"config-sms.yaml\"),       // 短信服务")
	fmt.Println("    config.WithFile(\"config-email.yaml\"),     // 邮件服务")
	fmt.Println("    config.WithFile(\"config-wechat.yaml\"),    // 微信服务")
	fmt.Println("    config.WithFile(\"config-alipay.yaml\"),    // 支付宝")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("场景 3: 功能模块化")
	fmt.Println("```go")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithFile(\"config-base.yaml\"),")
	fmt.Println("    config.WithFile(\"config-auth.yaml\"),      // 认证模块")
	fmt.Println("    config.WithFile(\"config-cache.yaml\"),     // 缓存模块")
	fmt.Println("    config.WithFile(\"config-search.yaml\"),    // 搜索模块")
	fmt.Println("    config.WithFile(\"config-analytics.yaml\"), // 分析模块")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	// ========================================
	// 最佳实践
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ 最佳实践:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("1. 配置文件命名规范")
	fmt.Println("   • config-{module}.yaml")
	fmt.Println("   • 清晰表明业务模块")
	fmt.Println("   • 例如: config-sms.yaml, config-email.yaml")
	fmt.Println()
	fmt.Println("2. 配置结构设计")
	fmt.Println("   • 每个业务模块有独立的配置段")
	fmt.Println("   • 避免配置键冲突")
	fmt.Println("   • 使用嵌套结构组织配置")
	fmt.Println()
	fmt.Println("3. 条件加载策略")
	fmt.Println("   • 根据环境加载不同模块")
	fmt.Println("   • 根据功能开关控制加载")
	fmt.Println("   • 提供配置验证机制")
	fmt.Println()
	fmt.Println("4. 团队协作")
	fmt.Println("   • 每个团队负责自己的业务配置")
	fmt.Println("   • 配置变更影响范围明确")
	fmt.Println("   • 便于独立开发和测试")
	fmt.Println()

	fmt.Println("✨ 所有示例执行完成！")
}
