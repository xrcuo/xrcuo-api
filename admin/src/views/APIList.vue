<template>
  <div class="api-list">
    <div class="card">
      <div class="card-title">🔗 API 接口列表</div>

      <div class="api-grid">
        <div v-for="api in apiList" :key="api.path" class="api-item">
          <div class="api-header">
            <span class="api-method" :class="'method-' + api.method.toLowerCase()">{{ api.method }}</span>
            <span class="api-path">{{ api.path }}</span>
          </div>
          <p class="api-desc">{{ api.description }}</p>
          <div class="api-meta">
            <span>路径: {{ api.path }}</span>
            <span>方法: {{ api.method }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-title">📖 接口文档</div>
      <div class="api-docs">
        <div v-for="api in apiList" :key="api.path + '-doc'" class="api-doc-item">
          <h4>{{ api.method }} {{ api.path }}</h4>
          <p>{{ api.description }}</p>
          <div class="api-example">
            <strong>示例请求:</strong>
            <code>{{ getExampleUrl(api) }}</code>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const apiList = ref([
  {
    path: '/api/ip',
    method: 'GET',
    description: '查询IP地理位置信息，支持批量查询'
  },
  {
    path: '/api/ipify',
    method: 'GET',
    description: '获取当前请求的公网IP地址'
  },
  {
    path: '/api/ping',
    method: 'GET',
    description: 'Ping命令执行，检测网络连通性'
  },
  {
    path: '/api/random',
    method: 'GET',
    description: '生成随机字符串或指定格式数据'
  },
  {
    path: '/api/mcpe',
    method: 'GET',
    description: 'Minecraft服务器状态查询(MCPE协议)'
  },
  {
    path: '/api/client',
    method: 'GET',
    description: '获取客户端信息(浏览器、操作系统等)'
  },
  {
    path: '/api/download',
    method: 'GET',
    description: '文件下载服务'
  },
  {
    path: '/api/stats',
    method: 'GET',
    description: '获取API使用统计数据'
  }
])

const getExampleUrl = (api) => {
  const baseUrl = 'http://localhost:8080'
  const examples = {
    '/api/ip': `${baseUrl}/api/ip?ip=114.114.114.114`,
    '/api/ipify': `${baseUrl}/api/ipify`,
    '/api/ping': `${baseUrl}/api/ping?target=www.baidu.com&count=3`,
    '/api/random': `${baseUrl}/api/random?length=16`,
    '/api/mcpe': `${baseUrl}/api/mcpe?host=play.hypixel.net`,
    '/api/client': `${baseUrl}/api/client`,
    '/api/download': `${baseUrl}/api/download/file.txt`,
    '/api/stats': `${baseUrl}/api/stats`
  }
  return examples[api.path] || `${baseUrl}${api.path}`
}
</script>

<style scoped>
.api-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.api-item {
  background: #fafafa;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid var(--border-color);
}

.api-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.api-method {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
  color: #fff;
}

.method-get { background: #52c41a; }
.method-post { background: #1890ff; }
.method-put { background: #faad14; }
.method-delete { background: #ff4d4f; }

.api-path {
  font-family: monospace;
  font-weight: 600;
  color: #333;
}

.api-desc {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 12px;
}

.api-meta {
  display: flex;
  gap: 16px;
  font-size: 0.85rem;
  color: #888;
}

.api-docs {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.api-doc-item {
  background: #fafafa;
  padding: 16px;
  border-radius: 8px;
}

.api-doc-item h4 {
  margin: 0 0 8px;
  font-family: monospace;
}

.api-doc-item p {
  color: #666;
  margin-bottom: 12px;
}

.api-example {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 6px;
}

.api-example code {
  display: block;
  margin-top: 8px;
  font-size: 0.85rem;
  word-break: break-all;
}
</style>