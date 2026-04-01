# Ping测试 API

## 功能描述

对指定目标进行网络Ping测试，返回目标IP、延迟、地理位置、运营商以及详细的Ping统计信息（发送/接收/丢失包数，延迟等）。

## 请求格式

```
GET /api/ping?target={目标地址}&timeout={超时时间}&count={Ping次数}
```

## 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `target` | string | 是 | 无 | 目标主机名或IP地址 |
| `timeout` | int | 否 | 3 | 超时时间（1-10秒） |
| `count` | int | 否 | 4 | Ping包数量（1-10个） |

## 响应格式

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "target": "www.baidu.com",
    "ip": "183.2.144.1",
    "delay": "12.34ms",
    "location": "中国 江苏 南京",
    "isp": "中国电信",
    "area": "中国 江苏 南京 中国电信",
    "ping_stats": {
      "sent": 4,
      "received": 4,
      "lost": 0,
      "lost_rate": 0,
      "min_delay": "11.23ms",
      "avg_delay": "12.34ms",
      "max_delay": "13.45ms",
      "std_dev": "0.56ms"
    }
  },
  "took": "125.123ms"
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `code` | int | 状态码（200成功，400参数错误，500服务器错误） |
| `msg` | string | 提示信息 |
| `data` | object | Ping测试结果数据 |
| `data.target` | string | 测试目标（域名/IP） |
| `data.ip` | string | 解析后的目标IP地址 |
| `data.delay` | string | 平均延迟 |
| `data.location` | string | 地理位置（国家+省份+城市） |
| `data.isp` | string | 运营商 |
| `data.area` | string | 完整地区信息 |
| `data.ping_stats` | object | 详细统计信息 |
| `data.ping_stats.sent` | int | 发送包数 |
| `data.ping_stats.received` | int | 接收包数 |
| `data.ping_stats.lost` | int | 丢失包数 |
| `data.ping_stats.lost_rate` | float | 丢包率（%） |
| `data.ping_stats.min_delay` | string | 最小延迟 |
| `data.ping_stats.avg_delay` | string | 平均延迟 |
| `data.ping_stats.max_delay` | string | 最大延迟 |
| `data.ping_stats.std_dev` | string | 标准差（稳定性指标） |
| `took` | string | 请求耗时 |

## 示例请求

```bash
curl YOUR_DOMAIN/api/ping?target=www.baidu.com&count=4
```

## 示例响应

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "target": "www.baidu.com",
    "ip": "183.2.144.1",
    "delay": "12.34ms",
    "location": "中国 江苏 南京",
    "isp": "中国电信",
    "area": "中国 江苏 南京 中国电信",
    "ping_stats": {
      "sent": 4,
      "received": 4,
      "lost": 0,
      "lost_rate": 0,
      "min_delay": "11.23ms",
      "avg_delay": "12.34ms",
      "max_delay": "13.45ms",
      "std_dev": "0.56ms"
    }
  },
  "took": "125.123ms"
}
```

## 错误响应

### 参数错误（目标为空）

```json
{
  "code": 400,
  "msg": "参数错误：目标（target）不能为空",
  "data": null,
  "took": "0.123ms"
}
```

### 参数错误（超时时间无效）

```json
{
  "code": 400,
  "msg": "参数错误：超时时间必须是1-10秒",
  "data": null,
  "took": "0.123ms"
}
```

### 目标解析失败

```json
{
  "code": 400,
  "msg": "目标解析失败：xxx",
  "data": null,
  "took": "0.123ms"
}
```

### Ping测试失败

```json
{
  "code": 500,
  "msg": "Ping测试失败：xxx",
  "data": null,
  "took": "5000.123ms"
}
```

## 注意事项

- 内网IP使用更短的Ping间隔（10ms），外网使用正常间隔（100ms）
- 延迟显示为"超时"表示该方向完全无响应
- 标准差越小说明延迟越稳定