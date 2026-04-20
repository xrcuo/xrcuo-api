package random

import (
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-lib/common"
	"github.com/xrcuo/xrcuo-lib/config"
)

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
	// 获取配置管理器实例
	cm := config.GetInstance()
	// 获取当前配置
	conf := cm.GetConfig()

	// 检查是否启用本地图片
	if !conf.RandomImage.LocalEnabled {
		return nil, nil
	}

	// 检查缓存是否有效
	if len(localImagesCache) > 0 && time.Since(lastCacheUpdate) < cacheDuration {
		return localImagesCache, nil
	}

	localPath := conf.RandomImage.LocalPath
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

	// 检查本地图片是否可用
	if len(images) > 0 && err == nil {
		// 随机选择一张本地图片
		index := rand.Intn(len(images))
		imagePath := images[index]
		// 获取配置管理器实例
		cm := config.GetInstance()
		// 获取当前配置
		conf := cm.GetConfig()
		fullPath := filepath.Join(conf.RandomImage.LocalPath, imagePath)

		// 记录请求日志（只记录关键信息）
		logrus.Debugf("本地随机图片请求: %s, IP: %s", imagePath, c.ClientIP())

		// 返回本地图片
		c.File(fullPath)
		return
	}

	// 本地图片不可用时返回错误
	logrus.Warnf("本地图片不可用，请求IP: %s", c.ClientIP())
	common.JSONResponse(c, http.StatusNotFound, map[string]string{
		"error": "本地图片不可用，请配置本地图片路径",
	})
}

// GetRandomImageInfoHandler 获取随机图片信息的处理函数
func GetRandomImageInfoHandler(c *gin.Context) {
	// 获取本地图片列表
	images, err := getLocalImages()

	// 检查本地图片是否可用
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

	// 本地图片不可用时返回错误
	logrus.Warnf("本地图片不可用，请求IP: %s", c.ClientIP())
	common.JSONResponse(c, http.StatusNotFound, map[string]string{
		"error": "本地图片不可用，请配置本地图片路径",
	})
}

// GetDmImgHandler 随机二次元壁纸处理函数
func GetDmImgHandler(c *gin.Context) {
	// 检查info参数，如果有则返回接口信息
	if c.Query("info") != "" {
		common.JSONResponse(c, http.StatusOK, DmImgInfoResponse{
			Name:   "随机二次元壁纸",
			Format: "image,json,xml,jsonp,text",
			Method: "GET,POST,PUT",
			Query: []DmImgQueryParam{
				{
					Param: "format",
					Value: "image",
					Des:   "跳转到图片",
				},
			},
		})
		return
	}

	// 生成随机图片URL
	url := "https://cache.aqco.top/static/api/img/dm/" + fmt.Sprintf("%d", rand.Intn(2494)+1) + ".jpg"
	format := c.Query("format")

	switch format {
	case "image":
		c.Redirect(http.StatusFound, url)
		return
	case "text":
		c.String(http.StatusOK, url)
		return
	default:
		common.JSONResponse(c, http.StatusOK, DmImgResult{
			Code: "200",
			Msg:  "请求成功",
			Data: DmImgData{
				URL: url,
			},
		})
	}
}
