# 配置说明

## 配置文件

项目使用 YAML 格式的配置文件，默认配置文件为 `config/default_config.yaml`。

## 完整配置示例

```yaml
# 服务器配置
server:
  port: ":8080"  # 服务监听端口
  mode: "debug"  # Gin运行模式（debug, release, test）
  json_format:
    enabled: false  # 是否启用格式化JSON响应

# API密钥配置
api_key:
  enabled: true  # 是否启用API密钥验证（默认启用，关闭后所有API无需密钥即可访问）

# IP2Region配置
ip2region:
  v4_db_path: "./ip2region_v4.xdb"  # IPv4数据库路径
  v6_db_path: "./ip2region_v6.xdb"  # IPv6数据库路径

# 日志配置
log:
  level: "info"  # 日志级别（debug, info, warn, error）
  file: "logs/app.log"  # 日志文件路径
  console_output: true  # 是否输出到控制台
  request_log: true  # 是否输出请求日志
  max_size: 10  # 单个日志文件最大大小（MB）
  max_backups: 5  # 保留的日志文件数量
  max_age: 7  # 日志文件保留天数

# 随机图片配置
random_image:
  local_enabled: false  # 是否启用本地图片
  local_path: "images/"  # 本地图片目录路径

# 下载配置
download:
  path: "downloads/"  # 下载文件目录路径

# 数据库配置
database:
  path: "./stats.db"  # SQLite数据库文件路径
  max_open_conns: 10  # 最大打开连接数
  max_idle_conns: 5  # 最大空闲连接数
```

## 主要配置项

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `server.port` | string | ":8080" | 服务监听端口 |
| `server.mode` | string | "debug" | Gin运行模式（debug, release, test） |
| `server.json_format.enabled` | bool | false | 是否启用格式化JSON响应 |
| `api_key.enabled` | bool | true | 是否启用API密钥验证 |
| `ip2region.v4_db_path` | string | "./ip2region_v4.xdb" | IPv4数据库路径 |
| `ip2region.v6_db_path` | string | "./ip2region_v6.xdb" | IPv6数据库路径 |
| `log.level` | string | "info" | 日志级别 |
| `log.request_log` | bool | true | 是否输出请求日志 |
| `random_image.local_enabled` | bool | false | 是否启用本地图片 |
| `random_image.local_path` | string | "images/" | 本地图片目录路径 |
| `download.path` | string | "downloads/" | 下载文件目录路径 |
| `database.path` | string | "./stats.db" | SQLite数据库文件路径 |
| `database.max_open_conns` | int | 10 | 最大打开连接数 |
| `database.max_idle_conns` | int | 5 | 最大空闲连接数 |

## 配置优先级

配置文件的优先级从高到低依次为：

1. 命令行参数（如果支持）
2. `config.yaml`（用户自定义配置）
3. `config/default_config.yaml`（默认配置）

## API密钥配置

```yaml
api_key:
  enabled: true  # true=启用验证，false=禁用验证
```

- `enabled: true`（默认）- 所有API请求需要有效的API密钥
- `enabled: false` - 所有API请求无需API密钥即可访问（适用于内网或不需要认证的场景）

## 日志配置

```yaml
log:
  level: "info"  # 日志级别
  file: "logs/app.log"  # 日志文件路径
  console_output: true  # 是否输出到控制台
  request_log: true  # 是否输出请求日志
  max_size: 10  # 单个日志文件最大大小（MB）
  max_backups: 5  # 保留的日志文件数量
  max_age: 7  # 日志文件保留天数
```

## 随机图片配置

```yaml
random_image:
  local_enabled: true  # 是否启用本地图片
  local_path: "images/"  # 本地图片目录路径
```

## 数据库配置

```yaml
database:
  path: "./stats.db"  # SQLite数据库文件路径
  max_open_conns: 10  # 最大打开连接数
  max_idle_conns: 5  # 最大空闲连接数
```