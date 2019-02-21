import Vue from 'vue'
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
// import VueWebsocket from 'vue-websocket'

import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap-vue/dist/bootstrap-vue.css"

Vue.use(VueRouter)
Vue.use(BootstrapVue)
// Vue.use(VueWebsocket, 'http://localhost:8081', {
//   reconnection: true,
//   path: '/ws',
//   transports: ['websocket']
// })

import Home from './pages/Home.vue'
import Register from './pages/Register.vue'
import Reconnect from './pages/Reconnect.vue'

const routes = [
  { path: '/', name: 'home', component: Home, },
  { path: '/register', name: 'register', component: Register, },
  { path: '/reconnect', name: 'reconnect', component: Reconnect, },
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
