package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/shima-park/agollo"
)

// ApolloProvider 实现 RemoteProvider 接口，用于从 Apollo 配置中心加载配置
type ApolloProvider struct {
	client    agollo.Agollo
	configKey string
}

// ApolloConfig Apollo 配置参数
type ApolloConfig struct {
	AppID     string
	ConfigKey string
	AccessKey string
	Cluster   string
	ServerURL string
}

// NewApolloProvider 创建 Apollo 配置提供者
func NewApolloProvider(cfg ApolloConfig) (*ApolloProvider, error) {
	if cfg.AppID == "" || cfg.ConfigKey == "" {
		return nil, fmt.Errorf("Apollo AppID and ConfigKey are required")
	}

	// 创建 Apollo 客户端
	client, err := agollo.New(
		cfg.ServerURL,
		cfg.AppID,
		agollo.Cluster(cfg.Cluster),
		agollo.AccessKey(cfg.AccessKey),
		agollo.FailTolerantOnBackupExists(),
		agollo.AutoFetchOnCacheMiss(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Apollo client: %w", err)
	}

	return &ApolloProvider{
		client:    client,
		configKey: cfg.ConfigKey,
	}, nil
}

// Load 从 Apollo 加载配置到 koanf
func (p *ApolloProvider) Load(ctx context.Context, k *koanf.Koanf) error {
	log.Println("[Apollo] Loading configuration from Apollo...")

	// 从 Apollo 获取配置
	configs := p.client.GetNameSpace(p.configKey)
	configStr, ok := configs["content"].(string)
	if !ok {
		return fmt.Errorf("invalid config content from Apollo: %+v", configs["content"])
	}

	if configStr == "" {
		return fmt.Errorf("empty config content from Apollo, configKey: %s", p.configKey)
	}

	// 解析 YAML 配置并加载到 koanf
	if err := k.Load(rawbytes.Provider([]byte(configStr)), yaml.Parser()); err != nil {
		return fmt.Errorf("failed to parse Apollo config: %w", err)
	}

	log.Printf("[Apollo] Successfully loaded configuration from Apollo (configKey: %s)", p.configKey)
	return nil
}

// Watch 监听 Apollo 配置变更（启动时加载一次，不进行热更新）
func (p *ApolloProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
	log.Println("[Apollo] Starting Apollo watcher (no hot reload)...")

	// 启动 Apollo 监听
	errorCh := p.client.Start()
	watchCh := p.client.Watch()

	// 启动 goroutine 处理 Apollo 事件，但不调用 onChange
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[Apollo] Context cancelled, stopping watcher")
				return
			case err := <-errorCh:
				if err != nil {
					log.Printf("[Apollo] Error from Apollo server: %v", err.Err)
				}
			case resp := <-watchCh:
				if resp.Error == nil {
					log.Printf("[Apollo] Configuration changed in Apollo (configKey: %s), but hot reload is disabled", p.configKey)
					// 不调用 onChange，因为配置热更新被禁用
				} else {
					log.Printf("[Apollo] Error watching Apollo changes: %v", resp.Error)
				}
			}
		}
	}()

	return nil
}

// Close 关闭 Apollo 客户端
func (p *ApolloProvider) Close() error {
	// Apollo client doesn't have a Close method, so we just return nil
	return nil
}
