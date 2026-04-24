package bootstrap

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/xrcuo/xrcuo-api/internal/models"
	"github.com/xrcuo/xrcuo-lib/db"
	"golang.org/x/term"
)

// InitAdminUser 初始化管理员用户（首次启动时）
func InitAdminUser() error {
	exists, err := models.UserExists()
	if err != nil {
		return fmt.Errorf("检查用户存在性失败: %v", err)
	}

	if exists {
		logrus.Info("管理员用户已存在，跳过初始化")
		return nil
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  首次启动，请设置管理员初始密码")
	fmt.Println("========================================")
	fmt.Println()

	var password string
	for {
		fmt.Print("请输入初始密码（至少6位）: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Printf("读取密码失败: %v，请重试\n", err)
			continue
		}

		password = strings.TrimSpace(string(bytePassword))
		if len(password) < 6 {
			fmt.Println("密码长度不能少于6位，请重新输入")
			continue
		}

		fmt.Print("请再次输入密码确认: ")
		byteConfirm, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Printf("读取确认密码失败: %v，请重试\n", err)
			continue
		}

		confirm := strings.TrimSpace(string(byteConfirm))
		if password != confirm {
			fmt.Println("两次输入的密码不一致，请重新输入")
			continue
		}

		break
	}

	if err := models.UserCreate("admin", password); err != nil {
		return fmt.Errorf("创建管理员用户失败: %v", err)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  管理员账户创建成功")
	fmt.Println("  用户名: admin")
	fmt.Println("========================================")
	fmt.Println()
	logrus.Info("管理员用户初始化完成")
	return nil
}

// InitDatabase 初始化数据库表结构
func InitDatabase() error {
	if db.GetDB() == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	if err := models.InitAllTables(); err != nil {
		return err
	}

	logrus.Info("数据库表初始化完成")
	return nil
}

// InitApiDocs 初始化默认API文档数据
func InitApiDocs() error {
	// 检查是否已有文档
	docs, err := models.ApiDocList("", "", "")
	if err != nil {
		return err
	}

	if len(docs) > 0 {
		logrus.Infof("已有 %d 条API文档，跳过初始化", len(docs))
		return nil
	}

	// 初始化默认数据
	defaultDocs := getDefaultApiDocs()
	if err := models.InitApiDocsFromData(defaultDocs); err != nil {
		return fmt.Errorf("初始化默认API文档失败: %v", err)
	}

	logrus.Info("默认API文档初始化完成")
	return nil
}

func getDefaultApiDocs() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "IP 地区查询",
			"path":        "/api/ip",
			"method":      "GET",
			"description": "查询指定 IP 地址的地理位置、运营商等信息",
			"category":    "IP 查询",
			"tags":        []string{"IP", "地区", "运营商"},
			"parameters": []map[string]interface{}{
				{"name": "ip", "type": "string", "required": true, "description": "要查询的 IP 地址", "example": "114.114.114.114"},
			},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses": []map[string]interface{}{
				{
					"statusCode":  200,
					"description": "查询成功",
					"contentType": "application/json",
					"body":        `{"code":200,"msg":"请求成功","data":{"ip":"114.114.114.114","location":"中国 江苏 南京","isp":"中国电信","area":"中国 江苏 南京 中国电信"},"took":"15.234ms"}`,
				},
			},
		},
		{
			"name":        "获取访问者 IP",
			"path":        "/api/ipify",
			"method":      "GET",
			"description": "获取当前访问者的真实 IP 地址",
			"category":    "IP 查询",
			"tags":        []string{"IP", "访问者"},
			"parameters":  []map[string]interface{}{},
			"headers": []map[string]interface{}{
				{"name": "X-Real-IP", "required": false, "description": "代理服务器转发的真实客户端 IP", "example": "203.0.113.1"},
			},
			"requestBody": nil,
			"responses": []map[string]interface{}{
				{
					"statusCode":  200,
					"description": "成功返回 IP 地址",
					"contentType": "text/plain",
					"body":        "203.0.113.45",
				},
			},
		},
		{
			"name":        "Ping 网络测试",
			"path":        "/api/ping",
			"method":      "GET",
			"description": "对指定目标执行 ICMP Ping 测试",
			"category":    "网络工具",
			"tags":        []string{"Ping", "网络", "延迟"},
			"parameters": []map[string]interface{}{
				{"name": "target", "type": "string", "required": true, "description": "目标地址", "example": "www.baidu.com"},
				{"name": "timeout", "type": "integer", "required": false, "description": "超时时间（秒）", "defaultValue": "3", "example": "5"},
				{"name": "count", "type": "integer", "required": false, "description": "Ping 包数", "defaultValue": "4", "example": "4"},
			},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
		{
			"name":        "获取客户端信息",
			"path":        "/api/client",
			"method":      "GET",
			"description": "获取当前访问者的客户端信息",
			"category":    "客户端信息",
			"tags":        []string{"客户端", "User-Agent", "IP"},
			"parameters":  []map[string]interface{}{},
			"headers": []map[string]interface{}{
				{"name": "User-Agent", "required": false, "description": "浏览器 User-Agent", "example": "Mozilla/5.0..."},
			},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
		{
			"name":        "MCPE 服务器状态查询",
			"path":        "/api/mcpe/status",
			"method":      "GET",
			"description": "查询 Minecraft PE 服务器状态",
			"category":    "游戏查询",
			"tags":        []string{"Minecraft", "MCPE", "游戏服务器"},
			"parameters": []map[string]interface{}{
				{"name": "server", "type": "string", "required": true, "description": "服务器地址", "example": "play.example.com"},
				{"name": "port", "type": "integer", "required": false, "description": "服务器端口", "defaultValue": "19132", "example": "19132"},
			},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
		{
			"name":        "随机图片",
			"path":        "/api/randomimage",
			"method":      "GET",
			"description": "从服务器随机返回一张图片",
			"category":    "工具",
			"tags":        []string{"图片", "随机"},
			"parameters":  []map[string]interface{}{},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
		{
			"name":        "服务统计页面",
			"path":        "/stats",
			"method":      "GET",
			"description": "返回服务访问统计的 HTML 页面",
			"category":    "统计",
			"tags":        []string{"统计", "监控"},
			"parameters":  []map[string]interface{}{},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
		{
			"name":        "服务统计 API",
			"path":        "/api/stats",
			"method":      "GET",
			"description": "返回服务访问统计的 JSON 数据",
			"category":    "统计",
			"tags":        []string{"统计", "API", "JSON"},
			"parameters":  []map[string]interface{}{},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
		{
			"name":        "一言",
			"path":        "/api/yiyan",
			"method":      "GET",
			"description": "返回一句随机的中文名言/诗句",
			"category":    "工具",
			"tags":        []string{"一言", "诗词", "随机"},
			"parameters":  []map[string]interface{}{},
			"headers":     []map[string]interface{}{},
			"requestBody": nil,
			"responses":   []map[string]interface{}{},
		},
	}
}
