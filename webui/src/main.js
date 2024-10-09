import {createApp, reactive} from 'vue'
// import { BModal, BButton } from 'bootstrap-vue-3'

import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Token from './components/Token.vue'
import ImageComponent from './components/ImageComponent.vue/index.js'
import FollowComponent from './components/FollowComponent.vue'
import CommentComponents from './components/CommentComponent.vue'
import ProfileImageComponent from './components/ProfileImageComponent.vue'
import './assets/dashboard.css'
import './assets/main.css'
import './assets/login.css'
import './assets/modal.css'
import './assets/dropdown.css'
console.log("PROVAAA");
const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Token", Token);
app.component("ImageComponent", ImageComponent);
app.component("FollowComponent",FollowComponent)
app.component("CommentComponents",CommentComponents)
app.component("ProfileImageComponent",ProfileImageComponent)
app.use(router)
app.mount('#app')

if (localStorage.getItem("Token") === null && window.location.pathname !== "/login"){
    location.replace("/login")
}