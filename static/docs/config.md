# 配置说明

## 配置文件

项目使用 YAML 格式的配置文件，默认配置文件为 `config/default_config.yaml`。

## 主要配置项

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `server.port` | string | ":8080" | 服务监听端口 |
| `server.mode` | string | "release" | Gin 运行模式（debug/release/test） |
| `database.dsn` | string | "sqlite3:./data.db" | 数据库连接字符串 |
| `rate_limit.enable` | bool | true | 是否启用速率限制 |
| `rate_limit.rate` | int | 100 | 每分钟请求次数限制 |
| `stats.enable` | bool | true | 是否启用统计功能 |

## 自定义配置

你可以通过创建 `config.yaml` 文件来覆盖默认配置：

```yaml
server:
  port: ":8080"
  mode: "release"

database:
  dsn: "sqlite3:./data.db"

rate_limit:
  enable: true
  rate: 100

stats:
  enable: true
```

## 配置优先级

配置文件的优先级从高到低依次为：

1. 命令行参数（如果支持）
2. `config.yaml`（用户自定义配置）
3. `config/default_config.yaml`（默认配置）

## 数据库配置

### SQLite

```yaml
database:
  dsn: "sqlite3:./data.db"
```

### MySQL

```yaml
database:
  dsn: "mysql://username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

### PostgreSQL

```yaml
database:
  dsn: "postgres://username:password@localhost:5432/dbname?sslmode=disable"
```

## 速率限制配置

```yaml
rate_limit:
  enable: true
  rate: 100  # 每分钟请求次数
  burst: 200  # 突发请求数
```

## 统计配置

```yaml
stats:
  enable: true
  retention_days: 30  # 统计数据保留天数
```