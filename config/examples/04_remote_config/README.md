# 示例 4: 远程配置（Apollo）

展示如何接入远程配置中心（以 Apollo 为例），实现配置的集中管理和动态推送。

## 运行示例

```bash
cd config/examples/04_remote_config
go run main.go
```

## 学习内容

1. **远程配置接入** - 实现 `RemoteProvider` 接口
2. **配置分层架构** - 本地配置作为兜底，远程配置优先
3. **动态推送** - 配置中心推送变更到客户端
4. **高可用性** - 远程配置失败时使用本地配置

## RemoteProvider 接口

```go
type RemoteProvider interface {
    Load(ctx context.Context, k *koanf.Koanf) error
    Watch(ctx context.Context, onChange func(map[string]interface{})) error
}
```

## 配置优先级

从低到高：
1. 本地配置文件 (config.yaml) - 默认配置/兜底配置
2. Apollo 远程配置 - 动态配置/运维配置 ⬅️ 最高优先级

## 实际应用场景

### 场景 1: 功能开关
运维人员在 Apollo 控制台修改功能开关，所有应用实例实时生效。

### 场景 2: 限流参数调整
动态调整 API 限流参数，无需重启应用。

### 场景 3: 灰度发布配置
为不同的应用实例推送不同的配置，实现灰度发布。

## Apollo 实际接入

本示例提供的是模拟实现，实际使用 Apollo 时：

### 1. 安装 Apollo SDK

```bash
go get github.com/apolloconfig/agollo/v4
```

### 2. 实现 RemoteProvider

```go
import "github.com/apolloconfig/agollo/v4"

type ApolloProvider struct {
    client agollo.Client
}

func (p *ApolloProvider) Load(ctx context.Context, k *koanf.Koanf) error {
    cache := p.client.GetDefaultConfigCache()
    for key, value := range cache {
        k.Set(key, value)
    }
    return nil
}

func (p *ApolloProvider) Watch(ctx context.Context, onChange func(map[string]interface{})) error {
    p.client.AddChangeListener(&apolloListener{onChange: onChange})
    return nil
}
```

### 3. 初始化配置

```go
apolloProvider := NewApolloProvider(&ApolloConfig{
    AppID:     "your-app-id",
    Cluster:   "default",
    Namespace: "application",
    ServerURL: "http://apollo-config.example.com",
})

config.Init(
    config.WithFile("config.yaml"),
    config.WithRemote(apolloProvider),
)
```

## 其他配置中心

同样的接口可以接入其他配置中心：

- **Nacos**: 阿里云的配置中心
- **Consul**: HashiCorp 的服务发现和配置管理工具
- **etcd**: 分布式键值存储
- **自建配置中心**: 实现 HTTP 接口即可

## 最佳实践

- ✅ 始终提供本地配置文件作为兜底
- ✅ 远程配置失败时降级到本地配置
- ✅ 监控远程配置的可用性
- ✅ 配置变更时记录审计日志
- ⚠️ 避免在远程配置中存储超大配置（如完整的白名单列表）

