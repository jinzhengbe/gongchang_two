import { format } from 'date-fns'
import { zhCN, enUS, ko, vi } from 'date-fns/locale'

const locales = {
  zh: zhCN,
  en: enUS,
  ko: ko,
  vi: vi
}

export function formatDate(date: string | Date, lang: string = 'zh'): string {
  const locale = locales[lang as keyof typeof locales] || zhCN
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale })
} 