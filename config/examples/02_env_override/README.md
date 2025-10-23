# 示例 2: 环境变量覆盖

展示如何使用环境变量覆盖文件配置，适用于不同环境（开发/测试/生产）。

## 运行示例

```bash
cd config/examples/02_env_override
go run main.go
```

## 学习内容

1. **环境变量覆盖** - 环境变量优先级高于文件配置
2. **环境变量命名规则** - `PREFIX_KEY_PATH`
3. **实际应用场景** - 不同环境使用不同配置
4. **敏感信息处理** - 密码等敏感信息通过环境变量传递

## 环境变量命名规则

格式：`PREFIX_KEY_PATH`

示例：
- `APP_SERVER_PORT=9090` → `server.port = 9090`
- `APP_DATABASE_HOST=prod.db.com` → `database.host = "prod.db.com"`
- `APP_APP_DEBUG=false` → `app.debug = false`

## 实际应用场景

### 开发环境
```bash
export APP_DATABASE_HOST=localhost
export APP_APP_DEBUG=true
```

### 生产环境
```bash
export APP_DATABASE_HOST=prod-db.example.com
export APP_DATABASE_PASSWORD=<从密钥管理系统获取>
export APP_APP_DEBUG=false
```

## 配置优先级

从低到高：
1. 配置文件 (config.yaml)
2. 环境变量 (APP_*) ⬅️ 优先级更高

## 最佳实践

- ✅ 在配置文件中设置默认值和开发环境配置
- ✅ 通过环境变量覆盖特定环境的配置
- ✅ 敏感信息（密码、密钥）始终使用环境变量
- ❌ 不要在配置文件中存储生产环境的敏感信息

