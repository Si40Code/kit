# 链路追踪示例

展示如何集成 OpenTelemetry 进行链路追踪。

## 功能

- 集成 OpenTelemetry
- 使用 stdout exporter（输出到终端）
- 在 handler 中创建子 span

> 注意：生产环境建议使用 Jaeger 或 OpenTelemetry Collector

## 运行

```bash
go run main.go
```

## 测试

```bash
curl http://localhost:8080/user/123
```

查看终端输出中的追踪信息。
