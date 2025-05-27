import './assets/main.css'
import 'element-plus/dist/index.css'

import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import i18n from '@/language/index.js'

createApp(App).use(ElementPlus).use(i18n).mount('#app')
