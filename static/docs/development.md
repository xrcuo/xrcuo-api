# 开发指南

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
│   └── ...          # 其他插件
├── static/          # 静态资源
├── templates/       # HTML 模板
├── config.yaml      # 配置文件
├── go.mod           # Go 模块
├── main.go          # 入口文件
└── README.md        # 项目文档
```

## 添加新插件

### 1. 创建插件目录

在 `plugin/` 目录下创建新的插件目录，例如 `myplugin`。

### 2. 实现 Plugin 接口

每个插件都需要实现 `Plugin` 接口，该接口定义在 `plugin/plugin.go` 文件中。

```go
package myplugin

import (
    "github.com/gin-gonic/gin"
    "github.com/xrcuo/xrcuo-api/plugin"
)

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

### 3. 注册插件

在 `main.go` 的 `registerRoutes` 函数中注册插件：

```go
// 注册各个插件
pluginManager.Register(ip.IPPlugin)
pluginManager.Register(ping.PingPlugin)
pluginManager.Register(random.RandomPlugin)
pluginManager.Register(client.ClientPlugin)
pluginManager.Register(ipify.IpifyPlugin)
pluginManager.Register(myplugin.MyPlugin) // 添加这一行
```

### 4. 测试插件

运行服务并测试新插件：

```bash
go run main.go
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/myplugin
```

## 插件开发规范

1. 每个插件应该有自己的目录，包含独立的代码文件
2. 插件应该实现 `Plugin` 接口
3. 插件路由应该挂载在 `/api` 路径下
4. 插件应该遵循 RESTful API 设计规范
5. 插件应该返回统一的响应格式

## 测试

### 单元测试

为插件编写单元测试，确保功能正常：

```go
package myplugin

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestMyPlugin(t *testing.T) {
    // 设置 Gin 模式为测试模式
    gin.SetMode(gin.TestMode)

    // 创建 Gin 引擎
    r := gin.Default()

    // 注册插件路由
    MyPlugin.Register(r.Group("/api"))

    // 创建测试请求
    req, _ := http.NewRequest("GET", "/api/myplugin", nil)
    req.Header.Set("X-API-Key", "test-key")

    // 执行请求
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    // 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"message":"Hello from my plugin"}`, w.Body.String())
}
```

### 集成测试

运行完整的服务并测试所有功能：

```bash
go run main.go
# 使用 curl 或其他工具测试 API
```

## 构建和部署

### 构建二进制文件

```bash
go build -o xrcuo-api main.go
```

### 构建跨平台二进制文件

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o xrcuo-api-linux main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o xrcuo-api-darwin main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o xrcuo-api-windows.exe main.go
```

### 部署

将构建好的二进制文件上传到服务器并运行：

```bash
./xrcuo-api
```