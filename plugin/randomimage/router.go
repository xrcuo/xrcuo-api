package randomimage

import "github.com/gin-gonic/gin"

var RandomImagePlugin = &randomImagePlugin{}

type randomImagePlugin struct{}

func (p *randomImagePlugin) Name() string {
	return "randomimage"
}

func (p *randomImagePlugin) Init() error {
	return nil
}

func (p *randomImagePlugin) RegisterRouter(group *gin.RouterGroup) {
	randomImageGroup := group.Group("/randomimage")
	{
		randomImageGroup.GET("", RandomImageHandler)
	}
}

func (p *randomImagePlugin) Cleanup() error {
	return nil
}

func RegisterRouter(group *gin.RouterGroup) {
	RandomImagePlugin.RegisterRouter(group)
}
