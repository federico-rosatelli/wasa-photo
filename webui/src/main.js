import {createApp, reactive} from 'vue'
import { BModal } from 'bootstrap-vue'

import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Token from './components/Token.vue'
import ImageComponent from './components/ImageComponent.vue'
import FollowComponent from './components/FollowComponent.vue'
import CommentComponents from './components/CommentComponent.vue'
import './assets/dashboard.css'
import './assets/main.css'
import './assets/login.css'
console.log("PROVAAA");
const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Token", Token);
app.component("ImageComponent", ImageComponent);
app.component("FollowComponent",FollowComponent)
app.component("CommentComponents",CommentComponents)
app.use(router)

app.component('b-modal', BModal)
app.mount('#app')
