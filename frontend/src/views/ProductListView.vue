<template>
  <div class="min-h-screen bg-gray-100">
    <!-- 顶部导航 -->
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-indigo-600">SewingMast</h1>
            </div>
          </div>
          <div class="flex items-center">
            <router-link
              to="/dashboard"
              class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium"
            >
              返回仪表盘
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- 标题和操作按钮 -->
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            产品管理
          </h2>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <router-link
            to="/products/create"
            class="ml-3 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            创建产品
          </router-link>
        </div>
      </div>

      <!-- 搜索和筛选 -->
      <div class="mt-4">
        <div class="flex flex-col sm:flex-row">
          <div class="flex-1">
            <div class="relative rounded-md shadow-sm">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                v-model="searchQuery"
                type="text"
                class="focus:ring-indigo-500 focus:border-indigo-500 block w-full pl-10 sm:text-sm border-gray-300 rounded-md"
                placeholder="搜索产品..."
              />
            </div>
          </div>
          <div class="mt-4 sm:mt-0 sm:ml-4">
            <select
              v-model="categoryFilter"
              class="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            >
              <option value="">所有类别</option>
              <option v-for="category in categories" :key="category" :value="category">
                {{ category }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- 产品列表 -->
      <div class="mt-8 flex flex-col">
        <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
            <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      产品名称
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      类别
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      价格
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      库存
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      状态
                    </th>
                    <th scope="col" class="relative px-6 py-3">
                      <span class="sr-only">操作</span>
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="product in filteredProducts" :key="product.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10">
                          <img class="h-10 w-10 rounded-full" :src="product.image" :alt="product.name" />
                        </div>
                        <div class="ml-4">
                          <div class="text-sm font-medium text-gray-900">
                            {{ product.name }}
                          </div>
                          <div class="text-sm text-gray-500">
                            {{ product.sku }}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ product.category }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      ¥{{ product.price }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ product.stock }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span :class="getStatusClass(product.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                        {{ product.status }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <router-link :to="'/products/' + product.id" class="text-indigo-600 hover:text-indigo-900 mr-4">
                        查看
                      </router-link>
                      <router-link :to="'/products/' + product.id + '/edit'" class="text-indigo-600 hover:text-indigo-900">
                        编辑
                      </router-link>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// 搜索和筛选
const searchQuery = ref('')
const categoryFilter = ref('')
const categories = ref(['服装', '配饰', '面料', '辅料'])

// 产品列表
const products = ref([])

// 获取状态样式
const getStatusClass = (status: string) => {
  const classes = {
    '在售': 'bg-green-100 text-green-800',
    '缺货': 'bg-red-100 text-red-800',
    '下架': 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

// 过滤后的产品列表
const filteredProducts = computed(() => {
  return products.value.filter(product => {
    const matchesSearch = product.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         product.sku.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory = !categoryFilter.value || product.category === categoryFilter.value
    return matchesSearch && matchesCategory
  })
})

// 获取产品列表
const fetchProducts = async () => {
  try {
    // TODO: 调用API获取数据
    // 模拟数据
    products.value = [
      {
        id: 1,
        name: '夏季连衣裙',
        sku: 'PRD-2024-001',
        category: '服装',
        price: 299.00,
        stock: 100,
        status: '在售',
        image: 'https://images.unsplash.com/photo-1515886657613-9f3515b0c78f?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60'
      },
      {
        id: 2,
        name: '休闲衬衫',
        sku: 'PRD-2024-002',
        category: '服装',
        price: 199.00,
        stock: 0,
        status: '缺货',
        image: 'https://images.unsplash.com/photo-1598033129183-c4f50c736f10?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60'
      },
      {
        id: 3,
        name: '棉麻面料',
        sku: 'PRD-2024-003',
        category: '面料',
        price: 89.00,
        stock: 500,
        status: '在售',
        image: 'https://images.unsplash.com/photo-1520013333831-acd5bdee5b1a?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60'
      }
    ]
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

onMounted(() => {
  fetchProducts()
})
</script> 