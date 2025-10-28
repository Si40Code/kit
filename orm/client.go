package orm

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Client ORM 客户端
type Client struct {
	*gorm.DB
	options *options
}

// New 创建一个新的 ORM 客户端
// dialector 是 GORM 的数据库方言，例如：
//   - mysql.Open(dsn)
//   - postgres.Open(dsn)
//   - sqlite.Open(dsn)
func New(dialector gorm.Dialector, opts ...Option) (*Client, error) {
	options := newOptions(opts...)

	// 创建 GORM 配置
	config := &gorm.Config{}

	// 配置日志
	if options.enableLog && options.logger != nil {
		config.Logger = newGormLogger(
			options.logger,
			options.slowThreshold,
			options.ignoreRecordNotFound,
		)
	}

	// 打开数据库连接
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 注册插件（用于 trace 和 metric）
	if options.enableTrace || options.enableMetric {
		if err := db.Use(&plugin{options: options}); err != nil {
			return nil, fmt.Errorf("failed to register plugin: %w", err)
		}
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(options.maxIdleConns)
	sqlDB.SetMaxOpenConns(options.maxOpenConns)
	sqlDB.SetConnMaxLifetime(options.connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(options.connMaxIdleTime)

	return &Client{
		DB:      db,
		options: options,
	}, nil
}

// WithContext 返回一个新的带有指定 context 的客户端实例
func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		DB:      c.DB.WithContext(ctx),
		options: c.options,
	}
}

// WithIgnoreRecordNotFound 返回一个新的客户端实例，该实例会忽略 RecordNotFound 错误
// 这个方法用于单次查询覆盖全局配置
func (c *Client) WithIgnoreRecordNotFound() *Client {
	return &Client{
		DB:      setIgnoreRecordNotFound(c.DB),
		options: c.options,
	}
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
