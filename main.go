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
	"github.com/xrcuo/xrcuo-api/plugin/download"
)

//go:embed static
//go:embed templates
//go:embed admin/dist
//go:embed docs
//go:embed docs/_sidebar.md
var embeddedFiles embed.FS

var globalPluginManager *plugin.PluginManager

func initApp() {
	config.Parse()
	log.InitLogger()

	config.GetInstance().RegisterUpdateCallback(func(newConfig *config.Config) {
		if dbDB := db.GetDB(); dbDB != nil {
			dbDB.SetMaxOpenConns(newConfig.Database.MaxOpenConns)
			dbDB.SetMaxIdleConns(newConfig.Database.MaxIdleConns)
			logrus.Info("数据库连接池配置已更新")
		}

		common.CloseIP2Region()
		if err := common.InitIP2Region(); err != nil {
			logrus.Errorf("IP2Region服务重新初始化失败: %v", err)
		} else {
			logrus.Info("IP2Region服务已重新初始化")
		}

		log.InitLogger()
	})

	config.GetInstance().WatchConfig()

	if err := db.InitDB(); err != nil {
		logrus.Fatalf("数据库初始化失败：%v", err)
	}

	if err := common.InitIP2Region(); err != nil {
		logrus.Fatalf("IP2Region数据库初始化失败：%v", err)
	}

	common.InitStats()
}

func setupGin() *gin.Engine {
	gin.SetMode(config.GetServerMode())
	r := gin.New()

	r.Use(common.RecoveryMiddleware())
	r.Use(common.RequestLoggerMiddleware())
	r.Use(common.CORSMiddleware())
	r.Use(common.RateLimitMiddleware())
	r.Use(common.PerformanceMiddleware())

	r.SetTrustedProxies(nil)

	return r
}

func setupTemplates(r *gin.Engine) {
	funcMap := template.FuncMap{
		"percentage": func(total, count int64) string {
			if total == 0 {
				return "0%"
			}
			return fmt.Sprintf("%d%%", int((float64(count)/float64(total))*100))
		},
	}

	tmpls, err := template.New("").Funcs(funcMap).ParseFS(embeddedFiles, "templates/*")
	if err != nil {
		logrus.Fatalf("加载模板失败：%v", err)
	}

	r.SetHTMLTemplate(tmpls)
}

func setupStaticFiles(r *gin.Engine) {
	r.Static("/images", "./images")

	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		logrus.Fatalf("获取static子目录失败：%v", err)
	}

	r.StaticFS("/static", http.FS(staticFS))

	docsFS, err := fs.Sub(embeddedFiles, "docs")
	if err != nil {
		logrus.Fatalf("获取docs子目录失败：%v", err)
	}
	r.StaticFS("/docs", http.FS(docsFS))

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(staticFS))
	})
}

func registerRoutes(r *gin.Engine) {
	pluginManager := plugin.NewPluginManager()
	pluginManager.RegisterBuiltinPlugins()

	if err := pluginManager.InitAll(); err != nil {
		logrus.Fatalf("插件初始化失败：%v", err)
	}

	globalPluginManager = pluginManager

	apiGroup := r.Group("/api")
	{
		apiGroup.Use(common.StatsMiddleware())
		apiGroup.Use(common.APIKeyMiddleware())
		pluginManager.RegisterAll(apiGroup)
	}

	downloadGroup := r.Group("/download")
	{
		downloadGroup.GET("/*filepath", download.DownloadHandler)
	}

	apiDownloadGroup := r.Group("/api/download")
	{
		apiDownloadGroup.GET("/*filepath", download.DownloadHandler)
	}

	authGroup := r.Group("/auth")
	{
		plugin.RegisterAPIRouter(authGroup)
	}

	r.GET("/stats", common.StatsHandler)
	r.GET("/api/stats", common.StatsAPIHandler)
	r.GET("/api_key", common.APIKeyHandler)

	adminFS, err := fs.Sub(embeddedFiles, "admin/dist")
	if err != nil {
		logrus.Fatalf("获取admin/dist子目录失败：%v", err)
	}
	r.StaticFS("/admin", http.FS(adminFS))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/")
	})
}

func startServer(r *gin.Engine) {
	port := config.GetServerPort()
	logrus.Infof("服务启动成功，监听地址：http://localhost%s", port)
	logrus.Infof("API文档：http://localhost%s/docs/", port)
	logrus.Infof("管理后台：http://localhost%s/admin/", port)
	logrus.Infof("IP接口示例：http://localhost%s/api/ip?ip=114.114.114.114", port)
	logrus.Infof("Ping接口示例：http://localhost%s/api/ping?target=www.baidu.com&count=3", port)

	if err := r.Run(port); err != nil {
		logrus.Fatalf("服务启动失败：%v", err)
	}
}

func main() {
	initApp()

	defer func() {
		common.CloseIP2Region()
		if err := db.CloseDB(); err != nil {
			logrus.Errorf("关闭数据库连接失败：%v", err)
		} else {
			logrus.Info("数据库连接已关闭")
		}
		if globalPluginManager != nil {
			globalPluginManager.CleanupAll()
		}
		config.GetInstance().StopWatching()
	}()

	r := setupGin()
	setupTemplates(r)
	setupStaticFiles(r)
	registerRoutes(r)
	startServer(r)
}
