import { createRouter, createWebHashHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import APIKeys from '../views/APIKeys.vue'
import APIList from '../views/APIList.vue'

const routes = [
    {
        path: '/',
        name: 'Dashboard',
        component: Dashboard
    },
    {
        path: '/api-keys',
        name: 'APIKeys',
        component: APIKeys
    },
    {
        path: '/api-list',
        name: 'APIList',
        component: APIList
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router