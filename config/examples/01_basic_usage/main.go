package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/Si40Code/go-pkg-sdk/config"
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

// AppConfig 完整应用配置
type AppConfig struct {
	App struct {
		Name    string `koanf:"name"`
		Version string `koanf:"version"`
		Debug   bool   `koanf:"debug"`
	} `koanf:"app"`
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
}

// ValidationError 校验错误
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("配置校验失败 [%s]: %s", e.Field, e.Message)
}

// ConfigValidator 配置校验器
type ConfigValidator struct {
	errors []ValidationError
}

func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		errors: make([]ValidationError, 0),
	}
}

// Required 检查必填字段
func (v *ConfigValidator) Required(path, fieldName string) *ConfigValidator {
	value := config.GetString(path)
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fieldName,
			Message: "不能为空",
		})
	}
	return v
}

// RequiredInt 检查必填整数字段
func (v *ConfigValidator) RequiredInt(path, fieldName string) *ConfigValidator {
	if !config.Exists(path) {
		v.errors = append(v.errors, ValidationError{
			Field:   fieldName,
			Message: "不能为空",
		})
	}
	return v
}

// Email 检查邮箱格式
func (v *ConfigValidator) Email(path, fieldName string) *ConfigValidator {
	email := config.GetString(path)
	if email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			v.errors = append(v.errors, ValidationError{
				Field:   fieldName,
				Message: "邮箱格式不正确",
			})
		}
	}
	return v
}

// URL 检查URL格式
func (v *ConfigValidator) URL(path, fieldName string) *ConfigValidator {
	url := config.GetString(path)
	if url != "" {
		urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
		if !urlRegex.MatchString(url) {
			v.errors = append(v.errors, ValidationError{
				Field:   fieldName,
				Message: "URL格式不正确",
			})
		}
	}
	return v
}

// Port 检查端口号范围
func (v *ConfigValidator) Port(path, fieldName string) *ConfigValidator {
	port := config.GetInt(path)
	if port < 1 || port > 65535 {
		v.errors = append(v.errors, ValidationError{
			Field:   fieldName,
			Message: "端口号必须在 1-65535 范围内",
		})
	}
	return v
}

// Host 检查主机地址格式
func (v *ConfigValidator) Host(path, fieldName string) *ConfigValidator {
	host := config.GetString(path)
	if host != "" {
		// 检查是否是有效的IP地址或主机名
		if net.ParseIP(host) == nil {
			// 如果不是IP，检查是否是有效的主机名
			if _, err := net.LookupHost(host); err != nil {
				v.errors = append(v.errors, ValidationError{
					Field:   fieldName,
					Message: "主机地址格式不正确",
				})
			}
		}
	}
	return v
}

// In 检查值是否在指定范围内
func (v *ConfigValidator) In(path, fieldName string, allowedValues []string) *ConfigValidator {
	value := config.GetString(path)
	if value != "" {
		found := false
		for _, allowed := range allowedValues {
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			v.errors = append(v.errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("值必须是以下之一: %s", strings.Join(allowedValues, ", ")),
			})
		}
	}
	return v
}

// MinLength 检查最小长度
func (v *ConfigValidator) MinLength(path, fieldName string, minLen int) *ConfigValidator {
	value := config.GetString(path)
	if len(value) < minLen {
		v.errors = append(v.errors, ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("长度不能少于 %d 个字符", minLen),
		})
	}
	return v
}

// MaxLength 检查最大长度
func (v *ConfigValidator) MaxLength(path, fieldName string, maxLen int) *ConfigValidator {
	value := config.GetString(path)
	if len(value) > maxLen {
		v.errors = append(v.errors, ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("长度不能超过 %d 个字符", maxLen),
		})
	}
	return v
}

// Validate 执行所有校验并返回错误
func (v *ConfigValidator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}

	var errorMessages []string
	for _, err := range v.errors {
		errorMessages = append(errorMessages, err.Error())
	}
	return fmt.Errorf("配置校验失败:\n%s", strings.Join(errorMessages, "\n"))
}

func main() {
	// 支持通过命令行参数指定配置格式
	format := flag.String("format", "yaml", "配置文件格式 (yaml/json/toml)")
	flag.Parse()

	fmt.Println("=== Config 基础用法示例 ===")
	fmt.Println()

	// 根据格式选择配置文件
	configFile := fmt.Sprintf("config.%s", *format)
	fmt.Printf("📄 使用配置文件: %s\n\n", configFile)

	// 初始化配置
	if err := config.Init(config.WithFile(configFile)); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	fmt.Println("✅ 配置初始化成功")
	fmt.Println()

	// 示例 1: 读取字符串配置
	fmt.Println("📖 示例 1: 读取字符串配置")
	appName := config.GetString("app.name")
	appVersion := config.GetString("app.version")
	fmt.Printf("  应用名称: %s\n", appName)
	fmt.Printf("  应用版本: %s\n\n", appVersion)

	// 示例 2: 读取整数配置
	fmt.Println("📖 示例 2: 读取整数配置")
	serverPort := config.GetInt("server.port")
	dbPort := config.GetInt("database.port")
	fmt.Printf("  服务器端口: %d\n", serverPort)
	fmt.Printf("  数据库端口: %d\n\n", dbPort)

	// 示例 3: 读取布尔配置
	fmt.Println("📖 示例 3: 读取布尔配置")
	debug := config.GetBool("app.debug")
	fmt.Printf("  调试模式: %v\n\n", debug)

	// 示例 4: 读取浮点数配置
	fmt.Println("📖 示例 4: 读取浮点数配置")
	timeout := config.GetFloat64("server.timeout")
	fmt.Printf("  服务器超时: %.1f 秒\n\n", timeout)

	// 示例 5: 读取字符串数组配置
	fmt.Println("📖 示例 5: 读取字符串数组配置")
	allowedHosts := config.GetStringSlice("server.allowed_hosts")
	fmt.Printf("  允许的主机列表: %v\n\n", allowedHosts)

	// 示例 6: 结构化读取 - 读取数据库配置
	fmt.Println("📖 示例 6: 结构化读取（Unmarshal）")
	var dbConfig DatabaseConfig
	if err := config.Unmarshal("database", &dbConfig); err != nil {
		log.Fatalf("解析数据库配置失败: %v", err)
	}
	fmt.Printf("  数据库配置:\n")
	fmt.Printf("    主机: %s\n", dbConfig.Host)
	fmt.Printf("    端口: %d\n", dbConfig.Port)
	fmt.Printf("    用户名: %s\n", dbConfig.Username)
	fmt.Printf("    数据库: %s\n\n", dbConfig.Database)

	// 示例 7: 读取整个配置
	fmt.Println("📖 示例 7: 读取完整配置")
	var appConfig AppConfig
	if err := config.Unmarshal("", &appConfig); err != nil {
		log.Fatalf("解析完整配置失败: %v", err)
	}
	fmt.Printf("  完整配置:\n")
	fmt.Printf("    应用名: %s\n", appConfig.App.Name)
	fmt.Printf("    服务器端口: %d\n", appConfig.Server.Port)
	fmt.Printf("    数据库主机: %s\n\n", appConfig.Database.Host)

	// 示例 8: 读取嵌套配置
	fmt.Println("📖 示例 8: 读取嵌套配置")
	dbHost := config.GetString("database.host")
	logLevel := config.GetString("log.level")
	fmt.Printf("  数据库主机: %s\n", dbHost)
	fmt.Printf("  日志级别: %s\n\n", logLevel)

	// 示例 9: 配置校验
	fmt.Println("📖 示例 9: 配置校验")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 创建校验器
	validator := NewConfigValidator()

	// 执行各种校验规则
	validator.
		Required("app.name", "应用名称").
		Required("database.host", "数据库主机").
		Required("database.username", "数据库用户名").
		Required("database.password", "数据库密码").
		RequiredInt("server.port", "服务器端口").
		RequiredInt("database.port", "数据库端口").
		Port("server.port", "服务器端口").
		Port("database.port", "数据库端口").
		Host("server.host", "服务器主机").
		Host("database.host", "数据库主机").
		In("log.level", "日志级别", []string{"debug", "info", "warn", "error"}).
		In("log.format", "日志格式", []string{"json", "text", "console"}).
		MinLength("app.name", "应用名称", 3).
		MaxLength("app.name", "应用名称", 50)

	// 如果有邮箱和URL配置，也进行校验
	if config.Exists("contact.email") {
		validator.Email("contact.email", "联系邮箱")
	}
	if config.Exists("api.base_url") {
		validator.URL("api.base_url", "API基础URL")
	}

	// 执行校验
	if err := validator.Validate(); err != nil {
		fmt.Printf("❌ 配置校验失败:\n%s\n\n", err)
	} else {
		fmt.Println("✅ 所有配置校验通过！")
		fmt.Println()
	}

	// 示例 10: 演示校验失败的情况
	fmt.Println("📖 示例 10: 演示校验失败情况")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 创建一个会失败的校验示例
	failValidator := NewConfigValidator()
	failValidator.
		Required("nonexistent.field", "不存在的字段").
		In("log.level", "日志级别", []string{"only_debug"}).
		Port("server.port", "服务器端口"). // 这个应该会通过
		Email("app.name", "应用名称")     // 这个会失败，因为应用名称不是邮箱格式

	if err := failValidator.Validate(); err != nil {
		fmt.Printf("❌ 预期的校验失败:\n%s\n\n", err)
	}

	fmt.Println("✨ 所有示例执行完成！")
	fmt.Println("\n💡 配置校验提示:")
	fmt.Println("   - 使用 NewConfigValidator() 创建校验器")
	fmt.Println("   - 链式调用各种校验方法")
	fmt.Println("   - 最后调用 Validate() 执行所有校验")
	fmt.Println("   - 支持必填字段、格式校验、范围校验等")
}
