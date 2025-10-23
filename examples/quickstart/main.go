package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Si40Code/go-pkg-sdk/config"
)

func main() {
	fmt.Println("=== go-pkg-sdk 快速开始示例 ===\n")

	// 设置一些环境变量（模拟实际环境）
	os.Setenv("APP_APP_ENV", "production")

	// 步骤 1: 初始化配置
	fmt.Println("📦 步骤 1: 初始化配置模块")
	if err := config.Init(
		config.WithFile("config.yaml"),
		config.WithEnv("APP_"),
		config.WithFileWatcher(),
	); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	fmt.Println("✅ 配置模块初始化成功\n")

	// 步骤 2: 读取配置
	fmt.Println("📖 步骤 2: 读取配置")
	appName := config.GetString("app.name")
	appEnv := config.GetString("app.env")
	serverPort := config.GetInt("server.port")

	fmt.Printf("  应用名称: %s\n", appName)
	fmt.Printf("  运行环境: %s\n", appEnv)
	fmt.Printf("  服务端口: %d\n\n", serverPort)

	// 步骤 3: 使用配置初始化应用组件
	fmt.Println("🔧 步骤 3: 使用配置初始化应用组件")

	// 模拟初始化数据库
	dbConfig := initDatabase()
	fmt.Printf("  ✅ 数据库已连接: %s:%d\n", dbConfig.Host, dbConfig.Port)

	// 模拟初始化缓存
	cacheConfig := initCache()
	fmt.Printf("  ✅ 缓存已初始化: TTL=%ds\n", cacheConfig.TTL)

	// 模拟初始化日志
	logLevel := config.GetString("log.level")
	fmt.Printf("  ✅ 日志系统已启动: level=%s\n\n", logLevel)

	// 步骤 4: 注册配置变更监听
	fmt.Println("👂 步骤 4: 注册配置变更监听")
	config.OnChange(func() {
		fmt.Println("\n🔔 检测到配置变更！")
		fmt.Printf("  新的日志级别: %s\n", config.GetString("log.level"))
		fmt.Printf("  新的端口: %d\n", config.GetInt("server.port"))
	})
	fmt.Println("  ✅ 配置变更监听已启动\n")

	// 步骤 5: 应用启动
	fmt.Println("🚀 步骤 5: 应用启动")
	fmt.Printf("  应用 %s 正在运行在端口 %d\n", appName, serverPort)
	fmt.Printf("  环境: %s\n", appEnv)
	fmt.Println()

	// 模拟应用运行
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 提示:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("这是一个快速开始示例，展示了:")
	fmt.Println("  • 如何初始化配置模块")
	fmt.Println("  • 如何读取配置")
	fmt.Println("  • 如何使用配置初始化应用组件")
	fmt.Println("  • 如何监听配置变更")
	fmt.Println()
	fmt.Println("更多示例请查看:")
	fmt.Println("  • config/examples/ - 配置模块的详细示例")
	fmt.Println("  • README.md - 完整文档")
	fmt.Println("  • ARCHITECTURE.md - 架构设计")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 模拟应用运行一段时间
	fmt.Println("⏳ 应用运行中... (5秒后退出)")
	time.Sleep(5 * time.Second)

	fmt.Println("\n👋 应用正常退出")
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Database string
}

// initDatabase 初始化数据库
func initDatabase() DatabaseConfig {
	return DatabaseConfig{
		Host:     config.GetString("database.host"),
		Port:     config.GetInt("database.port"),
		Username: config.GetString("database.username"),
		Database: config.GetString("database.database"),
	}
}

// CacheConfig 缓存配置
type CacheConfig struct {
	TTL int
}

// initCache 初始化缓存
func initCache() CacheConfig {
	return CacheConfig{
		TTL: config.GetInt("cache.ttl"),
	}
}
