import { RouteRecordRaw } from 'vue-router'

const orderRoutes: RouteRecordRaw[] = [
  {
    path: '/orders',
    name: 'OrderList',
    component: () => import('@/views/OrderListView.vue'),
    meta: {
      title: '订单列表',
      requiresAuth: true
    }
  },
  {
    path: '/orders/create',
    name: 'OrderCreate',
    component: () => import('@/views/OrderCreateView.vue'),
    meta: {
      title: '创建订单',
      requiresAuth: true
    }
  },
  {
    path: '/orders/:id',
    name: 'OrderDetail',
    component: () => import('@/views/OrderDetailView.vue'),
    meta: {
      title: '订单详情',
      requiresAuth: true
    },
    props: true
  },
  {
    path: '/orders/:id/edit',
    name: 'OrderEdit',
    component: () => import('@/views/OrderCreateView.vue'),
    meta: {
      title: '编辑订单',
      requiresAuth: true
    },
    props: true
  }
]

export default orderRoutes 