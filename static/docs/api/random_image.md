# 随机图片 API

## 功能描述

获取随机图片，支持本地图片和远程图片源。可配置优先使用本地图片。

## 接口一：获取随机图片

返回一张随机图片（直接返回图片内容或重定向到图片URL）。

### 请求格式

```
GET /api/random/image
```

### 请求参数

无

### 响应格式

返回图片内容（Content-Type: image/jpeg、image/png、image/gif、image/webp 等）

### 示例请求

```bash
curl -H "X-API-Key: your-api-key" "http://localhost:8080/api/random/image"
```

---

## 接口二：获取随机图片信息

返回随机图片的信息（JSON格式），包含图片URL和来源。

### 请求格式

```
GET /api/random/image/info
```

### 请求参数

无

### 响应格式

```json
{
  "code": 200,
  "data": {
    "url": "/images/sample.jpg",
    "provider": "local"
  }
}
```

### 响应字段说明

| 字段名 | 类型 | 描述 |
|-------|------|------|
| `code` | int | 状态码 |
| `data` | object | 图片信息数据 |
| `data.url` | string | 图片URL地址 |
| `data.provider` | string | 图片来源：local（本地）/ picsum.photos / unsplash.com / random.imagecdn.app |

### 示例请求

```bash
curl -H "X-API-Key: your-api-key" "http://localhost:8080/api/random/image/info"
```

### 示例响应

```json
{
  "code": 200,
  "data": {
    "url": "/images/photo.jpg",
    "provider": "local"
  }
}
```

## 配置说明

在 `config.yaml` 中可以配置：

```yaml
random_image:
  local_enabled: true    # 是否启用本地图片
  local_path: "images/"  # 本地图片目录路径
```

- 当 `local_enabled` 为 `true` 且本地图片目录有图片时，优先返回本地图片
- 当本地图片不可用时，自动切换到远程图片源

## 注意事项

- 支持的图片格式：jpg、jpeg、png、gif、webp
- 本地图片列表会缓存5分钟，减少文件系统访问
- 远程图片来源：picsum.photos、unsplash.com、random.imagecdn.app