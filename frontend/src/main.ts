import { createApp } from 'vue'
import App from './App.vue'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import axios from 'axios'
import VueAxios from 'vue-axios'

import dayjs from 'dayjs'
import isLeapYear from 'dayjs/plugin/isLeapYear' // 导入插件
import 'dayjs/locale/zh-cn' // 导入本地化语言
dayjs.extend(isLeapYear)
dayjs.locale('zh-cn')

const app = createApp(App)
app.use(ElementPlus, {
  locale: zhCn,
})
app.use(VueAxios, axios)
app.config.globalProperties.$dayjs = dayjs
app.mount('#app')