# MCPE服务器查询 API

## 功能描述

查询Minecraft PE（基岩版）服务器状态信息，包括在线人数、服务器版本、延迟等。

## 请求格式

```
GET /api/mcpe/status?server={服务器地址}&port={端口}
```

## 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `server` | string | 是 | 无 | 服务器地址（IP地址或域名） |
| `port` | int | 否 | 19132 | 服务器端口 |

## 响应格式

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "server_ip": "play.example.com",
    "port": 19132,
    "online": 10,
    "max_players": 100,
    "version": "1.20.30",
    "motd": "§l§a生存服务器",
    "ping_time": "25ms",
    "time": "2024-01-01T12:00:00Z"
  },
  "took": "25.123ms"
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `code` | int | 状态码（200成功，400参数错误，500服务器错误） |
| `msg` | string | 提示信息 |
| `data` | object | 服务器状态数据 |
| `data.server_ip` | string | 服务器地址 |
| `data.port` | int | 服务器端口 |
| `data.online` | int | 当前在线玩家数 |
| `data.max_players` | int | 最大玩家数 |
| `data.version` | string | 服务器游戏版本 |
| `data.motd` | string | 服务器描述（Motd） |
| `data.ping_time` | string | 服务器延迟 |
| `data.time` | string | 查询时间 |
| `took` | string | 请求耗时 |

## 示例请求

```bash
curl "YOUR_DOMAIN/api/mcpe/status?server=play.example.com&port=19132"
```

## 示例响应

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "server_ip": "play.example.com",
    "port": 19132,
    "online": 10,
    "max_players": 100,
    "version": "1.20.30",
    "motd": "§l§a生存服务器",
    "ping_time": "25ms",
    "time": "2024-01-01T12:00:00Z"
  },
  "took": "25.123ms"
}
```

## 错误响应

### 参数错误（服务器地址为空）

```json
{
  "code": 400,
  "msg": "参数错误：服务器地址（server）不能为空",
  "data": null,
  "took": "0.123ms"
}
```

### 参数错误（端口无效）

```json
{
  "code": 400,
  "msg": "参数错误：端口必须是1-65535之间的整数",
  "data": null,
  "took": "0.123ms"
}
```

### 服务器查询失败

```json
{
  "code": 500,
  "msg": "MCPE服务器查询失败：连接超时",
  "data": null,
  "took": "5000.123ms"
}
```

## 注意事项

- 默认端口 19132 是 Minecraft PE 的标准端口
- 如果服务器使用了自定义端口，请确保在请求中指定正确的端口
- 连接超时时间为 5 秒