<template>
  <div class="notifications-page">
    <h1 class="page-title">消息通知</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0" class="empty">暂无消息</div>
    <ul v-else class="notif-list">
      <li
        v-for="n in list"
        :key="n.id"
        class="card notif-item"
        :class="{ unread: !isRead(n) }"
      >
        <router-link
          v-if="postId(n)"
          :to="{ name: 'PostDetail', params: { id: postId(n) } }"
          class="notif-link"
          @click="markAsReadIfUnread(n)"
        >
          <span class="notif-text">{{ notifText(n) }}</span>
          <time class="notif-time">{{ formatTime(n.created_at) }}</time>
        </router-link>
        <div v-else class="notif-link" @click="markAsReadIfUnread(n)">
          <span class="notif-text">{{ notifText(n) }}</span>
          <time class="notif-time">{{ formatTime(n.created_at) }}</time>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useNotifStore } from '../stores/notif'
import { notifications as notifApi } from '../api'

const list = ref([])
const loading = ref(true)
const POLL_INTERVAL_MS = 3000
let pollTimer = null

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

function isRead(n) {
  if (!n.read_at) return false
  return n.read_at !== '1970-01-01T00:00:00Z' && n.read_at !== '1970-01-01 00:00:00'
}

function notifText(n) {
  if (n.type === 'like') return '有人赞了你的帖子'
  if (n.type === 'comment') return '有人评论了你的帖子'
  return n.type || '新消息'
}

function postId(n) {
  if (n.type === 'like' && n.related_id) return String(n.related_id)
  if (n.type === 'comment' && n.payload && n.payload.post_id) return String(n.payload.post_id)
  return null
}

async function fetchList() {
  try {
    const { data } = await notifApi.list({ limit: 50 })
    list.value = data.list || []
  } catch (e) {
    list.value = []
  } finally {
    loading.value = false
  }
}

const notifStore = useNotifStore()

async function markAsReadIfUnread(n) {
  if (isRead(n)) return
  try {
    await notifApi.markRead(n.id)
    notifStore.decrementUnread()
    const idx = list.value.findIndex((x) => x.id === n.id)
    if (idx !== -1) {
      list.value[idx] = { ...list.value[idx], read_at: new Date().toISOString() }
    }
  } catch (_) {
    // 忽略失败，下次轮询会同步
  }
}

onMounted(() => {
  fetchList()
  notifStore.fetchUnreadCount()
  pollTimer = setInterval(() => {
    fetchList()
    notifStore.fetchUnreadCount()
  }, POLL_INTERVAL_MS)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
.notifications-page {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 20px;
  width: 100%;
  min-width: 0;
}
.notif-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 12px; }
.notif-item { display: block; transition: background var(--transition); padding: 16px; }
.notif-item.unread { background: var(--color-primary-light); border-color: rgba(129, 216, 207, 0.4); }
.notif-link { display: block; text-decoration: none; color: inherit; }
.notif-text { font-weight: 500; color: var(--color-text); font-size: clamp(0.9375rem, 2.5vw, 1rem); }
.notif-time { display: block; margin-top: 4px; font-size: 12px; color: var(--color-text-muted); }
@media (max-width: 768px) {
  .notifications-page { padding: 0 16px; }
  .notif-item { padding: 12px; }
}
@media (max-width: 375px) {
  .notifications-page { padding: 0 12px; }
  .notif-item { padding: 10px; }
  .notif-time { font-size: 11px; }
}
</style>
