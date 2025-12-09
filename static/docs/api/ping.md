# Ping测试 API

## 功能描述

对指定目标进行网络Ping测试，返回Ping的结果，包括延迟、TTL等信息。

## 请求格式

```
GET /api/ping?target={目标地址}&count={Ping次数}
```

## 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `target` | string | 是 | 无 | 目标主机名或IP地址 |
| `count` | int | 否 | 3 | Ping次数（1-10） |

## 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "target": "www.baidu.com",
    "count": 3,
    "results": [
      {"seq": 1, "ttl": 54, "time": "12.345ms"},
      {"seq": 2, "ttl": 54, "time": "11.234ms"},
      {"seq": 3, "ttl": 54, "time": "13.456ms"}
    ],
    "avg_time": "12.345ms"
  }
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `target` | string | 测试目标地址 |
| `count` | int | 实际Ping次数 |
| `results` | array | Ping结果列表 |
| `results[].seq` | int | 序列号 |
| `results[].ttl` | int | 生存时间 |
| `results[].time` | string | 延迟时间 |
| `avg_time` | string | 平均延迟时间 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/ping?target=www.baidu.com&count=3
```

## 示例响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "target": "www.baidu.com",
    "count": 3,
    "results": [
      {"seq": 1, "ttl": 54, "time": "12.345ms"},
      {"seq": 2, "ttl": 54, "time": "11.234ms"},
      {"seq": 3, "ttl": 54, "time": "13.456ms"}
    ],
    "avg_time": "12.345ms"
  }
}
```