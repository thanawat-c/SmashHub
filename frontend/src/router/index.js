import { createRouter, createWebHistory } from 'vue-router'
import LiffView from '../views/liff/LiffView.vue'
import AdminDashboard from '../views/admin/AdminDashboard.vue'
import DemoView from '../views/demo/DemoView.vue'

const routes = [
  {
    path: '/',
    redirect: '/liff'
  },
  {
    path: '/liff',
    name: 'liff',
    component: LiffView,
    meta: { title: 'SmashHub - จองคอร์ทแบดมินตัน (LINE LIFF)' }
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminDashboard,
    meta: { title: 'SmashHub - ระบบจัดการหลังร้าน & POS' }
  },
  {
    path: '/demo',
    name: 'demo',
    component: DemoView,
    meta: { title: 'SmashHub - Split-screen Test' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }
  next()
})

export default router
