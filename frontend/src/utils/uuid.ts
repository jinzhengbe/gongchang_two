// 检查浏览器是否支持 crypto.randomUUID
const hasCrypto = typeof window !== 'undefined' && 
                 window.crypto && 
                 typeof window.crypto.randomUUID === 'function'

// 生成 UUID 的函数
function generateUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

// 导出函数，优先使用浏览器的实现
export const uuid = hasCrypto ? window.crypto.randomUUID : generateUUID 