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
	"github.com/xrcuo/xrcuo-api/common"
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

// 本地图片列表缓存
var (
	localImagesCache []string
	lastCacheUpdate  time.Time
	cacheDuration    = 5 * time.Minute // 缓存有效期5分钟
)

// 初始化随机数生成器（只初始化一次）
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 获取本地图片文件列表，带缓存
func getLocalImages() ([]string, error) {
	// 检查是否启用本地图片
	if !config.Conf.RandomImage.LocalEnabled {
		return nil, nil
	}

	// 检查缓存是否有效
	if len(localImagesCache) > 0 && time.Since(lastCacheUpdate) < cacheDuration {
		return localImagesCache, nil
	}

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

	if err != nil {
		return nil, err
	}

	// 更新缓存
	localImagesCache = images
	lastCacheUpdate = time.Now()

	return images, nil
}

// GetRandomImageHandler 获取随机图片的处理函数
func GetRandomImageHandler(c *gin.Context) {
	// 获取本地图片列表
	images, err := getLocalImages()

	// 优先使用本地图片（如果启用且有图片）
	if len(images) > 0 && err == nil {
		// 随机选择一张本地图片
		index := rand.Intn(len(images))
		imagePath := images[index]
		fullPath := filepath.Join(config.Conf.RandomImage.LocalPath, imagePath)

		// 记录请求日志（只记录关键信息）
		logrus.Debugf("本地随机图片请求: %s, IP: %s", imagePath, c.ClientIP())

		// 返回本地图片
		c.File(fullPath)
		return
	}

	// 如果本地图片不可用，使用远程图片提供者
	index := rand.Intn(len(imageProviders))
	imageURL := imageProviders[index]

	// 记录请求日志（只记录关键信息）
	logrus.Debugf("随机图片请求: %s, IP: %s", imageURL, c.ClientIP())

	// 重定向到随机图片URL
	c.Redirect(http.StatusFound, imageURL)
}

// GetRandomImageInfoHandler 获取随机图片信息的处理函数
func GetRandomImageInfoHandler(c *gin.Context) {
	// 获取本地图片列表
	images, err := getLocalImages()

	// 优先使用本地图片（如果启用且有图片）
	if len(images) > 0 && err == nil {
		// 随机选择一张本地图片
		index := rand.Intn(len(images))
		imagePath := images[index]

		// 返回本地图片信息
		common.JSONResponse(c, http.StatusOK, ImageResponse{
			URL:      "/images/" + imagePath, // 本地图片的访问路径
			Provider: "local",
		})
		return
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
	common.JSONResponse(c, http.StatusOK, ImageResponse{
		URL:      imageURL,
		Provider: provider,
	})
}
