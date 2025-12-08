package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/common"
	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/plugin"
	"github.com/xrcuo/xrcuo-api/plugin/api_key"
	"github.com/xrcuo/xrcuo-api/plugin/client"
	"github.com/xrcuo/xrcuo-api/plugin/ip"
	"github.com/xrcuo/xrcuo-api/plugin/ipify"
	"github.com/xrcuo/xrcuo-api/plugin/ping"
	"github.com/xrcuo/xrcuo-api/plugin/random"
)

//go:embed static
//go:embed templates
var embeddedFiles embed.FS

// 初始化应用
func initApp() {
	// 解析配置文件
	config.Parse()

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		logrus.Fatalf("数据库初始化失败：%v", err)
	}

	// 预加载IP2Region数据库
	if err := common.InitIP2Region(); err != nil {
		logrus.Fatalf("IP2Region数据库初始化失败：%v", err)
	}

	// 初始化统计信息
	common.InitStats()
}

// 设置Gin引擎和中间件
func setupGin() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.GetServerMode())

	// 创建Gin引擎实例
	r := gin.Default()

	// 信任所有代理，确保能正确获取客户端真实IP
	r.SetTrustedProxies(nil)

	// 添加全局中间件
	r.Use(common.RequestLoggerMiddleware()) // 请求日志中间件
	r.Use(common.CORSMiddleware())          // 跨域中间件
	r.Use(common.RateLimitMiddleware())     // 速率限制中间件

	return r
}

// 设置模板
func setupTemplates(r *gin.Engine) {
	// 设置模板自定义函数
	funcMap := template.FuncMap{
		"percentage": func(total, count int64) string {
			if total == 0 {
				return "0%"
			}
			return fmt.Sprintf("%d%%", int((float64(count)/float64(total))*100))
		},
	}

	// 从嵌入式文件系统加载模板
	tmpls, err := template.New("").Funcs(funcMap).ParseFS(embeddedFiles, "templates/*")
	if err != nil {
		logrus.Fatalf("加载模板失败：%v", err)
	}

	// 设置Gin使用解析后的模板
	r.SetHTMLTemplate(tmpls)
}

// 设置静态文件服务
func setupStaticFiles(r *gin.Engine) {
	// 添加静态文件服务，用于提供本地图片（如果需要）
	r.Static("/images", "./images")

	// 从嵌入式文件系统中获取static子目录
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		logrus.Fatalf("获取static子目录失败：%v", err)
	}

	// 使用嵌入式文件系统提供静态资源
	r.StaticFS("/static", http.FS(staticFS))

	// 使用嵌入式文件系统提供favicon.ico
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(staticFS))
	})
}

// 注册路由
func registerRoutes(r *gin.Engine) {
	// 创建插件管理器
	pluginManager := plugin.NewPluginManager()

	// 注册各个插件
	pluginManager.Register(ip.IPPlugin)
	pluginManager.Register(ping.PingPlugin)
	pluginManager.Register(random.RandomPlugin)
	pluginManager.Register(client.ClientPlugin)
	pluginManager.Register(ipify.IpifyPlugin)
	// 后续新增插件，只需在这里添加注册语句即可

	// 注册API根路由（所有插件路由都挂载在/api下）
	apiGroup := r.Group("/api")
	{
		// 为所有API注册统计中间件
		apiGroup.Use(common.StatsMiddleware())
		// 为所有API注册API密钥验证中间件
		apiGroup.Use(common.APIKeyMiddleware())
		// 使用插件管理器注册所有插件路由
		pluginManager.RegisterAll(apiGroup)
	}

	// 注册API密钥管理路由（不需要API密钥验证）
	authGroup := r.Group("/auth")
	{
		// 注册API密钥管理路由
		api_key.RegisterRouter(authGroup)
	}

	// 添加统计信息展示页面路由
	r.GET("/stats", common.StatsHandler)
	// 添加统计信息API路由，返回JSON格式数据
	r.GET("/api/stats", common.StatsAPIHandler)
	// 添加API密钥管理页面路由
	r.GET("/api_key", common.APIKeyHandler)
}

// 启动服务
func startServer(r *gin.Engine) {
	port := config.GetServerPort()
	logrus.Infof("服务启动成功，监听地址：http://localhost%s", port)
	logrus.Infof("IP接口示例：http://localhost%s/api/ip?ip=114.114.114.114", port)
	logrus.Infof("Ping接口示例：http://localhost%s/api/ping?target=www.baidu.com&count=3", port)
	logrus.Infof("统计页面：http://localhost%s/stats", port)

	if err := r.Run(port); err != nil {
		logrus.Fatalf("服务启动失败：%v", err)
	}
}

func main() {
	// 初始化应用
	initApp()

	// 确保应用退出时关闭数据库连接
	defer func() {
		if err := db.CloseDB(); err != nil {
			logrus.Errorf("关闭数据库连接失败：%v", err)
		} else {
			logrus.Info("数据库连接已关闭")
		}
	}()

	// 设置Gin引擎和中间件
	r := setupGin()

	// 设置模板
	setupTemplates(r)

	// 设置静态文件服务
	setupStaticFiles(r)

	// 注册路由
	registerRoutes(r)

	// 启动服务
	startServer(r)
}
