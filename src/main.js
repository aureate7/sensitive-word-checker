// src/main.js
import { createApp } from 'vue'
import App from './App.vue'

import router from './router'

import ElementPlus from 'element-plus'
import 'bootstrap-icons/font/bootstrap-icons.min.css'
import 'element-plus/dist/index.css'
import './style.css'     // 如果没有这个文件，可以先删掉这一行

const app = createApp(App)

// 生产环境配置
if (process.env.NODE_ENV === 'production') {
  // 禁止生产环境中的生产提示
  app.config.productionTip = false;
  // 禁止生产环境中的 Vue DevTools
  app.config.devtools = false;
  // 禁止任何 Vue 的警告和日志输出
  app.config.silent = true;
  
  // 禁用 console 输出
  console.log = () => {};   // 禁用 console.log
  console.warn = () => {};  // 禁用 console.warn
  console.error = () => {}; // 禁用 console.error
}
// # 在命令行中设置环境变量，或者通过 .env 文件
// NODE_ENV=production npm run build

app.use(router)       
app.use(ElementPlus)  // ElementPlus
app.mount('#app')
