# 示例 06: 忽略安全证书验证

本示例演示如何使用 `WithInsecureSkipVerify()` 选项来跳过 SSL/TLS 证书验证。

## ⚠️ 安全警告

**不建议在生产环境使用** `WithInsecureSkipVerify()` 选项，因为它会跳过证书验证，使你的应用容易受到中间人攻击。

## 使用场景

- 开发环境访问使用自签名证书的内部服务
- 测试环境
- 本地开发时快速调试

## 代码示例

```go
// 创建客户端，启用跳过证书验证
client := httpclient.New(
    httpclient.WithLogger(logger.Default()),
    httpclient.WithInsecureSkipVerify(), // 关键配置
)

// 现在可以访问使用自签名证书的 HTTPS 服务
resp, err := client.R(context.Background()).
    Get("https://internal-service.example.com/api")
```

## 实现原理

在 `client.go` 中的实现逻辑：

```go
// TLS 配置
if options.insecureSkipVerify {
    transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
```

当启用 `WithInsecureSkipVerify()` 时，会将 `InsecureSkipVerify: true` 传递给 `tls.Config`，从而跳过 TLS 证书验证。

## 运行示例

```bash
cd 06_insecure_skip_verify && go run main.go
```

## 更好的替代方案

对于生产环境，建议：

1. **使用正确的证书**：为你的服务配置有效的 SSL 证书
2. **自定义证书池**：如果需要使用自定义 CA，可以修改 TLS 配置：

```go
// 加载自定义 CA 证书
certPool := x509.NewCertPool()
certFile, _ := os.ReadFile("custom-ca.crt")
certPool.AppendCertsFromPEM(certFile)

transport.TLSClientConfig = &tls.Config{
    RootCAs: certPool,
}
```

3. **仅允许特定证书**：使用客户端证书认证

