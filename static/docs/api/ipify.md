# 获取公网IP API

## 功能描述

获取当前客户端的公网IP地址。

## 请求格式

```
GET /api/ipify
```

## 请求参数

无

## 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "123.123.123.123"
  }
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `ip` | string | 客户端的公网IP地址 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/ipify
```

## 示例响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "123.123.123.123"
  }
}
```