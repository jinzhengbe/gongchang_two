import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import router from './router'
import './assets/main.css'

// 创建 i18n 实例
const i18n = createI18n({
  legacy: false,
  locale: 'zh',
  fallbackLocale: 'en',
  messages: {
    zh: {
      common: {
        login: '登录',
        register: '注册',
        logout: '退出',
        search: '搜索',
        submit: '提交',
        cancel: '取消',
        confirm: '确认',
        delete: '删除',
        edit: '编辑',
        create: '创建',
        update: '更新',
        success: '成功',
        error: '错误',
        warning: '警告',
        info: '信息'
      },
      auth: {
        username: '用户名',
        password: '密码',
        email: '邮箱',
        role: '角色',
        loginSuccess: '登录成功',
        registerSuccess: '注册成功',
        logoutSuccess: '退出成功'
      },
      product: {
        name: '产品名称',
        description: '产品描述',
        category: '产品类别',
        price: '价格',
        stock: '库存',
        status: '状态',
        createProduct: '创建产品',
        updateProduct: '更新产品',
        deleteProduct: '删除产品',
        productList: '产品列表',
        productDetail: '产品详情'
      },
      order: {
        title: '订单标题',
        description: '订单描述',
        status: '订单状态',
        createOrder: '创建订单',
        updateOrder: '更新订单',
        deleteOrder: '删除订单',
        orderList: '订单列表',
        orderDetail: '订单详情'
      }
    },
    en: {
      // English translations
    }
  }
})

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.use(i18n)

app.mount('#app') 