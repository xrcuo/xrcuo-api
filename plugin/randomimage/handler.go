package randomimage

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	imageDir = "images"
)

var extensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomImageHandler(c *gin.Context) {
	startTime := time.Now()

	availableImages, err := scanImages()
	if err != nil {
		logrus.Errorf("扫描图片目录失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "扫描图片目录失败",
			"took": time.Since(startTime).String(),
		})
		return
	}

	if len(availableImages) == 0 {
		logrus.Warnf("未找到任何图片文件，目录：%s", imageDir)
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "未找到图片",
			"took": time.Since(startTime).String(),
		})
		return
	}

	randomIndex := rand.Intn(len(availableImages))
	imagePath := availableImages[randomIndex]

	logrus.Infof("返回随机图片：%s（共%d张）", imagePath, len(availableImages))
	c.File(imagePath)
}

func scanImages() ([]string, error) {
	var images []string

	entries, err := os.ReadDir(imageDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		for _, allowedExt := range extensions {
			if ext == allowedExt {
				images = append(images, filepath.Join(imageDir, entry.Name()))
				break
			}
		}
	}

	return images, nil
}
