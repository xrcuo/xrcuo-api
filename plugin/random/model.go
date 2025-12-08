package random

// ImageResponse 随机图片响应模型
type ImageResponse struct {
	URL      string `json:"url"`
	Provider string `json:"provider"`
}
