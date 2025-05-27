import { createI18n } from 'vue-i18n'
import enUS from './locales/en-US'
import zhCN from './locales/zh-CN'

const i18n = createI18n({
  legacy: false,
  locale: 'en-us', // 默认显示语言
  messages: {
    'en-us': enUS,
    'zh-cn': zhCN
  }
})

export default i18n
