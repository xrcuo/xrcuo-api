# API 密钥管理

## 功能描述

API密钥用于验证请求的合法性，防止API被滥用。系统支持生成、查看和删除API密钥。

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