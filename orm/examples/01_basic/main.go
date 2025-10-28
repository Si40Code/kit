package main

import (
	"context"
	"fmt"
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

func main() {
	// 1. 初始化 logger
	logger.Init(
		logger.WithLevel(logger.InfoLevel),
		logger.WithFormat(logger.ConsoleFormat),
		logger.WithStdout(),
	)
	defer logger.Sync()

	ctx := context.Background()

	// 2. 创建 ORM 客户端（请替换为你的数据库连接信息）
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := orm.New(
		mysql.Open(dsn),
		orm.WithLogger(logger.Default()),
		orm.WithSlowThreshold(100*time.Millisecond), // 100ms 为慢查询
	)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", "error", err)
	}
	defer client.Close()

	logger.Info(ctx, "Successfully connected to database")

	// 3. 自动迁移（创建表）
	if err := client.AutoMigrate(&User{}); err != nil {
		logger.Fatal(ctx, "Failed to migrate database", "error", err)
	}

	// 4. 创建用户
	logger.Info(ctx, "Creating user...")
	user := User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   25,
	}

	if err := client.WithContext(ctx).Create(&user).Error; err != nil {
		logger.Error(ctx, "Failed to create user", "error", err)
	} else {
		logger.Info(ctx, "User created successfully", "id", user.ID)
	}

	// 5. 查询单个用户
	logger.Info(ctx, "Querying user by ID...")
	var foundUser User
	if err := client.WithContext(ctx).First(&foundUser, user.ID).Error; err != nil {
		logger.Error(ctx, "Failed to find user", "error", err)
	} else {
		logger.Info(ctx, "User found", "name", foundUser.Name, "email", foundUser.Email)
	}

	// 6. 查询所有用户
	logger.Info(ctx, "Querying all users...")
	var users []User
	if err := client.WithContext(ctx).Find(&users).Error; err != nil {
		logger.Error(ctx, "Failed to find users", "error", err)
	} else {
		logger.Info(ctx, "Users found", "count", len(users))
	}

	// 7. 更新用户
	logger.Info(ctx, "Updating user...")
	if err := client.WithContext(ctx).Model(&foundUser).Update("age", 26).Error; err != nil {
		logger.Error(ctx, "Failed to update user", "error", err)
	} else {
		logger.Info(ctx, "User updated successfully")
	}

	// 8. 条件查询
	logger.Info(ctx, "Querying users with condition...")
	var youngUsers []User
	if err := client.WithContext(ctx).Where("age < ?", 30).Find(&youngUsers).Error; err != nil {
		logger.Error(ctx, "Failed to find young users", "error", err)
	} else {
		logger.Info(ctx, "Young users found", "count", len(youngUsers))
	}

	// 9. 删除用户
	logger.Info(ctx, "Deleting user...")
	if err := client.WithContext(ctx).Delete(&foundUser).Error; err != nil {
		logger.Error(ctx, "Failed to delete user", "error", err)
	} else {
		logger.Info(ctx, "User deleted successfully")
	}

	// 10. 测试查询不存在的记录（会报错）
	logger.Info(ctx, "Querying non-existent user (will return error)...")
	var notFoundUser User
	if err := client.WithContext(ctx).First(&notFoundUser, 99999).Error; err != nil {
		logger.Warn(ctx, "User not found (expected error)", "error", err)
	}

	// 11. 使用 WithIgnoreRecordNotFound 忽略 RecordNotFound 错误
	logger.Info(ctx, "Querying non-existent user (with ignore flag)...")
	if err := client.WithContext(ctx).WithIgnoreRecordNotFound().First(&notFoundUser, 99999).Error; err != nil {
		logger.Error(ctx, "Unexpected error", "error", err)
	} else {
		logger.Info(ctx, "Query completed without error (RecordNotFound ignored)")
	}

	fmt.Println("\n✅ All operations completed successfully!")
}
