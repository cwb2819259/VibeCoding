<template>
  <div v-if="loading" class="detail-loading">
    <span class="detail-loading-spinner" aria-hidden="true" />
    <span class="detail-loading-text">加载中...</span>
  </div>
  <div v-else-if="post" class="detail-page">
    <article class="detail-card">
      <!-- 作者区 -->
      <header class="detail-header">
        <div class="author-avatar" aria-hidden="true">
          <span class="author-initial">{{ authorInitial }}</span>
        </div>
        <div class="author-info">
          <span class="author-name">{{ post.user?.nickname || '用户' }}</span>
          <time class="post-time" :datetime="post.created_at">{{ formatTime(post.created_at) }}</time>
        </div>
      </header>

      <!-- 正文 -->
      <div v-if="post.content" class="detail-content">
        {{ post.content }}
      </div>

      <!-- 媒体：多图横向滑动 / 单图或视频 -->
      <div v-if="post.media && post.media.length" class="detail-media">
        <div
          v-if="post.media.length > 1"
          class="media-scroll-wrap"
          @scroll="onMediaScroll"
        >
          <div class="media-scroll">
            <div v-for="(m, i) in post.media" :key="i" class="media-slide">
              <img
                v-if="m.type === 'image'"
                :src="mediaUrl(m.url)"
                :alt="`图片 ${i + 1}`"
                class="media-item media-img"
                loading="lazy"
              />
              <video
                v-else
                :src="mediaUrl(m.url)"
                class="media-item media-video-el"
                controls
                playsinline
              />
            </div>
          </div>
        </div>
        <div v-else class="media-single">
          <img
            v-if="post.media[0].type === 'image'"
            :src="mediaUrl(post.media[0].url)"
            :alt="'图片 1'"
            class="media-item media-img"
            loading="lazy"
          />
          <video
            v-else
            :src="mediaUrl(post.media[0].url)"
            class="media-item media-video-el"
            controls
            playsinline
          />
        </div>
        <div v-if="post.media.length > 1" class="media-dots">
          <span
            v-for="(_, i) in post.media"
            :key="i"
            class="dot"
            :class="{ active: currentSlide === i }"
            aria-hidden="true"
          />
        </div>
      </div>

      <!-- 操作栏 -->
      <div class="detail-actions">
        <button
          type="button"
          class="action-btn action-like"
          :class="{ liked: liked }"
          :aria-pressed="liked"
          @click="toggleLike"
        >
          <Icons :name="liked ? 'heartSolid' : 'heart'" size="md" />
          <span>{{ likeCount }}</span>
        </button>
        <span class="action-stat" aria-label="评论数">
          <Icons name="chat" size="md" />
          <span>{{ commentCount }}</span>
        </span>
      </div>
    </article>

    <!-- 评论区 -->
    <section class="comments-card" aria-labelledby="comments-heading">
      <h2 id="comments-heading" class="comments-heading">评论 <span class="comments-count">{{ comments.length }}</span></h2>
      <div v-if="authStore.token" class="comment-form">
        <textarea
          v-model="commentText"
          placeholder="写一条评论..."
          rows="3"
          class="comment-textarea"
          aria-label="评论内容"
          @focus="commentFocused = true"
          @blur="commentFocused = false"
        />
        <button
          type="button"
          class="btn btn-primary btn-send"
          :disabled="!commentText.trim()"
          @click="submitComment"
        >
          发送
        </button>
      </div>
      <ul class="comment-list" aria-label="评论列表">
        <li v-for="c in comments" :key="c.id" class="comment-item">
          <div class="comment-avatar" aria-hidden="true">
            <span class="comment-initial">{{ commentInitial(c) }}</span>
          </div>
          <div class="comment-body">
            <strong class="comment-author">{{ c.user?.nickname || '用户' }}</strong>
            <span class="comment-time">{{ formatTime(c.created_at) }}</span>
            <p class="comment-content">{{ c.content }}</p>
          </div>
        </li>
      </ul>
      <p v-if="comments.length === 0 && !loading" class="comments-empty">暂无评论，来抢沙发吧</p>
    </section>
  </div>
  <div v-else class="detail-empty">
    <p>未找到帖子</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { posts as postsApi, comments as commentsApi, likes as likesApi } from '../api'
import Icons from '../components/Icons.vue'
import { mediaUrl } from '../utils/media'

const route = useRoute()
const authStore = useAuthStore()
const post = ref(null)
const likeCount = ref(0)
const commentCount = ref(0)
const liked = ref(false)
const comments = ref([])
const commentText = ref('')
const commentFocused = ref(false)
const loading = ref(true)
const currentSlide = ref(0)

const id = computed(() => route.params.id)
const authorInitial = computed(() => {
  const name = post.value?.user?.nickname || '用'
  return name.charAt(0).toUpperCase()
})

function commentInitial(c) {
  const name = c.user?.nickname || '用'
  return name.charAt(0).toUpperCase()
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

function onMediaScroll(e) {
  const el = e.target
  if (!el.scrollWidth) return
  const index = Math.round(el.scrollLeft / el.clientWidth)
  const inner = el.firstElementChild
  const count = inner ? inner.children.length : 0
  currentSlide.value = Math.min(Math.max(0, index), count - 1)
}

async function load() {
  try {
    const { data } = await postsApi.get(id.value)
    post.value = data.post
    likeCount.value = data.like_count ?? 0
    commentCount.value = data.comment_count ?? 0
    liked.value = !!data.liked
    const { data: cr } = await commentsApi.list(id.value, { limit: 50 })
    comments.value = cr.list || []
  } catch (e) {
    post.value = null
  } finally {
    loading.value = false
  }
}

async function toggleLike() {
  if (!authStore.token) return
  try {
    const res = liked.value
      ? await likesApi.unlike(id.value)
      : await likesApi.like(id.value)
    liked.value = !!res?.data?.liked
    if (typeof res?.data?.like_count === 'number') likeCount.value = res.data.like_count
  } catch (e) {
    console.error(e)
  }
}

async function submitComment() {
  if (!commentText.value.trim()) return
  try {
    await commentsApi.create(id.value, { content: commentText.value.trim() })
    commentText.value = ''
    const { data } = await commentsApi.list(id.value, { limit: 50 })
    comments.value = data.list || []
    commentCount.value = comments.length
  } catch (e) {
    console.error(e)
  }
}

onMounted(load)
</script>

<style scoped>
.detail-page {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 20px;
  width: 100%;
  min-width: 0;
}
.detail-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  min-height: 200px;
  color: var(--color-text-muted);
}
.detail-loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.detail-loading-text { font-size: 14px; }
.detail-empty {
  text-align: center;
  padding: 48px 20px;
  color: var(--color-text-muted);
}
.detail-empty p { margin: 0; font-size: 1rem; }

/* 帖子卡片 */
.detail-card {
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  border: 1px solid var(--color-border);
  padding: 24px;
  margin-bottom: var(--section-gap);
  transition: box-shadow var(--transition), border-color var(--transition);
}
.detail-card:hover {
  box-shadow: var(--shadow-hover);
  border-color: rgba(129, 216, 207, 0.3);
}

/* 作者区 */
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}
.author-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--color-primary-light);
  color: var(--color-primary-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.author-initial {
  font-family: var(--font-heading);
  font-size: 1.125rem;
  font-weight: 700;
}
.author-info { min-width: 0; }
.author-name {
  display: block;
  font-weight: 600;
  color: var(--color-text);
  font-size: 1rem;
}
.post-time {
  display: block;
  font-size: 13px;
  color: var(--color-text-muted);
  margin-top: 2px;
}

/* 正文 */
.detail-content {
  font-size: 1.0625rem;
  line-height: 1.7;
  color: var(--color-text);
  white-space: pre-wrap;
  word-break: break-word;
  margin-bottom: 24px;
}

/* 媒体 */
.detail-media {
  position: relative;
  margin-bottom: 24px;
  border-radius: var(--radius);
  overflow: hidden;
  background: #0f0f0f;
}
.media-scroll-wrap {
  overflow-x: auto;
  overflow-y: hidden;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
}
.media-scroll-wrap::-webkit-scrollbar { display: none; }
.media-scroll {
  display: flex;
  scroll-snap-type: x mandatory;
}
.media-slide {
  flex: 0 0 100%;
  scroll-snap-align: start;
  scroll-snap-stop: always;
}
.media-single { display: block; }
.media-item {
  display: block;
  width: 100%;
  max-height: min(420px, 75vh);
  object-fit: contain;
}
.media-video-el { max-height: min(400px, 70vh); }
.media-dots {
  position: absolute;
  bottom: 12px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  z-index: 1;
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.4);
  transition: background var(--transition);
}
.dot.active { background: #fff; }

/* 操作栏 */
.detail-actions {
  display: flex;
  align-items: center;
  gap: 24px;
  padding-top: 20px;
  border-top: 1px solid var(--color-border);
}
.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-text-muted);
  font-size: 15px;
  cursor: pointer;
  transition: color var(--transition), background var(--transition);
}
.action-btn:hover { color: var(--color-primary); background: var(--color-primary-light); }
.action-btn:focus-visible { outline: 2px solid var(--color-primary); outline-offset: 2px; }
.action-btn.liked { color: #e11d48; }
.action-btn.liked:hover { background: rgba(225, 29, 72, 0.1); }
.action-stat {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-muted);
  font-size: 15px;
}

/* 评论区 */
.comments-card {
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  border: 1px solid var(--color-border);
  padding: 24px;
  margin-bottom: var(--section-gap);
}
.comments-heading {
  font-family: var(--font-heading);
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 20px;
  color: var(--color-text);
}
.comments-count {
  font-weight: 500;
  color: var(--color-text-muted);
  font-size: 0.9375rem;
}
.comment-form {
  margin-bottom: 24px;
}
.comment-textarea {
  width: 100%;
  min-height: 88px;
  padding: 14px 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: 15px;
  line-height: 1.5;
  resize: vertical;
  box-sizing: border-box;
  transition: border-color var(--transition);
  margin-bottom: 12px;
}
.comment-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
}
.comment-textarea::placeholder { color: var(--color-text-muted); }
.btn-send { min-width: 100px; }
.comment-list { list-style: none; padding: 0; margin: 0; }
.comment-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid var(--color-border);
}
.comment-item:last-child { border-bottom: none; }
.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--color-bg);
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 600;
}
.comment-body { min-width: 0; flex: 1; }
.comment-author {
  font-size: 14px;
  color: var(--color-text);
  margin-right: 8px;
}
.comment-time {
  font-size: 12px;
  color: var(--color-text-muted);
}
.comment-content {
  margin: 8px 0 0;
  font-size: 14px;
  line-height: 1.5;
  color: var(--color-text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
}
.comments-empty {
  margin: 0;
  padding: 24px 0;
  text-align: center;
  color: var(--color-text-muted);
  font-size: 14px;
}

@media (max-width: 768px) {
  .detail-page { padding: 0 16px; }
  .detail-card { padding: 20px; }
  .detail-header { margin-bottom: 16px; }
  .author-avatar { width: 40px; height: 40px; }
  .author-initial { font-size: 1rem; }
  .detail-content { font-size: 1rem; margin-bottom: 20px; }
  .media-item { max-height: min(320px, 65vh); }
  .media-video-el { max-height: min(300px, 60vh); }
  .detail-actions { padding-top: 16px; gap: 20px; }
  .comments-card { padding: 20px; }
  .comments-heading { margin-bottom: 16px; font-size: 1.125rem; }
}
@media (max-width: 375px) {
  .detail-page { padding: 0 12px; }
  .detail-card { padding: 16px; }
  .author-avatar { width: 36px; height: 36px; }
  .comment-avatar { width: 32px; height: 32px; font-size: 12px; }
  .comment-item { padding: 12px 0; }
}
@media (prefers-reduced-motion: reduce) {
  .detail-loading-spinner { animation: none; border-top-color: var(--color-primary); opacity: 0.8; }
}
</style>
