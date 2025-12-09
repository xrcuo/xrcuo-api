# 随机数生成 API

## 功能描述

根据指定的范围，生成一个随机整数。

## 请求格式

```
GET /api/random?min={最小值}&max={最大值}
```

## 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `min` | int | 否 | 0 | 随机数最小值 |
| `max` | int | 否 | 100 | 随机数最大值 |

## 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "random": 42
  }
}
```

## 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `random` | int | 生成的随机整数 |

## 示例请求

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/random?min=1&max=100
```

## 示例响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "random": 42
  }
}
```