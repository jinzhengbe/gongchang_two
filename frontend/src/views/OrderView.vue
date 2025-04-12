<template>
  <div class="order-container">
    <div class="order-nav">
      <h2>订单管理</h2>
      <div class="nav-actions">
        <button class="btn-create" @click="navigateToCreate">创建订单</button>
        <div class="search-bar">
          <input
            type="text"
            v-model="searchQuery"
            placeholder="搜索订单号/客户/产品"
            @input="handleSearch"
          />
          <select v-model="statusFilter" @change="handleStatusFilter">
            <option value="">全部状态</option>
            <option v-for="status in statuses" :key="status" :value="status">
              {{ status }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <div class="order-list">
      <div class="order-header">
        <div>订单号</div>
        <div>客户名称</div>
        <div>产品名称</div>
        <div>数量</div>
        <div>总价</div>
        <div>状态</div>
        <div>创建时间</div>
      </div>

      <div v-if="filteredOrders.length === 0" class="no-orders">
        没有找到符合条件的订单
      </div>

      <div v-else class="order-items">
        <div 
          v-for="order in filteredOrders" 
          :key="order.id" 
          class="order-item"
          @click="navigateToDetail(order.id)"
        >
          <div>{{ order.orderNumber }}</div>
          <div>{{ order.customerName }}</div>
          <div>{{ order.productName }}</div>
          <div>{{ order.quantity }}</div>
          <div>¥{{ order.totalPrice.toFixed(2) }}</div>
          <div :class="['status', orderService.getStatusClass(order.status)]">
            {{ order.status }}
          </div>
          <div>{{ orderService.formatDate(order.createdAt) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { OrderService, OrderStatus } from '@/services/orderService'
import '@/styles/order.css'

const router = useRouter()
const orderService = new OrderService()
const searchQuery = ref('')
const statusFilter = ref<OrderStatus | ''>('')
const statuses = orderService.getStatuses()
const filteredOrders = ref<Order[]>([])

const navigateToCreate = () => {
  router.push('/orders/create')
}

const navigateToDetail = (orderId: string) => {
  router.push(`/orders/${orderId}`)
}

const handleSearch = () => {
  orderService.setSearchQuery(searchQuery.value)
  filteredOrders.value = orderService.getFilteredOrders()
}

const handleStatusFilter = () => {
  orderService.setStatusFilter(statusFilter.value)
  filteredOrders.value = orderService.getFilteredOrders()
}

onMounted(async () => {
  await orderService.fetchOrders()
  filteredOrders.value = orderService.getFilteredOrders()
})
</script>

<style scoped>
.order-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.order-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.nav-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.btn-create {
  padding: 0.5rem 1rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
}

.search-bar {
  display: flex;
  gap: 1rem;
}

.search-bar input {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  width: 300px;
}

.search-bar select {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  min-width: 150px;
}

.order-list {
  background-color: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.order-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  padding: 1rem;
  background-color: #f3f4f6;
  font-weight: 500;
  border-bottom: 1px solid #e5e7eb;
}

.order-items {
  display: flex;
  flex-direction: column;
}

.order-item {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
  cursor: pointer;
  transition: background-color 0.2s;
}

.order-item:hover {
  background-color: #f9fafb;
}

.no-orders {
  padding: 2rem;
  text-align: center;
  color: #6b7280;
}

.status {
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.status.pending {
  background-color: #fef3c7;
  color: #92400e;
}

.status.processing {
  background-color: #dbeafe;
  color: #1e40af;
}

.status.completed {
  background-color: #dcfce7;
  color: #166534;
}
</style> 