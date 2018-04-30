import Vue from 'vue'
import Router from 'vue-router'
import top from '@/components/top'
import regist from '@/components/register'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'top',
      component: top
    },
    {
      path: '/regist',
      name: 'regist',
      component: regist
    }
  ]
})
