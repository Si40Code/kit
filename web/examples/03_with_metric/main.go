package main

import (
	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMetricRecorder Prometheus 指标记录器
type PrometheusMetricRecorder struct {
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewPrometheusMetricRecorder() *PrometheusMetricRecorder {
	recorder := &PrometheusMetricRecorder{
		requestCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latencies in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	prometheus.MustRegister(recorder.requestCounter)
	prometheus.MustRegister(recorder.requestDuration)

	return recorder
}

func (r *PrometheusMetricRecorder) RecordRequest(data web.MetricData) {
	r.requestCounter.WithLabelValues(
		data.Method,
		data.Path,
		string(rune(data.Status)),
	).Inc()

	r.requestDuration.WithLabelValues(
		data.Method,
		data.Path,
	).Observe(data.Duration.Seconds())
}

func main() {
	// 创建 Prometheus recorder
	metricRecorder := NewPrometheusMetricRecorder()

	// 创建服务器
	server := web.New(
		web.WithMode(web.ReleaseMode),
		web.WithServiceName("metric-example"),
		web.WithMetric(metricRecorder),
		web.WithSkipPaths("/metrics"), // 不记录 metrics 端点的日志
	)

	engine := server.Engine()

	// Prometheus metrics 端点
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 业务路由
	engine.GET("/api/hello", func(c *gin.Context) {
		web.Success(c, gin.H{"message": "Hello World"})
	})

	server.RunWithGracefulShutdown(":8080")
}
