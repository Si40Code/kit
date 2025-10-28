# 示例 5: 文件操作和日志优化

这个示例演示了 httpclient 如何智能地处理文件上传/下载的日志记录，以及如何保护敏感信息。

## 功能特性

- 文件上传的智能日志记录
- 文件下载的智能日志记录
- 敏感头信息自动过滤
- 大文件响应的日志截断

## 运行示例

```bash
cd examples/05_file_operations
go run main.go
```

## 代码说明

### 1. 文件上传

```go
resp, err := client.R(context.Background()).
    SetFile("file", tmpFile.Name()).
    SetFormData(map[string]string{
        "description": "测试文件上传",
    }).
    Post("https://httpbin.org/post")
```

**日志输出：**
```json
{
  "level": "info",
  "msg": "HTTP request started",
  "http.method": "POST",
  "http.url": "https://httpbin.org/post",
  "http.request.content_type": "multipart/form-data; boundary=...",
  "http.request.body": "[multipart/form-data, size: 1234 bytes]"
}
```

### 2. 文件下载

```go
resp, err := client.R(context.Background()).
    Get("https://httpbin.org/image/png")
```

**日志输出：**
```json
{
  "level": "info",
  "msg": "HTTP request completed successfully",
  "http.status_code": 200,
  "http.response.content_type": "image/png",
  "http.response.body": "[file download, size: 8090 bytes]"
}
```

### 3. 敏感头信息过滤

```go
resp, err := client.R(context.Background()).
    SetHeader("Authorization", "Bearer super-secret-token").
    SetHeader("X-API-Key", "my-api-key").
    Get("https://httpbin.org/get")
```

**日志输出：**
```json
{
  "http.request.headers": {
    "Authorization": "******",
    "X-API-Key": "******",
    "User-Agent": "MyApp/1.0"
  }
}
```

### 4. 大文件响应截断

```go
client := httpclient.New(
    httpclient.WithLogger(logger.Default()),
    httpclient.WithMaxBodyLogSize(100), // 只记录前 100 字节
)
```

**日志输出：**
```json
{
  "http.response.body": "first 100 bytes...",
  "http.response.body_truncated": true,
  "http.response.body_original_size": 5000
}
```

## 智能检测机制

### 文件上传检测

自动检测以下情况：
- Content-Type 包含 `multipart/`
- 使用 SetFile() 或 SetFiles() 方法

### 文件下载检测

自动检测以下情况：
- Content-Disposition 包含 `attachment` 或 `inline`
- Content-Type 为常见文件类型：
  - `application/octet-stream`
  - `application/pdf`
  - `application/zip`
  - `image/*`
  - `video/*`
  - `audio/*`
  - Office 文档类型

### 敏感头信息过滤

自动过滤以下头信息：
- `Authorization`
- `Cookie` / `Set-Cookie`
- `X-API-Key` / `API-Key` / `ApiKey`
- `Token` / `Access-Token` / `Refresh-Token`
- `X-Auth-Token` / `X-CSRF-Token`
- `Password` / `Secret`

## 配置选项

### 1. 日志大小限制

```go
client := httpclient.New(
    httpclient.WithMaxBodyLogSize(10 * 1024), // 10KB
)
```

### 2. 禁用日志

```go
client := httpclient.New(
    httpclient.WithDisableLog(),
)
```

## 日志字段说明

### 请求日志字段

| 字段 | 说明 | 示例 |
|------|------|------|
| http.method | HTTP 方法 | `POST` |
| http.url | 请求 URL | `https://api.example.com/upload` |
| http.request.headers | 请求头（已过滤） | `{"User-Agent": "..."}` |
| http.request.query_params | 查询参数 | `{"page": "1"}` |
| http.request.form_data | 表单数据 | `{"name": "test"}` |
| http.request.content_type | Content-Type | `multipart/form-data` |
| http.request.content_length | 内容长度 | `1234` |
| http.request.body | 请求体 | `[multipart/form-data]` 或实际内容 |
| http.request.body_truncated | 是否截断 | `true` |
| http.request.body_original_size | 原始大小 | `5000` |

### 响应日志字段

| 字段 | 说明 | 示例 |
|------|------|------|
| http.status_code | 状态码 | `200` |
| http.status | 状态信息 | `200 OK` |
| http.response.headers | 响应头（已过滤） | `{"Content-Type": "..."}` |
| http.response.content_type | Content-Type | `image/png` |
| http.response.content_length | 内容长度 | `8090` |
| http.response.content_disposition | Content-Disposition | `attachment; filename="file.pdf"` |
| http.response.body | 响应体 | `[file download]` 或实际内容 |
| http.response.body_truncated | 是否截断 | `true` |
| http.response.body_original_size | 原始大小 | `10000` |

## 最佳实践

### 1. 文件上传

```go
// ✅ 推荐：适当的日志大小限制
client := httpclient.New(
    httpclient.WithLogger(logger.Default()),
    httpclient.WithMaxBodyLogSize(10 * 1024),
)

// 上传文件
resp, err := client.R(ctx).
    SetFile("file", filePath).
    Post(url)
```

### 2. 文件下载

```go
// ✅ 推荐：下载大文件时可以考虑禁用日志
client := httpclient.New(
    httpclient.WithLogger(logger.Default()),
    httpclient.WithMaxBodyLogSize(1024), // 只记录元信息
)

resp, err := client.R(ctx).Get(downloadURL)
```

### 3. 敏感信息处理

```go
// ✅ 自动过滤敏感头信息
resp, err := client.R(ctx).
    SetHeader("Authorization", "Bearer " + token).
    Get(url)

// 日志中会显示：
// "Authorization": "******"
```

### 4. 生产环境配置

```go
client := httpclient.New(
    httpclient.WithLogger(logger.Default()),
    httpclient.WithMaxBodyLogSize(5 * 1024), // 5KB
    httpclient.WithTrace("my-service"),
    httpclient.WithMetric(recorder),
)
```

## 性能考虑

1. **日志大小限制**：避免记录过大的文件内容影响性能
2. **敏感信息过滤**：最小的性能开销，字符串比较
3. **文件检测**：基于 Content-Type，几乎无开销

## 注意事项

⚠️ **文件上传**：multipart/form-data 的内容不会记录到日志中，只记录元信息

⚠️ **文件下载**：二进制文件内容不会记录到日志中，只记录大小和类型

⚠️ **敏感信息**：所有匹配的敏感头信息都会被自动替换为 `******`

⚠️ **日志截断**：当 body 超过限制时，会自动截断并设置 `body_truncated=true`

## 输出示例

```
=== 示例 1: 文件上传（日志会显示元信息） ===
{"level":"info","msg":"HTTP request started","http.method":"POST","http.request.body":"[multipart/form-data, size: 1234 bytes]"}
✓ 文件上传成功，状态码: 200

=== 示例 2: 文件下载（日志会显示元信息） ===
{"level":"info","msg":"HTTP request completed successfully","http.response.body":"[file download, size: 8090 bytes]"}
✓ 文件下载成功，状态码: 200

=== 示例 3: 敏感头信息过滤 ===
{"level":"info","http.request.headers":{"Authorization":"******","X-API-Key":"******","User-Agent":"MyApp/1.0"}}
✓ 请求成功，状态码: 200

=== 示例 4: 大文件响应（日志会截断） ===
{"http.response.body_truncated":true,"http.response.body_original_size":1500}
✓ 请求成功，状态码: 200
```

