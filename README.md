

一个基于 Golang 和 Gin 框架的轻量级、插件化 API 服务。

## 🌟 功能特性

- **插件化架构**：支持动态添加和管理 API 插件
- **API 密钥管理**：支持生成、验证和管理 API 密钥
- **统计功能**：实时统计 API 请求次数和响应时间
- **多种 API 插件**：内置 IP 查询、Ping 测试、随机数生成等实用插件
- **跨域支持**：内置 CORS 中间件
- **速率限制**：防止 API 滥用
- **请求日志**：详细记录每个请求的信息
- **嵌入式资源**：静态文件和模板嵌入到二进制文件中
- **YAML 配置**：灵活的配置管理

## 🛠️ 技术栈

- **语言**：Golang
- **Web 框架**：Gin
- **配置管理**：YAML
- **数据库**：支持多种数据库（通过 GORM）
- **IP 库**：IP2Region
- **日志**：logrus

## 📦 插件列表

| 插件名称 | 功能描述 | 示例请求 |
|---------|---------|---------|
| **ip** | IP 地址信息查询 | `GET /api/ip?ip=114.114.114.114` |
| **ping** | 网络 Ping 测试 | `GET /api/ping?target=www.baidu.com&count=3` |
| **random** | 随机数生成 | `GET /api/random?min=1&max=100` |
| **client** | 客户端信息获取 | `GET /api/client` |
| **ipify** | 获取客户端公网 IP | `GET /api/ipify` |

## 🚀 快速开始

### 环境要求

- Golang 1.18 或更高版本

### 安装和运行

1. **克隆仓库**
   ```bash
   git clone https://github.com/xrcuo/xrcuo-api.git
   cd xrcuo-api
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **运行服务**
   ```bash
   go run main.go
   ```

4. **访问服务**
   - API 文档：http://localhost:8080
   - 统计页面：http://localhost:8080/stats
   - API 密钥管理：http://localhost:8080/api_key

### 构建二进制文件

```bash
go build -o xrcuo-api main.go
./xrcuo-api
```

## ⚙️ 配置说明

### 配置文件

项目使用 YAML 格式的配置文件，默认配置文件为 `config/default_config.yaml`。

### 主要配置项

| 配置项 | 类型 | 默认值 | 描述 |
|-------|------|-------|------|
| `server.port` | string | ":8080" | 服务监听端口 |
| `server.mode` | string | "release" | Gin 运行模式（debug/release/test） |
| `database.dsn` | string | "sqlite3:./data.db" | 数据库连接字符串 |
| `rate_limit.enable` | bool | true | 是否启用速率限制 |
| `rate_limit.rate` | int | 100 | 每分钟请求次数限制 |
| `stats.enable` | bool | true | 是否启用统计功能 |

### 自定义配置

你可以通过创建 `config.yaml` 文件来覆盖默认配置：

```yaml
server:
  port: ":8080"
  mode: "release"

database:
  dsn: "sqlite3:./data.db"
```

## 🔑 API 密钥管理

### 生成 API 密钥

1. 访问 http://localhost:8080/api_key
2. 点击 "生成新密钥" 按钮
3. 复制生成的 API 密钥

### 使用 API 密钥

在请求头中添加 `X-API-Key` 字段：

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/ip?ip=114.114.114.114
```

或者作为查询参数：

```bash
curl http://localhost:8080/api/ip?ip=114.114.114.114&api_key=your-api-key
```

## 📊 统计功能

### 查看统计信息

1. 访问 http://localhost:8080/stats
2. 查看 API 请求次数、响应时间等统计信息
3. 支持按 API 路径和状态码筛选

### API 统计接口

```bash
curl http://localhost:8080/api/stats
```

## 📝 API 文档

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

### 插件 API 详情

#### IP 查询 API

```
GET /api/ip?ip=114.114.114.114
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "114.114.114.114",
    "country": "中国",
    "region": "江苏",
    "city": "南京",
    "isp": "江苏省南京市 电信"
  }
}
```

#### Ping 测试 API

```
GET /api/ping?target=www.baidu.com&count=3
```

**参数：**
- `target`：目标主机名或 IP 地址
- `count`：Ping 次数（默认 3）

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "target": "www.baidu.com",
    "count": 3,
    "results": [
      {"seq": 1, "ttl": 54, "time": "12.345ms"},
      {"seq": 2, "ttl": 54, "time": "11.234ms"},
      {"seq": 3, "ttl": 54, "time": "13.456ms"}
    ],
    "avg_time": "12.345ms"
  }
}
```

#### 随机数生成 API

```
GET /api/random?min=1&max=100
```

**参数：**
- `min`：最小值（默认 0）
- `max`：最大值（默认 100）

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "random": 42
  }
}
```

## 📁 项目结构

```
xrcuo-api/
├── common/          # 公共工具和中间件
├── config/          # 配置管理
├── db/              # 数据库操作
├── models/          # 数据模型
├── plugin/          # 插件目录
│   ├── ip/          # IP 查询插件
│   ├── ping/        # Ping 测试插件
│   ├── random/      # 随机数插件
│   └── ...          # 其他插件
├── static/          # 静态资源
├── templates/       # HTML 模板
├── config.yaml      # 配置文件
├── go.mod           # Go 模块
├── main.go          # 入口文件
└── README.md        # 项目文档
```

## 🔧 开发指南

### 添加新插件

1. 在 `plugin/` 目录下创建新的插件目录
2. 实现 `Plugin` 接口
3. 在 `main.go` 的 `registerRoutes` 函数中注册插件

**插件示例：**

```go
package myplugin

import "github.com/gin-gonic/gin"

// MyPlugin 定义插件
var MyPlugin = &plugin.Plugin{
    Name:        "myplugin",
    Description: "我的插件",
    Register: func(rg *gin.RouterGroup) {
        rg.GET("/myplugin", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "Hello from my plugin",
            })
        })
    },
}
```

### 测试插件

```bash
go run main.go
curl http://localhost:8080/api/myplugin?api_key=your-api-key
```

## 🚀 部署方式

### 二进制部署

1. 构建二进制文件
   ```bash
   GOOS=linux GOARCH=amd64 go build -o xrcuo-api-linux main.go
   ```

2. 上传到服务器
   ```bash
   scp xrcuo-api-linux user@server:/path/to/directory
   ```

3. 运行服务
   ```bash
   ./xrcuo-api-linux
   ```

### Docker 部署

（待添加 Docker 支持）

## 📋 许可证

MIT License

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 `git checkout -b feature/AmazingFeature`
3. 提交更改 `git commit -m 'Add some AmazingFeature'`
4. 推送到分支 `git push origin feature/AmazingFeature`
5. 提交 Pull Request

## 📞 联系方式

如有问题或建议，欢迎提交 Issue 或 Pull Request。

---

**Star ⭐ 支持一下吧！**