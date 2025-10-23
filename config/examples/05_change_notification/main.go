package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/silin/go-pkg-sdk/config"
)

func main() {
	fmt.Println("=== 配置变更通知和日志示例 ===\n")

	// 创建测试配置文件
	createTestConfigFile("config.yaml")

	// 初始化配置 - 启用文件监控
	if err := config.Init(
		config.WithFile("config.yaml"),
		config.WithFileWatcher(),
	); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	fmt.Println("✅ 配置初始化成功\n")

	// 示例 1: 单个组件监听配置变更
	fmt.Println("📝 示例 1: 单个组件监听配置")
	setupDatabaseListener()
	fmt.Println()

	// 示例 2: 多个组件同时监听
	fmt.Println("📝 示例 2: 多个组件同时监听")
	setupHTTPServerListener()
	setupCacheListener()
	setupLoggerListener()
	fmt.Println()

	// 示例 3: 配置变更日志
	fmt.Println("📝 示例 3: 配置变更日志")
	fmt.Println("  配置变更时会自动输出 JSON 格式的日志")
	fmt.Println("  日志格式:")
	fmt.Println(`  {`)
	fmt.Println(`    "type": "config_change",`)
	fmt.Println(`    "source": "file",`)
	fmt.Println(`    "key": "server.port",`)
	fmt.Println(`    "old": "8080",`)
	fmt.Println(`    "new": "9090",`)
	fmt.Println(`    "change": "UPDATE",`)
	fmt.Println(`    "timestamp": "2024-01-01T12:00:00Z"`)
	fmt.Println(`  }`)
	fmt.Println()

	// 示例 4: 敏感信息脱敏
	fmt.Println("📝 示例 4: 敏感信息自动脱敏")
	fmt.Println("  包含以下关键词的配置会自动脱敏:")
	fmt.Println("    • password")
	fmt.Println("    • secret")
	fmt.Println("    • token")
	fmt.Println("    • key")
	fmt.Println()
	fmt.Println("  例如:")
	fmt.Println(`  {`)
	fmt.Println(`    "key": "database.password",`)
	fmt.Println(`    "old": "******",  ← 自动脱敏`)
	fmt.Println(`    "new": "******"   ← 自动脱敏`)
	fmt.Println(`  }`)
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 测试说明:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("程序会自动修改配置文件来演示配置变更通知")
	fmt.Println()
	fmt.Println("观察以下内容:")
	fmt.Println("1. 各个组件收到的变更通知")
	fmt.Println("2. 控制台输出的配置变更日志（JSON 格式）")
	fmt.Println("3. 敏感信息（password）的脱敏处理")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 自动修改配置文件进行演示
	fmt.Println("⏳ 等待 3 秒后开始演示...\n")
	time.Sleep(3 * time.Second)

	// 第一次修改
	fmt.Println("📝 第一次修改配置...")
	updateConfig1()
	time.Sleep(3 * time.Second)

	// 第二次修改
	fmt.Println("\n📝 第二次修改配置...")
	updateConfig2()
	time.Sleep(3 * time.Second)

	// 第三次修改（包含敏感信息）
	fmt.Println("\n📝 第三次修改配置（包含敏感信息）...")
	updateConfig3()
	time.Sleep(2 * time.Second)

	fmt.Println("\n✨ 演示完成！")
	fmt.Println("\n💡 提示:")
	fmt.Println("  • 配置变更日志会记录所有变更历史")
	fmt.Println("  • 可以将日志输出到文件或日志系统进行审计")
	fmt.Println("  • 敏感信息会自动脱敏，保护安全")
}

// 数据库组件监听配置
func setupDatabaseListener() {
	config.OnChange(func() {
		host := config.GetString("database.host")
		port := config.GetInt("database.port")
		fmt.Printf("  [数据库连接池] 配置已更新: %s:%d\n", host, port)
	})
	fmt.Println("  ✅ 数据库组件已注册监听器")
}

// HTTP 服务器监听配置
func setupHTTPServerListener() {
	config.OnChange(func() {
		port := config.GetInt("server.port")
		fmt.Printf("  [HTTP 服务器] 端口配置已更新: %d\n", port)
	})
	fmt.Println("  ✅ HTTP 服务器已注册监听器")
}

// 缓存组件监听配置
func setupCacheListener() {
	config.OnChange(func() {
		ttl := config.GetInt("cache.ttl")
		fmt.Printf("  [缓存管理器] TTL 配置已更新: %d 秒\n", ttl)
	})
	fmt.Println("  ✅ 缓存管理器已注册监听器")
}

// 日志组件监听配置
func setupLoggerListener() {
	config.OnChange(func() {
		level := config.GetString("log.level")
		fmt.Printf("  [日志系统] 日志级别已更新: %s\n", level)
	})
	fmt.Println("  ✅ 日志系统已注册监听器")
}

// 创建初始配置文件
func createTestConfigFile(path string) {
	content := `app:
  name: notification-demo
  version: 1.0.0

server:
  host: 0.0.0.0
  port: 8080

database:
  host: localhost
  port: 3306
  password: initial-secret

cache:
  ttl: 300

log:
  level: info
`
	os.WriteFile(path, []byte(content), 0644)
}

// 第一次配置更新
func updateConfig1() {
	content := `app:
  name: notification-demo
  version: 1.0.0

server:
  host: 0.0.0.0
  port: 9090

database:
  host: localhost
  port: 3306
  password: initial-secret

cache:
  ttl: 300

log:
  level: info
`
	os.WriteFile("config.yaml", []byte(content), 0644)
}

// 第二次配置更新
func updateConfig2() {
	content := `app:
  name: notification-demo
  version: 1.0.0

server:
  host: 0.0.0.0
  port: 9090

database:
  host: prod-db.example.com
  port: 3307
  password: initial-secret

cache:
  ttl: 600

log:
  level: debug
`
	os.WriteFile("config.yaml", []byte(content), 0644)
}

// 第三次配置更新（包含敏感信息）
func updateConfig3() {
	content := `app:
  name: notification-demo
  version: 1.0.0

server:
  host: 0.0.0.0
  port: 9090

database:
  host: prod-db.example.com
  port: 3307
  password: new-super-secret-password

cache:
  ttl: 600

log:
  level: debug
`
	os.WriteFile("config.yaml", []byte(content), 0644)
}
