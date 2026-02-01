import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { public: true } },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    meta: { requireAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: () => import('../views/Dashboard.vue') },
      { path: 'content', name: 'Content', component: () => import('../views/Content.vue') },
      { path: 'content/:id', name: 'ContentDetail', component: () => import('../views/ContentDetail.vue') },
      { path: 'users', name: 'Users', component: () => import('../views/Users.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes,
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
