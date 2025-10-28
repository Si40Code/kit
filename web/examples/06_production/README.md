# 生产环境配置示例

展示生产环境的最佳实践配置。

## 功能

- 生产环境优化配置
- 健康检查和就绪检查
- 慢请求监控
- 链路追踪集成
- 优雅关闭

## 运行

```bash
# 开发环境
go run main.go

# 模拟生产环境
APP_ENV=production go run main.go
```

## 测试

```bash
# 健康检查
curl http://localhost:8080/health

# 就绪检查
curl http://localhost:8080/ready

# 获取用户列表
curl http://localhost:8080/api/users

# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```
