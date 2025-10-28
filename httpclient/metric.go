package httpclient

import "time"

// MetricData HTTP 请求的详细指标数据
type MetricData struct {
	// 基础信息
	Method     string // HTTP 方法
	Host       string // 主机名
	Path       string // 请求路径
	StatusCode int    // HTTP 状态码

	// 时间指标（来自 Resty TraceInfo）
	TotalTime    time.Duration // 总耗时
	DNSLookup    time.Duration // DNS 查询时间
	TCPConn      time.Duration // TCP 连接时间
	TLSHandshake time.Duration // TLS 握手时间
	ServerTime   time.Duration // 服务器处理时间
	ResponseTime time.Duration // 响应传输时间

	// 连接信息
	IsConnReused   bool          // 连接是否复用
	IsConnWasIdle  bool          // 连接是否空闲
	ConnIdleTime   time.Duration // 连接空闲时间
	RequestAttempt int           // 请求尝试次数
	RemoteAddr     string        // 远程地址
}

// MetricRecorder HTTP 请求指标记录器接口
type MetricRecorder interface {
	// RecordRequest 记录 HTTP 请求指标
	// 实现者可以将数据发送到 Prometheus、SigNoz 或其他监控系统
	RecordRequest(data MetricData)
}
