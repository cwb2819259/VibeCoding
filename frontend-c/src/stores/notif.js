import { defineStore } from 'pinia'
import { ref } from 'vue'
import { notifications as notifApi } from '../api'

const POLL_INTERVAL_MS = 3000

export const useNotifStore = defineStore('notif', () => {
  const unreadCount = ref(0)
  let pollTimer = null

  async function fetchUnreadCount() {
    const token = localStorage.getItem('c_token')
    if (!token) {
      unreadCount.value = 0
      return
    }
    try {
      const { data } = await notifApi.unreadCount()
      unreadCount.value = data?.count ?? 0
    } catch {
      unreadCount.value = 0
    }
  }

  function startPolling() {
    if (pollTimer) return
    fetchUnreadCount()
    pollTimer = setInterval(fetchUnreadCount, POLL_INTERVAL_MS)
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    unreadCount.value = 0
  }

  /** 单条标已读后本地未读数 -1，与后端 markRead 配合使用 */
  function decrementUnread() {
    if (unreadCount.value > 0) unreadCount.value -= 1
  }

  return { unreadCount, fetchUnreadCount, startPolling, stopPolling, decrementUnread }
})
