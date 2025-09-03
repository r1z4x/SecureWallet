import { createI18n } from 'vue-i18n'
import en from './locales/en.js'
import tr from './locales/tr.js'
import es from './locales/es.js'

const messages = {
  en,
  tr,
  es
}

// Tarayıcı dilini algıla
const getDefaultLocale = () => {
  const savedLocale = localStorage.getItem('locale')
  if (savedLocale && messages[savedLocale]) {
    return savedLocale
  }
  
  const browserLocale = navigator.language.split('-')[0]
  return messages[browserLocale] ? browserLocale : 'en'
}

export default createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: 'en',
  messages,
  globalInjection: true
})
