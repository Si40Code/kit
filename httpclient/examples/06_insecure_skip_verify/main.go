package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Si40Code/kit/httpclient"
	"github.com/Si40Code/kit/logger"
)

func main() {
	// 初始化 logger
	err := logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	fmt.Println("=== 示例: 忽略证书验证 ===")

	// 创建客户端，启用 WithInsecureSkipVerify 选项
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithInsecureSkipVerify(), // 忽略 SSL 证书验证
	)

	// 测试请求到自签名证书的 HTTPS 服务器
	// 注意：这个示例使用 httpbin.org，它可以正常工作
	// 如果你需要测试跳过证书验证，可以访问使用自签名证书的内部服务
	resp, err := client.R(context.Background()).
		Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("状态码: %d\n", resp.StatusCode())
	fmt.Printf("响应体前 200 字符: %s\n", resp.String()[:200])
}
