# 统计功能

## 功能描述

统计功能用于记录和展示API的使用情况，包括请求次数、响应时间、状态码分布等信息。

## 查看统计信息

1. 访问 http://localhost:8080/stats
2. 查看API请求的统计数据
3. 可以按API路径和状态码筛选统计结果

## 统计API

### 请求格式

```
GET /api/stats
```

### 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `path` | string | 否 | 无 | 按API路径筛选 |
| `status_code` | int | 否 | 无 | 按状态码筛选 |
| `start_time` | string | 否 | 无 | 开始时间（格式：2023-01-01T00:00:00Z） |
| `end_time` | string | 否 | 无 | 结束时间（格式：2023-01-01T23:59:59Z） |

### 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_requests": 1000,
    "avg_response_time": "12.345ms",
    "status_code_distribution": {
      "200": 950,
      "401": 30,
      "429": 20
    },
    "top_paths": [
      {"path": "/api/ip", "count": 500},
      {"path": "/api/ping", "count": 300},
      {"path": "/api/random", "count": 200}
    ]
  }
}
```

### 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `total_requests` | int | 总请求次数 |
| `avg_response_time` | string | 平均响应时间 |
| `status_code_distribution` | object | 状态码分布 |
| `top_paths` | array | 热门API路径 |
| `top_paths[].path` | string | API路径 |
| `top_paths[].count` | int | 请求次数 |

### 示例请求

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/stats
```

### 示例响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_requests": 1000,
    "avg_response_time": "12.345ms",
    "status_code_distribution": {
      "200": 950,
      "401": 30,
      "429": 20
    },
    "top_paths": [
      {"path": "/api/ip", "count": 500},
      {"path": "/api/ping", "count": 300},
      {"path": "/api/random", "count": 200}
    ]
  }
}
```

## 统计数据存储

统计数据默认存储在SQLite数据库中，数据库文件为`data.db`。可以通过修改配置文件中的`database.dsn`字段来使用其他数据库。