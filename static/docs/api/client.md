# 客户端信息获取 API

## 功能描述

获取当前请求客户端的详细信息，包括IP地址、浏览器信息、操作系统等。

## 请求格式

```
GET /api/client
```

## 请求参数

无

## 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "127.0.0.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "browser": "Chrome",
    "browser_version": "91.0.4472.124",
    "os": "Windows",
    "os_version": "10.0"
  }
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `ip` | string | 客户端IP地址 |
| `user_agent` | string | 完整的User-Agent字符串 |
| `browser` | string | 浏览器名称 |
| `browser_version` | string | 浏览器版本 |
| `os` | string | 操作系统名称 |
| `os_version` | string | 操作系统版本 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/client
```

## 示例响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "ip": "127.0.0.1",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "browser": "Chrome",
    "browser_version": "91.0.4472.124",
    "os": "Windows",
    "os_version": "10.0"
  }
}
```