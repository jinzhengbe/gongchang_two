import { createI18n } from 'vue-i18n'
import { useStorage } from '@vueuse/core'
import { format } from 'date-fns'
import { zhCN, enUS, ko, vi } from 'date-fns/locale'

// 导入语言包
import zhCNMessages from './locales/zh-CN'
import enUSMessages from './locales/en-US'
import koKRMessages from './locales/ko-KR'
import viVNMessages from './locales/vi-VN'

// 支持的语言类型
type SupportedLocale = 'zh-CN' | 'en-US' | 'ko-KR' | 'vi-VN'

// 日期格式化本地化配置
const dateLocales = {
  'zh-CN': zhCN,
  'en-US': enUS,
  'ko-KR': ko,
  'vi-VN': vi
}

// 获取浏览器语言
const getBrowserLanguage = (): SupportedLocale => {
  const lang = navigator.language
  const supportedLanguages: SupportedLocale[] = ['zh-CN', 'en-US', 'ko-KR', 'vi-VN']
  const defaultLang: SupportedLocale = 'en-US'

  if (lang.startsWith('zh')) return 'zh-CN'
  if (lang.startsWith('ko')) return 'ko-KR'
  if (lang.startsWith('vi')) return 'vi-VN'
  
  return supportedLanguages.includes(lang as SupportedLocale) ? (lang as SupportedLocale) : defaultLang
}

// 从本地存储获取语言设置
const storedLanguage = useStorage<SupportedLocale>('app-language', getBrowserLanguage())

// 创建i18n实例
const i18n = createI18n({
  legacy: false,
  locale: storedLanguage.value,
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCNMessages,
    'en-US': enUSMessages,
    'ko-KR': koKRMessages,
    'vi-VN': viVNMessages
  },
  numberFormats: {
    'zh-CN': {
      currency: {
        style: 'currency',
        currency: 'CNY',
        notation: 'standard'
      },
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 2
      }
    },
    'en-US': {
      currency: {
        style: 'currency',
        currency: 'USD',
        notation: 'standard'
      },
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 2
      }
    },
    'ko-KR': {
      currency: {
        style: 'currency',
        currency: 'KRW',
        notation: 'standard'
      },
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 2
      }
    },
    'vi-VN': {
      currency: {
        style: 'currency',
        currency: 'VND',
        notation: 'standard'
      },
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 2
      }
    }
  }
})

// 格式化日期
export const formatDate = (date: Date | string | number, formatStr: string = 'yyyy-MM-dd'): string => {
  const locale = dateLocales[i18n.global.locale.value as SupportedLocale]
  return format(new Date(date), formatStr, { locale })
}

// 格式化货币
export const formatCurrency = (value: number, locale?: SupportedLocale): string => {
  const currentLocale = locale || (i18n.global.locale.value as SupportedLocale)
  return new Intl.NumberFormat(currentLocale, {
    style: 'currency',
    currency: getCurrencyByLocale(currentLocale)
  }).format(value)
}

// 格式化数字
export const formatNumber = (value: number, minimumFractionDigits: number = 2): string => {
  return new Intl.NumberFormat(i18n.global.locale.value as SupportedLocale, {
    minimumFractionDigits,
    maximumFractionDigits: minimumFractionDigits
  }).format(value)
}

// 格式化百分比
export const formatPercent = (value: number, minimumFractionDigits: number = 2): string => {
  return new Intl.NumberFormat(i18n.global.locale.value as SupportedLocale, {
    style: 'percent',
    minimumFractionDigits,
    maximumFractionDigits: minimumFractionDigits
  }).format(value)
}

// 根据地区获取货币代码
const getCurrencyByLocale = (locale: SupportedLocale): string => {
  const currencyMap: Record<SupportedLocale, string> = {
    'zh-CN': 'CNY',
    'en-US': 'USD',
    'ko-KR': 'KRW',
    'vi-VN': 'VND'
  }
  return currencyMap[locale]
}

// 切换语言
export const setLanguage = (locale: SupportedLocale): void => {
  if (i18n.global.locale.value === locale) return
  
  i18n.global.locale.value = locale
  useStorage('app-language', locale)
  document.querySelector('html')?.setAttribute('lang', locale)
  
  // 触发自定义事件，用于通知其他组件语言已更改
  window.dispatchEvent(new CustomEvent('language-changed', { detail: locale }))
}

// 获取支持的语言列表
export const getSupportedLanguages = (): Array<{
  code: SupportedLocale
  name: string
  nativeName: string
}> => {
  return [
    { code: 'zh-CN', name: 'Chinese', nativeName: '简体中文' },
    { code: 'en-US', name: 'English', nativeName: 'English' },
    { code: 'ko-KR', name: 'Korean', nativeName: '한국어' },
    { code: 'vi-VN', name: 'Vietnamese', nativeName: 'Tiếng Việt' }
  ]
}

// 获取当前语言
export const getCurrentLanguage = (): SupportedLocale => i18n.global.locale.value as SupportedLocale

export default i18n 