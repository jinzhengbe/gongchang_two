import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Home from '@/views/Home.vue'
import orderRoutes from './order'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/designer-orders',
    name: 'DesignerOrders',
    component: () => import('@/views/DesignerOrders.vue')
  },
  {
    path: '/factory-list',
    name: 'FactoryList',
    component: () => import('@/views/FactoryList.vue')
  },
  {
    path: '/fabric-categories',
    name: 'FabricCategories',
    component: () => import('@/views/FabricCategories.vue')
  },
  {
    path: '/login/:role',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  ...orderRoutes
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 