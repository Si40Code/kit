# 示例 1: 基本用法

这个示例演示了如何使用 httpclient 包进行基本的 HTTP 请求。

## 功能特性

- 基本的 GET/POST 请求
- 查询参数设置
- 请求体设置
- 自动日志记录

## 运行示例

```bash
cd examples/01_basic
go run main.go
```

## 代码说明

### 1. 创建客户端

```go
client := httpclient.New(
    httpclient.WithLogger(logger.L()),
)
```

### 2. 发起 GET 请求

```go
resp, err := client.R(context.Background()).
    Get("https://httpbin.org/get")
```

### 3. 发起 POST 请求

```go
resp, err := client.R(context.Background()).
    SetBody(user).
    Post("https://httpbin.org/post")
```

### 4. 设置查询参数

```go
resp, err := client.R(context.Background()).
    SetQueryParams(map[string]string{
        "page": "1",
        "pageSize": "10",
    }).
    Get("https://httpbin.org/get")
```

## 配置选项

- `WithLogger`: 设置日志记录器
- `WithTimeout`: 设置请求超时时间
- `WithMaxBodyLogSize`: 设置最大日志 body 大小

## 输出示例

程序会输出详细的请求日志，包括：
- 请求方法和 URL
- 请求头和查询参数
- 请求体和响应体
- 响应状态码

