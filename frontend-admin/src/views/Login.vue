<template>
  <div class="login-page">
    <div class="login-bg" aria-hidden="true">
      <div class="login-bg-gradient" />
      <div class="login-bg-orb login-bg-orb-1" />
      <div class="login-bg-orb login-bg-orb-2" />
      <div class="login-bg-orb login-bg-orb-3" />
    </div>
    <div class="login-container">
      <div class="login-brand">
        <div class="login-brand-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
          </svg>
        </div>
        <h1 class="login-title">VibeCoding</h1>
        <p class="login-subtitle">社区后台管理</p>
      </div>
      <el-card class="login-card" shadow="hover">
        <template #header>
          <span class="login-card-title">管理员登录</span>
        </template>
        <el-form :model="form" @submit.prevent="onSubmit" class="login-form" label-position="top">
          <el-form-item label="用户名" class="login-form-item">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              size="large"
              clearable
              autocomplete="username"
              class="login-input"
            >
              <template #prefix>
                <span class="input-icon" aria-hidden="true">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="密码" class="login-form-item">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password
              clearable
              autocomplete="current-password"
              @keyup.enter="onSubmit"
              class="login-input"
            >
              <template #prefix>
                <span class="input-icon" aria-hidden="true">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item class="login-submit-item">
            <el-button
              type="primary"
              native-type="submit"
              :loading="loading"
              class="login-btn"
              size="large"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
      <p class="login-footer">仅限授权管理员使用</p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import { auth as authApi } from '../api'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const form = reactive({ username: 'admin', password: 'admin123' })
const loading = ref(false)

async function onSubmit() {
  loading.value = true
  try {
    const { data } = await authApi.login(form.username, form.password)
    auth.setLogin(data.token, data.admin)
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  font-family: 'Fira Sans', -apple-system, BlinkMacSystemFont, sans-serif;
}
.login-bg {
  position: fixed;
  inset: 0;
  background: #faf5ff;
  z-index: 0;
  overflow: hidden;
}
.login-bg-gradient {
  position: absolute;
  inset: -50%;
  background: linear-gradient(135deg, #faf5ff 0%, #f3e8ff 25%, #ede9fe 50%, #e9d5ff 75%, #f3e8ff 100%);
  background-size: 400% 400%;
  animation: gradientMove 18s ease infinite;
}
.login-bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.35;
  animation: orbFloat 20s ease-in-out infinite;
}
.login-bg-orb-1 {
  width: 320px;
  height: 320px;
  background: #c4b5fd;
  top: -10%;
  left: -5%;
  animation-delay: 0s;
}
.login-bg-orb-2 {
  width: 280px;
  height: 280px;
  background: #a78bfa;
  top: 50%;
  right: -8%;
  animation-delay: -7s;
  animation-duration: 22s;
}
.login-bg-orb-3 {
  width: 240px;
  height: 240px;
  background: #ddd6fe;
  bottom: -15%;
  left: 20%;
  animation-delay: -14s;
  animation-duration: 24s;
}
@keyframes gradientMove {
  0%, 100% { transform: translate(0, 0) scale(1); background-position: 0% 50%; }
  50% { transform: translate(2%, 2%) scale(1.02); background-position: 100% 50%; }
}
@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -40px) scale(1.05); }
  66% { transform: translate(-20px, 25px) scale(0.98); }
}
.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
}
.login-brand {
  text-align: center;
  margin-bottom: 32px;
}
.login-brand-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #7c3aed;
  color: #fff;
  border-radius: 14px;
  box-shadow: 0 4px 14px rgba(124, 58, 237, 0.4);
}
.login-brand-icon svg {
  width: 28px;
  height: 28px;
}
.login-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #4c1d95;
  margin: 0 0 4px;
  letter-spacing: -0.02em;
}
.login-subtitle {
  font-size: 0.9375rem;
  color: #6b21a8;
  margin: 0;
  font-weight: 500;
}
.login-card {
  border-radius: 16px;
  border: 1px solid rgba(124, 58, 237, 0.12);
  overflow: hidden;
}
.login-card :deep(.el-card__header) {
  padding: 20px 24px 16px;
  border-bottom: 1px solid #f3e8ff;
  font-weight: 600;
  font-size: 1.0625rem;
  color: #4c1d95;
}
.login-card :deep(.el-card__body) {
  padding: 24px;
}
.login-form :deep(.el-form-item) {
  margin-bottom: 22px;
}
.login-form-item :deep(.el-form-item__label) {
  display: block;
  width: 100%;
  text-align: left;
  padding-bottom: 8px;
  color: #5b21b6;
  font-weight: 500;
  font-size: 0.9375rem;
  line-height: 1.4;
}
.login-form-item :deep(.el-form-item__content) {
  display: block;
  width: 100%;
  margin-left: 0 !important;
}
.login-form-item :deep(.el-input) {
  width: 100%;
}
.login-form-item :deep(.el-input__wrapper) {
  width: 100%;
  min-width: 0;
}
.login-submit-item {
  margin-bottom: 0;
  margin-top: 12px;
}
.login-submit-item :deep(.el-form-item__content) {
  margin-left: 0 !important;
}
.login-btn {
  width: 100%;
  height: 44px;
  font-weight: 600;
  font-size: 1rem;
  border-radius: 10px;
  background: #7c3aed;
  border-color: #7c3aed;
  transition: background 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}
.login-btn:hover,
.login-btn:focus {
  background: #6d28d9;
  border-color: #6d28d9;
  box-shadow: 0 4px 14px rgba(124, 58, 237, 0.35);
}
.login-input :deep(.el-input__wrapper) {
  border-radius: 10px;
  transition: box-shadow 0.2s ease, border-color 0.2s ease;
}
.login-input :deep(.el-input__wrapper:hover),
.login-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.2);
}
.input-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #7c3aed;
  opacity: 0.8;
}
.input-icon svg {
  width: 18px;
  height: 18px;
}
.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 0.8125rem;
  color: #7c3aed;
  opacity: 0.85;
}
@media (prefers-reduced-motion: reduce) {
  .login-btn,
  .login-input :deep(.el-input__wrapper) { transition: none; }
  .login-bg-gradient,
  .login-bg-orb { animation: none; }
}
</style>
