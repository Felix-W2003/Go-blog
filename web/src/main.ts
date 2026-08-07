import { createApp } from 'vue'                      // 1. 导入核心函数
import '@/assets/base.css'                          // 2. 导入全局样式
import * as ElementPlusIconsVue from '@element-plus/icons-vue' // 3. 导入图标库
import App from './App.vue'                        // 4. 导入根组件
import router from './router'                      // 5. 导入路由配置
import { pinia } from "@/stores";                  // 6. 导入状态管理

const app = createApp(App)   // 1. 创建 Vue 应用实例

// 2. 注册所有 Element Plus 图标为全局组件
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

// 3. 注册插件（Pinia 和 Router）
app.use(pinia).use(router)

// 4. 挂载到 DOM 元素 #app 上
app.mount('#app')