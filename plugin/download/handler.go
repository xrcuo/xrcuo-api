package download

import (
	"net/http"
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

	fullPath := filepath.Join(downloadPath, filePath)
	logrus.Debugf("下载请求: %s, IP: %s", fullPath, c.ClientIP())

	c.File(fullPath)
}
