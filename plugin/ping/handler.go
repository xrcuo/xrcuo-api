package ping

import (
	"fmt"
	"net/http"
	"time"

	"go-boot/common"

	"github.com/gin-gonic/gin"
	"github.com/go-ping/ping"
)

// PingHandler Ping测试处理函数
func PingHandler(c *gin.Context) {
	startTime := time.Now()
	response := &Response{Code: 200, Msg: "请求成功"}

	// 统一响应出口
	defer func() {
		response.Took = time.Since(startTime).String()
		c.JSON(http.StatusOK, response)
	}()

	// 1. 获取并校验参数
	target := c.Query("target")
	if target == "" {
		response.Code = 400
		response.Msg = "参数错误：目标（target）不能为空"
		return
	}

	// 超时时间（默认3秒，1-10秒）
	timeoutSec := c.Query("timeout")
	timeout := 3 * time.Second
	if sec := common.StrToInt(timeoutSec, 3); sec >= 1 && sec <= 10 {
		timeout = time.Duration(sec) * time.Second
	} else {
		response.Code = 400
		response.Msg = "参数错误：超时时间必须是1-10秒"
		return
	}

	// Ping包数（默认4个，1-10个）
	countStr := c.Query("count")
	count := 4
	if cnt := common.StrToInt(countStr, 4); cnt >= 1 && cnt <= 10 {
		count = cnt
	} else {
		response.Code = 400
		response.Msg = "参数错误：Ping包数必须是1-10之间的整数"
		return
	}

	// 2. 解析目标（域名→IP）
	ipAddr, err := common.ResolveTarget(target)
	if err != nil {
		response.Code = 400
		response.Msg = "目标解析失败：" + err.Error()
		return
	}

	// 3. 执行Ping测试
	pingStats, err := doPing(ipAddr, timeout, count)
	if err != nil {
		response.Code = 500
		response.Msg = "Ping测试失败：" + err.Error()
		return
	}

	// 4. 查询地区信息
	regionParts, err := common.GetRegionByIP(ipAddr)
	if err != nil {
		response.Msg = "Ping成功，但地区查询失败：" + err.Error()
		regionParts = common.RegionParts{}
	}

	// 5. 构造响应数据
	locationParts := []string{regionParts.Country, regionParts.Province, regionParts.City}
	location := common.JoinNonEmpty(locationParts, "")
	area := common.JoinNonEmpty(append(locationParts, regionParts.Isp), "")

	// 格式化延迟显示
	avgDelay := formatDelay(pingStats.AvgRtt)
	minDelay := formatDelay(pingStats.MinRtt)
	maxDelay := formatDelay(pingStats.MaxRtt)
	stdDev := formatDelay(pingStats.StdDevRtt)

	response.Data = &Data{
		Target:   target,
		IP:       ipAddr,
		Delay:    avgDelay,
		Location: location,
		Isp:      regionParts.Isp,
		Area:     area,
		PingStats: &PingStats{
			Sent:     pingStats.PacketsSent,
			Received: pingStats.PacketsRecv,
			Lost:     pingStats.PacketsSent - pingStats.PacketsRecv,
			LostRate: float64(pingStats.PacketsSent-pingStats.PacketsRecv) / float64(pingStats.PacketsSent) * 100,
			MinDelay: minDelay,
			AvgDelay: avgDelay,
			MaxDelay: maxDelay,
			StdDev:   stdDev,
		},
	}
}

// doPing 执行ICMP Ping测试（适配内外网间隔）
func doPing(ip string, timeout time.Duration, count int) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return nil, fmt.Errorf("Pinger初始化失败：%v", err)
	}

	// 配置Ping参数
	pinger.Count = count
	pinger.Timeout = timeout
	pinger.SetPrivileged(true) // Windows需要管理员权限

	// 优化：内网缩短间隔（10ms），外网正常间隔（100ms）
	if common.IsPrivateIP(ip) {
		pinger.Interval = 10 * time.Millisecond
	} else {
		pinger.Interval = 100 * time.Millisecond
	}

	// 执行Ping（阻塞）
	if err := pinger.Run(); err != nil {
		return nil, fmt.Errorf("Ping执行失败：%v", err)
	}

	return pinger.Statistics(), nil
}

// formatDelay 格式化延迟（微秒→毫秒，保留2位小数）
func formatDelay(delay time.Duration) string {
	if delay == 0 {
		return "超时"
	}
	return fmt.Sprintf("%.2fms", delay.Seconds()*1000)
}
