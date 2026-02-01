<template>
  <div class="page-content">
    <div class="detail-header">
      <el-button text class="back-btn" @click="goBack">
        <span class="back-arrow">←</span> 返回列表
      </el-button>
      <h1 class="page-title">帖子详情</h1>
    </div>

    <el-card v-loading="loading" class="detail-card" shadow="hover">
      <template v-if="error">
        <p class="detail-error">{{ error }}</p>
      </template>
      <template v-else-if="post">
        <div class="detail-meta">
          <span class="meta-item"><strong>ID</strong> {{ post.id }}</span>
          <span class="meta-item"><strong>用户ID</strong> {{ post.user_id }}</span>
          <span v-if="post.user" class="meta-item"><strong>昵称</strong> {{ post.user.nickname || '-' }}</span>
          <el-tag :type="statusTagType(post.status)" size="small" class="status-tag">{{ statusLabel(post.status) }}</el-tag>
          <span class="meta-item"><strong>发布时间</strong> {{ formatTime(post.created_at) }}</span>
        </div>

        <div class="detail-content-section">
          <h3 class="section-label">正文</h3>
          <div class="detail-content">{{ post.content || '（无文字）' }}</div>
        </div>

        <div v-if="post.media && post.media.length" class="detail-media-section">
          <h3 class="section-label">图片 / 视频</h3>
          <div class="media-grid">
            <div v-for="(m, i) in post.media" :key="m.id || i" v-show="true" class="media-item" :class="isVideo(m) ? 'media-video' : 'media-image'">
              <template v-if="isVideo(m)">
                <video :src="mediaUrl(m.url)" controls playsinline class="media-video-el" />
                <span class="media-type-tag">视频</span>
              </template>
              <img v-else :src="mediaUrl(m.url)" :alt="`图片 ${i + 1}`" class="media-img" loading="lazy" />
            </div>
          </div>
        </div>

        <div class="detail-actions">
          <el-button v-if="post.status !== 'normal'" type="success" size="default" @click="setStatus(post.id, 'normal')">恢复</el-button>
          <el-button v-if="post.status === 'normal'" size="default" @click="setStatus(post.id, 'hidden')">隐藏</el-button>
          <el-button v-if="post.status === 'normal'" size="default" @click="setStatus(post.id, 'flagged')">标记</el-button>
          <el-button type="danger" size="default" @click="del(post.id)">删除</el-button>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { posts as postsApi, mediaUrl } from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const post = ref(null)
const loading = ref(false)
const error = ref('')

const id = computed(() => route.params.id)

function statusLabel(s) {
  const map = { normal: '正常', hidden: '隐藏', flagged: '标记' }
  return map[s] || s
}
function statusTagType(s) {
  const map = { normal: 'success', hidden: 'warning', flagged: 'danger' }
  return map[s] || 'info'
}
function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}
function isVideo(m) {
  if (!m || !m.url) return false
  const u = (m.url || '').toLowerCase()
  return u.includes('video') || u.endsWith('.mp4') || u.endsWith('.webm')
}

function goBack() {
  router.push({ name: 'Content' })
}

async function load() {
  if (!id.value) return
  loading.value = true
  error.value = ''
  try {
    const { data } = await postsApi.get(id.value)
    post.value = data
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
    post.value = null
  } finally {
    loading.value = false
  }
}

async function setStatus(postId, status) {
  try {
    await postsApi.updateStatus(postId, status)
    ElMessage.success('已更新')
    post.value = { ...post.value, status }
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '失败')
  }
}

async function del(postId) {
  if (!confirm('确定删除该帖子？')) return
  try {
    await postsApi.delete(postId)
    ElMessage.success('已删除')
    router.push({ name: 'Content' })
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '失败')
  }
}

onMounted(load)
</script>

<style scoped>
.page-content { max-width: 900px; }
.detail-header {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}
.back-btn {
  color: var(--admin-text-secondary);
  font-size: 0.9375rem;
  padding: 4px 0;
}
.back-btn:hover { color: var(--admin-primary); }
.back-arrow { margin-right: 4px; }
.page-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--admin-text);
  letter-spacing: -0.02em;
}
.detail-card {
  border-radius: var(--admin-card-radius);
  border: 1px solid var(--admin-border);
  box-shadow: var(--admin-card-shadow);
  overflow: hidden;
}
.detail-card :deep(.el-card__body) { padding: 24px; }
.detail-error {
  color: var(--el-color-danger);
  margin: 0;
}
.detail-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px 24px;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--admin-border);
}
.meta-item {
  font-size: 0.9375rem;
  color: var(--admin-text-secondary);
}
.meta-item strong { color: var(--admin-text); margin-right: 6px; }
.status-tag { margin-right: 0; }
.section-label {
  margin: 0 0 10px;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--admin-text-secondary);
}
.detail-content-section { margin-bottom: 24px; }
.detail-content {
  font-size: 1rem;
  line-height: 1.6;
  color: var(--admin-text);
  white-space: pre-wrap;
  word-break: break-word;
  padding: 16px;
  background: var(--admin-bg);
  border-radius: var(--admin-card-radius);
  border: 1px solid var(--admin-border);
}
.detail-media-section { margin-bottom: 24px; }
.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}
.media-item {
  position: relative;
  border-radius: var(--admin-card-radius);
  overflow: hidden;
  border: 1px solid var(--admin-border);
  background: #111;
}
.media-img {
  display: block;
  width: 100%;
  height: auto;
  max-height: 320px;
  object-fit: contain;
}
.media-video { min-height: 120px; }
.media-video-el {
  display: block;
  width: 100%;
  max-height: 360px;
  object-fit: contain;
}
.media-type-tag {
  position: absolute;
  top: 8px;
  left: 8px;
  font-size: 11px;
  padding: 2px 8px;
  background: rgba(0,0,0,0.6);
  color: #fff;
  border-radius: 4px;
}
.detail-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--admin-border);
}
</style>
