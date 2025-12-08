package random

import (
	"io/fs"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/config"
)

// 随机图片API提供者列表
var imageProviders = []string{
	"https://picsum.photos/800/600",
	"https://source.unsplash.com/random/800x600",
	"https://random.imagecdn.app/800/600",
}

// 支持的图片扩展名
var supportedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 获取本地图片文件列表
func getLocalImages() ([]string, error) {
	localPath := config.Conf.RandomImage.LocalPath
	var images []string

	// 遍历本地图片目录
	err := filepath.Walk(localPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理文件，不处理目录
		if !info.IsDir() {
			// 检查文件扩展名是否为图片
			ext := strings.ToLower(filepath.Ext(path))
			if supportedImageExtensions[ext] {
				// 转换为相对路径，用于URL访问
				relPath, err := filepath.Rel(localPath, path)
				if err != nil {
					return err
				}
				// 使用正斜杠作为路径分隔符
				relPath = strings.ReplaceAll(relPath, "\\", "/")
				images = append(images, relPath)
			}
		}
		return nil
	})

	return images, err
}

// 检查本地图片目录是否存在且有图片
func hasLocalImages() bool {
	if !config.Conf.RandomImage.LocalEnabled {
		return false
	}

	images, err := getLocalImages()
	if err != nil {
		logrus.Warnf("获取本地图片失败: %v", err)
		return false
	}

	return len(images) > 0
}

// GetRandomImageHandler 获取随机图片的处理函数
func GetRandomImageHandler(c *gin.Context) {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 优先使用本地图片（如果启用且有图片）
	if hasLocalImages() {
		images, err := getLocalImages()
		if err == nil && len(images) > 0 {
			// 随机选择一张本地图片
			index := rand.Intn(len(images))
			imagePath := images[index]
			fullPath := filepath.Join(config.Conf.RandomImage.LocalPath, imagePath)

			// 记录请求日志
			logrus.WithFields(logrus.Fields{
				"provider":   "local_random_image",
				"image_path": imagePath,
				"client_ip":  c.ClientIP(),
			}).Info("本地随机图片请求")

			// 返回本地图片
			logrus.WithFields(logrus.Fields{
				"full_path":  fullPath,
				"local_path": config.Conf.RandomImage.LocalPath,
				"image_path": imagePath,
			}).Info("尝试返回本地图片")
			c.File(fullPath)
			return
		}
	}

	// 如果本地图片不可用，使用远程图片提供者
	index := rand.Intn(len(imageProviders))
	imageURL := imageProviders[index]

	// 记录请求日志
	logrus.WithFields(logrus.Fields{
		"provider":  "random_image",
		"image_url": imageURL,
		"client_ip": c.ClientIP(),
	}).Info("随机图片请求")

	// 重定向到随机图片URL
	c.Redirect(http.StatusFound, imageURL)
}

// GetRandomImageInfoHandler 获取随机图片信息的处理函数
func GetRandomImageInfoHandler(c *gin.Context) {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 优先使用本地图片（如果启用且有图片）
	if hasLocalImages() {
		images, err := getLocalImages()
		if err == nil && len(images) > 0 {
			// 随机选择一张本地图片
			index := rand.Intn(len(images))
			imagePath := images[index]

			// 返回本地图片信息
			c.JSON(http.StatusOK, ImageResponse{
				URL:      "/images/" + imagePath, // 本地图片的访问路径
				Provider: "local",
			})
			return
		}
	}

	// 如果本地图片不可用，使用远程图片提供者
	index := rand.Intn(len(imageProviders))
	imageURL := imageProviders[index]
	provider := "random"

	if index == 0 {
		provider = "picsum.photos"
	} else if index == 1 {
		provider = "unsplash.com"
	} else if index == 2 {
		provider = "random.imagecdn.app"
	}

	// 返回远程图片信息
	c.JSON(http.StatusOK, ImageResponse{
		URL:      imageURL,
		Provider: provider,
	})
}
