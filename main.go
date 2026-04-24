package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/cmd"
	"github.com/xrcuo/xrcuo-api/internal/bootstrap"
	"github.com/xrcuo/xrcuo-lib/common"
	"github.com/xrcuo/xrcuo-lib/config"
	"github.com/xrcuo/xrcuo-lib/db"
)

func main() {
	cmd.InitApp()

	// 初始化数据库表
	if err := bootstrap.InitDatabase(); err != nil {
		logrus.Fatalf("数据库表初始化失败: %v", err)
	}

	// 初始化管理员用户（首次启动时提示设置密码）
	if err := bootstrap.InitAdminUser(); err != nil {
		logrus.Fatalf("管理员用户初始化失败: %v", err)
	}

	// 初始化默认API文档数据
	if err := bootstrap.InitApiDocs(); err != nil {
		logrus.Warnf("默认API文档初始化失败: %v", err)
	}

	defer func() {
		common.CloseIP2Region()
		if err := db.CloseDB(); err != nil {
			logrus.Errorf("关闭数据库连接失败：%v", err)
		} else {
			logrus.Info("数据库连接已关闭")
		}
		if cmd.GlobalPluginManager != nil {
			cmd.GlobalPluginManager.CleanupAll()
		}
		config.GetInstance().StopWatching()
	}()

	r := cmd.SetupGin()
	cmd.SetupStaticFiles(r)
	cmd.SetupRoutes(r)
	cmd.StartServer(r)
}
