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
    meta: { title: 'SmashHub - Badminton Court Booking (LINE LIFF)' }
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminDashboard,
    meta: { title: 'SmashHub - Backoffice & POS System' }
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
  // Prevent non-admin users from accessing /admin backoffice
  if (to.path === '/admin') {
    const savedUser = localStorage.getItem('smashhub_user')
    if (savedUser) {
      try {
        const user = JSON.parse(savedUser)
        if (user && user.role !== 'ADMIN') {
          return next('/liff')
        }
      } catch (e) {
        // ignore parse error
      }
    }
  }
  next()
})

export default router
