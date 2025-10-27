package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Si40Code/kit/config"
)

func main() {
	fmt.Println("=== 多配置文件示例 ===\n")

	// ========================================
	// 示例 1: 使用 WithFile 加载多个文件
	// ========================================
	fmt.Println("📝 示例 1: 使用多个 WithFile 加载配置")
	fmt.Println("  配置加载顺序:")
	fmt.Println("  1. config-base.yaml     (基础配置)")
	fmt.Println("  2. config-dev.yaml      (开发环境)")
	fmt.Println("  后加载的配置会覆盖先加载的配置\n")

	if err := config.Init(
		config.WithFile("config-base.yaml"),
		config.WithFile("config-dev.yaml"),
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	fmt.Println("📖 查看配置值:")
	fmt.Printf("  app.name = %s (来自 base)\n", config.GetString("app.name"))
	fmt.Printf("  app.env = %s (来自 dev)\n", config.GetString("app.env"))
	fmt.Printf("  server.port = %d (dev 覆盖了 base 的 8080)\n", config.GetInt("server.port"))
	fmt.Printf("  database.host = %s (dev 覆盖了 base 的 localhost)\n", config.GetString("database.host"))
	fmt.Printf("  log.level = %s (来自 base)\n\n", config.GetString("log.level"))

	// ========================================
	// 示例 2: 使用 WithFiles 一次加载多个文件
	// ========================================
	fmt.Println("📝 示例 2: 使用 WithFiles 一次加载多个文件")

	// 根据环境变量决定加载哪个配置
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	configFiles := []string{
		"config-base.yaml",
		fmt.Sprintf("config-%s.yaml", env),
	}

	fmt.Printf("  环境: %s\n", env)
	fmt.Printf("  配置文件: %v\n\n", configFiles)

	if err := config.Init(config.WithFiles(configFiles...)); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功\n")

	// ========================================
	// 示例 3: 分层配置架构
	// ========================================
	fmt.Println("📝 示例 3: 分层配置架构\n")

	fmt.Println("  典型的配置分层:")
	fmt.Println("  ┌─────────────────────────────────────┐")
	fmt.Println("  │ 1. config-base.yaml                 │  ← 基础配置")
	fmt.Println("  │    - 所有环境通用的配置             │")
	fmt.Println("  │    - 默认值                         │")
	fmt.Println("  ├─────────────────────────────────────┤")
	fmt.Println("  │ 2. config-{env}.yaml                │  ← 环境配置")
	fmt.Println("  │    - config-dev.yaml                │")
	fmt.Println("  │    - config-test.yaml               │")
	fmt.Println("  │    - config-prod.yaml               │")
	fmt.Println("  ├─────────────────────────────────────┤")
	fmt.Println("  │ 3. config-local.yaml (可选)         │  ← 本地配置")
	fmt.Println("  │    - 开发者本地覆盖                 │")
	fmt.Println("  │    - 不提交到 Git                   │")
	fmt.Println("  └─────────────────────────────────────┘")
	fmt.Println()

	// ========================================
	// 示例 4: 完整的配置加载策略
	// ========================================
	fmt.Println("📝 示例 4: 完整的配置加载策略\n")

	defaults := map[string]interface{}{
		"server.port": 8080,
		"log.level":   "info",
	}

	configFiles = []string{
		"config-base.yaml",
		"config-dev.yaml",
	}

	// 如果有本地配置文件，也加载它
	if _, err := os.Stat("config-local.yaml"); err == nil {
		configFiles = append(configFiles, "config-local.yaml")
		fmt.Println("  ✅ 检测到 config-local.yaml，将加载本地配置")
	} else {
		fmt.Println("  ℹ️  未检测到 config-local.yaml")
	}

	fmt.Println()
	fmt.Println("  配置优先级（从低到高）:")
	fmt.Println("  1. 默认值 (代码中)")
	fmt.Println("  2. config-base.yaml")
	fmt.Println("  3. config-dev.yaml")
	fmt.Println("  4. config-local.yaml (如果存在)")
	fmt.Println("  5. 环境变量")
	fmt.Println("  6. 远程配置")
	fmt.Println()

	if err := config.Init(
		config.WithDefaults(defaults),
		config.WithFiles(configFiles...),
		config.WithEnv("APP_"),
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 完整配置加载成功\n")

	// ========================================
	// 实际应用场景
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 实际应用场景:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Println("场景 1: 多环境部署")
	fmt.Println("```go")
	fmt.Println("env := os.Getenv(\"ENV\")")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithFile(\"config-base.yaml\"),")
	fmt.Println("    config.WithFile(fmt.Sprintf(\"config-%s.yaml\", env)),")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("场景 2: 功能模块化配置")
	fmt.Println("```go")
	fmt.Println("config.Init(")
	fmt.Println("    config.WithFiles(")
	fmt.Println("        \"config-base.yaml\",")
	fmt.Println("        \"config-database.yaml\",   // 数据库配置")
	fmt.Println("        \"config-redis.yaml\",      // Redis 配置")
	fmt.Println("        \"config-mq.yaml\",         // 消息队列配置")
	fmt.Println("    ),")
	fmt.Println(")")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("场景 3: 团队协作")
	fmt.Println("  • config-base.yaml    → 提交到 Git")
	fmt.Println("  • config-dev.yaml     → 提交到 Git")
	fmt.Println("  • config-local.yaml   → 添加到 .gitignore")
	fmt.Println()
	fmt.Println("  开发者可以在 config-local.yaml 中设置个人配置")
	fmt.Println("  不会影响其他团队成员")
	fmt.Println()

	// ========================================
	// 最佳实践
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ 最佳实践:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("1. 配置分层")
	fmt.Println("   • base - 通用配置")
	fmt.Println("   • env - 环境特定配置")
	fmt.Println("   • local - 本地覆盖（不提交）")
	fmt.Println()
	fmt.Println("2. 配置文件命名")
	fmt.Println("   • 使用统一的命名规范")
	fmt.Println("   • config-{layer}.yaml")
	fmt.Println("   • 清晰表明配置用途")
	fmt.Println()
	fmt.Println("3. Git 管理")
	fmt.Println("   • 提交通用配置和环境模板")
	fmt.Println("   • 忽略本地配置和敏感信息")
	fmt.Println("   • 提供 .example 文件作为参考")
	fmt.Println()
	fmt.Println("4. 文档化")
	fmt.Println("   • 说明每个配置文件的用途")
	fmt.Println("   • 记录配置项的含义和默认值")
	fmt.Println("   • 提供示例配置")
	fmt.Println()

	fmt.Println("✨ 所有示例执行完成！")
}


