package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-api/common"
	"github.com/xrcuo/xrcuo-api/config"
	"github.com/xrcuo/xrcuo-api/db"
	"github.com/xrcuo/xrcuo-api/plugin/ip"
	"github.com/xrcuo/xrcuo-api/plugin/ping"
)

func main() {
	// 1. 初始化数据库
	if err := db.InitDB(); err != nil {
		panic("数据库初始化失败：" + err.Error())
	}

	// 2. 初始化统计信息
	common.InitStats()

	// 2. 初始化Gin引擎（生产环境改为gin.ReleaseMode）
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

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

	// 3. 注册API根路由（所有插件路由都挂载在/api下）
	apiGroup := r.Group("/api")
	{
		// 为所有API注册统计中间件
		apiGroup.Use(common.StatsMiddleware())
		// 注册各个插件的路由（插件化核心：按需启用/禁用）
		ip.RegisterRouter(apiGroup)   // 启用IP插件
		ping.RegisterRouter(apiGroup) // 启用Ping插件
		// 后续新增插件，只需在这里添加注册语句即可
	}

	// 4. 添加统计信息展示页面路由
	r.GET("/stats", common.StatsHandler)

	// 5. 启动服务
	println("服务启动成功，监听地址：http://localhost" + config.ServerPort)
	println("IP接口示例：http://localhost" + config.ServerPort + "/api/ip?ip=114.114.114.114")
	println("Ping接口示例：http://localhost" + config.ServerPort + "/api/ping?target=www.baidu.com&count=3")
	println("统计页面：http://localhost" + config.ServerPort + "/stats")

	if err := r.Run(config.ServerPort); err != nil {
		panic("服务启动失败：" + err.Error())
	}
}
