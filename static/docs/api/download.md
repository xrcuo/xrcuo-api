# 文件下载 API

## 功能描述

提供文件直链下载功能，可以从配置的 downloads 目录下载文件。

## 请求格式

```
GET /download?file={filename}
# 或
GET /download/{filepath}
```

## 请求参数

| 参数名 | 类型 | 必填 | 描述 |
|-------|------|------|------|
| `file` | string | 是（方式一） | 要下载的文件名 |

## 响应格式

直接返回文件内容（如果文件存在），否则返回JSON错误。

成功时：返回文件内容，Content-Type 根据文件类型自动设置

失败时：
```json
{
  "error": "缺少file参数"
}
```

## 示例请求

```bash
# 方式一：使用查询参数
curl -O "http://localhost:8080/download?file=ccmsi.lua"

# 方式二：使用路径参数
curl -O http://localhost:8080/download/ccmsi.lua
```

## 错误响应

### 缺少参数

```json
{
  "error": "缺少file参数"
}
```

### 文件不存在

返回 HTTP 404 状态码

## 注意事项

- 文件必须存放在配置的 `downloads/` 目录下
- 不需要 API 密钥即可访问
- 支持任意文件类型下载