# 文件下载 API

## 功能描述

提供文件直链下载功能，可以从本地 downloads 目录下载文件。

## 请求格式

```
GET /download/{filename}
# 或
GET /api/download/{filename}
```

## 请求参数

无（文件路径通过URL路径传递）

## 响应格式

直接返回文件内容（如果文件存在），否则返回404错误。

## 示例请求

```bash
# 方式一（推荐）
curl -O http://localhost:8080/download/ccmsi.lua

# 方式二
curl -O http://localhost:8080/api/download/ccmsi.lua
```

## 注意事项

- 文件必须存放在 `downloads/` 目录下
- 不需要 API 密钥即可访问
- 支持任意文件类型下载
