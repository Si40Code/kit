package provider

import (
	"context"

	"github.com/knadh/koanf/v2"
)

// RemoteProvider 远程配置提供者接口
// 可以实现 Apollo、Nacos、Consul 等远程配置中心的接入
type RemoteProvider interface {
	Load(ctx context.Context, k *koanf.Koanf) error
	Watch(ctx context.Context, onChange func(map[string]interface{})) error
}
