<template>
  <div class="file-preview">
    <el-image
      v-if="isImage"
      :src="fileUrl"
      :preview-src-list="[fileUrl]"
      fit="cover"
    />
    <div v-else class="file-icon">
      <el-icon :size="48">
        <Document v-if="isDocument" />
        <VideoCamera v-else-if="isVideo" />
        <Headset v-else-if="isAudio" />
        <Files v-else />
      </el-icon>
      <span class="file-name">{{ fileName }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Document, VideoCamera, Headset, Files } from '@element-plus/icons-vue'

const props = defineProps<{
  fileUrl: string
  fileName: string
  fileType: string
}>()

const isImage = computed(() => {
  return /^image\//.test(props.fileType)
})

const isDocument = computed(() => {
  return /^application\/(pdf|msword|vnd\.openxmlformats-officedocument\.wordprocessingml\.document)$/.test(props.fileType)
})

const isVideo = computed(() => {
  return /^video\//.test(props.fileType)
})

const isAudio = computed(() => {
  return /^audio\//.test(props.fileType)
})
</script>

<style scoped>
.file-preview {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  overflow: hidden;
}

.file-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.file-name {
  font-size: 14px;
  color: #606266;
  word-break: break-all;
  text-align: center;
  max-width: 200px;
}
</style> 