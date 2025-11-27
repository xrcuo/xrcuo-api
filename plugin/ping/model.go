package ping

// Response Ping接口统一响应结构
type Response struct {
	Code int    `json:"code"` // 状态码
	Msg  string `json:"msg"`  // 提示信息
	Data *Data  `json:"data"` // 核心数据
	Took string `json:"took"` // 总耗时
}

// Data Ping核心数据（含延迟+地区）
type Data struct {
	Target    string     `json:"target"`     // 目标（域名/IP）
	IP        string     `json:"ip"`         // 解析后的IP
	Delay     string     `json:"delay"`      // 平均延迟（友好显示）
	Location  string     `json:"location"`   // 地区（国家+省份+城市）
	Isp       string     `json:"isp"`        // 运营商
	Area      string     `json:"area"`       // 完整地区信息
	PingStats *PingStats `json:"ping_stats"` // Ping详细统计
}

// PingStats Ping测试详细统计
type PingStats struct {
	Sent     int     `json:"sent"`      // 发送包数
	Received int     `json:"received"`  // 接收包数
	Lost     int     `json:"lost"`      // 丢失包数
	LostRate float64 `json:"lost_rate"` // 丢包率（%）
	MinDelay string  `json:"min_delay"` // 最小延迟
	AvgDelay string  `json:"avg_delay"` // 平均延迟
	MaxDelay string  `json:"max_delay"` // 最大延迟
	StdDev   string  `json:"std_dev"`   // 标准差（稳定性）
}
