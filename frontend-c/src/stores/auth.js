import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('c_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('c_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)

  function setLogin(t, u) {
    token.value = t
    user.value = u
    if (t) localStorage.setItem('c_token', t)
    if (u) localStorage.setItem('c_user', JSON.stringify(u))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('c_token')
    localStorage.removeItem('c_user')
  }

  return { token, user, isLoggedIn, setLogin, logout }
})
