# 指标监控示例

展示如何集成 Prometheus 指标监控。

## 功能

- 集成 Prometheus 指标收集
- 自定义 MetricRecorder 实现
- 记录请求数量和延迟

## 运行

```bash
go run main.go
```

## 访问 metrics 端点

```bash
curl http://localhost:8080/metrics
```

## 测试

```bash
# 发送一些请求
for i in {1..10}; do
  curl http://localhost:8080/api/hello
  sleep 0.5
done

# 查看 metrics
curl http://localhost:8080/metrics | grep http_requests
```
