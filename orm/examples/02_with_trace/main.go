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

func main() {
	// 1. 初始化 logger
	logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
		logger.WithTrace("orm-example"), // 启用 trace 集成
	)
	defer logger.Sync()

	// 2. 初始化 OpenTelemetry（使用 stdout exporter 用于演示）
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	defer tp.Shutdown(context.Background())

	otel.SetTracerProvider(tp)

	ctx := context.Background()
	logger.Info(ctx, "OpenTelemetry initialized")

	// 3. 创建 ORM 客户端（请替换为你的数据库连接信息）
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := orm.New(
		mysql.Open(dsn),
		orm.WithLogger(logger.Default()),
		orm.WithTrace("orm-client"), // 启用 trace
		orm.WithSlowThreshold(100*time.Millisecond),
	)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", "error", err)
	}
	defer client.Close()

	logger.Info(ctx, "Successfully connected to database")

	// 4. 自动迁移
	if err := client.AutoMigrate(&User{}); err != nil {
		logger.Fatal(ctx, "Failed to migrate database", "error", err)
	}

	// 5. 创建业务操作的 span
	tracer := otel.Tracer("business-logic")
	ctx, span := tracer.Start(ctx, "CreateAndQueryUser")
	defer span.End()

	// 6. 在业务 span 下执行数据库操作
	logger.Info(ctx, "Starting business operation: CreateAndQueryUser")

	// 创建用户（会自动创建子 span）
	user := User{
		Name:  "Bob",
		Email: "bob@example.com",
		Age:   30,
	}

	if err := client.WithContext(ctx).Create(&user).Error; err != nil {
		logger.Error(ctx, "Failed to create user", "error", err)
		return
	}
	logger.Info(ctx, "User created", "id", user.ID)

	// 查询用户（会自动创建另一个子 span）
	var foundUser User
	if err := client.WithContext(ctx).First(&foundUser, user.ID).Error; err != nil {
		logger.Error(ctx, "Failed to find user", "error", err)
		return
	}
	logger.Info(ctx, "User found", "name", foundUser.Name)

	// 更新用户（又一个子 span）
	if err := client.WithContext(ctx).Model(&foundUser).Update("age", 31).Error; err != nil {
		logger.Error(ctx, "Failed to update user", "error", err)
		return
	}
	logger.Info(ctx, "User updated")

	// 7. 创建另一个独立的业务操作
	ctx2, span2 := tracer.Start(context.Background(), "BatchQueryUsers")
	defer span2.End()

	logger.Info(ctx2, "Starting business operation: BatchQueryUsers")

	var users []User
	if err := client.WithContext(ctx2).Where("age > ?", 25).Find(&users).Error; err != nil {
		logger.Error(ctx2, "Failed to find users", "error", err)
		return
	}
	logger.Info(ctx2, "Users found", "count", len(users))

	// 删除用户
	if err := client.WithContext(ctx2).Delete(&foundUser).Error; err != nil {
		logger.Error(ctx2, "Failed to delete user", "error", err)
		return
	}

	fmt.Println("\n✅ All operations completed successfully!")
	fmt.Println("📊 Check the trace output above to see the span hierarchy:")
	fmt.Println("   - Business spans (CreateAndQueryUser, BatchQueryUsers)")
	fmt.Println("   - Database spans (DB INSERT, DB SELECT, DB UPDATE, DB DELETE)")
	fmt.Println("   - Each database span includes:")
	fmt.Println("     * db.operation (SELECT/INSERT/UPDATE/DELETE)")
	fmt.Println("     * db.table (table name)")
	fmt.Println("     * db.statement (SQL)")
	fmt.Println("     * db.rows_affected (affected rows)")
	fmt.Println("     * db.duration_ms (execution time)")
}
