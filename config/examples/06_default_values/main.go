package main

import (
	"fmt"
	"log"

	"github.com/Si40Code/kit/config"
)

// AppConfig 应用配置结构体
type AppConfig struct {
	App struct {
		Name    string `koanf:"name"`
		Version string `koanf:"version"`
		Debug   bool   `koanf:"debug"`
		Timeout int    `koanf:"timeout"`
	} `koanf:"app"`
	Server struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	} `koanf:"server"`
	Database struct {
		Host     string `koanf:"host"`
		Port     int    `koanf:"port"`
		MaxConns int    `koanf:"max_conns"`
		MinConns int    `koanf:"min_conns"`
	} `koanf:"database"`
}

func main() {
	fmt.Println("=== 默认值功能示例 ===\n")

	// ========================================
	// 方式 1: 使用 WithDefaults 设置 Map 默认值
	// ========================================
	fmt.Println("📝 方式 1: 使用 Map 设置默认值")

	defaults := map[string]interface{}{
		"app.name":           "default-app",
		"app.version":        "0.0.1",
		"app.debug":          false,
		"app.timeout":        30,
		"server.host":        "0.0.0.0",
		"server.port":        8080,
		"database.host":      "localhost",
		"database.port":      3306,
		"database.max_conns": 100,
		"database.min_conns": 10,
	}

	if err := config.Init(
		config.WithDefaults(defaults),  // 设置默认值
		config.WithFile("config.yaml"), // 文件配置会覆盖默认值
	); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	fmt.Println("✅ 配置初始化成功（使用 Map 默认值）\n")

	// 示例 1: 配置文件中有的值会覆盖默认值
	fmt.Println("📖 示例 1: 配置文件覆盖默认值")
	fmt.Printf("  默认值: app.name = \"default-app\"\n")
	fmt.Printf("  配置文件: app.name = \"my-app\"\n")
	fmt.Printf("  ✨ 实际值: %s （配置文件覆盖了默认值）\n\n", config.GetString("app.name"))

	// 示例 2: 配置文件中没有的值使用默认值
	fmt.Println("📖 示例 2: 使用默认值")
	fmt.Printf("  默认值: app.timeout = 30\n")
	fmt.Printf("  配置文件: 未设置\n")
	fmt.Printf("  ✨ 实际值: %d （使用默认值）\n\n", config.GetInt("app.timeout"))

	fmt.Println("📖 示例 3: 数据库连接配置")
	fmt.Printf("  database.max_conns = %d （使用默认值）\n", config.GetInt("database.max_conns"))
	fmt.Printf("  database.min_conns = %d （使用默认值）\n\n", config.GetInt("database.min_conns"))

	// ========================================
	// 方式 2: 使用 GetXxxOr 在读取时指定默认值
	// ========================================
	fmt.Println("📝 方式 2: 使用 GetXxxOr 方法指定默认值\n")

	// 示例 4: GetStringOr - 存在的配置
	fmt.Println("📖 示例 4: GetStringOr - 配置存在")
	appName := config.GetStringOr("app.name", "fallback-name")
	fmt.Printf("  app.name 存在，值为: %s\n\n", appName)

	// 示例 5: GetStringOr - 不存在的配置
	fmt.Println("📖 示例 5: GetStringOr - 配置不存在")
	logFile := config.GetStringOr("log.file", "/var/log/app.log")
	fmt.Printf("  log.file 不存在，使用默认值: %s\n\n", logFile)

	// 示例 6: GetIntOr - 不存在的配置
	fmt.Println("📖 示例 6: GetIntOr - 配置不存在")
	maxRetry := config.GetIntOr("http.max_retry", 3)
	fmt.Printf("  http.max_retry 不存在，使用默认值: %d\n\n", maxRetry)

	// 示例 7: GetBoolOr - 不存在的配置
	fmt.Println("📖 示例 7: GetBoolOr - 配置不存在")
	enableCache := config.GetBoolOr("cache.enabled", true)
	fmt.Printf("  cache.enabled 不存在，使用默认值: %v\n\n", enableCache)

	// 示例 8: GetFloat64Or - 不存在的配置
	fmt.Println("📖 示例 8: GetFloat64Or - 配置不存在")
	cacheRatio := config.GetFloat64Or("cache.ratio", 0.75)
	fmt.Printf("  cache.ratio 不存在，使用默认值: %.2f\n\n", cacheRatio)

	// 示例 9: GetStringSliceOr - 不存在的配置
	fmt.Println("📖 示例 9: GetStringSliceOr - 配置不存在")
	allowedIPs := config.GetStringSliceOr("security.allowed_ips", []string{"127.0.0.1", "::1"})
	fmt.Printf("  security.allowed_ips 不存在，使用默认值: %v\n\n", allowedIPs)

	// 示例 10: 检查配置是否存在
	fmt.Println("📖 示例 10: 使用 Exists 检查配置是否存在")
	if config.Exists("app.name") {
		fmt.Println("  ✅ app.name 存在")
	}
	if !config.Exists("not.existing.key") {
		fmt.Println("  ❌ not.existing.key 不存在")
	}
	fmt.Println()

	// ========================================
	// 配置优先级说明
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 配置优先级（从低到高）:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("  1. 默认值 (WithDefaults)")
	fmt.Println("  2. 配置文件 (WithFile)")
	fmt.Println("  3. 环境变量 (WithEnv)")
	fmt.Println("  4. 远程配置 (WithRemote)")
	fmt.Println()
	fmt.Println("  高优先级的配置会覆盖低优先级的配置")
	fmt.Println()

	// ========================================
	// 实际应用场景
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 实际应用场景:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("场景 1: 应用启动时提供合理的默认值")
	fmt.Println("  - 即使没有配置文件，应用也能正常启动")
	fmt.Println("  - 降低配置门槛，开箱即用")
	fmt.Println()
	fmt.Println("场景 2: 平滑升级，新增配置有默认值")
	fmt.Println("  - 新版本增加新配置项")
	fmt.Println("  - 老版本配置文件仍可使用")
	fmt.Println("  - 避免因缺少配置导致的错误")
	fmt.Println()
	fmt.Println("场景 3: 简化配置文件")
	fmt.Println("  - 只需在配置文件中覆盖非默认值")
	fmt.Println("  - 配置文件更简洁易读")
	fmt.Println()

	// ========================================
	// 最佳实践
	// ========================================
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ 最佳实践:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("1. 使用 WithDefaults 设置全局默认值")
	fmt.Println("   适合: 应用的基础配置")
	fmt.Println()
	fmt.Println("2. 使用 GetXxxOr 设置局部默认值")
	fmt.Println("   适合: 某些可选的配置项")
	fmt.Println()
	fmt.Println("3. 结合使用两种方式")
	fmt.Println("   WithDefaults: 提供基础默认值")
	fmt.Println("   GetXxxOr: 处理特殊情况")
	fmt.Println()

	fmt.Println("✨ 所有示例执行完成！")
}
