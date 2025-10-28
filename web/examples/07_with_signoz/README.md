# SigNoz 集成示例

展示如何集成 SigNoz 进行统一的日志、追踪和指标监控。

## 功能

- 集成 SigNoz 进行追踪
- 集成 SigNoz 进行指标监控
- 自定义日志记录器
- 统一的遥测数据收集

## 前置条件

确保 SigNoz 实例运行在 `47.83.197.11`。

SigNoz 是一个开源的 APM 平台，提供日志、追踪和指标的统一收集。

## 运行

```bash
go run main.go
```

**✅ 现在直接运行就可以看到 Trace、Metrics 和 Logs！**

使用了 SigNoz 官方的 [zap_otlp](https://github.com/SigNoz/zap_otlp) 库，日志通过 OTLP 协议直接发送到 SigNoz，无需额外配置。

服务启动后，会自动每 10 秒发送一批测试请求，这样你可以立即在 SigNoz 中看到追踪、指标和日志数据。

## 访问 SigNoz UI

打开浏览器访问：http://47.83.197.11

等待几秒后，你应该能看到：
- **Traces** - 自动生成的追踪数据
- **Metrics** - 请求指标
- **Logs** - 日志数据（使用 `send_logs.sh` 运行时）

## 测试

### 1. 健康检查
```bash
curl http://localhost:8080/health
```

### 2. 获取用户列表
```bash
curl http://localhost:8080/api/users
```

### 3. 获取单个用户
```bash
curl http://localhost:8080/api/users/1
```

### 4. 创建用户
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"David","email":"david@example.com"}'
```

### 5. 生成流量（用于查看追踪和指标）
```bash
# 运行多次以生成追踪数据
for i in {1..20}; do
  curl http://localhost:8080/api/users
  curl http://localhost:8080/api/users/$((i % 3 + 1))
  sleep 0.5
done
```

## 在 SigNoz 中查看

### 1. Traces（追踪）
导航到 **Traces** 页面，查看请求的详细追踪链。

### 2. Logs（日志）
✅ 在 **Logs** 页面直接查看：
- 使用 [zap_otlp](https://github.com/SigNoz/zap_otlp) 通过 OTLP 协议发送
- 自动包含 trace_id 和 span_id，可关联追踪
- 点击日志中的 trace_id 可跳转到对应的追踪详情
- 在 Traces 页面也能看到关联的日志

### 3. Metrics（指标）

#### 方式 1: 通过 Prometheus 端点查看（推荐）

访问应用的 metrics 端点：
```bash
curl http://localhost:8080/metrics
```

你会看到类似这样的指标：
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/users",service="signoz-example",status="200"} 10

# HELP http_request_duration_seconds HTTP request latencies in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/api/users",service="signoz-example",le="0.005"} 5
```

#### 方式 2: 在 SigNoz 中查看

1. 在 SigNoz UI 中，导航到 **Dashboard** 或 **Metrics Explorer**
2. 搜索指标名称：
   - `http_requests_total` - 请求总数
   - `http_request_duration_seconds` - 请求延迟
3. 按标签过滤：
   - `service="signoz-example"`
   - `method="GET"`
   - `path="/api/users"`

#### 配置 SigNoz 抓取 Prometheus 指标

如果 SigNoz 还没有抓取你的指标，需要配置 Prometheus 抓取任务。在 SigNoz 的配置中添加：

```yaml
scrape_configs:
  - job_name: 'signoz-example'
    static_configs:
      - targets: ['your-app-host:8080']
    metrics_path: '/metrics'
```

## 配置说明

- SigNoz 地址：`47.83.197.11`
- 服务名称：`signoz-example`（会在所有遥测数据中使用）
- 追踪使用 OTLP HTTP 协议（端口 4318）
- 日志使用 OTLP gRPC 协议（端口 4317）
- 指标使用 Prometheus 格式

## 技术实现

使用 SigNoz 官方的 [zap_otlp](https://github.com/SigNoz/zap_otlp) 库：

- **Zap Logger**: 高性能的结构化日志库
- **OTLP Encoder**: 将日志编码为 OTLP 格式
- **OTLP Syncer**: 通过 gRPC 批量发送日志到 SigNoz
- **Trace Context**: 自动关联日志和追踪

## 生产环境建议

1. 使用环境变量配置 SigNoz 地址
2. 配置采样率以减少追踪开销
3. 调整日志批量大小（默认 100 条）
4. 设置合理的批量发送间隔（默认 5 秒）
5. 配置日志级别，避免过多日志

### 环境变量配置

```bash
# 配置 SigNoz 端点
export SIGNOZ_ENDPOINT="your-signoz-endpoint:4317"

# 配置服务名称
export SERVICE_NAME="your-service-name"

# 如果使用 TLS
export OTEL_EXPORTER_OTLP_INSECURE="false"
```
