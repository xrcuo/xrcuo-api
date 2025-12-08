package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/common"
	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/plugin/api_key"
	"github.com/xrcuo/xrcuo-api/plugin/client"
	"github.com/xrcuo/xrcuo-api/plugin/ip"
	"github.com/xrcuo/xrcuo-api/plugin/ipify"
	"github.com/xrcuo/xrcuo-api/plugin/ping"
	"github.com/xrcuo/xrcuo-api/plugin/random"
)

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
	// 设置Gin模式（生产环境改为gin.ReleaseMode）
	gin.SetMode(gin.DebugMode)

	// 创建Gin引擎实例
	r := gin.Default()

	// 信任所有代理，确保能正确获取客户端真实IP
	r.SetTrustedProxies(nil)

	// 添加全局中间件
	r.Use(common.RequestLoggerMiddleware()) // 请求日志中间件
	r.Use(common.CORSMiddleware())          // 跨域中间件

	return r
}

// 设置模板
func setupTemplates(r *gin.Engine) {
	// 设置模板自定义函数
	r.SetFuncMap(template.FuncMap{
		"percentage": func(total, count int64) int {
			if total == 0 {
				return 0
			}
			return int((float64(count) / float64(total)) * 100)
		},
	})

	// 加载模板文件
	r.LoadHTMLGlob("templates/*")
}

// 设置静态文件服务
func setupStaticFiles(r *gin.Engine) {
	// 添加静态文件服务，用于提供本地图片
	r.Static("/images", "./images")

	// 添加静态文件服务，用于提供其他静态资源
	r.Static("/static", "./static")

	// 直接映射favicon.ico文件
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
}

// 注册路由
func registerRoutes(r *gin.Engine) {
	// 注册API根路由（所有插件路由都挂载在/api下）
	apiGroup := r.Group("/api")
	{
		// 为所有API注册统计中间件
		apiGroup.Use(common.StatsMiddleware())
		// 为所有API注册API密钥验证中间件
		apiGroup.Use(common.APIKeyMiddleware())
		// 注册各个插件的路由（插件化核心：按需启用/禁用）
		ip.RegisterRouter(apiGroup)     // 启用IP插件
		ping.RegisterRouter(apiGroup)   // 启用Ping插件
		random.RegisterRouter(apiGroup) // 启用随机图片插件
		client.RegisterRouter(apiGroup) // 启用客户端信息插件
		ipify.RegisterRouter(apiGroup)  // 启用IP获取插件
		// 后续新增插件，只需在这里添加注册语句即可
	}

	// 注册API密钥管理路由（不需要API密钥验证）
	authGroup := r.Group("/auth")
	{
		// 注册API密钥管理路由
		api_key.RegisterRouter(authGroup)
	}

	// 添加统计信息展示页面路由
	r.GET("/stats", common.StatsHandler)
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
