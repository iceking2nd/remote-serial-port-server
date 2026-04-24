import { createI18n } from 'vue-i18n'
import enUS from './locales/en-US'
import zhCN from './locales/zh-CN'
import zhTW from './locales/zh-TW'
import frFR from './locales/fr-FR'
import deDE from './locales/de-DE'
import ruRU from './locales/ru-RU'
import plPL from './locales/pl-PL'

const supportedLocales = ['en-us', 'zh-cn', 'zh-tw', 'fr-fr', 'de-de', 'ru-ru', 'pl-pl']

function detectLocale() {
  const browserLangs = navigator.languages || [navigator.language]
  for (const lang of browserLangs) {
    const lower = lang.toLowerCase()
    if (supportedLocales.includes(lower)) {
      return lower
    }
    const base = lower.split('-')[0]
    const matched = supportedLocales.find(l => l.startsWith(base))
    if (matched) {
      return matched
    }
  }
  return 'en-us'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  messages: {
    'en-us': enUS,
    'zh-cn': zhCN,
    'zh-tw': zhTW,
    'fr-fr': frFR,
    'de-de': deDE,
    'ru-ru': ruRU,
    'pl-pl': plPL
  }
})

export default i18n
