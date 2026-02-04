<template>
  <div class="login-page">
    <div class="card login-card">
      <h1 class="page-title">登录</h1>
      <form @submit.prevent="onSubmit">
        <div v-if="errorMessage" class="form-error" role="alert">
          {{ errorMessage }}
        </div>
        <div class="form-group">
          <label for="login-phone">手机号</label>
          <input id="login-phone" v-model="phone" type="tel" placeholder="请输入手机号" required autocomplete="tel" @input="errorMessage = ''" />
        </div>
        <div class="form-group form-group-code">
          <label for="login-code">验证码（mock 固定 123456）</label>
          <div class="code-row">
            <input id="login-code" v-model="code" type="text" placeholder="123456" autocomplete="one-time-code" class="code-input" @input="errorMessage = ''" />
            <button type="button" class="btn btn-code" :disabled="codeCountdown > 0" @click="onGetCode">
              {{ codeCountdown > 0 ? `${codeCountdown}s 后重试` : '获取验证码' }}
            </button>
          </div>
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
import { ref, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { auth as authApi } from '../api'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const phone = ref('')
const code = ref('123456')
const loading = ref(false)
const codeCountdown = ref(0)
const errorMessage = ref('')
let codeTimer = null

function onGetCode() {
  if (codeCountdown.value > 0) return
  codeCountdown.value = 30
  codeTimer = setInterval(() => {
    codeCountdown.value -= 1
    if (codeCountdown.value <= 0 && codeTimer) {
      clearInterval(codeTimer)
      codeTimer = null
    }
  }, 1000)
}

onUnmounted(() => {
  if (codeTimer) clearInterval(codeTimer)
})

// 中国大陆手机号：1 开头，11 位数字
const PHONE_REG = /^1[3-9]\d{9}$/

function isValidPhone(val) {
  return PHONE_REG.test(String(val).trim())
}

async function onSubmit() {
  errorMessage.value = ''
  const phoneVal = phone.value.trim()
  if (!phoneVal) {
    errorMessage.value = '请填写手机号'
    return
  }
  if (!isValidPhone(phoneVal)) {
    errorMessage.value = '请填写正确的 11 位手机号（以 1 开头）'
    return
  }
  loading.value = true
  try {
    const { data } = await authApi.login(phoneVal, code.value || '123456')
    auth.setLogin(data.token, data.user)
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    errorMessage.value = e.response?.data?.error || '登录失败，请稍后重试'
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
.form-error {
  padding: 10px 14px;
  margin-bottom: 16px;
  background: #fef2f2;
  color: #b91c1c;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.4;
}
.btn-block { width: 100%; }
.form-group-code .code-row { display: flex; gap: 10px; align-items: stretch; }
.form-group-code .code-input { flex: 1; min-width: 0; }
.form-group-code .btn-code { flex-shrink: 0; white-space: nowrap; padding: 0 14px; }
.form-group-code .btn-code:disabled { opacity: 0.6; cursor: not-allowed; }
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
