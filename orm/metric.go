package orm

import "time"

// MetricData 数据库查询的详细指标数据
type MetricData struct {
	// 基础信息
	Operation string // 查询类型：SELECT, INSERT, UPDATE, DELETE, RAW 等
	Table     string // 表名
	SQL       string // 完整 SQL 语句

	// 性能指标
	Duration     time.Duration // 查询耗时
	RowsAffected int64         // 影响/返回的行数

	// 错误信息
	Error error // 错误（如果有）
}

// MetricRecorder 数据库查询指标记录器接口
type MetricRecorder interface {
	// RecordQuery 记录数据库查询指标
	// 实现者可以将数据发送到 Prometheus、SigNoz 或其他监控系统
	RecordQuery(data MetricData)
}
