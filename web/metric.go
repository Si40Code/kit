package web

import "time"

// MetricData 指标数据
type MetricData struct {
	Method   string
	Path     string
	Status   int
	Duration time.Duration
}

// MetricRecorder 指标记录器接口
type MetricRecorder interface {
	// RecordRequest 记录请求指标
	RecordRequest(data MetricData)
}
