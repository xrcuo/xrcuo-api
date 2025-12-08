package client

// Response 统一响应结构
type Response struct {
	Code int    `json:"code"` // 状态码（200成功，400参数错，500服务器错）
	Msg  string `json:"msg"`  // 提示信息
	Data *Data  `json:"data"` // 核心数据
	Took string `json:"took"` // 总耗时
}

// Data 客户端信息核心数据
type Data struct {
	IP              string `json:"ip"`              // 客户端IP地址
	Location        string `json:"location"`        // 地理位置（国家+省份+城市）
	ISP             string `json:"isp"`             // 运营商
	Area            string `json:"area"`            // 完整信息（国家+省份+城市+运营商）
	OS              string `json:"os"`              // 操作系统
	Browser         string `json:"browser"`         // 浏览器
	BrowserVersion  string `json:"browser_version"` // 浏览器版本
}
