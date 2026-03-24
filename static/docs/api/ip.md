# IP地址查询 API

## 功能描述

根据提供的IP地址，返回该IP的详细信息，包括地理位置（国家、省份、城市）和运营商信息。

## 请求格式

```
GET /api/ip?ip={ip地址}
```

## 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `ip` | string | 是 | 无 | 要查询的IP地址 |

## 响应格式

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "ip": "114.114.114.114",
    "location": "中国 江苏 南京",
    "isp": "江苏省南京市 电信",
    "area": "中国 江苏 南京 江苏省南京市 电信"
  },
  "took": "25.123ms"
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `code` | int | 状态码（200成功，400参数错误，500服务器错误） |
| `msg` | string | 提示信息 |
| `data` | object | IP信息数据 |
| `data.ip` | string | 查询的IP地址 |
| `data.location` | string | 地理位置（国家+省份+城市） |
| `data.isp` | string | 运营商信息 |
| `data.area` | string | 完整地区信息（国家+省份+城市+运营商） |
| `took` | string | 请求耗时 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" YOUR_DOMAIN/api/ip?ip=114.114.114.114
```

## 示例响应

```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "ip": "114.114.114.114",
    "location": "中国 江苏 南京",
    "isp": "江苏省南京市 电信",
    "area": "中国 江苏 南京 江苏省南京市 电信"
  },
  "took": "25.123ms"
}
```

## 错误响应

### 参数错误（IP为空）

```json
{
  "code": 400,
  "msg": "参数错误：IP地址不能为空",
  "data": null,
  "took": "0.123ms"
}
```

### 参数错误（IP格式无效）

```json
{
  "code": 400,
  "msg": "参数错误：无效的IP地址格式",
  "data": null,
  "took": "0.123ms"
}
```

### 服务器错误

```json
{
  "code": 500,
  "msg": "查询失败：xxx",
  "data": null,
  "took": "25.123ms"
}
```