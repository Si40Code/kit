package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("=== 示例 1: 文件上传（日志会显示元信息） ===")
	example1FileUpload()

	fmt.Println("\n=== 示例 2: 文件下载（日志会显示元信息） ===")
	example2FileDownload()

	fmt.Println("\n=== 示例 3: 敏感头信息过滤 ===")
	example3SensitiveHeaders()

	fmt.Println("\n=== 示例 4: 大文件响应（日志会截断） ===")
	example4LargeResponse()
}

// 示例 1: 文件上传
func example1FileUpload() {
	// 创建客户端
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
	)

	// 创建一个临时文件
	content := []byte("This is a test file content for upload demo")
	tmpFile, err := os.CreateTemp("", "upload-test-*.txt")
	if err != nil {
		log.Printf("创建临时文件失败: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		log.Printf("写入临时文件失败: %v", err)
		return
	}

	// 上传文件
	resp, err := client.R(context.Background()).
		SetFile("file", tmpFile.Name()).
		SetFormData(map[string]string{
			"description": "测试文件上传",
			"category":    "demo",
		}).
		Post("https://httpbin.org/post")
	if err != nil {
		log.Printf("文件上传失败: %v", err)
		return
	}

	fmt.Printf("✓ 文件上传成功，状态码: %d\n", resp.StatusCode())
	fmt.Println("注意：日志中只显示 [multipart/form-data] 而不是完整的文件内容")
}

// 示例 2: 文件下载
func example2FileDownload() {
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
	)

	// 下载一个图片文件
	resp, err := client.R(context.Background()).
		Get("https://httpbin.org/image/png")
	if err != nil {
		log.Printf("文件下载失败: %v", err)
		return
	}

	fmt.Printf("✓ 文件下载成功，状态码: %d\n", resp.StatusCode())
	fmt.Printf("✓ 文件大小: %d bytes\n", len(resp.Body()))
	fmt.Println("注意：日志中只显示 [file download, size: X bytes] 而不是二进制内容")
}

// 示例 3: 敏感头信息过滤
func example3SensitiveHeaders() {
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
	)

	// 发送带敏感信息的请求
	resp, err := client.R(context.Background()).
		SetHeader("Authorization", "Bearer super-secret-token-12345").
		SetHeader("X-API-Key", "my-api-key-67890").
		SetHeader("User-Agent", "MyApp/1.0").
		Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("✓ 请求成功，状态码: %d\n", resp.StatusCode())
	fmt.Println("注意：日志中 Authorization 和 X-API-Key 会显示为 ******")
}

// 示例 4: 大文件响应（日志会截断）
func example4LargeResponse() {
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithMaxBodyLogSize(100), // 只记录前 100 字节
	)

	// 生成一个大的 JSON 响应
	resp, err := client.R(context.Background()).
		SetBody(map[string]string{
			"data": strings.Repeat("A", 1000), // 1000 字符
		}).
		Post("https://httpbin.org/post")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("✓ 请求成功，状态码: %d\n", resp.StatusCode())
	fmt.Printf("✓ 响应大小: %d bytes\n", len(resp.Body()))
	fmt.Println("注意：日志中的 body 会被截断，并标记 body_truncated=true")
}
