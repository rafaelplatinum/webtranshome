import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { vCan } from './composables/usePermission'
import { vReveal } from './composables/useAnimasi'
import './assets/main.css'

const app = createApp(App)

// Pinia dipasang sebelum router karena route guard memakai store auth.
app.use(createPinia())
app.use(router)
app.directive('can', vCan)
app.directive('reveal', vReveal)

app.mount('#app')
