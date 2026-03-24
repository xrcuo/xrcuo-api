<template>
  <div id="app" class="app-container">
    <nav class="sidebar">
      <div class="logo">
        <h2>Xrcuo API</h2>
        <p>管理后台</p>
      </div>
      <ul class="nav-menu">
        <li>
          <router-link to="/" class="nav-link">
            <span class="icon">📊</span>
            <span>仪表盘</span>
          </router-link>
        </li>
        <li>
          <router-link to="/api-keys" class="nav-link">
            <span class="icon">🔑</span>
            <span>API Key管理</span>
          </router-link>
        </li>
        <li>
          <router-link to="/api-list" class="nav-link">
            <span class="icon">🔗</span>
            <span>API接口列表</span>
          </router-link>
        </li>
      </ul>
    </nav>
    <main class="main-content">
      <header class="top-bar">
        <h1>{{ pageTitle }}</h1>
        <div class="time">{{ currentTime }}</div>
      </header>
      <div class="content">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const currentTime = ref('')

const pageTitle = computed(() => {
  const titles = {
    '/': '仪表盘',
    '/api-keys': 'API Key管理',
    '/api-list': 'API接口列表'
  }
  return titles[route.path] || '管理后台'
})

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleString('zh-CN')
}

let timer
onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style scoped>
.app-container {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 240px;
  background: #1a1a2e;
  color: #fff;
  padding: 20px;
  position: fixed;
  height: 100vh;
}

.logo {
  text-align: center;
  margin-bottom: 40px;
}

.logo h2 {
  font-size: 1.5rem;
  margin: 0;
}

.logo p {
  font-size: 0.85rem;
  color: #888;
  margin: 5px 0 0;
}

.nav-menu {
  list-style: none;
  padding: 0;
}

.nav-menu li {
  margin-bottom: 8px;
}

.nav-link {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  color: #ccc;
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.3s;
}

.nav-link:hover,
.nav-link.router-link-active {
  background: #16213e;
  color: #fff;
}

.nav-link .icon {
  margin-right: 12px;
  font-size: 1.2rem;
}

.main-content {
  flex: 1;
  margin-left: 240px;
  background: #f5f6fa;
  min-height: 100vh;
}

.top-bar {
  background: #fff;
  padding: 20px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
}

.top-bar h1 {
  font-size: 1.5rem;
  margin: 0;
  color: #333;
}

.time {
  color: #888;
  font-size: 0.9rem;
}

.content {
  padding: 30px;
}
</style>