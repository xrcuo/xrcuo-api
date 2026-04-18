package mcpe

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xrcuo/xrcuo-lib/common"
)

// MCPEHandler 处理Minecraft PE服务器查询请求
func MCPEHandler(c *gin.Context) {
	startTime := time.Now()
	response := &Response{Code: 200, Msg: "请求成功"}

	// 统一响应出口
	defer func() {
		response.Took = time.Since(startTime).String()
		common.JSONResponse(c, http.StatusOK, response)
	}()

	// 1. 获取并校验参数
	server := c.Query("server")
	if server == "" {
		response.Code = 400
		response.Msg = "参数错误：服务器地址（server）不能为空"
		return
	}

	portStr := c.Query("port")
	port := 19132 // 默认Minecraft PE端口
	if portStr != "" {
		portInt := common.StrToInt(portStr, 19132)
		if portInt < 1 || portInt > 65535 {
			response.Code = 400
			response.Msg = "参数错误：端口必须是1-65535之间的整数"
			return
		}
		port = portInt
	}

	// 2. 构建服务器地址
	host := fmt.Sprintf("%s:%d", server, port)

	// 3. 执行MCPE服务器查询
	mcpeData, err := queryMCPEServer(host)
	if err != nil {
		response.Code = 500
		response.Msg = "MCPE服务器查询失败：" + err.Error()
		return
	}

	// 4. 构造响应数据
	response.Data = &Data{
		ServerIP:   server,
		Port:       port,
		Online:     mcpeData.Online,
		MaxPlayers: mcpeData.MaxPlayers,
		Version:    mcpeData.Version,
		Motd:       mcpeData.Motd,
		PingTime:   fmt.Sprintf("%dms", mcpeData.PingTime),
		Time:       time.Now(),
	}
}

// MCPEData MCPE服务器查询结果
type MCPEData struct {
	Online     int    `json:"online"`
	MaxPlayers int    `json:"max_players"`
	Version    string `json:"version"`
	Motd       string `json:"motd"`
	PingTime   int64  `json:"ping_time"`
}

// queryMCPEServer 查询MCPE服务器
func queryMCPEServer(host string) (*MCPEData, error) {
	// 创建UDP连接
	conn, err := net.Dial("udp", host)
	if err != nil {
		return nil, fmt.Errorf("UDP连接失败：%v", err)
	}
	defer conn.Close()

	// 组成发送数据
	PacketID := []byte{0x01} // Packet ID
	// 获取当前时间戳
	ClientSendTime := make([]byte, 8) // 客户端发送时间
	binary.BigEndian.PutUint64(ClientSendTime, uint64(time.Now().UnixMilli()))
	Magic := []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD} // Magic Number
	ClientID := []byte{0x12, 0x34, 0x56, 0x78, 0x00} // 客户端ID
	ClientGUID := make([]byte, 8) // 客户端GUID
	binary.BigEndian.PutUint64(ClientGUID, 0)
	// 组合数据
	SendData := append(PacketID, ClientSendTime...)
	SendData = append(SendData, Magic...)
	SendData = append(SendData, ClientID...)
	SendData = append(SendData, ClientGUID...)

	// 发送数据
	startTime := time.Now().UnixNano() / 1e6 // 记录发送时间
	_, err = conn.Write(SendData)
	if err != nil {
		return nil, fmt.Errorf("发送数据包失败：%v", err)
	}

	// 接收数据
	UDPdata := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // 设置读取五秒超时
	// 读取数据
	_, err = conn.Read(UDPdata)
	if err != nil {
		return nil, fmt.Errorf("接收响应失败：%v", err)
	}
	endTime := time.Now().UnixNano() / 1e6 // 记录接收时间

	// 解析服务器响应
	ServerInfo := UDPdata[33:] // 服务器信息

	// 按;分割数据
	MotdData := strings.Split(string(ServerInfo), ";")

	// 检查数据完整性
	if len(MotdData) < 9 {
		return nil, fmt.Errorf("响应数据不完整")
	}

	// 解析数据
	motd := MotdData[1] // 服务器MOTD line 1
	version := MotdData[3] // 服务器游戏版本
	playerCount := MotdData[4] // 在线人数
	maxPlayerCount := MotdData[5] // 最大在线人数

	// 转换数据类型
	playerCountInt, err := strconv.Atoi(playerCount)
	if err != nil {
		return nil, fmt.Errorf("解析在线人数失败：%v", err)
	}

	maxPlayerCountInt, err := strconv.Atoi(maxPlayerCount)
	if err != nil {
		return nil, fmt.Errorf("解析最大在线人数失败：%v", err)
	}

	// 计算延迟
	pingTime := endTime - startTime

	return &MCPEData{
		Online:     playerCountInt,
		MaxPlayers: maxPlayerCountInt,
		Version:    version,
		Motd:       motd,
		PingTime:   pingTime,
	}, nil
}
