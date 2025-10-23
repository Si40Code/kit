package main

import (
	"fmt"
	"log"
	"os"

	"github.com/silin/go-pkg-sdk/config"
)

func main() {
	fmt.Println("=== 环境变量覆盖配置示例 ===\n")

	// 设置一些环境变量用于演示
	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_DATABASE_HOST", "prod-db.example.com")
	os.Setenv("APP_DATABASE_PASSWORD", "prod-secret")
	os.Setenv("APP_APP_DEBUG", "false")

	fmt.Println("📝 设置的环境变量:")
	fmt.Println("  APP_SERVER_PORT=9090")
	fmt.Println("  APP_DATABASE_HOST=prod-db.example.com")
	fmt.Println("  APP_DATABASE_PASSWORD=prod-secret")
	fmt.Println("  APP_APP_DEBUG=false")
	fmt.Println()

	// 初始化配置：文件 + 环境变量
	if err := config.Init(
		config.WithFile("config.yaml"),
		config.WithEnv("APP_"), // 环境变量前缀
	); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	fmt.Println("✅ 配置初始化成功（文件 + 环境变量）\n")

	// 示例 1: 环境变量覆盖数值配置
	fmt.Println("📖 示例 1: 环境变量覆盖端口配置")
	serverPort := config.GetInt("server.port")
	fmt.Printf("  配置文件中: server.port = 8080\n")
	fmt.Printf("  环境变量: APP_SERVER_PORT = 9090\n")
	fmt.Printf("  ✨ 实际值: %d （环境变量生效）\n\n", serverPort)

	// 示例 2: 环境变量覆盖字符串配置
	fmt.Println("📖 示例 2: 环境变量覆盖数据库主机")
	dbHost := config.GetString("database.host")
	fmt.Printf("  配置文件中: database.host = localhost\n")
	fmt.Printf("  环境变量: APP_DATABASE_HOST = prod-db.example.com\n")
	fmt.Printf("  ✨ 实际值: %s （环境变量生效）\n\n", dbHost)

	// 示例 3: 环境变量覆盖敏感配置
	fmt.Println("📖 示例 3: 环境变量覆盖敏感信息（推荐做法）")
	dbPassword := config.GetString("database.password")
	fmt.Printf("  配置文件中: database.password = dev-password\n")
	fmt.Printf("  环境变量: APP_DATABASE_PASSWORD = prod-secret\n")
	fmt.Printf("  ✨ 实际值: %s （环境变量生效）\n", dbPassword)
	fmt.Println("  💡 提示: 生产环境的敏感信息应该通过环境变量传递，而不是写在配置文件中\n")

	// 示例 4: 环境变量覆盖布尔配置
	fmt.Println("📖 示例 4: 环境变量覆盖布尔配置")
	debug := config.GetBool("app.debug")
	fmt.Printf("  配置文件中: app.debug = true\n")
	fmt.Printf("  环境变量: APP_APP_DEBUG = false\n")
	fmt.Printf("  ✨ 实际值: %v （环境变量生效）\n\n", debug)

	// 示例 5: 没有环境变量时使用文件配置
	fmt.Println("📖 示例 5: 没有对应环境变量时使用文件配置")
	appName := config.GetString("app.name")
	fmt.Printf("  配置文件中: app.name = %s\n", appName)
	fmt.Printf("  环境变量: 无 APP_APP_NAME\n")
	fmt.Printf("  ✨ 实际值: %s （使用文件配置）\n\n", appName)

	// 示例 6: 实际应用场景 - 根据环境切换配置
	fmt.Println("📖 示例 6: 实际应用场景")
	fmt.Println("  场景: 使用相同的配置文件，通过环境变量区分不同环境")
	fmt.Println()
	fmt.Println("  开发环境:")
	fmt.Println("    export APP_DATABASE_HOST=localhost")
	fmt.Println("    export APP_APP_DEBUG=true")
	fmt.Println()
	fmt.Println("  测试环境:")
	fmt.Println("    export APP_DATABASE_HOST=test-db.example.com")
	fmt.Println("    export APP_APP_DEBUG=true")
	fmt.Println()
	fmt.Println("  生产环境:")
	fmt.Println("    export APP_DATABASE_HOST=prod-db.example.com")
	fmt.Println("    export APP_DATABASE_PASSWORD=<从密钥管理系统获取>")
	fmt.Println("    export APP_APP_DEBUG=false")
	fmt.Println()

	// 优先级说明
	fmt.Println("📋 配置优先级（从低到高）:")
	fmt.Println("  1. 文件配置 (config.yaml)")
	fmt.Println("  2. 环境变量 (APP_*) ⬅️ 优先级更高")
	fmt.Println()

	fmt.Println("✨ 所有示例执行完成！")
	fmt.Println()
	fmt.Println("💡 最佳实践:")
	fmt.Println("  • 在配置文件中设置默认值和开发环境配置")
	fmt.Println("  • 通过环境变量覆盖特定环境的配置")
	fmt.Println("  • 敏感信息（密码、密钥）始终使用环境变量")
	fmt.Println("  • 环境变量命名规范: <PREFIX>_<KEY_PATH>")
	fmt.Println("    例如: APP_DATABASE_HOST -> database.host")
}
