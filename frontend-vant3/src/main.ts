import { createApp } from 'vue'
import App from './App.vue'

import axios from 'axios'
import VueAxios from 'vue-axios'

import Vant from 'vant';
import 'vant/lib/index.css';

const app = createApp(App);
app.use(Vant);
app.use(VueAxios, axios);
app.mount('#app');