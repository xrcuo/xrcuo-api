package random

// ImageResponse 随机图片响应模型
type ImageResponse struct {
	URL      string `json:"url"`
	Provider string `json:"provider"`
}

// DmImgQueryParam 接口查询参数
type DmImgQueryParam struct {
	Param string `json:"param"`
	Value string `json:"value"`
	Des   string `json:"des"`
}

// DmImgInfoResponse 接口信息响应
type DmImgInfoResponse struct {
	Name   string          `json:"name"`
	Format string          `json:"format"`
	Method string          `json:"method"`
	Query  []DmImgQueryParam `json:"query"`
}

// DmImgData 壁纸数据
type DmImgData struct {
	URL string `json:"url"`
}

// DmImgResult 接口响应结果
type DmImgResult struct {
	Code string    `json:"code"`
	Msg  string    `json:"msg"`
	Data DmImgData `json:"data"`
}
