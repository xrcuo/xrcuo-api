package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"github.com/xrcuo/xrcuo-api/config"
)

func InitLogger() {
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()
	if cfg == nil {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logrus.SetOutput(os.Stdout)
		return
	}

	levelStr := cfg.Log.Level
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", levelStr)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if cfg.Log.File != "" {
		logDir := filepath.Dir(cfg.Log.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Warnf("创建日志目录失败: %v, 仅输出到控制台", err)
		} else {
			if cfg.Log.NewFileOnStartup {
				rotateExistingLog(cfg.Log.File)
			}

			fileLogger := &lumberjack.Logger{
				Filename:   cfg.Log.File,
				MaxSize:    cfg.Log.MaxSize,
				MaxBackups: cfg.Log.MaxBackups,
				MaxAge:     cfg.Log.MaxAge,
				Compress:   cfg.Log.Compress,
			}

			if cfg.Log.ConsoleOutput {
				logrus.SetOutput(io.MultiWriter(os.Stdout, fileLogger))
			} else {
				logrus.SetOutput(fileLogger)
			}

			logrus.Infof("日志文件已配置: %s", cfg.Log.File)
			if cfg.Log.Compress {
				logrus.Info("日志压缩功能已启用")
			}
			if cfg.Log.NewFileOnStartup {
				logrus.Info("每次启动创建新日志文件功能已启用")
			}
		}
	} else {
		logrus.SetOutput(os.Stdout)
	}

	logrus.Debugf("日志级别已设置为: %s", level)
}

func rotateExistingLog(logPath string) {
	if _, err := os.Stat(logPath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		ext := filepath.Ext(logPath)
		base := logPath[:len(logPath)-len(ext)]
		newLogPath := fmt.Sprintf("%s-%s%s", base, timestamp, ext)

		if err := os.Rename(logPath, newLogPath); err == nil {
			logrus.Infof("旧日志文件已重命名: %s -> %s", logPath, newLogPath)
		} else {
			logrus.Warnf("重命名旧日志文件失败: %v", err)
		}
	}
}

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
