# 客户端信息获取 API

## 功能描述

获取当前请求客户端的详细信息，包括IP地址、地理位置、运营商、操作系统、浏览器信息等。

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
  "msg": "请求成功",
  "data": {
    "ip": "127.0.0.1",
    "location": "中国 江苏 南京",
    "isp": "江苏省南京市 电信",
    "area": "中国 江苏 南京 江苏省南京市 电信",
    "os": "Windows 10",
    "browser": "Google Chrome",
    "browser_version": "91.0.4472.124"
  },
  "took": "25.123ms"
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `code` | int | 状态码（200成功，400参数错误，500服务器错误） |
| `msg` | string | 提示信息 |
| `data` | object | 客户端信息数据 |
| `data.ip` | string | 客户端IP地址 |
| `data.location` | string | 地理位置（国家+省份+城市） |
| `data.isp` | string | 运营商信息 |
| `data.area` | string | 完整地区信息（国家+省份+城市+运营商） |
| `data.os` | string | 操作系统名称 |
| `data.browser` | string | 浏览器名称 |
| `data.browser_version` | string | 浏览器版本 |
| `took` | string | 请求耗时 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" YOUR_DOMAIN/api/client
```

## 示例响应

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "ip": "127.0.0.1",
    "location": "中国 江苏 南京",
    "isp": "江苏省南京市 电信",
    "area": "中国 江苏 南京 江苏省南京市 电信",
    "os": "Windows 10",
    "browser": "Google Chrome",
    "browser_version": "91.0.4472.124"
  },
  "took": "25.123ms"
}
```

## 错误响应

### 服务器错误

```json
{
  "code": 500,
  "msg": "请求失败：xxx",
  "data": null,
  "took": "25.123ms"
}
```

## 注意事项

- 如果无法获取地理位置信息，`location`、`isp`、`area` 字段可能为空字符串
- 操作系统和浏览器信息通过解析 User-Agent 头获取
- 支持通过 `X-Real-IP` 或 `X-Forwarded-For` 头获取真实客户端IP（用于代理场景）