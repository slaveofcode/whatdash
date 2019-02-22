import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import VueTimeago from 'vue-timeago'

import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap-vue/dist/bootstrap-vue.css"

Vue.use(VueRouter)
Vue.use(BootstrapVue)
Vue.use(VueTimeago, {
  locale: 'en',
  locales: {
    id: require('date-fns/locale/id')
  }
})

import Home from './pages/Home.vue'
import Register from './pages/Register.vue'
import Reconnect from './pages/Reconnect.vue'
import ChatWindow from './pages/ChatWindow.vue'

const routes = [
  { path: '/', name: 'home', component: Home, },
  { path: '/register', name: 'register', component: Register, },
  { path: '/reconnect', name: 'reconnect', component: Reconnect, },
  { path: '/chat/:id', name: 'chat', component: ChatWindow, },
]

const router = new VueRouter({ routes, })

new Vue({
  router,
  template: `
    <transition>
      <router-view></router-view>
    </transition>
  `,
}).$mount('#app')
