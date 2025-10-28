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

	fmt.Println("=== 示例 1: 基本 GET 请求 ===")
	example1BasicGet()

	fmt.Println("\n=== 示例 2: POST 请求 ===")
	example2Post()

	fmt.Println("\n=== 示例 3: 带查询参数的请求 ===")
	example3QueryParams()

	fmt.Println("\n=== 示例 4: 自定义配置 ===")
	example4CustomConfig()
}

// 示例 1: 基本 GET 请求
func example1BasicGet() {
	// 创建客户端
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
	)

	// 发起请求
	resp, err := client.R(context.Background()).
		Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}
	fmt.Printf("状态码: %d\n", resp.StatusCode())
	fmt.Printf("响应体: %s\n", resp.String())
}

// 示例 2: POST 请求
func example2Post() {
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
	)

	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := User{
		Name:  "Alice",
		Email: "alice@example.com",
	}

	resp, err := client.R(context.Background()).
		SetBody(user).
		Post("https://httpbin.org/post")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("状态码: %d\n", resp.StatusCode())
	fmt.Printf("响应体: %s\n", resp.String())
}

// 示例 3: 带查询参数的请求
func example3QueryParams() {
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
	)

	resp, err := client.R(context.Background()).
		SetQueryParams(map[string]string{
			"page":     "1",
			"pageSize": "10",
		}).
		Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("状态码: %d\n", resp.StatusCode())
	fmt.Printf("响应体: %s\n", resp.String())
}

// 示例 4: 自定义配置
func example4CustomConfig() {
	client := httpclient.New(
		httpclient.WithLogger(logger.Default()),
		httpclient.WithTimeout(10*1000000000), // 10 seconds
		httpclient.WithMaxBodyLogSize(1024),   // 1KB
	)

	resp, err := client.R(context.Background()).
		SetHeader("User-Agent", "MyApp/1.0").
		Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("状态码: %d\n", resp.StatusCode())
}
