<template>
  <div class="file-uploader">
    <div
      class="upload-area"
      :class="{ 'is-dragover': isDragover }"
      @dragenter.prevent="isDragover = true"
      @dragleave.prevent="isDragover = false"
      @dragover.prevent
      @drop.prevent="handleDrop"
    >
      <el-upload
        ref="uploadRef"
        :action="uploadUrl"
        :headers="headers"
        :data="uploadData"
        :on-success="handleSuccess"
        :on-error="handleError"
        :before-upload="beforeUpload"
        :on-progress="handleProgress"
        :show-file-list="false"
        :multiple="multiple"
        :accept="accept"
        :auto-upload="false"
        drag
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          {{ $t('fileUploader.dropFile') }} <em>{{ $t('fileUploader.clickToUpload') }}</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            {{ $t('fileUploader.tip', { size: maxSize / 1024 / 1024 }) }}
          </div>
        </template>
      </el-upload>
    </div>

    <div v-if="uploadQueue.length > 0" class="upload-queue">
      <div class="queue-header">
        <span>{{ $t('fileUploader.queue', { count: uploadQueue.length }) }}</span>
        <el-button
          type="primary"
          size="small"
          :loading="isUploading"
          @click="handleStartUpload"
        >
          {{ $t('fileUploader.startUpload') }}
        </el-button>
      </div>
      <div class="queue-list">
        <div v-for="file in uploadQueue" :key="file.uid" class="queue-item">
          <div class="queue-info">
            <span class="file-name">{{ file.name }}</span>
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
          </div>
          <el-button
            type="danger"
            link
            @click="handleRemoveFromQueue(file)"
          >
            <el-icon><delete /></el-icon>
          </el-button>
        </div>
      </div>
    </div>

    <div v-if="uploadingFiles.length > 0" class="upload-progress">
      <div v-for="file in uploadingFiles" :key="file.name" class="progress-item">
        <div class="progress-info">
          <span class="file-name">{{ file.name }}</span>
          <span class="progress-percentage">{{ file.percentage }}%</span>
        </div>
        <el-progress
          :percentage="file.percentage"
          :status="file.status"
          :stroke-width="4"
        />
      </div>
    </div>

    <div v-if="fileList.length > 0" class="file-list">
      <div v-for="file in fileList" :key="file.id" class="file-item">
        <div class="file-icon">
          <el-icon v-if="isImageFile(file)" class="image-icon"><picture /></el-icon>
          <el-icon v-else-if="isPDFFile(file)" class="pdf-icon"><document /></el-icon>
          <el-icon v-else-if="isWordFile(file)" class="word-icon"><document /></el-icon>
          <el-icon v-else-if="isExcelFile(file)" class="excel-icon"><document /></el-icon>
          <el-icon v-else class="default-icon"><document /></el-icon>
        </div>
        <div class="file-info">
          <div class="file-name">{{ file.name }}</div>
          <div class="file-meta">
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
            <span class="file-type">{{ getFileType(file) }}</span>
          </div>
        </div>
        <div class="file-actions">
          <el-button
            type="primary"
            link
            @click="handlePreview(file)"
          >
            <el-icon><view /></el-icon>
          </el-button>
          <el-button
            type="danger"
            link
            @click="handleRemove(file)"
          >
            <el-icon><delete /></el-icon>
          </el-button>
        </div>
      </div>
    </div>

    <file-preview
      v-if="previewFile"
      :file="previewFile"
      @close="previewFile = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { UploadFilled, Document, Delete, View, Picture } from '@element-plus/icons-vue'
import type { UploadInstance, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'
import FilePreview from './FilePreview.vue'

const props = defineProps<{
  orderId: string
  multiple?: boolean
  accept?: string
  maxSize?: number
}>()

const emit = defineEmits<{
  (e: 'upload-success', file: any): void
  (e: 'upload-error', error: any): void
  (e: 'remove', file: any): void
}>()

const { t } = useI18n()
const uploadRef = ref<UploadInstance>()
const isDragover = ref(false)
const fileList = ref<any[]>([])
const previewFile = ref<any>(null)
const uploadingFiles = ref<Array<{
  name: string
  percentage: number
  status: 'success' | 'exception' | 'warning' | ''
}>>([])
const uploadQueue = ref<File[]>([])
const isUploading = ref(false)

const uploadUrl = computed(() => `/api/files/upload`)
const headers = computed(() => ({
  'Authorization': `Bearer ${localStorage.getItem('token')}`
}))
const uploadData = computed(() => ({
  orderId: props.orderId
}))

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / 1024 / 1024).toFixed(2) + ' MB'
  return (size / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

const handleDrop = (e: DragEvent) => {
  isDragover.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    handleFiles(files)
  }
}

const handleFiles = (files: FileList) => {
  Array.from(files).forEach(file => {
    if (props.accept && !file.type.match(props.accept)) {
      ElMessage.error(t('fileUploader.invalidType'))
      return
    }
    if (props.maxSize && file.size > props.maxSize) {
      ElMessage.error(t('fileUploader.tooLarge'))
      return
    }
    uploadQueue.value.push(file)
  })
}

const handleStartUpload = async () => {
  if (isUploading.value) return
  
  isUploading.value = true
  try {
    for (const file of uploadQueue.value) {
      await uploadRef.value?.upload(file)
    }
    uploadQueue.value = []
  } catch (error) {
    ElMessage.error(t('fileUploader.uploadError'))
  } finally {
    isUploading.value = false
  }
}

const handleRemoveFromQueue = (file: File) => {
  const index = uploadQueue.value.findIndex(f => f.uid === file.uid)
  if (index > -1) {
    uploadQueue.value.splice(index, 1)
  }
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  if (props.accept && !file.type.match(props.accept)) {
    ElMessage.error(t('fileUploader.invalidType'))
    return false
  }
  if (props.maxSize && file.size > props.maxSize) {
    ElMessage.error(t('fileUploader.tooLarge'))
    return false
  }
  return true
}

const handleProgress: UploadProps['onProgress'] = (event, file) => {
  const index = uploadingFiles.value.findIndex(f => f.name === file.name)
  if (index === -1) {
    uploadingFiles.value.push({
      name: file.name,
      percentage: 0,
      status: ''
    })
  } else {
    uploadingFiles.value[index].percentage = Math.round(event.percent || 0)
  }
}

const handleSuccess: UploadProps['onSuccess'] = (response, file) => {
  const index = uploadingFiles.value.findIndex(f => f.name === file.name)
  if (index > -1) {
    uploadingFiles.value[index].status = 'success'
    setTimeout(() => {
      uploadingFiles.value.splice(index, 1)
    }, 1000)
  }

  fileList.value.push({
    id: response.id,
    name: file.name,
    size: file.size,
    url: response.url
  })
  emit('upload-success', response)
  ElMessage.success(t('fileUploader.success'))
}

const handleError: UploadProps['onError'] = (error, file) => {
  const index = uploadingFiles.value.findIndex(f => f.name === file.name)
  if (index > -1) {
    uploadingFiles.value[index].status = 'exception'
    setTimeout(() => {
      uploadingFiles.value.splice(index, 1)
    }, 2000)
  }

  emit('upload-error', error)
  ElMessage.error(t('fileUploader.error'))
}

const handleRemove = (file: any) => {
  const index = fileList.value.indexOf(file)
  if (index > -1) {
    fileList.value.splice(index, 1)
    emit('remove', file)
  }
}

const handlePreview = (file: any) => {
  previewFile.value = file
}

const isImageFile = (file: any) => {
  return file.type?.startsWith('image/')
}

const isPDFFile = (file: any) => {
  return file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')
}

const isWordFile = (file: any) => {
  return file.type === 'application/msword' || 
         file.type === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document' ||
         file.name.toLowerCase().endsWith('.doc') ||
         file.name.toLowerCase().endsWith('.docx')
}

const isExcelFile = (file: any) => {
  return file.type === 'application/vnd.ms-excel' ||
         file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
         file.name.toLowerCase().endsWith('.xls') ||
         file.name.toLowerCase().endsWith('.xlsx')
}

const getFileType = (file: any) => {
  if (isImageFile(file)) return '图片'
  if (isPDFFile(file)) return 'PDF'
  if (isWordFile(file)) return 'Word'
  if (isExcelFile(file)) return 'Excel'
  return '文件'
}
</script>

<style scoped>
.file-uploader {
  width: 100%;
}

.upload-area {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.upload-area.is-dragover {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
}

.el-upload {
  width: 100%;
}

.el-upload__text {
  margin: 10px 0;
}

.el-upload__tip {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.file-list {
  margin-top: 20px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  margin-bottom: 8px;
  background-color: var(--el-bg-color);
  transition: all 0.3s;
}

.file-item:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.file-icon {
  margin-right: 12px;
  font-size: 24px;
}

.image-icon {
  color: var(--el-color-success);
}

.pdf-icon {
  color: var(--el-color-danger);
}

.word-icon {
  color: var(--el-color-primary);
}

.excel-icon {
  color: var(--el-color-success);
}

.default-icon {
  color: var(--el-text-color-secondary);
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.file-actions {
  display: flex;
  gap: 8px;
  margin-left: 12px;
}

.file-actions .el-button {
  padding: 4px;
}

.upload-progress {
  margin-top: 20px;
}

.progress-item {
  margin-bottom: 10px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.progress-percentage {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.upload-queue {
  margin-top: 20px;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  padding: 12px;
}

.queue-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.queue-list {
  max-height: 200px;
  overflow-y: auto;
}

.queue-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid var(--el-border-color);
}

.queue-item:last-child {
  border-bottom: none;
}

.queue-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.queue-info .file-name {
  font-size: 14px;
}

.queue-info .file-size {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style> 