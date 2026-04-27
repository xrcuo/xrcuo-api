package cmd

import (
	"io"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/internal/handler"
	"github.com/xrcuo/xrcuo-api/plugin"
	"github.com/xrcuo/xrcuo-lib/common"
	"github.com/xrcuo/xrcuo-lib/config"
)

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

	// 一言 API 代理接口
	r.GET("/api/yiyan", func(c *gin.Context) {
		resp, err := http.Get("https://api.kekc.cn/api/yien")
		if err != nil {
			logrus.Errorf("请求一言 API 失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"cn": "人生若只如初见，何事秋风悲画扇。",
			})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("读取一言 API 响应失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"cn": "人生若只如初见，何事秋风悲画扇。",
			})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	})

	webFS, _ := fs.Sub(embeddedFiles, "web")

	r.GET("/", func(c *gin.Context) {
		siteConfig := config.GetInstance().GetConfig()

		tmplContent, err := fs.ReadFile(webFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read template: %v", err)
			return
		}

		tmpl, err := template.New("index.html").Parse(string(tmplContent))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to parse template: %v", err)
			return
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(c.Writer, gin.H{
			"Title":     siteConfig.Site.Title,
			"Name":      siteConfig.Site.Name,
			"Motto":     siteConfig.Site.Motto,
			"AvatarURL": siteConfig.Site.AvatarURL,
			"ICP":       siteConfig.Site.ICP,
			"Copyright": siteConfig.Site.Copyright,
			"Links":     siteConfig.Site.Links,
			"Contact":   siteConfig.Site.Contact,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to execute template: %v", err)
		}
	})

	// API 文档静态文件服务 - /docs/ 路径
	docsFS, _ := fs.Sub(embeddedFiles, "docs")
	r.StaticFS("/docs/", http.FS(docsFS))

	// 重定向 /docs 到 /docs/
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/")
	})

	// 后台管理 API 路由
	adminGroup := r.Group("/admin")
	{
		adminGroup.POST("/login", handler.Login)
		adminGroup.GET("/user", handler.AuthMiddleware(), handler.GetCurrentUser)
		adminGroup.POST("/password", handler.AuthMiddleware(), handler.ChangePassword)

		// API 文档管理
		adminGroup.GET("/api-docs", handler.ApiDocList)
		adminGroup.GET("/api-docs/:id", handler.ApiDocGet)
		adminGroup.POST("/api-docs", handler.AuthMiddleware(), handler.ApiDocCreate)
		adminGroup.PUT("/api-docs", handler.AuthMiddleware(), handler.ApiDocUpdate)
		adminGroup.DELETE("/api-docs/:id", handler.AuthMiddleware(), handler.ApiDocDelete)
		adminGroup.GET("/api-docs/categories", handler.ApiDocCategories)

		// 系统配置管理
		adminGroup.GET("/config", handler.AuthMiddleware(), handler.SystemConfigGet)
		adminGroup.POST("/config", handler.AuthMiddleware(), handler.SystemConfigSet)
		adminGroup.GET("/configs", handler.AuthMiddleware(), handler.SystemConfigList)

		// 系统资源监控
		adminGroup.GET("/metrics", handler.AuthMiddleware(), handler.GetSystemMetrics)
		adminGroup.GET("/metrics/history", handler.AuthMiddleware(), handler.GetSystemMetricsHistory)

		// 备案号管理
		adminGroup.GET("/icp", handler.ICPGet)
		adminGroup.POST("/icp", handler.AuthMiddleware(), handler.ICPSet)
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path[0] == '/' {
			path = path[1:]
		}
		data, err := fs.ReadFile(webFS, path)
		if err != nil {
			c.String(http.StatusNotFound, "File not found: %s", path)
			return
		}
		contentType := "text/plain"
		switch {
		case strings.HasSuffix(path, ".html"):
			contentType = "text/html; charset=utf-8"
		case strings.HasSuffix(path, ".css"):
			contentType = "text/css; charset=utf-8"
		case strings.HasSuffix(path, ".js"):
			contentType = "application/javascript; charset=utf-8"
		case strings.HasSuffix(path, ".png"):
			contentType = "image/png"
		case strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg"):
			contentType = "image/jpeg"
		case strings.HasSuffix(path, ".ico"):
			contentType = "image/x-icon"
		}
		c.Data(http.StatusOK, contentType, data)
	})
}
