package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"github.com/xrcuo/xrcuo-api/config"
)

// InitLogger 初始化日志配置
func InitLogger() {
	configManager := config.GetInstance()
	cfg := configManager.GetConfig()
	if cfg == nil {
		// 使用默认配置
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logrus.SetOutput(os.Stdout)
		return
	}

	// 设置日志级别
	levelStr := cfg.Log.Level
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Warnf("无效的日志级别: %s, 使用默认级别: info", levelStr)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 如果配置了日志文件路径，则配置日志输出
	if cfg.Log.File != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.Log.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.Warnf("创建日志目录失败: %v, 仅输出到控制台", err)
		} else {
			// 创建日志文件输出
			fileLogger := &lumberjack.Logger{
				Filename:   cfg.Log.File,
				MaxSize:    cfg.Log.MaxSize,
				MaxBackups: cfg.Log.MaxBackups,
				MaxAge:     cfg.Log.MaxAge,
			}

			if cfg.Log.ConsoleOutput {
				// 同时输出到控制台和文件
				logrus.SetOutput(io.MultiWriter(os.Stdout, fileLogger))
			} else {
				// 只输出到文件
				logrus.SetOutput(fileLogger)
			}

			logrus.Infof("日志文件已配置: %s", cfg.Log.File)
		}
	} else {
		// 只输出到控制台
		logrus.SetOutput(os.Stdout)
	}

	logrus.Debugf("日志级别已设置为: %s", level)
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

// Debug 调试日志
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
