import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/', name: 'Landing', component: () => import('../views/Landing.vue'), meta: {} },
  { path: '/posts', name: 'Posts', component: () => import('../views/Home.vue'), meta: {} },
  { path: '/post/:id', name: 'PostDetail', component: () => import('../views/PostDetail.vue') },
  { path: '/publish', name: 'Publish', component: () => import('../views/Publish.vue'), meta: { requireAuth: true } },
  { path: '/user', name: 'User', component: () => import('../views/User.vue'), meta: { requireAuth: true } },
  { path: '/notifications', name: 'Notifications', component: () => import('../views/Notifications.vue'), meta: { requireAuth: true } },
  { path: '/search', name: 'Search', component: () => import('../views/Search.vue') },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { hideNav: false } },
]

const router = createRouter({
  history: createWebHistory('/c/'),
  routes,
  linkActiveClass: 'active',
  linkExactActiveClass: 'active',
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.meta.requireAuth && !auth.token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
