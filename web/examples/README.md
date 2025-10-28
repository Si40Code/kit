# Web 模块示例

本文档包含所有 web 模块的使用示例。

## 示例列表

### 01_basic - 基础用法
最简单的使用方式，展示如何创建一个基本的 HTTP 服务器。

### 02_with_trace - 链路追踪
展示如何集成 OpenTelemetry 进行链路追踪。

### 03_with_metric - 指标监控
展示如何集成 Prometheus 指标监控。

### 04_file_upload - 文件上传
展示如何处理文件上传，包括单文件和多文件上传。

### 05_custom_logger - 自定义日志
展示如何实现自定义日志记录器。

### 06_production - 生产环境配置
展示生产环境的最佳实践配置。

### 07_with_signoz - SigNoz 集成
展示如何集成 SigNoz 进行统一的日志、追踪和指标监控。

## 快速开始

每个示例都可以独立运行：

```bash
# 进入示例目录
cd examples/01_basic

# 运行示例
go run main.go
```

## 选择指南

| 需求 | 推荐示例 |
|------|---------|
| 快速上手 | 01_basic |
| 需要追踪 | 02_with_trace |
| 需要指标 | 03_with_metric |
| 处理文件上传 | 04_file_upload |
| 自定义日志 | 05_custom_logger |
| 生产部署 | 06_production |
| 使用 SigNoz | 07_with_signoz |

## 注意事项

1. 某些示例需要依赖外部服务（如 Jaeger、SigNoz），请先阅读对应的 README
2. 文件上传示例需要创建 `uploads/` 目录
3. 生产环境示例建议使用环境变量配置
