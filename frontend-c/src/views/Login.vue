<template>
  <div class="login-page">
    <div class="card login-card">
      <h1 class="page-title">登录</h1>
      <form @submit.prevent="onSubmit">
        <div class="form-group">
          <label for="login-phone">手机号</label>
          <input id="login-phone" v-model="phone" type="tel" placeholder="请输入手机号" required autocomplete="tel" />
        </div>
        <div class="form-group">
          <label for="login-code">验证码（mock 固定 123456）</label>
          <input id="login-code" v-model="code" type="text" placeholder="123456" autocomplete="one-time-code" />
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      <p class="login-hint">演示：验证码固定为 123456</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { auth as authApi } from '../api'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const phone = ref('')
const code = ref('123456')
const loading = ref(false)

async function onSubmit() {
  loading.value = true
  try {
    const { data } = await authApi.login(phone.value, code.value || '123456')
    auth.setLogin(data.token, data.user)
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    alert(e.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  max-width: 400px;
  margin: 48px auto;
  padding: 0 20px;
  width: 100%;
  min-width: 0;
  padding-left: max(20px, env(safe-area-inset-left));
  padding-right: max(20px, env(safe-area-inset-right));
}
.login-card { padding: 32px; }
.btn-block { width: 100%; }
.login-hint { margin-top: 20px; color: var(--color-text-muted); font-size: 13px; }
@media (max-width: 768px) {
  .login-page { margin: 32px auto; padding: 0 16px; padding-left: max(16px, env(safe-area-inset-left)); padding-right: max(16px, env(safe-area-inset-right)); }
  .login-card { padding: 24px; }
  .login-hint { font-size: 12px; }
}
@media (max-width: 375px) {
  .login-page { margin: 24px auto; padding: 0 12px; padding-left: max(12px, env(safe-area-inset-left)); padding-right: max(12px, env(safe-area-inset-right)); }
  .login-card { padding: 20px; }
  .login-hint { margin-top: 16px; font-size: 12px; }
}
</style>
