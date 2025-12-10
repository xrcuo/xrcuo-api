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
	"github.com/xrcuo/xrcuo-api/log"
	"github.com/xrcuo/xrcuo-api/plugin"
)

//go:embed static
//go:embed templates
var embeddedFiles embed.FS

// 全局插件管理器实例，用于在程序退出时清理资源
var globalPluginManager *plugin.PluginManager

// initApp 初始化应用程序
// 功能：
// 1. 解析配置文件
// 2. 初始化日志配置
// 3. 启动配置文件监听
// 4. 初始化数据库连接
// 5. 预加载IP2Region数据库
// 6. 初始化统计信息
func initApp() {
	// 解析配置文件
	config.Parse()

	// 初始化日志配置
	log.InitLogger()

	// 注册配置更新回调
	config.GetInstance().RegisterUpdateCallback(func(newConfig *config.Config) {
		// 更新数据库连接池配置
		if dbDB := db.GetDB(); dbDB != nil {
			dbDB.SetMaxOpenConns(newConfig.Database.MaxOpenConns)
			dbDB.SetMaxIdleConns(newConfig.Database.MaxIdleConns)
			logrus.Info("数据库连接池配置已更新")
		}

		// 重新初始化IP2Region服务
		common.CloseIP2Region()
		if err := common.InitIP2Region(); err != nil {
			logrus.Errorf("IP2Region服务重新初始化失败: %v", err)
		} else {
			logrus.Info("IP2Region服务已重新初始化")
		}

		// 重新初始化日志配置
		log.InitLogger()
	})

	// 启动配置文件监听，实现配置热重载
	config.GetInstance().WatchConfig()

	// 初始化数据库连接和表结构
	if err := db.InitDB(); err != nil {
		logrus.Fatalf("数据库初始化失败：%v", err)
	}

	// 预加载IP2Region数据库，用于IP地址查询
	if err := common.InitIP2Region(); err != nil {
		logrus.Fatalf("IP2Region数据库初始化失败：%v", err)
	}

	// 初始化统计信息，用于记录API调用次数和性能指标
	common.InitStats()
}

// 设置Gin引擎和中间件
func setupGin() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.GetServerMode())

	// 创建Gin引擎实例（不使用默认中间件，手动添加）
	r := gin.New()

	// 添加自定义的Recovery中间件（替换默认的Recovery中间件）
	r.Use(common.RecoveryMiddleware())
	// 添加请求日志中间件
	r.Use(common.RequestLoggerMiddleware())
	// 添加跨域中间件
	r.Use(common.CORSMiddleware())
	// 添加速率限制中间件
	r.Use(common.RateLimitMiddleware())
	// 添加性能监控中间件
	r.Use(common.PerformanceMiddleware())

	// 信任所有代理，确保能正确获取客户端真实IP
	r.SetTrustedProxies(nil)

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

	// 使用静态文件服务提供docsify文档，直接映射到static/docs目录
	r.StaticFS("/docs", http.Dir("./static/docs"))

	// 使用嵌入式文件系统提供favicon.ico
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(staticFS))
	})
}

// 注册路由
func registerRoutes(r *gin.Engine) {
	// 创建插件管理器
	pluginManager := plugin.NewPluginManager()

	// 注册所有内置插件
	pluginManager.RegisterBuiltinPlugins()

	// 初始化所有插件
	if err := pluginManager.InitAll(); err != nil {
		logrus.Fatalf("插件初始化失败：%v", err)
	}

	// 将插件管理器添加到全局变量，以便在程序退出时清理资源
	globalPluginManager = pluginManager

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
		plugin.RegisterAPIRouter(authGroup)
	}

	// 添加统计信息展示页面路由
	r.GET("/stats", common.StatsHandler)
	// 添加统计信息API路由，返回JSON格式数据
	r.GET("/api/stats", common.StatsAPIHandler)
	// 添加API密钥管理页面路由
	r.GET("/api_key", common.APIKeyHandler)

	// 根路径重定向到docs
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/")
	})
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

	// 确保应用退出时关闭资源
	defer func() {
		// 关闭IP2Region服务
		common.CloseIP2Region()
		// 关闭数据库连接
		if err := db.CloseDB(); err != nil {
			logrus.Errorf("关闭数据库连接失败：%v", err)
		} else {
			logrus.Info("数据库连接已关闭")
		}
		// 清理所有插件资源
		if globalPluginManager != nil {
			globalPluginManager.CleanupAll()
		}
		// 停止配置文件监听
		config.GetInstance().StopWatching()
		// 停止API密钥缓存清理任务
		common.StopAPICacheCleanup()
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
