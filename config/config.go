package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/Si40Code/go-pkg-sdk/config/provider"
	"github.com/knadh/koanf/v2"
)

var (
	k               *koanf.Koanf
	mu              sync.RWMutex
	changeCallbacks []func()
	lastSnapshot    map[string]interface{} = make(map[string]interface{})
)

func Init(opts ...Option) error {
	options := newOptions(opts...)

	k = koanf.New(".")

	// 1. 加载默认配置（最低优先级）
	if options.defaults != nil {
		if err := provider.LoadDefaults(k, options.defaults); err != nil {
			return fmt.Errorf("load default config failed: %w", err)
		}
	}

	// 2. 加载文件配置（按顺序加载，后面的覆盖前面的）
	for _, filePath := range options.filePaths {
		if err := provider.LoadFile(k, filePath); err != nil {
			return fmt.Errorf("load file config failed (%s): %w", filePath, err)
		}
	}

	// 3. 加载环境变量配置
	if options.useEnv {
		if err := provider.LoadEnv(k, options.envPrefix); err != nil {
			return fmt.Errorf("load env config failed: %w", err)
		}
	}

	// 4. 加载远程配置（最高优先级）
	if options.remoteProvider != nil {
		if err := options.remoteProvider.Load(context.Background(), k); err != nil {
			return fmt.Errorf("load remote config failed: %w", err)
		}
		go options.remoteProvider.Watch(context.Background(), func(newCfg map[string]interface{}) {
			mu.Lock()
			LogConfigDiff("apollo", lastSnapshot, newCfg)
			k.Load(provider.MapProvider(newCfg), nil)
			lastSnapshot = cloneMap(k.Raw())
			mu.Unlock()
			notifyChange()
		})
	}

	// 5. 启动文件监控（监控所有配置文件）
	if options.watchFile {
		for _, filePath := range options.filePaths {
			startWatcher(filePath)
		}
	}

	lastSnapshot = cloneMap(k.Raw())
	return nil
}

func GetString(path string) string {
	mu.RLock()
	defer mu.RUnlock()
	return k.String(path)
}

func GetInt(path string) int {
	mu.RLock()
	defer mu.RUnlock()
	return k.Int(path)
}

func GetBool(path string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return k.Bool(path)
}

func GetFloat64(path string) float64 {
	mu.RLock()
	defer mu.RUnlock()
	return k.Float64(path)
}

func GetStringSlice(path string) []string {
	mu.RLock()
	defer mu.RUnlock()
	return k.Strings(path)
}

// GetStringOr 读取字符串配置，如果不存在则返回默认值
func GetStringOr(path string, defaultValue string) string {
	mu.RLock()
	defer mu.RUnlock()
	if !k.Exists(path) {
		return defaultValue
	}
	return k.String(path)
}

// GetIntOr 读取整数配置，如果不存在则返回默认值
func GetIntOr(path string, defaultValue int) int {
	mu.RLock()
	defer mu.RUnlock()
	if !k.Exists(path) {
		return defaultValue
	}
	return k.Int(path)
}

// GetBoolOr 读取布尔配置，如果不存在则返回默认值
func GetBoolOr(path string, defaultValue bool) bool {
	mu.RLock()
	defer mu.RUnlock()
	if !k.Exists(path) {
		return defaultValue
	}
	return k.Bool(path)
}

// GetFloat64Or 读取浮点数配置，如果不存在则返回默认值
func GetFloat64Or(path string, defaultValue float64) float64 {
	mu.RLock()
	defer mu.RUnlock()
	if !k.Exists(path) {
		return defaultValue
	}
	return k.Float64(path)
}

// GetStringSliceOr 读取字符串数组配置，如果不存在则返回默认值
func GetStringSliceOr(path string, defaultValue []string) []string {
	mu.RLock()
	defer mu.RUnlock()
	if !k.Exists(path) {
		return defaultValue
	}
	return k.Strings(path)
}

// Exists 检查配置键是否存在
func Exists(path string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return k.Exists(path)
}

func Unmarshal(path string, out interface{}) error {
	mu.RLock()
	defer mu.RUnlock()
	return k.Unmarshal(path, out)
}

func OnChange(cb func()) {
	changeCallbacks = append(changeCallbacks, cb)
}

func notifyChange() {
	for _, cb := range changeCallbacks {
		cb()
	}
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
