package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/cmd"
	"github.com/xrcuo/xrcuo-lib/common"
	"github.com/xrcuo/xrcuo-lib/config"
	"github.com/xrcuo/xrcuo-lib/db"
)

func main() {
	cmd.InitApp()

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
