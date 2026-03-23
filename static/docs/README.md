# Xrcuo API 文档

一个基于 Golang 和 Gin 框架的轻量级、插件化 API 服务。

## 功能特性

- **插件化架构**：支持动态添加和管理 API 插件
- **API 密钥管理**：支持生成、验证和管理 API 密钥
- **统计功能**：实时统计 API 请求次数和响应时间
- **多种 API 插件**：内置 IP 查询、Ping 测试、随机数生成、MCPE服务器查询等实用插件
- **跨域支持**：内置 CORS 中间件
- **速率限制**：防止 API 滥用
- **请求日志**：详细记录每个请求的信息
- **嵌入式资源**：静态文件和模板嵌入到二进制文件中
- **YAML 配置**：灵活的配置管理
- **配置文件热重载**：修改配置无需重启服务
- **性能监控**：实时监控 QPS、响应时间等性能指标

## 快速开始

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

3. **编译运行**
   ```bash
   go build -o xrcuo-api.exe
   ./xrcuo-api.exe
   ```

4. **访问服务**
   - API 文档：http://localhost:8080/docs
   - 统计页面：http://localhost:8080/stats
   - API 密钥管理：http://localhost:8080/api_key

## API 列表

所有需要 API 密钥的接口都挂载在 `/api/` 路径下，文件下载接口无需密钥。

| 接口 | 路径 | 描述 | 需要密钥 |
|------|------|------|----------|
| IP查询 | `/api/ip` | 查询IP归属地信息 | ✅ |
| Ping测试 | `/api/ping` | 对目标进行Ping测试 | ✅ |
| 随机数生成 | `/api/random` | 生成指定范围随机整数 | ✅ |
| 随机图片 | `/api/random/image` | 获取随机图片 | ✅ |
| 随机图片信息 | `/api/random/image/info` | 获取随机图片信息 | ✅ |
| 客户端信息 | `/api/client` | 获取客户端详细信息 | ✅ |
| 公网IP | `/api/ipify` | 获取客户端公网IP | ✅ |
| MCPE查询 | `/api/mcpe/status` | 查询Minecraft PE服务器状态 | ✅ |
| 文件下载 | `/download/{file}` | 下载本地文件 | ❌ |
| API下载 | `/api/download/{file}` | 下载本地文件 | ❌ |

## API 密钥管理

### 生成 API 密钥

1. 访问 http://localhost:8080/api_key
2. 点击 "生成新密钥" 按钮
3. 复制生成的 API 密钥

### 使用 API 密钥

在请求头中添加 `Authorization` 字段：

```bash
curl -H "Authorization: your-api-key" http://localhost:8080/api/ip?ip=114.114.114.114
```

或者作为查询参数：

```bash
curl http://localhost:8080/api/ip?ip=114.114.114.114&api_key=your-api-key
```

## 通用 API 格式

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

## 项目结构

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
│   ├── mcpe/        # MCPE服务器查询插件
│   ├── client/      # 客户端信息插件
│   ├── ipify/       # 公网IP插件
│   ├── api_key/     # API密钥管理插件
│   └── download/    # 文件下载插件
├── static/          # 静态资源
├── templates/       # HTML 模板
├── config.yaml      # 配置文件
├── go.mod           # Go 模块
├── main.go          # 入口文件
└── README.md        # 项目文档
```

## 配置说明

配置文件 `config.yaml` 包含以下配置项：

- **server**: 服务端口和模式
- **database**: 数据库配置
- **ip2region**: IP地址库配置
- **log**: 日志配置
- **random_image**: 随机图片配置
- **download**: 下载目录配置

详细配置说明请参阅 [配置说明](config.md)
