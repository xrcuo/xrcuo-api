<template>
  <div class="dashboard">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>加载中...</span>
    </div>

    <template v-else>
      <div class="stat-grid">
        <div class="stat-card">
          <div class="stat-icon blue">🌐</div>
          <div class="stat-info">
            <h3>{{ stats.total_calls || 0 }}</h3>
            <p>总调用次数</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon green">📅</div>
          <div class="stat-info">
            <h3>{{ stats.daily_calls || 0 }}</h3>
            <p>今日调用</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon orange">⏱️</div>
          <div class="stat-info">
            <h3>{{ stats.hourly_calls || 0 }}</h3>
            <p>每小时调用</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon purple">🔗</div>
          <div class="stat-info">
            <h3>{{ Object.keys(stats.path_calls || {}).length }}</h3>
            <p>API路径数量</p>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-title">📈 HTTP方法统计</div>
        <div class="method-stats">
          <div v-for="(count, method) in stats.method_calls" :key="method" class="method-item">
            <div class="method-info">
              <span class="method-name">{{ method }}</span>
              <span class="method-count">{{ count }} 次</span>
            </div>
            <div class="progress">
              <div class="progress-bar" :style="{ width: getPercentage(count) + '%' }"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-title">🔗 API路径统计</div>
        <table>
          <thead>
            <tr>
              <th>API路径</th>
              <th>调用次数</th>
              <th>占比</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(count, path) in stats.path_calls" :key="path">
              <td>{{ path }}</td>
              <td>{{ count }}</td>
              <td>
                <div class="progress" style="max-width: 200px;">
                  <div class="progress-bar" :style="{ width: getPercentage(count) + '%' }"></div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="card">
        <div class="card-title">📋 最近调用记录</div>
        <table>
          <thead>
            <tr>
              <th>时间</th>
              <th>路径</th>
              <th>方法</th>
              <th>IP</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="call in stats.last_call_details" :key="call.timestamp + call.path">
              <td>{{ formatTime(call.timestamp) }}</td>
              <td>{{ call.path }}</td>
              <td><span class="badge badge-success">{{ call.method }}</span></td>
              <td>{{ call.ip }}</td>
              <td>
                <span :class="getStatusClass(call.status_code)">{{ call.status_code }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const loading = ref(true)
const stats = ref({})

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/stats')
    stats.value = response.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

const getPercentage = (count) => {
  if (!stats.value.total_calls) return 0
  return Math.min(100, Math.round((count / stats.value.total_calls) * 100))
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

const getStatusClass = (code) => {
  if (code === 200) return 'badge badge-success'
  if (code >= 400 && code < 500) return 'badge badge-warning'
  return 'badge badge-error'
}

onMounted(() => {
  fetchStats()
  setInterval(fetchStats, 5000)
})
</script>

<style scoped>
.method-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.method-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.method-info {
  display: flex;
  justify-content: space-between;
}

.method-name {
  font-weight: 600;
  color: #333;
}

.method-count {
  color: #666;
}
</style>