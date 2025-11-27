package main

import (
	"go-boot/config"
	"go-boot/plugin/ip"
	"go-boot/plugin/ping"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化Gin引擎（生产环境改为gin.ReleaseMode）
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 2. 注册API根路由（所有插件路由都挂载在/api下）
	apiGroup := r.Group("/api")
	{
		// 注册各个插件的路由（插件化核心：按需启用/禁用）
		ip.RegisterRouter(apiGroup)   // 启用IP插件
		ping.RegisterRouter(apiGroup) // 启用Ping插件
		// 后续新增插件，只需在这里添加注册语句即可
	}

	// 3. 启动服务
	println("服务启动成功，监听地址：http://localhost" + config.ServerPort)
	println("IP接口示例：http://localhost" + config.ServerPort + "/api/ip?ip=114.114.114.114")
	println("Ping接口示例：http://localhost" + config.ServerPort + "/api/ping?target=www.baidu.com&count=3")

	if err := r.Run(config.ServerPort); err != nil {
		panic("服务启动失败：" + err.Error())
	}
}
