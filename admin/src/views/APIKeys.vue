<template>
  <div class="api-keys">
    <div class="card">
      <div class="card-header">
        <div class="card-title">🔑 API Key 管理</div>
        <button class="btn btn-primary" @click="showCreateModal = true">
          + 创建新API Key
        </button>
      </div>

      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>

      <div v-else-if="apiKeys.length === 0" class="empty-state">
        <p>暂无API Key，请创建一个新的API Key</p>
      </div>

      <table v-else>
        <thead>
          <tr>
            <th>名称</th>
            <th>API Key</th>
            <th>使用情况</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="key in apiKeys" :key="key.id">
            <td>{{ key.name }}</td>
            <td>
              <div class="api-key-cell">
                <code>{{ key.key }}</code>
                <button class="copy-btn" @click="copyKey(key.key)">复制</button>
              </div>
            </td>
            <td>
              <div class="usage-info">
                <span>{{ key.current_usage }} / {{ key.is_permanent ? '∞' : key.max_usage }}</span>
                <div class="progress" v-if="!key.is_permanent">
                  <div class="progress-bar" :style="{ width: getUsagePercent(key) + '%' }"></div>
                </div>
              </div>
            </td>
            <td>
              <span :class="key.is_permanent ? 'badge badge-success' : 'badge badge-warning'">
                {{ key.is_permanent ? '永久有效' : '限制使用' }}
              </span>
            </td>
            <td>{{ formatTime(key.created_at) }}</td>
            <td>
              <button class="btn btn-danger btn-sm" @click="deleteKey(key.id, key.name)">
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>创建新API Key</h3>
          <button class="modal-close" @click="showCreateModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">名称</label>
            <input v-model="newKey.name" type="text" class="input" placeholder="输入API Key名称">
          </div>
          <div class="form-group">
            <label class="form-label">最大使用次数</label>
            <input v-model.number="newKey.max_usage" type="number" class="input" placeholder="0表示无限制" min="0">
          </div>
          <div class="form-group">
            <div class="checkbox-group">
              <input v-model="newKey.is_permanent" type="checkbox" id="isPermanent">
              <label for="isPermanent">永久有效</label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn" @click="showCreateModal = false">取消</button>
          <button class="btn btn-primary" @click="createKey">创建</button>
        </div>
      </div>
    </div>

    <div v-if="showSuccessModal" class="success-toast" @click="showSuccessModal = false">
      {{ successMessage }}
    </div>

    <div v-if="createdKeyData" class="modal-overlay" @click.self="closeCreatedModal">
      <div class="modal">
        <div class="modal-header">
          <h3>✅ API Key 创建成功</h3>
          <button class="modal-close" @click="closeCreatedModal">&times;</button>
        </div>
        <div class="modal-body">
          <p style="margin-bottom: 16px;">请妥善保存您的API Key，只显示一次！</p>
          <div class="api-key-display">{{ createdKeyData.key }}</div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" @click="copyAndClose">复制并关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const loading = ref(true)
const apiKeys = ref([])
const showCreateModal = ref(false)
const showSuccessModal = ref(false)
const successMessage = ref('')
const createdKeyData = ref(null)

const newKey = ref({
  name: '',
  max_usage: 100,
  is_permanent: false
})

const fetchAPIKeys = async () => {
  try {
    const response = await axios.get('/auth/api_key')
    apiKeys.value = response.data.api_keys || []
  } catch (error) {
    console.error('获取API Keys失败:', error)
  } finally {
    loading.value = false
  }
}

const createKey = async () => {
  try {
    const response = await axios.post('/auth/api_key', {
      name: newKey.value.name,
      max_usage: newKey.value.max_usage,
      is_permanent: newKey.value.is_permanent
    })

    if (response.data.api_key) {
      createdKeyData.value = response.data.api_key
      showCreateModal.value = false
      newKey.value = { name: '', max_usage: 100, is_permanent: false }
      fetchAPIKeys()
    }
  } catch (error) {
    console.error('创建API Key失败:', error)
    showToast('创建失败')
  }
}

const deleteKey = async (id, name) => {
  if (!confirm(`确定要删除API Key "${name}"吗？此操作不可恢复。`)) {
    return
  }

  try {
    await axios.delete(`/auth/api_key/${id}`)
    showToast('删除成功')
    fetchAPIKeys()
  } catch (error) {
    console.error('删除API Key失败:', error)
    showToast('删除失败')
  }
}

const copyKey = (key) => {
  navigator.clipboard.writeText(key)
  showToast('已复制到剪贴板')
}

const copyAndClose = () => {
  if (createdKeyData.value) {
    navigator.clipboard.writeText(createdKeyData.value.key)
  }
  closeCreatedModal()
  showToast('已复制到剪贴板')
}

const closeCreatedModal = () => {
  createdKeyData.value = null
}

const showToast = (message) => {
  successMessage.value = message
  showSuccessModal.value = true
  setTimeout(() => {
    showSuccessModal.value = false
  }, 2000)
}

const getUsagePercent = (key) => {
  if (key.is_permanent || !key.max_usage) return 0
  return Math.min(100, Math.round((key.current_usage / key.max_usage) * 100))
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  fetchAPIKeys()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.api-key-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.api-key-cell code {
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.85rem;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.usage-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.usage-info .progress {
  max-width: 120px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.85rem;
}
</style>