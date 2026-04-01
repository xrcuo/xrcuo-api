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
  type: "sqlite"  # 数据库类型：sqlite, mysql, postgresql
  path: "./stats.db"  # SQLite数据库文件路径（仅sqlite使用）
  host: "localhost"  # 数据库主机（mysql/postgresql使用）
  port: 3306  # 数据库端口（mysql默认3306，postgresql默认5432）
  user: ""  # 数据库用户名（mysql/postgresql使用）
  password: ""  # 数据库密码（mysql/postgresql使用）
  dbname: "xrcuo_api"  # 数据库名称（mysql/postgresql使用）
  max_open_conns: 10  # 最大打开连接数
  max_idle_conns: 5  # 最大空闲连接数
```

## 主要配置项

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `server.port` | string | ":8080" | 服务监听端口 |
| `server.mode` | string | "debug" | Gin运行模式（debug, release, test） |
| `server.json_format.enabled` | bool | false | 是否启用格式化JSON响应 |
| `ip2region.v4_db_path` | string | "./ip2region_v4.xdb" | IPv4数据库路径 |
| `ip2region.v6_db_path` | string | "./ip2region_v6.xdb" | IPv6数据库路径 |
| `log.level` | string | "info" | 日志级别 |
| `log.request_log` | bool | true | 是否输出请求日志 |
| `random_image.local_enabled` | bool | false | 是否启用本地图片 |
| `random_image.local_path` | string | "images/" | 本地图片目录路径 |
| `download.path` | string | "downloads/" | 下载文件目录路径 |
| `database.type` | string | "sqlite" | 数据库类型：sqlite, mysql, postgresql |
| `database.path` | string | "./stats.db" | SQLite数据库文件路径 |
| `database.host` | string | "localhost" | 数据库主机（mysql/postgresql使用） |
| `database.port` | int | 3306 | 数据库端口（mysql默认3306，postgresql默认5432） |
| `database.user` | string | "" | 数据库用户名（mysql/postgresql使用） |
| `database.password` | string | "" | 数据库密码（mysql/postgresql使用） |
| `database.dbname` | string | "xrcuo_api" | 数据库名称（mysql/postgresql使用） |
| `database.max_open_conns` | int | 10 | 最大打开连接数 |
| `database.max_idle_conns` | int | 5 | 最大空闲连接数 |

## 配置优先级

配置文件的优先级从高到低依次为：

1. 命令行参数（如果支持）
2. `config.yaml`（用户自定义配置）
3. `config/default_config.yaml`（默认配置）

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

项目支持三种数据库：SQLite（默认）、MySQL 和 PostgreSQL。

## SQLite 配置（默认）

```yaml
database:
  type: "sqlite"  # 数据库类型
  path: "./stats.db"  # SQLite数据库文件路径
  max_open_conns: 10  # 最大打开连接数
  max_idle_conns: 5  # 最大空闲连接数
```

## MySQL 配置

```yaml
database:
  type: "mysql"  # 数据库类型
  host: "localhost"  # 数据库主机
  port: 3306  # 数据库端口（MySQL默认3306）
  user: "root"  # 数据库用户名
  password: "your_password"  # 数据库密码
  dbname: "xrcuo_api"  # 数据库名称
  max_open_conns: 10  # 最大打开连接数
  max_idle_conns: 5  # 最大空闲连接数
```

使用 MySQL 前需要：
1. 安装 MySQL 服务
2. 创建数据库：`CREATE DATABASE xrcuo_api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
3. 配置正确的用户名和密码

## PostgreSQL 配置

```yaml
database:
  type: "postgresql"  # 数据库类型
  host: "localhost"  # 数据库主机
  port: 5432  # 数据库端口（PostgreSQL默认5432）
  user: "postgres"  # 数据库用户名
  password: "your_password"  # 数据库密码
  dbname: "xrcuo_api"  # 数据库名称
  max_open_conns: 10  # 最大打开连接数
  max_idle_conns: 5  # 最大空闲连接数
```

使用 PostgreSQL 前需要：
1. 安装 PostgreSQL 服务
2. 创建数据库：`CREATE DATABASE xrcuo_api;`
3. 配置正确的用户名和密码