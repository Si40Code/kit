# 基础用法示例

最简单的使用方式，展示如何创建一个基本的 HTTP 服务器。

## 功能

- 创建一个简单的 web 服务器
- 注册 GET 和 POST 路由
- 使用统一响应格式

## 运行

```bash
go run main.go
```

## 测试

```bash
# 测试 ping
curl http://localhost:8080/ping

# 测试 echo
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"name":"John","age":30}'
```
