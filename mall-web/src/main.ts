import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// Vant UI
import { Lazyload } from 'vant'
import '@vant/touch-emulator'

// Global styles
import 'virtual:windi.css'
import './styles/index.scss'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Lazyload)

app.mount('#app')