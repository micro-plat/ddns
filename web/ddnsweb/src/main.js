import "jquery"
import Vue from 'vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import 'view-design/dist/styles/iview.css';
import store from './utility/store'
import DateFilter from './utility/filter';
import mavonEditor from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'
import ViewUI from 'view-design';
import htmlToPdf from '@/utility/htmlToPdf'
import { EnumUtility, EnumFilter } from 'qxnw-enum';
import App from './App.vue'
import VueClipboard from 'vue-clipboard2'

import router from './utility/router'
import {
  get,
  post
} from './utility/http'



Vue.prototype.$get = get;
Vue.prototype.$post = post;
Vue.use(mavonEditor)
Vue.use(ElementUI);
Vue.use(VueClipboard)
Vue.config.productionTip = false
console.log("当前环境：", process.env.NODE_ENV)
Vue.use(ViewUI);
Vue.use(htmlToPdf)

// console.log("当前环境：", process.env.NODE_ENV)
// console.log("当前环境：", process.env.VUE_APP_API_URL)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')