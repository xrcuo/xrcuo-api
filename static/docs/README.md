# Xrcuo API 文档

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
   - API 文档：http://localhost:8080/docs
   - 统计页面：http://localhost:8080/stats
   - API 密钥管理：http://localhost:8080/api_key

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

## 📝 通用 API 格式

### 请求格式

```
GET /api/{plugin-name}?{params}
```

### 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
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