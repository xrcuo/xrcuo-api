package download

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
)

func DownloadHandler(c *gin.Context) {
	filePath := c.Query("file")
	if filePath == "" {
		filePath = c.Param("filepath")
		if filePath != "" {
			filePath = strings.TrimPrefix(filePath, "/")
		}
	}

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少file参数",
		})
		return
	}

	downloadPath := config.GetDownloadPath()

	cleanFilePath := filepath.Clean(filePath)
	fullPath := filepath.Join(downloadPath, cleanFilePath)

	absDownloadPath, err := filepath.Abs(downloadPath)
	if err != nil {
		logrus.Errorf("获取下载目录绝对路径失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		logrus.Errorf("获取文件绝对路径失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件路径"})
		return
	}

	if !strings.HasPrefix(absFullPath, absDownloadPath+string(os.PathSeparator)) && absFullPath != absDownloadPath {
		logrus.Warnf("路径遍历攻击尝试: %s, IP: %s", filePath, c.ClientIP())
		c.JSON(http.StatusForbidden, gin.H{"error": "禁止访问"})
		return
	}

	fileInfo, err := os.Stat(absFullPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
			return
		}
		logrus.Errorf("获取文件信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能下载目录"})
		return
	}

	logrus.Debugf("下载请求: %s, IP: %s", absFullPath, c.ClientIP())

	c.File(absFullPath)
}
