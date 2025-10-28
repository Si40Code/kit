package orm

import (
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	spanKey        = "orm:span"
	startTimeKey   = "orm:start_time"
	ignoreNotFound = "orm:ignore_not_found"
)

// plugin ORM 插件，用于集成 trace 和 metric
type plugin struct {
	options *options
}

// Name 返回插件名称
func (p *plugin) Name() string {
	return "kit-orm-plugin"
}

// Initialize 初始化插件，注册所有回调
func (p *plugin) Initialize(db *gorm.DB) error {
	// 注册 Create 操作的 before 和 after 回调
	if err := db.Callback().Create().Before("gorm:create").Register("kit:before_create", p.before); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:create").Register("kit:after_create", p.after); err != nil {
		return err
	}

	// 注册 Query 操作的 before 和 after 回调
	if err := db.Callback().Query().Before("gorm:query").Register("kit:before_query", p.before); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("kit:after_query", p.after); err != nil {
		return err
	}

	// 注册 Update 操作的 before 和 after 回调
	if err := db.Callback().Update().Before("gorm:update").Register("kit:before_update", p.before); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("kit:after_update", p.after); err != nil {
		return err
	}

	// 注册 Delete 操作的 before 和 after 回调
	if err := db.Callback().Delete().Before("gorm:delete").Register("kit:before_delete", p.before); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("kit:after_delete", p.after); err != nil {
		return err
	}

	// 注册 Row 操作的 before 和 after 回调
	if err := db.Callback().Row().Before("gorm:row").Register("kit:before_row", p.before); err != nil {
		return err
	}
	if err := db.Callback().Row().After("gorm:row").Register("kit:after_row", p.after); err != nil {
		return err
	}

	// 注册 Raw 操作的 before 和 after 回调
	if err := db.Callback().Raw().Before("gorm:raw").Register("kit:before_raw", p.before); err != nil {
		return err
	}
	if err := db.Callback().Raw().After("gorm:raw").Register("kit:after_raw", p.after); err != nil {
		return err
	}

	return nil
}

// before 在数据库操作前执行
func (p *plugin) before(db *gorm.DB) {
	if db.Statement == nil || db.Statement.Context == nil {
		return
	}

	ctx := db.Statement.Context

	// 记录开始时间
	db.InstanceSet(startTimeKey, time.Now())

	// 如果启用了 trace，创建 span
	if p.options.enableTrace {
		operation := getOperation(db)
		newCtx, span := createSpan(ctx, p.options.serviceName, operation)

		// 保存 span 到实例变量中
		db.InstanceSet(spanKey, span)

		// 更新 context
		db.Statement.Context = newCtx
	}
}

// after 在数据库操作后执行
func (p *plugin) after(db *gorm.DB) {
	if db.Statement == nil {
		return
	}

	// 获取开始时间
	startTime, ok := db.InstanceGet(startTimeKey)
	if !ok {
		return
	}

	duration := time.Since(startTime.(time.Time))
	operation := getOperation(db)
	table := getTable(db)
	sql := db.Statement.SQL.String()
	rowsAffected := db.Statement.RowsAffected
	err := db.Error

	// 处理 ignoreRecordNotFound
	if err == gorm.ErrRecordNotFound {
		// 检查是否配置了单次忽略
		if val, exists := db.InstanceGet(ignoreNotFound); exists && val.(bool) {
			err = nil
			db.Error = nil
		} else if p.options.ignoreRecordNotFound {
			// 或者全局配置了忽略
			err = nil
			db.Error = nil
		}
	}

	// Trace: 更新 span
	if p.options.enableTrace {
		if spanVal, ok := db.InstanceGet(spanKey); ok {
			span := spanVal.(trace.Span)

			attrs := map[string]interface{}{
				"db.operation":     operation,
				"db.table":         table,
				"db.statement":     sql,
				"db.rows_affected": rowsAffected,
				"db.duration_ms":   duration.Milliseconds(),
			}
			setSpanAttributes(span, attrs)

			if err != nil {
				markSpanError(span, err, fmt.Sprintf("database %s failed", operation))
			} else {
				markSpanSuccess(span)
			}

			span.End()
		}
	}

	// Metric: 记录指标
	if p.options.enableMetric && p.options.metricRecorder != nil {
		p.options.metricRecorder.RecordQuery(MetricData{
			Operation:    operation,
			Table:        table,
			SQL:          sql,
			Duration:     duration,
			RowsAffected: rowsAffected,
			Error:        err,
		})
	}
}

// getOperation 获取操作类型
func getOperation(db *gorm.DB) string {
	if db.Statement == nil {
		return "UNKNOWN"
	}

	sql := strings.ToUpper(db.Statement.SQL.String())

	switch {
	case strings.HasPrefix(sql, "SELECT"):
		return "SELECT"
	case strings.HasPrefix(sql, "INSERT"):
		return "INSERT"
	case strings.HasPrefix(sql, "UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(sql, "DELETE"):
		return "DELETE"
	case strings.HasPrefix(sql, "CREATE"):
		return "CREATE"
	case strings.HasPrefix(sql, "ALTER"):
		return "ALTER"
	case strings.HasPrefix(sql, "DROP"):
		return "DROP"
	default:
		return "RAW"
	}
}

// getTable 获取表名
func getTable(db *gorm.DB) string {
	if db.Statement == nil {
		return ""
	}

	if db.Statement.Table != "" {
		return db.Statement.Table
	}

	if db.Statement.Schema != nil {
		return db.Statement.Schema.Table
	}

	return ""
}

// setIgnoreRecordNotFound 在当前会话中设置忽略 RecordNotFound 错误
func setIgnoreRecordNotFound(db *gorm.DB) *gorm.DB {
	return db.InstanceSet(ignoreNotFound, true)
}
