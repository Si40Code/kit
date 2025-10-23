package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/Si40Code/go-pkg-sdk/config"
)

func main() {
	fmt.Println("=== 远程配置（Apollo）示例 ===\n")

	// 创建 Apollo 配置提供者
	apolloProvider := NewMockApolloProvider(&ApolloConfig{
		AppID:     "example-app",
		Cluster:   "default",
		Namespace: "application",
		ServerURL: "http://apollo-config.example.com",
	})

	fmt.Println("📝 Apollo 配置:")
	fmt.Println("  AppID: example-app")
	fmt.Println("  Cluster: default")
	fmt.Println("  Namespace: application")
	fmt.Println("  Server: http://apollo-config.example.com")
	fmt.Println()

	// 初始化配置：本地文件 + 远程配置
	if err := config.Init(
		config.WithFile("config.yaml"),    // 本地配置作为兜底
		config.WithRemote(apolloProvider), // 远程配置（优先级更高）
	); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	fmt.Println("✅ 配置初始化成功（本地文件 + Apollo 远程配置）\n")

	// 示例 1: 读取配置（远程配置优先）
	fmt.Println("📖 示例 1: 读取配置（远程优先）")
	printCurrentConfig()

	// 示例 2: 监听远程配置变更
	fmt.Println("\n📝 示例 2: 监听远程配置变更")
	config.OnChange(func() {
		fmt.Println("\n🔔 远程配置已更新！")
		fmt.Println("📖 新的配置:")
		printCurrentConfig()
	})
	fmt.Println("  ✅ 监听器已注册\n")

	// 示例 3: 配置分层架构
	fmt.Println("📝 示例 3: 配置分层架构说明")
	fmt.Println("  配置优先级（从低到高）:")
	fmt.Println("    1. 本地配置文件 (config.yaml) - 默认配置/兜底配置")
	fmt.Println("    2. Apollo 远程配置 - 动态配置/运维配置")
	fmt.Println()
	fmt.Println("  💡 优点:")
	fmt.Println("    • 远程配置失败时仍能使用本地配置")
	fmt.Println("    • 可以动态调整配置而不用重启应用")
	fmt.Println("    • 便于运维人员统一管理配置")
	fmt.Println()

	// 使用说明
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 测试说明:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("这是一个模拟的 Apollo 配置示例")
	fmt.Println()
	fmt.Println("在实际使用中:")
	fmt.Println("1. 从 Apollo 服务器加载配置")
	fmt.Println("2. 长轮询监听配置变更")
	fmt.Println("3. 配置变更时自动推送到客户端")
	fmt.Println()
	fmt.Println("本示例会模拟配置推送（每30秒一次）")
	fmt.Println()
	fmt.Println("按 Ctrl+C 退出程序")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("⏳ 监听远程配置变更...\n")

	<-sigChan
	fmt.Println("\n👋 程序退出")
}

func printCurrentConfig() {
	appName := config.GetString("app.name")
	serverPort := config.GetInt("server.port")
	dbHost := config.GetString("database.host")
	logLevel := config.GetString("log.level")
	featureEnabled := config.GetBool("feature.new_feature_enabled")

	fmt.Printf("  应用名称: %s\n", appName)
	fmt.Printf("  服务器端口: %d\n", serverPort)
	fmt.Printf("  数据库地址: %s\n", dbHost)
	fmt.Printf("  日志级别: %s\n", logLevel)
	fmt.Printf("  新功能开关: %v\n", featureEnabled)
}

// ============ Apollo Provider 实现示例 ============

// ApolloConfig Apollo 配置
type ApolloConfig struct {
	AppID     string
	Cluster   string
	Namespace string
	ServerURL string
}

// MockApolloProvider 模拟的 Apollo 配置提供者
type MockApolloProvider struct {
	config *ApolloConfig
}

// NewMockApolloProvider 创建 Apollo 配置提供者
func NewMockApolloProvider(cfg *ApolloConfig) *MockApolloProvider {
	return &MockApolloProvider{config: cfg}
}

// Load 加载远程配置
func (p *MockApolloProvider) Load(ctx context.Context, k *koanf.Koanf) error {
	fmt.Println("  [Apollo] 正在从远程服务器加载配置...")

	// 模拟从 Apollo 加载配置
	// 实际实现中，这里会调用 Apollo SDK 或 HTTP API
	time.Sleep(100 * time.Millisecond)

	fmt.Println("  [Apollo] 配置加载成功")
	return nil
}

// Watch 监听配置变更
func (p *MockApolloProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
	fmt.Println("  [Apollo] 开始监听配置变更...")

	go func() {
		// 模拟 Apollo 的长轮询机制
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count++
				fmt.Printf("\n  [Apollo] 检测到配置变更 (模拟推送 #%d)\n", count)

				// 模拟新的配置数据
				newConfig := map[string]interface{}{
					"server.port":                 9090 + count,
					"log.level":                   "debug",
					"feature.new_feature_enabled": count%2 == 0,
				}

				// 通知配置变更
				onChange(newConfig)
			}
		}
	}()

	return nil
}

// 注意：这只是一个示例实现
// 实际使用 Apollo 时，应该使用官方 SDK:
// import "github.com/apolloconfig/agollo/v4"
//
// type ApolloProvider struct {
//     client agollo.Client
// }
//
// func (p *ApolloProvider) Load(ctx context.Context, k *koanf.Koanf) error {
//     cache := p.client.GetDefaultConfigCache()
//     for key, value := range cache {
//         k.Set(key, value)
//     }
//     return nil
// }
//
// func (p *ApolloProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
//     p.client.AddChangeListener(&apolloListener{onChange: onChange})
//     return nil
// }
