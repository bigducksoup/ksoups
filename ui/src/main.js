import { createApp } from 'vue'
import App from './App.vue'
import './index.css'
import router from  './router'
import 'vfonts/Lato.css'
import {init} from "./graph/custom.js";
init()

createApp(App).use(router).mount('#app')
