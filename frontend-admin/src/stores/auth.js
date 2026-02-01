import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const admin = ref(JSON.parse(localStorage.getItem('admin_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)

  function setLogin(t, u) {
    token.value = t
    admin.value = u
    if (t) localStorage.setItem('admin_token', t)
    if (u) localStorage.setItem('admin_user', JSON.stringify(u))
  }

  function logout() {
    token.value = ''
    admin.value = null
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_user')
  }

  return { token, admin, isLoggedIn, setLogin, logout }
})
