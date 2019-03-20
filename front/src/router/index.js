import Vue from 'vue'
import Router from 'vue-router'
import top from '@/components/top'
import regist from '@/components/register'
import welcome from '@/components/welcome'
import goodbye from '@/components/goodbye'
import question from '@/components/question'
import forgot from '@/components/forgot'

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
    },
    {
      path: '/welcome',
      name: 'welcome',
      component: welcome
    },
    {
      path: '/goodbye',
      name: 'goodbye',
      component: goodbye
    },
    {
      path: '/question',
      name: 'question',
      component: question
    },
    {
      path: '/forgot',
      name: 'forgot',
      component: forgot
    }
  ]
})
