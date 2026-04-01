# xrcuo-api

一个基于 Golang 和 Gin 框架的轻量级、插件化 API 服务。

## 功能特性

- **插件化架构**：支持动态添加和管理 API 插件
- **统计功能**：实时统计 API 请求次数和响应时间
- **多种 API 插件**：内置 IP 查询、Ping 测试、随机数生成等实用插件
- **跨域支持**：内置 CORS 中间件
- **速率限制**：防止 API 滥用
- **请求日志**：详细记录每个请求的信息
- **嵌入式资源**：静态文件和模板嵌入到二进制文件中
- **YAML 配置**：灵活的配置管理

## 技术栈

- **语言**：Golang
- **Web 框架**：Gin
- **配置管理**：YAML
- **数据库**：支持 SQLite（默认）、MySQL、PostgreSQL
- **IP 库**：IP2Region
- **日志**：logrus

## 插件列表

| 插件名称 | 功能描述 | 示例请求 |
|---------|---------|---------|
| **ip** | IP 地址信息查询 | `GET /api/ip?ip=114.114.114.114` |
| **ping** | 网络 Ping 测试 | `GET /api/ping?target=www.baidu.com&count=3` |
| **random** | 随机数生成 | `GET /api/random?min=1&max=100` |
| **client** | 客户端信息获取 | `GET /api/client` |
| **ipify** | 获取客户端公网 IP | `GET /api/ipify` |

## 快速开始

### 环境要求

- Golang 1.18 或更高版本

### 安装和运行

```bash
# 克隆仓库
git clone https://github.com/xrcuo/xrcuo-api.git
cd xrcuo-api

# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

### 访问服务

- API 文档：http://localhost:8080
- 统计页面：http://localhost:8080/stats

### 构建二进制文件

```bash
go build -o xrcuo-api main.go
./xrcuo-api
```

## 配置说明

项目使用 YAML 格式的配置文件，默认配置文件为 `config/default_config.yaml`。

### 数据库配置

项目支持三种数据库：SQLite（默认）、MySQL 和 PostgreSQL。通过修改 `config.yaml` 中的 `database.type` 字段即可切换。

#### SQLite 配置（默认）

```yaml
database:
  type: "sqlite"
  path: "./stats.db"
  max_open_conns: 10
  max_idle_conns: 5
```

#### MySQL 配置

```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "xrcuo_api"
  max_open_conns: 10
  max_idle_conns: 5
```

使用前需创建数据库：
```sql
CREATE DATABASE xrcuo_api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### PostgreSQL 配置

```yaml
database:
  type: "postgresql"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password"
  dbname: "xrcuo_api"
  max_open_conns: 10
  max_idle_conns: 5
```

使用前需创建数据库：
```sql
CREATE DATABASE xrcuo_api;
```

### 主要配置项

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `server.port` | string | ":8080" | 服务监听端口 |
| `server.mode` | string | "debug" | Gin 运行模式（debug/release/test） |
| `database.type` | string | "sqlite" | 数据库类型：sqlite, mysql, postgresql |
| `database.path` | string | "./stats.db" | SQLite 数据库文件路径 |
| `database.host` | string | "localhost" | 数据库主机（MySQL/PostgreSQL） |
| `database.port` | int | 3306 | 数据库端口 |
| `database.user` | string | "" | 数据库用户名（MySQL/PostgreSQL） |
| `database.password` | string | "" | 数据库密码（MySQL/PostgreSQL） |
| `database.dbname` | string | "xrcuo_api" | 数据库名称（MySQL/PostgreSQL） |

### 自定义配置

你可以通过创建 `config.yaml` 文件来覆盖默认配置：

```yaml
server:
  port: ":8080"
  mode: "release"

database:
  type: "sqlite"
  path: "./stats.db"
```

## 统计功能

### 查看统计信息

1. 访问 http://localhost:8080/stats
2. 查看 API 请求次数、响应时间等统计信息
3. 支持按 API 路径和状态码筛选

### API 统计接口

```bash
curl http://localhost:8080/api/stats
```

## API 文档

### 通用 API 格式

#### 请求格式

```
GET /api/{plugin-name}?{params}
```

#### 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## 项目结构

```
xrcuo-api/
├── common/          # 公共工具和中间件
├── config/          # 配置管理
├── db/              # 数据库操作
├── models/          # 数据模型
├── plugin/          # 插件目录
├── static/          # 静态资源
├── templates/       # HTML 模板
├── main.go          # 入口文件
└── README.md        # 项目文档
```

## 开发指南

### 添加新插件

1. 在 `plugin/` 目录下创建新的插件目录
2. 实现 `Plugin` 接口
3. 在 `plugin/plugin.go` 的 `RegisterBuiltinPlugins` 函数中注册插件

### 测试插件

```bash
go run main.go
curl http://localhost:8080/api/myplugin
```

## 部署方式

### 二进制部署

```bash
# 构建二进制文件
GOOS=linux GOARCH=amd64 go build -o xrcuo-api-linux main.go

# 上传到服务器
scp xrcuo-api-linux user@server:/path/to/directory

# 运行服务
./xrcuo-api-linux
```

### Docker 部署

（待添加 Docker 支持）

## 许可证

MIT License

## 贡献

1. Fork 本仓库
2. 创建特性分支 `git checkout -b feature/AmazingFeature`
3. 提交更改 `git commit -m 'Add some AmazingFeature'`
4. 推送到分支 `git push origin feature/AmazingFeature`
5. 提交 Pull Request

---

**Star ⭐ 支持一下吧！**
