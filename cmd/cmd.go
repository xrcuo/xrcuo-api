package cmd

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/common"
	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/log"
	"github.com/xrcuo/xrcuo-api/plugin"
)

//go:embed static
//go:embed docs
//go:embed docs/_sidebar.md
var embeddedFiles embed.FS

var GlobalPluginManager *plugin.PluginManager

func InitApp() {
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

func SetupGin() *gin.Engine {
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

func SetupStaticFiles(r *gin.Engine) {
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

func SetupRoutes(r *gin.Engine) {
	pluginManager := plugin.NewPluginManager()
	pluginManager.RegisterBuiltinPlugins()

	if err := pluginManager.InitAll(); err != nil {
		logrus.Fatalf("插件初始化失败：%v", err)
	}

	GlobalPluginManager = pluginManager

	apiGroup := r.Group("/api")
	{
		apiGroup.Use(common.StatsMiddleware())
		pluginManager.RegisterAll(apiGroup)
	}

	r.GET("/stats", common.StatsHandler)
	r.GET("/api/stats", common.StatsAPIHandler)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/")
	})
}

func StartServer(r *gin.Engine) {

	port := config.GetServerPort()
	logrus.Infof("服务启动成功，监听地址：http://localhost%s", port)
	logrus.Infof("API文档：http://localhost%s/docs/", port)
	logrus.Infof("IP地址查询：http://localhost%s/api/ip?ip=114.114.114.114", port)
	logrus.Infof("获取访问者IP：http://localhost%s/api/ipify", port)
	logrus.Infof("Ping接口示例：http://localhost%s/api/ping?target=www.baidu.com&count=3", port)

	// 获取并打印本地网络接口IP
	if localIPs, err := common.GetLocalIPs(); err == nil && len(localIPs) > 0 {
		var publicIPs []string
		var privateIPs []string

		for _, ip := range localIPs {
			if common.IsPrivateIP(ip) {
				privateIPs = append(privateIPs, ip)
			} else {
				publicIPs = append(publicIPs, ip)
			}
		}

		if len(publicIPs) > 0 {
			logrus.Infof("本地公网IP：%v", publicIPs)
		}
		if len(privateIPs) > 0 {
			logrus.Infof("本地内网IP：%v", privateIPs)
		}
	} else if err != nil {
		logrus.Warnf("获取本地IP失败：%v", err)
	}

	// 获取并打印外网IP（使用国内API）
	if publicIP, err := common.GetPublicIP(); err == nil {
		logrus.Infof("外网IP地址：%s", publicIP)
	} else {
		logrus.Warnf("获取外网IP失败：%v", err)
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("服务启动失败：%v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// 监听信号：SIGINT (Ctrl+C), SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("正在关闭服务...")

	// 设置5秒超时用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("服务强制关闭：%v", err)
	}

	logrus.Info("服务已优雅关闭")
}
