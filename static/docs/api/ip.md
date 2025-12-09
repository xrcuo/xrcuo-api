# IP查询 API

## 功能描述

根据提供的IP地址，返回该IP的详细信息，包括国家、地区、城市和ISP等信息。

## 请求格式

```
GET /api/ip?ip={ip地址}
```

## 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `ip` | string | 否 | 客户端IP | 要查询的IP地址 |

## 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "114.114.114.114",
    "country": "中国",
    "region": "江苏",
    "city": "南京",
    "isp": "江苏省南京市 电信"
  }
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `ip` | string | 查询的IP地址 |
| `country` | string | 国家名称 |
| `region` | string | 地区/省份名称 |
| `city` | string | 城市名称 |
| `isp` | string | 互联网服务提供商信息 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/ip?ip=114.114.114.114
```

## 示例响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "114.114.114.114",
    "country": "中国",
    "region": "江苏",
    "city": "南京",
    "isp": "江苏省南京市 电信"
  }
}
```