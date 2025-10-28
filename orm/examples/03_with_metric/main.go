package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Si40Code/kit/logger"
	"github.com/Si40Code/kit/orm"
	"gorm.io/driver/mysql"
)

// User 用户模型
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex"`
	Age       int
}

// SimpleMetricRecorder 简单的 Metric 记录器示例
type SimpleMetricRecorder struct {
	mu      sync.Mutex
	metrics []orm.MetricData
}

func NewSimpleMetricRecorder() *SimpleMetricRecorder {
	return &SimpleMetricRecorder{
		metrics: make([]orm.MetricData, 0),
	}
}

// RecordQuery 实现 MetricRecorder 接口
func (r *SimpleMetricRecorder) RecordQuery(data orm.MetricData) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics = append(r.metrics, data)

	// 实时打印 metric 信息
	fmt.Printf("📊 Metric: %s on %s, duration=%dms, rows=%d",
		data.Operation,
		data.Table,
		data.Duration.Milliseconds(),
		data.RowsAffected,
	)

	if data.Error != nil {
		fmt.Printf(", error=%v", data.Error)
	}
	fmt.Println()
}

// PrintSummary 打印统计摘要
func (r *SimpleMetricRecorder) PrintSummary() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.metrics) == 0 {
		fmt.Println("No metrics recorded")
		return
	}

	fmt.Println("\n" + repeatString("=").repeat(60))
	fmt.Println("📈 Metrics Summary")
	fmt.Println(repeatString("=").repeat(60))

	// 按操作类型统计
	operationStats := make(map[string]struct {
		count    int
		totalMs  int64
		errorCnt int
	})

	// 按表统计
	tableStats := make(map[string]struct {
		count   int
		totalMs int64
	})

	var totalDuration time.Duration
	var errorCount int

	for _, m := range r.metrics {
		// 操作统计
		stat := operationStats[m.Operation]
		stat.count++
		stat.totalMs += m.Duration.Milliseconds()
		if m.Error != nil {
			stat.errorCnt++
			errorCount++
		}
		operationStats[m.Operation] = stat

		// 表统计
		if m.Table != "" {
			tstat := tableStats[m.Table]
			tstat.count++
			tstat.totalMs += m.Duration.Milliseconds()
			tableStats[m.Table] = tstat
		}

		totalDuration += m.Duration
	}

	// 打印总体统计
	fmt.Printf("\nTotal Queries: %d\n", len(r.metrics))
	fmt.Printf("Total Duration: %dms\n", totalDuration.Milliseconds())
	fmt.Printf("Average Duration: %dms\n", totalDuration.Milliseconds()/int64(len(r.metrics)))
	fmt.Printf("Error Count: %d\n", errorCount)

	// 打印操作类型统计
	fmt.Println("\nBy Operation Type:")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("%-10s | %8s | %10s | %8s\n", "Operation", "Count", "Total(ms)", "Avg(ms)")
	fmt.Println("--------------------------------------------------")
	for op, stat := range operationStats {
		avg := stat.totalMs / int64(stat.count)
		fmt.Printf("%-10s | %8d | %10d | %8d", op, stat.count, stat.totalMs, avg)
		if stat.errorCnt > 0 {
			fmt.Printf(" (errors: %d)", stat.errorCnt)
		}
		fmt.Println()
	}

	// 打印表统计
	if len(tableStats) > 0 {
		fmt.Println("\nBy Table:")
		fmt.Println("--------------------------------------------------")
		fmt.Printf("%-10s | %8s | %10s | %8s\n", "Table", "Count", "Total(ms)", "Avg(ms)")
		fmt.Println("--------------------------------------------------")
		for table, stat := range tableStats {
			avg := stat.totalMs / int64(stat.count)
			fmt.Printf("%-10s | %8d | %10d | %8d\n", table, stat.count, stat.totalMs, avg)
		}
	}

	fmt.Println(repeatString("=").repeat(60))
}

// String repeat helper
type repeatString string

func (s repeatString) repeat(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += string(s)
	}
	return result
}

func main() {
	// 1. 初始化 logger
	logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
	)
	defer logger.Sync()

	ctx := context.Background()

	// 2. 创建 Metric 记录器
	metricRecorder := NewSimpleMetricRecorder()

	// 3. 创建 ORM 客户端（请替换为你的数据库连接信息）
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := orm.New(
		mysql.Open(dsn),
		orm.WithLogger(logger.Default()),
		orm.WithMetric(metricRecorder), // 启用 metric
		orm.WithSlowThreshold(50*time.Millisecond),
	)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", "error", err)
	}
	defer client.Close()

	logger.Info(ctx, "Successfully connected to database with metric enabled")

	// 4. 自动迁移
	if err := client.AutoMigrate(&User{}); err != nil {
		logger.Fatal(ctx, "Failed to migrate database", "error", err)
	}

	// 5. 执行各种数据库操作，观察 metric 输出
	fmt.Println("\n🔄 Starting database operations...")

	// 创建多个用户
	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 25},
		{Name: "Bob", Email: "bob@example.com", Age: 30},
		{Name: "Charlie", Email: "charlie@example.com", Age: 35},
	}

	for _, user := range users {
		if err := client.WithContext(ctx).Create(&user).Error; err != nil {
			logger.Error(ctx, "Failed to create user", "error", err)
		}
		time.Sleep(10 * time.Millisecond) // 模拟间隔
	}

	// 查询操作
	var foundUsers []User
	client.WithContext(ctx).Find(&foundUsers)
	client.WithContext(ctx).Where("age > ?", 25).Find(&foundUsers)
	client.WithContext(ctx).First(&foundUsers[0], foundUsers[0].ID)

	// 更新操作
	client.WithContext(ctx).Model(&foundUsers[0]).Update("age", 26)
	client.WithContext(ctx).Model(&foundUsers[1]).Updates(User{Age: 31, Name: "Bob Updated"})

	// 删除操作
	client.WithContext(ctx).Delete(&foundUsers[2])

	// 测试错误情况
	var notFound User
	client.WithContext(ctx).First(&notFound, 99999) // 这会产生 RecordNotFound 错误

	// 6. 打印 metric 摘要
	metricRecorder.PrintSummary()

	fmt.Println("\n✅ All operations completed!")
	fmt.Println("\n💡 In production, you would send these metrics to:")
	fmt.Println("   - Prometheus (using prometheus client)")
	fmt.Println("   - SigNoz (using OTLP metric exporter)")
	fmt.Println("   - Custom monitoring system")
}
