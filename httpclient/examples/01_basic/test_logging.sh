#!/bin/bash

# 格式化 JSON 日志输出
echo "=== 测试 HTTP 请求日志 ==="
echo ""

# 运行程序并提取 JSON 日志
go run main.go 2>&1 | grep "^\{" | jq -c 'select(.msg == "HTTP request started")' | head -1 | jq .

echo ""
echo "=== 上述日志包含请求信息："
echo "  - http.method: HTTP 方法"
echo "  - http.url: 请求 URL"
echo "  - http.request.body: 请求体内容"
echo "  - http.request.headers: 请求头（如果设置）"
echo "  - http.request.query_params: 查询参数（如果设置）"
echo "  - http.request.form_data: 表单数据（如果设置）"
echo ""

echo "=== 测试 HTTP 响应日志 ==="
go run main.go 2>&1 | grep "^\{" | jq -c 'select(.msg == "HTTP request completed successfully")' | head -1 | jq .

echo ""
echo "=== 上述日志包含响应信息："
echo "  - http.status_code: 状态码"
echo "  - http.response.body: 响应体内容"
echo "  - http.response.headers: 响应头"
echo "  - http.total_time_ms: 总耗时"
echo "  - http.dns_lookup_ms: DNS 查询时间"
echo "  - http.tcp_conn_ms: TCP 连接时间"
echo "  - http.tls_handshake_ms: TLS 握手时间"
echo "  - http.server_time_ms: 服务器处理时间"
echo "  - http.conn_reused: 连接是否复用"
echo ""

