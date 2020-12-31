import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    // {
    //   path: '/',
    //   name: 'Home',
    //   component: () => import('../pages/Home.vue')
    // },
    {
      path: '/',
      name: 'menu',
      component: () => import('../pages/menu.vue'),
      meta: {
        name: "用户权限系统"
      },
      children: [{
        path: 'home',
        name: 'Home',
        component: () => import('../pages/Home.vue')
      }]
    },
  ]
})