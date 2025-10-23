package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/silin/go-pkg-sdk/config"
)

func main() {
	fmt.Println("=== 配置文件监控示例 ===\n")

	// 初始化配置 - 启用文件监控
	if err := config.Init(
		config.WithFile("config.yaml"),
		config.WithFileWatcher(), // 启用文件监控
	); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	fmt.Println("✅ 配置初始化成功（已启用文件监控）\n")

	// 读取初始配置
	fmt.Println("📖 当前配置:")
	printCurrentConfig()

	// 示例 1: 注册配置变更回调
	fmt.Println("\n📝 示例 1: 注册配置变更回调")
	config.OnChange(func() {
		fmt.Println("\n🔔 配置已变更！回调被触发")
		fmt.Println("📖 新的配置:")
		printCurrentConfig()
	})
	fmt.Println("  ✅ 回调已注册\n")

	// 示例 2: 注册多个回调（模拟不同组件监听配置）
	fmt.Println("📝 示例 2: 注册多个回调函数")

	// 数据库连接池监听配置
	config.OnChange(func() {
		fmt.Println("  [数据库连接池] 检测到配置变更，准备重新连接...")
		dbHost := config.GetString("database.host")
		dbPort := config.GetInt("database.port")
		fmt.Printf("  [数据库连接池] 新地址: %s:%d\n", dbHost, dbPort)
	})

	// 日志系统监听配置
	config.OnChange(func() {
		fmt.Println("  [日志系统] 检测到配置变更，更新日志级别...")
		logLevel := config.GetString("log.level")
		fmt.Printf("  [日志系统] 新级别: %s\n", logLevel)
	})

	// HTTP 服务器监听配置
	config.OnChange(func() {
		fmt.Println("  [HTTP 服务器] 检测到配置变更，检查是否需要重启...")
		serverPort := config.GetInt("server.port")
		fmt.Printf("  [HTTP 服务器] 当前端口: %d\n", serverPort)
	})

	fmt.Println("  ✅ 多个回调已注册\n")

	// 使用说明
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 测试说明:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("程序正在监控配置文件变更...")
	fmt.Println()
	fmt.Println("📝 请尝试以下操作:")
	fmt.Println()
	fmt.Println("1. 修改 config.yaml 文件内容")
	fmt.Println("   例如: 将 server.port 从 8080 改为 9090")
	fmt.Println()
	fmt.Println("2. 保存文件")
	fmt.Println()
	fmt.Println("3. 观察控制台输出的配置变更通知")
	fmt.Println()
	fmt.Println("4. 按 Ctrl+C 退出程序")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 模拟应用运行中的配置检查
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("\n⏰ 定期检查当前配置...")
			printCurrentConfig()
		}
	}()

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("⏳ 等待配置文件变更...\n")

	<-sigChan
	fmt.Println("\n👋 程序退出")
}

func printCurrentConfig() {
	appName := config.GetString("app.name")
	serverPort := config.GetInt("server.port")
	dbHost := config.GetString("database.host")
	dbPort := config.GetInt("database.port")
	logLevel := config.GetString("log.level")
	debug := config.GetBool("app.debug")

	fmt.Printf("  应用名称: %s\n", appName)
	fmt.Printf("  服务器端口: %d\n", serverPort)
	fmt.Printf("  数据库地址: %s:%d\n", dbHost, dbPort)
	fmt.Printf("  日志级别: %s\n", logLevel)
	fmt.Printf("  调试模式: %v\n", debug)
}
