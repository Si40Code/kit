# 自定义日志示例

展示如何实现自定义日志记录器。

## 功能

- 实现 Logger 接口
- 自定义日志格式
- 设置日志大小限制

## 运行

```bash
go run main.go
```

## 测试

```bash
# 测试 GET
curl http://localhost:8080/hello

# 测试 POST（会记录日志）
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{"key":"value","data":"test"}'
```
