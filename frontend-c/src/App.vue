<template>
  <div id="app">
    <!-- 桌面端：单行导航 -->
    <header class="nav nav-desktop" v-if="!hideNav">
      <router-link to="/" class="logo">VibeCoding社区</router-link>
      <nav class="nav-links" aria-label="主导航">
        <router-link to="/" class="nav-link" exact-active-class="active"> <Icons name="home" size="sm" /> 首页 </router-link>
        <router-link to="/posts" class="nav-link" exact-active-class="active"> <Icons name="feed" size="sm" /> 帖子 </router-link>
        <router-link to="/search" class="nav-link" exact-active-class="active"> <Icons name="search" size="sm" /> 搜索 </router-link>
        <template v-if="token">
          <router-link to="/publish" class="nav-link" exact-active-class="active"> <Icons name="pen" size="sm" /> 发帖 </router-link>
          <router-link to="/user" class="nav-link" exact-active-class="active"> <Icons name="user" size="sm" /> 我的 </router-link>
          <router-link to="/notifications" class="nav-link nav-link-notif" exact-active-class="active">
            <Icons name="bell" size="sm" /> 消息
            <span v-if="notifStore.unreadCount > 0" class="notif-badge">{{ notifStore.unreadCount > 99 ? '99+' : notifStore.unreadCount }}</span>
          </router-link>
          <a href="#" class="nav-link" @click.prevent="logout" aria-label="退出登录"> <Icons name="logout" size="sm" /> 退出 </a>
        </template>
        <router-link v-else to="/login" class="nav-link" exact-active-class="active">登录</router-link>
      </nav>
    </header>

    <!-- 手机端：顶部栏 + 底部 Tab -->
    <template v-if="!hideNav">
      <header class="nav-mobile-top">
        <router-link to="/" class="nav-mobile-logo">VibeCoding社区</router-link>
        <router-link to="/search" class="nav-mobile-icon" aria-label="搜索">
          <Icons name="search" size="sm" />
        </router-link>
      </header>
      <nav class="nav-mobile-bottom" aria-label="主导航">
        <router-link to="/" class="tab-item" exact-active-class="active">
          <Icons name="home" size="sm" /> <span class="tab-label">首页</span>
        </router-link>
        <router-link to="/posts" class="tab-item" exact-active-class="active">
          <Icons name="feed" size="sm" /> <span class="tab-label">帖子</span>
        </router-link>
        <router-link v-if="token" to="/publish" class="tab-item tab-item-publish" exact-active-class="active">
          <Icons name="pen" size="sm" /> <span class="tab-label">发帖</span>
        </router-link>
        <template v-if="token">
          <router-link to="/notifications" class="tab-item tab-item-notif" exact-active-class="active">
            <Icons name="bell" size="sm" />
            <span class="tab-label">消息</span>
            <span v-if="notifStore.unreadCount > 0" class="tab-badge">{{ notifStore.unreadCount > 99 ? '99+' : notifStore.unreadCount }}</span>
          </router-link>
          <router-link to="/user" class="tab-item" exact-active-class="active">
            <Icons name="user" size="sm" /> <span class="tab-label">我的</span>
          </router-link>
        </template>
        <template v-else>
          <router-link to="/search" class="tab-item" exact-active-class="active">
            <Icons name="search" size="sm" /> <span class="tab-label">搜索</span>
          </router-link>
          <router-link to="/login" class="tab-item" exact-active-class="active">
            <Icons name="user" size="sm" /> <span class="tab-label">登录</span>
          </router-link>
        </template>
      </nav>
    </template>

    <main :class="['main-content', route.name === 'Landing' ? 'main-content--full' : 'container']">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useNotifStore } from './stores/notif'
import Icons from './components/Icons.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const notifStore = useNotifStore()

const token = computed(() => auth.token)
const hideNav = computed(() => route.meta.hideNav)

watch(token, (t) => {
  if (t) notifStore.startPolling()
  else notifStore.stopPolling()
}, { immediate: true })

function logout() {
  auth.logout()
  notifStore.stopPolling()
  router.push('/')
}
</script>

<style scoped>
.nav-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  min-height: 44px;
  padding: 0 6px;
  border-radius: 8px;
  transition: background var(--transition);
}
.nav-link:hover { background: var(--color-primary-light); }
.nav-link-notif { position: relative; }
.notif-badge {
  position: absolute;
  top: -4px;
  right: -6px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  font-size: 11px;
  line-height: 16px;
  text-align: center;
  color: #fff;
  background: #e11d48;
  border-radius: 8px;
}
.main-content {
  padding-top: 24px;
  padding-bottom: 48px;
  padding-bottom: calc(48px + env(safe-area-inset-bottom));
  width: 100%;
  min-width: 0;
}
/* 首页全屏：PC 画幅下去掉窄容器，由 Landing 自控宽度 */
.main-content--full {
  max-width: none;
  padding-left: 0;
  padding-right: 0;
}

/* 手机端：顶部栏 + 底部 Tab；PC 画幅下不展示 */
.nav-mobile-top,
.nav-mobile-bottom { display: none; }
@media (min-width: 769px) {
  .nav-mobile-top,
  .nav-mobile-bottom { display: none !important; }
}
.nav-mobile-top {
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  padding-left: max(16px, env(safe-area-inset-left));
  padding-right: max(16px, env(safe-area-inset-right));
  background: var(--color-card);
  border-bottom: 1px solid var(--color-border);
  min-height: 48px;
}
.nav-mobile-logo {
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 1rem;
  color: var(--color-text);
  text-decoration: none;
  transition: color var(--transition);
}
.nav-mobile-logo:hover { color: var(--color-primary); }
.nav-mobile-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  transition: color var(--transition), background var(--transition);
}
.nav-mobile-icon:hover { color: var(--color-primary); background: var(--color-primary-light); }
.nav-mobile-bottom {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  align-items: stretch;
  justify-content: space-around;
  padding: 8px 0;
  padding-bottom: max(8px, env(safe-area-inset-bottom));
  background: var(--color-card);
  border-top: 1px solid var(--color-border);
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.06);
}
.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 6px 4px;
  color: var(--color-text-muted);
  text-decoration: none;
  font-size: 10px;
  transition: color var(--transition);
  min-width: 0;
  cursor: pointer;
}
.tab-item .icon-wrap { flex-shrink: 0; }
.tab-item:hover { color: var(--color-primary); }
.tab-item.active { color: var(--color-primary); font-weight: 600; }
.tab-label { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 100%; }
.tab-item-publish { color: var(--color-primary); }
.tab-item-publish .icon-wrap { width: 22px; height: 22px; }
.tab-item-publish.active { color: var(--color-primary-dark); }
.tab-item-notif { position: relative; }
.tab-badge {
  position: absolute;
  top: 0;
  left: calc(50% + 10px);
  min-width: 14px;
  height: 14px;
  padding: 0 4px;
  font-size: 10px;
  line-height: 14px;
  text-align: center;
  color: #fff;
  background: #e11d48;
  border-radius: 7px;
}
@media (max-width: 768px) {
  .nav-mobile-top,
  .nav-mobile-bottom { display: flex; }
  .main-content { padding-top: 12px; padding-bottom: 72px; padding-bottom: calc(72px + env(safe-area-inset-bottom)); }
}
</style>
