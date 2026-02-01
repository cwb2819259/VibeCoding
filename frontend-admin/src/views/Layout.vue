<template>
  <el-container class="admin-layout">
    <el-aside width="220px" class="admin-aside">
      <div class="admin-brand">
        <div class="admin-brand-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
          </svg>
        </div>
        <span class="admin-brand-text">VibeCoding</span>
        <span class="admin-brand-sub">后台管理</span>
      </div>
      <el-menu
        router
        :default-active="$route.path"
        class="admin-menu"
        background-color="transparent"
        text-color="#c4b5fd"
        active-text-color="#7c3aed"
      >
        <el-menu-item index="/" class="admin-menu-item">
          <span class="menu-icon-wrap" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
          </span>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/content" class="admin-menu-item">
          <span class="menu-icon-wrap" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
          </span>
          <span>内容管理</span>
        </el-menu-item>
        <el-menu-item index="/users" class="admin-menu-item">
          <span class="menu-icon-wrap" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
          </span>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
      <div class="admin-footer">
        <el-button text class="admin-logout" @click="logout">
          <span class="menu-icon-wrap" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          </span>
          退出登录
        </el-button>
      </div>
    </el-aside>
    <el-main class="admin-main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-layout { min-height: 100vh; }
.admin-aside {
  background: linear-gradient(180deg, #4c1d95 0%, #5b21b6 50%, #4c1d95 100%);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 20px rgba(76, 29, 149, 0.15);
}
.admin-brand {
  padding: 24px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
.admin-brand-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  border-radius: 12px;
}
.admin-brand-icon svg { width: 22px; height: 22px; }
.admin-brand-text {
  font-size: 1.125rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.02em;
}
.admin-brand-sub {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.75);
  font-weight: 500;
}
.admin-menu {
  flex: 1;
  padding: 16px 12px 0;
  border-right: none;
}
.admin-menu-item {
  margin-bottom: 4px;
  border-radius: 10px;
  height: 44px;
  line-height: 44px;
  transition: background var(--admin-transition), color var(--admin-transition);
}
.admin-menu-item:hover {
  background: rgba(255, 255, 255, 0.08) !important;
  color: #e9d5ff !important;
}
.admin-menu-item.is-active {
  background: rgba(255, 255, 255, 0.95) !important;
  color: var(--admin-primary) !important;
  font-weight: 600;
}
.menu-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  vertical-align: middle;
}
.menu-icon-wrap svg { width: 18px; height: 18px; }
.admin-footer {
  padding: 16px 12px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
.admin-logout {
  width: 100%;
  justify-content: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.875rem;
  transition: color var(--admin-transition);
}
.admin-logout:hover {
  color: #fff;
}
.admin-main {
  background: var(--admin-bg-main);
  padding: 24px;
  min-height: 100vh;
}
@media (prefers-reduced-motion: reduce) {
  .admin-menu-item,
  .admin-logout { transition: none; }
}
</style>
