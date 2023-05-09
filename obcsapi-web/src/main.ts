import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { create,NForm,NFormItem,NInput,NInputNumber,NButton,NSwitch } from "naive-ui";

const app = createApp(App)

app.use(createPinia())
app.use(router)

const naive = create({
    components: [NForm,NFormItem,NInput,NInputNumber,NButton,NSwitch]
  })
  
app.use(naive)


app.mount('#app')
