import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'
import Integrations from './components/Integrations.vue'
import LoginForm from './components/LoginForm.vue'
import PluginDetails from './components/PluginDetails.vue'
import Rules from './components/Rules.vue'
import Settings from './components/Settings.vue'
import { BootstrapVue } from 'bootstrap-vue'
import store from './store/store'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(BootstrapVue);
Vue.use(VueRouter);

const routes = [
  { name: 'home', path: '/', redirect: '/integrations' },
  { name: 'integrations', path: '/integrations', component: Integrations },
  { name: 'rules', path: '/rules', component: Rules },
  { name: 'settings', path: '/settings', component: Settings },
  { name: 'login', path: '/login', component: LoginForm },
  { name: 'add-plugin', path: '/plugin', component: PluginDetails },
  { name: 'plugin', path: '/plugin/:name', component: PluginDetails }
];

export const router = new VueRouter({
  routes, mode: 'history'
});

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
