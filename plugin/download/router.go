package download

import (
	"github.com/gin-gonic/gin"
)

var DownloadPlugin = &downloadPlugin{}

type downloadPlugin struct{}

func (p *downloadPlugin) Name() string {
	return "download"
}

func (p *downloadPlugin) Init() error {
	return nil
}

func (p *downloadPlugin) RegisterRouter(group *gin.RouterGroup) {
	group.GET("/download", DownloadHandler)
	group.GET("/download/*filepath", DownloadHandler)
}

func (p *downloadPlugin) Cleanup() error {
	return nil
}

func RegisterRouter(router *gin.RouterGroup) {
	DownloadPlugin.RegisterRouter(router)
}
