# API 密钥管理

## 功能描述

API密钥用于验证请求的合法性，防止API被滥用。系统支持生成、查看和删除API密钥。

## 启用/禁用API密钥验证

API密钥验证功能可以通过配置文件开关控制：

```yaml
# config.yaml
api_key:
  enabled: true  # true=启用验证，false=禁用验证（所有API无需密钥即可访问）
```

- `enabled: true`（默认）- 所有API请求需要有效的API密钥
- `enabled: false` - 所有API请求无需API密钥即可访问（适用于内网或不需要认证的场景）

修改配置后需要重启服务生效。

## 生成API密钥

1. 访问 http://localhost:8080/api_key
2. 点击 "生成新密钥" 按钮
3. 复制生成的API密钥

## 使用API密钥

### 在请求头中使用

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/ip?ip=114.114.114.114
```

### 在查询参数中使用

```bash
curl http://localhost:8080/api/ip?ip=114.114.114.114&api_key=your-api-key
```

## 管理API密钥

1. 访问 http://localhost:8080/api_key
2. 查看所有已生成的API密钥
3. 点击 "删除" 按钮删除不再使用的API密钥

## API密钥限制

- 每个API密钥每分钟最多可以发送100个请求（可通过配置文件修改）
- API密钥是大小写敏感的
- 请妥善保管您的API密钥，避免泄露

## 错误处理

当API密钥无效或已过期时，系统将返回以下错误响应：

```json
{
  "code": 401,
  "message": "Invalid API Key",
  "data": null
}
```

当请求次数超过限制时，系统将返回以下错误响应：

```json
{
  "code": 429,
  "message": "Rate Limit Exceeded",
  "data": null
}
```

## REST API 接口

### 获取所有API密钥

```
GET /api/api_key
```

#### 响应格式

```json
{
  "code": 200,
  "data": {
    "api_keys": [
      {
        "id": 1,
        "name": "测试密钥",
        "key": "xxxx.xxxx.xxxx",
        "max_usage": 1000,
        "current_usage": 50,
        "is_permanent": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `api_keys` | array | API密钥列表 |
| `api_keys[].id` | int | 密钥ID |
| `api_keys[].name` | string | 密钥名称 |
| `api_keys[].key` | string | API密钥（仅创建时返回完整密钥） |
| `api_keys[].max_usage` | int | 最大使用次数（0表示无限制） |
| `api_keys[].current_usage` | int | 当前已使用次数 |
| `api_keys[].is_permanent` | bool | 是否永久有效 |
| `api_keys[].created_at` | string | 创建时间 |
| `api_keys[].updated_at` | string | 更新时间 |

### 创建新的API密钥

```
POST /api/api_key
```

#### 请求体

```json
{
  "name": "新密钥名称",
  "max_usage": 1000,
  "is_permanent": true
}
```

| 参数名 | 类型 | 必填 | 默认值 | 描述 |
|-------|------|------|-------|------|
| `name` | string | 是 | 无 | 密钥名称 |
| `max_usage` | int | 否 | 0 | 最大使用次数（0表示无限制） |
| `is_permanent` | bool | 否 | false | 是否永久有效 |

#### 响应格式

```json
{
  "code": 201,
  "data": {
    "api_key": {
      "id": 2,
      "name": "新密钥名称",
      "key": "abc123.def456.ghi789",
      "max_usage": 1000,
      "current_usage": 0,
      "is_permanent": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 删除API密钥

```
DELETE /api/api_key/:id
```

#### 响应格式

```json
{
  "code": 200,
  "data": {
    "message": "API密钥删除成功"
  }
}
```

### 错误响应

#### 参数无效

```json
{
  "code": 400,
  "data": {
    "error": "请求参数无效"
  }
}
```

#### 内部错误

```json
{
  "code": 500,
  "data": {
    "error": "创建API密钥失败"
  }
}
```