package ip

import (
	"net"
	"net/http"
	"time"

	"go-boot/common"

	"github.com/gin-gonic/gin"
)

// SearchRegionHandler IP地区查询处理函数
func SearchRegionHandler(c *gin.Context) {
	startTime := time.Now()
	response := &Response{Code: 200, Msg: "请求成功"}

	// 统一响应出口（确保took字段必赋值）
	defer func() {
		response.Took = time.Since(startTime).String()
		c.JSON(http.StatusOK, response)
	}()

	// 1. 获取并校验IP参数
	ip := c.Query("ip")
	if ip == "" {
		response.Code = 400
		response.Msg = "参数错误：IP地址不能为空"
		return
	}

	// 校验IP格式
	if net.ParseIP(ip) == nil {
		response.Code = 400
		response.Msg = "参数错误：无效的IP地址格式"
		return
	}

	// 2. 调用公共工具查询地区
	regionParts, err := common.GetRegionByIP(ip)
	if err != nil {
		response.Code = 500
		response.Msg = "查询失败：" + err.Error()
		return
	}

	// 3. 构造响应数据
	locationParts := []string{regionParts.Country, regionParts.Province, regionParts.City}
	location := common.JoinNonEmpty(locationParts, "")
	area := common.JoinNonEmpty(append(locationParts, regionParts.Isp), "")

	response.Data = &Data{
		IP:       ip,
		Location: location,
		Isp:      regionParts.Isp,
		Area:     area,
	}
}
