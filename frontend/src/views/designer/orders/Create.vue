<template>
  <div class="create-order-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('order.create.title') }}</span>
        </div>
      </template>

      <el-form
        :model="orderForm"
        :rules="rules"
        ref="orderFormRef"
        label-width="120px"
        class="order-form"
      >
        <el-form-item :label="$t('order.create.basicInfo')" class="form-section-title">
          <div class="form-section-content">
            <el-form-item prop="title" :label="$t('order.create.title')">
              <el-input v-model="orderForm.title" />
            </el-form-item>

            <el-form-item prop="description" :label="$t('order.create.description')">
              <el-input
                v-model="orderForm.description"
                type="textarea"
                :rows="3"
              />
            </el-form-item>

            <el-form-item prop="quantity" :label="$t('order.create.quantity')">
              <el-input-number
                v-model="orderForm.quantity"
                :min="1"
                :max="10000"
              />
            </el-form-item>

            <el-form-item prop="unitPrice" :label="$t('order.create.unitPrice')">
              <el-input-number
                v-model="orderForm.unitPrice"
                :precision="2"
                :step="0.1"
                :min="0"
              />
            </el-form-item>

            <el-form-item prop="dueDate" :label="$t('order.create.dueDate')">
              <el-date-picker
                v-model="orderForm.dueDate"
                type="date"
                :placeholder="$t('order.create.selectDate')"
              />
            </el-form-item>
          </div>
        </el-form-item>

        <el-form-item :label="$t('order.create.requirements')" class="form-section-title">
          <div class="form-section-content">
            <el-form-item prop="requirements" :label="$t('order.create.requirements')">
              <el-input
                v-model="orderForm.requirements"
                type="textarea"
                :rows="4"
              />
            </el-form-item>

            <el-form-item :label="$t('order.create.attachments')">
              <el-upload
                class="upload-demo"
                action="/api/upload"
                :on-success="handleUploadSuccess"
                :on-error="handleUploadError"
                :before-upload="beforeUpload"
                multiple
                :limit="5"
              >
                <el-button type="primary">{{ $t('order.create.upload') }}</el-button>
                <template #tip>
                  <div class="el-upload__tip">
                    {{ $t('order.create.uploadTip') }}
                  </div>
                </template>
              </el-upload>
            </el-form-item>
          </div>
        </el-form-item>

        <el-form-item :label="$t('order.create.factory')" class="form-section-title">
          <div class="form-section-content">
            <el-form-item prop="factoryId" :label="$t('order.create.selectFactory')">
              <el-select
                v-model="orderForm.factoryId"
                :placeholder="$t('order.create.selectFactory')"
                filterable
                remote
                :remote-method="searchFactories"
                :loading="factoryLoading"
              >
                <el-option
                  v-for="factory in factories"
                  :key="factory.id"
                  :label="factory.name"
                  :value="factory.id"
                >
                  <span>{{ factory.name }}</span>
                  <span class="factory-info">
                    {{ factory.capacity }} | {{ factory.location }}
                  </span>
                </el-option>
              </el-select>
            </el-form-item>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitOrder" :loading="submitting">
            {{ $t('order.create.submit') }}
          </el-button>
          <el-button @click="cancel">{{ $t('order.create.cancel') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default defineComponent({
  name: 'CreateOrder',
  setup() {
    const router = useRouter()
    const { t } = useI18n()
    const orderFormRef = ref()
    const submitting = ref(false)
    const factoryLoading = ref(false)
    const factories = ref([])

    const orderForm = reactive({
      title: '',
      description: '',
      quantity: 1,
      unitPrice: 0,
      dueDate: '',
      requirements: '',
      factoryId: null,
      attachments: []
    })

    const rules = {
      title: [{ required: true, message: t('order.create.rules.title'), trigger: 'blur' }],
      description: [{ required: true, message: t('order.create.rules.description'), trigger: 'blur' }],
      quantity: [{ required: true, message: t('order.create.rules.quantity'), trigger: 'blur' }],
      unitPrice: [{ required: true, message: t('order.create.rules.unitPrice'), trigger: 'blur' }],
      dueDate: [{ required: true, message: t('order.create.rules.dueDate'), trigger: 'change' }],
      requirements: [{ required: true, message: t('order.create.rules.requirements'), trigger: 'blur' }],
      factoryId: [{ required: true, message: t('order.create.rules.factory'), trigger: 'change' }]
    }

    const searchFactories = async (query: string) => {
      if (query.length < 2) return
      factoryLoading.value = true
      try {
        const response = await axios.get('/api/factories/search', {
          params: { query }
        })
        factories.value = response.data
      } catch (error) {
        console.error('Failed to search factories:', error)
      } finally {
        factoryLoading.value = false
      }
    }

    const handleUploadSuccess = (response: any, file: any) => {
      orderForm.attachments.push({
        name: file.name,
        url: response.url
      })
      ElMessage.success(t('order.create.uploadSuccess'))
    }

    const handleUploadError = () => {
      ElMessage.error(t('order.create.uploadError'))
    }

    const beforeUpload = (file: any) => {
      const isLt10M = file.size / 1024 / 1024 < 10
      if (!isLt10M) {
        ElMessage.error(t('order.create.fileTooLarge'))
        return false
      }
      return true
    }

    const submitOrder = async () => {
      try {
        await orderFormRef.value.validate()
        submitting.value = true

        const response = await axios.post('/api/orders', {
          ...orderForm,
          totalPrice: orderForm.quantity * orderForm.unitPrice
        })

        ElMessage.success(t('order.create.success'))
        router.push('/designer/orders')
      } catch (error) {
        ElMessage.error(t('order.create.error'))
      } finally {
        submitting.value = false
      }
    }

    const cancel = () => {
      router.push('/designer/orders')
    }

    return {
      orderForm,
      orderFormRef,
      rules,
      submitting,
      factoryLoading,
      factories,
      searchFactories,
      handleUploadSuccess,
      handleUploadError,
      beforeUpload,
      submitOrder,
      cancel
    }
  }
})
</script>

<style scoped>
.create-order-container {
  padding: 20px;
}

.form-section-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 10px;
}

.form-section-content {
  padding-left: 20px;
}

.factory-info {
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}

.el-upload__tip {
  color: #909399;
  font-size: 12px;
  margin-top: 5px;
}
</style> 