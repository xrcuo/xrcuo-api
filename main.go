package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/common"
	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/plugin/ip"
	"github.com/xrcuo/xrcuo-api/plugin/ping"
	"github.com/xrcuo/xrcuo-api/plugin/random"
)

func main() {
	// 1. 解析配置文件
	config.Parse()

	// 2. 初始化数据库
	if err := db.InitDB(); err != nil {
		logrus.Fatalf("数据库初始化失败：%v", err)
	}

	// 确保应用退出时关闭数据库连接
	defer func() {
		if err := db.CloseDB(); err != nil {
			logrus.Errorf("关闭数据库连接失败：%v", err)
		}
	}()

	// 3. 预加载IP2Region数据库
	if err := common.InitIP2Region(); err != nil {
		logrus.Fatalf("IP2Region数据库初始化失败：%v", err)
	}

	// 4. 初始化统计信息
	common.InitStats()

	// 4. 初始化Gin引擎（生产环境改为gin.ReleaseMode）
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	// 添加全局中间件
	r.Use(common.RequestLoggerMiddleware()) // 请求日志中间件
	r.Use(common.CORSMiddleware())          // 跨域中间件

	// 设置模板路径和自定义函数
	r.SetFuncMap(template.FuncMap{
		"percentage": func(total, count int64) int {
			if total == 0 {
				return 0
			}
			return int((float64(count) / float64(total)) * 100)
		},
	})
	r.LoadHTMLGlob("templates/*")

	// 添加静态文件服务，用于提供本地图片
	r.Static("/images", "./images")
	// 添加静态文件服务，用于提供其他静态资源
	r.Static("/static", "./static")
	// 直接映射favicon.ico文件
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// 3. 注册API根路由（所有插件路由都挂载在/api下）
	apiGroup := r.Group("/api")
	{
		// 为所有API注册统计中间件
		apiGroup.Use(common.StatsMiddleware())
		// 注册各个插件的路由（插件化核心：按需启用/禁用）
		ip.RegisterRouter(apiGroup)     // 启用IP插件
		ping.RegisterRouter(apiGroup)   // 启用Ping插件
		random.RegisterRouter(apiGroup) // 启用随机图片插件
		// 后续新增插件，只需在这里添加注册语句即可
	}

	// 4. 添加统计信息展示页面路由
	r.GET("/stats", common.StatsHandler)

	// 5. 启动服务
	port := config.GetServerPort()
	logrus.Infof("服务启动成功，监听地址：http://localhost%s", port)
	logrus.Infof("IP接口示例：http://localhost%s/api/ip?ip=114.114.114.114", port)
	logrus.Infof("Ping接口示例：http://localhost%s/api/ping?target=www.baidu.com&count=3", port)
	logrus.Infof("统计页面：http://localhost%s/stats", port)

	if err := r.Run(port); err != nil {
		logrus.Fatalf("服务启动失败：%v", err)
	}
}
