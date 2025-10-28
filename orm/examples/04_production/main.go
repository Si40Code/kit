package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Si40Code/kit/logger"
	"github.com/Si40Code/kit/orm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:100;not null;index"`
	Email     string `gorm:"size:100;uniqueIndex"`
	Age       int    `gorm:"index"`
}

// ProductionMetricRecorder 生产环境 Metric 记录器
type ProductionMetricRecorder struct {
	// 在实际生产中，这里应该是 Prometheus/SigNoz 客户端
}

func NewProductionMetricRecorder() *ProductionMetricRecorder {
	return &ProductionMetricRecorder{}
}

func (r *ProductionMetricRecorder) RecordQuery(data orm.MetricData) {
	// 在生产环境中，这里应该：
	// 1. 发送到 Prometheus
	// 2. 发送到 SigNoz
	// 3. 发送到自定义监控系统

	// 示例：只打印关键信息
	if data.Error != nil || data.Duration > 100*time.Millisecond {
		fmt.Printf("⚠️  [Metric] %s on %s: %dms, rows=%d, error=%v\n",
			data.Operation,
			data.Table,
			data.Duration.Milliseconds(),
			data.RowsAffected,
			data.Error,
		)
	}
}

func main() {
	// ========================================
	// 1. 初始化 Logger (生产配置)
	// ========================================
	logger.Init(
		logger.WithLevel(logger.InfoLevel),   // 生产环境使用 Info 级别
		logger.WithFormat(logger.JSONFormat), // JSON 格式便于日志收集
		logger.WithStdout(),                  // 输出到 stdout
		logger.WithFile(
			"/var/log/app/app.log",       // 同时写入文件
			logger.WithFileMaxSize(100),  // 100MB
			logger.WithFileMaxAge(7),     // 保留 7 天
			logger.WithFileMaxBackups(3), // 保留 3 个备份
		),
		logger.WithTrace("my-service"), // 启用 trace 集成
		logger.WithCaller(true),        // 记录调用者信息
	)
	defer logger.Sync()

	// ========================================
	// 2. 初始化 OpenTelemetry (生产配置)
	// ========================================
	// 注意：在实际生产中，应该使用 OTLP exporter 连接到 SigNoz/Jaeger
	// import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	//
	// exporter, err := otlptracehttp.New(
	//     context.Background(),
	//     otlptracehttp.WithEndpoint("signoz:4318"),
	//     otlptracehttp.WithInsecure(),
	// )

	// 这里为了演示使用 stdout exporter
	exporter, err := stdouttrace.New()
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.TraceIDRatioBased(1.0)), // 生产环境建议 0.1 (10% 采样)
	)
	defer tp.Shutdown(context.Background())

	otel.SetTracerProvider(tp)

	ctx := context.Background()
	logger.Info(ctx, "Application starting", "version", "1.0.0")

	// ========================================
	// 3. 创建 ORM 客户端 (生产配置)
	// ========================================
	// 生产环境建议从配置文件或环境变量读取
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

	client, err := orm.New(
		mysql.Open(dsn),
		// 日志配置
		orm.WithLogger(logger.Default()),
		orm.WithSlowThreshold(200*time.Millisecond), // 200ms 为慢查询
		orm.WithIgnoreRecordNotFoundError(),         // 全局忽略 RecordNotFound

		// Trace 配置
		orm.WithTrace("orm-client"),

		// Metric 配置
		orm.WithMetric(NewProductionMetricRecorder()),

		// 连接池配置（重要！）
		orm.WithMaxIdleConns(10),                // 最大空闲连接数
		orm.WithMaxOpenConns(100),               // 最大打开连接数
		orm.WithConnMaxLifetime(time.Hour),      // 连接最大生命周期
		orm.WithConnMaxIdleTime(10*time.Minute), // 连接最大空闲时间
	)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", "error", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			logger.Error(ctx, "Failed to close database", "error", err)
		}
	}()

	logger.Info(ctx, "Database connected successfully",
		"max_idle_conns", 10,
		"max_open_conns", 100,
	)

	// ========================================
	// 4. 自动迁移
	// ========================================
	if err := client.AutoMigrate(&User{}); err != nil {
		logger.Fatal(ctx, "Failed to migrate database", "error", err)
	}

	// ========================================
	// 5. 业务逻辑示例
	// ========================================

	// 创建业务 span
	tracer := otel.Tracer("business-service")
	ctx, span := tracer.Start(ctx, "UserRegistration")
	defer span.End()

	logger.Info(ctx, "Starting user registration flow")

	// 5.1 创建用户（带重试逻辑）
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   28,
	}

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if err := client.WithContext(ctx).Create(&user).Error; err != nil {
			logger.Warn(ctx, "Create user failed, retrying",
				"attempt", i+1,
				"max_retries", maxRetries,
				"error", err,
			)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		break
	}

	if user.ID == 0 {
		logger.Error(ctx, "Failed to create user after retries")
		span.End()
		return
	}

	logger.Info(ctx, "User created successfully",
		"user_id", user.ID,
		"email", user.Email,
	)

	// 5.2 查询用户（使用忽略 RecordNotFound）
	var foundUser User
	if err := client.WithContext(ctx).First(&foundUser, user.ID).Error; err != nil {
		logger.Error(ctx, "Failed to find user", "error", err)
		return
	}

	logger.Info(ctx, "User found",
		"user_id", foundUser.ID,
		"name", foundUser.Name,
	)

	// 5.3 更新用户（事务）
	err = client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新年龄
		if err := tx.Model(&foundUser).Update("age", 29).Error; err != nil {
			return err
		}

		// 记录操作日志（示例）
		logger.Info(ctx, "User updated in transaction",
			"user_id", foundUser.ID,
			"old_age", user.Age,
			"new_age", 29,
		)

		return nil
	})
	if err != nil {
		logger.Error(ctx, "Transaction failed", "error", err)
		return
	}

	// 5.4 批量查询
	var users []User
	if err := client.WithContext(ctx).
		Where("age > ?", 25).
		Order("created_at DESC").
		Limit(10).
		Find(&users).Error; err != nil {
		logger.Error(ctx, "Failed to query users", "error", err)
		return
	}

	logger.Info(ctx, "Batch query completed",
		"count", len(users),
	)

	// 5.5 清理（删除测试数据）
	if err := client.WithContext(ctx).Delete(&foundUser).Error; err != nil {
		logger.Error(ctx, "Failed to delete user", "error", err)
	}

	logger.Info(ctx, "User registration flow completed successfully")

	// ========================================
	// 6. 健康检查示例
	// ========================================
	if err := checkDatabaseHealth(client, ctx); err != nil {
		logger.Error(ctx, "Database health check failed", "error", err)
	} else {
		logger.Info(ctx, "Database health check passed")
	}

	fmt.Println("\n✅ Production example completed successfully!")
	fmt.Println("\n📝 Best Practices Demonstrated:")
	fmt.Println("   ✓ Structured logging with JSON format")
	fmt.Println("   ✓ OpenTelemetry trace integration")
	fmt.Println("   ✓ Metric monitoring")
	fmt.Println("   ✓ Connection pool configuration")
	fmt.Println("   ✓ Slow query detection")
	fmt.Println("   ✓ Error handling with retries")
	fmt.Println("   ✓ Transaction support")
	fmt.Println("   ✓ Health check")
	fmt.Println("   ✓ Graceful shutdown")
}

// checkDatabaseHealth 数据库健康检查
func checkDatabaseHealth(client *orm.Client, ctx context.Context) error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 检查连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// 检查连接池状态
	stats := sqlDB.Stats()
	logger.Info(ctx, "Database connection pool stats",
		"open_connections", stats.OpenConnections,
		"in_use", stats.InUse,
		"idle", stats.Idle,
		"wait_count", stats.WaitCount,
		"wait_duration_ms", stats.WaitDuration.Milliseconds(),
	)

	// 如果等待时间过长，可能需要调整连接池配置
	if stats.WaitCount > 100 {
		logger.Warn(ctx, "High connection pool wait count, consider increasing max_open_conns",
			"wait_count", stats.WaitCount,
		)
	}

	return nil
}
