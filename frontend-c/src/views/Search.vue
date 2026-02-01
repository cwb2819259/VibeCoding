<template>
  <div class="search-page">
    <h1 class="page-title">搜索 / 话题</h1>
    <div class="search-form card">
      <div class="form-group" style="margin-bottom: 0;">
        <label for="search-keyword">关键词</label>
        <input
          id="search-keyword"
          v-model="keyword"
          type="search"
          placeholder="输入关键词搜索"
          @keyup.enter="search"
          aria-label="搜索关键词"
        />
      </div>
      <button type="button" class="btn btn-primary" @click="search">搜索</button>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="list.length === 0 && searched" class="empty">无结果</div>
    <div v-else-if="list.length" class="feed-list">
      <article
        v-for="p in list"
        :key="p.id"
        class="card post-card"
        @click="goDetail(p.id)"
      >
        <div class="post-meta">
          <span class="author">{{ p.user?.nickname || '用户' + (p.user_id || '').toString().slice(-4) }}</span>
          <span class="time">{{ formatTime(p.created_at) }}</span>
        </div>
        <p v-if="p.content" class="post-content">{{ p.content }}</p>
        <!-- 与帖子列表一致：多图横向滑动，视频可播放 -->
        <div v-if="p.media && p.media.length" class="post-media-wrap" @click.stop>
          <div
            class="post-media-scroll"
            :class="{ 'has-multi': p.media.length > 1 }"
            :data-post-id="p.id"
            @scroll="onMediaScroll"
          >
            <div v-for="(m, i) in p.media" :key="m.id || i" class="post-media-slide">
              <img
                v-if="m.type === 'image'"
                :src="mediaUrl(m.url)"
                :alt="`图片 ${i + 1}`"
                class="post-thumb"
                loading="lazy"
                @click="goDetail(p.id)"
              />
              <video
                v-else
                :src="mediaUrl(m.url)"
                class="post-thumb post-video"
                controls
                playsinline
                muted
                @click.stop
              />
            </div>
          </div>
          <div v-if="p.media.length > 1" class="media-dots">
            <span v-for="(_, i) in p.media" :key="i" class="dot" :class="{ active: currentSlide[p.id] === i }" />
          </div>
        </div>
        <div class="post-actions">
          <span class="action-item"><Icons name="heart" size="sm" /> {{ p.like_count ?? 0 }}</span>
          <span class="action-item"><Icons name="chat" size="sm" /> {{ p.comment_count ?? 0 }}</span>
        </div>
      </article>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { posts as postsApi } from '../api'
import Icons from '../components/Icons.vue'
import { mediaUrl } from '../utils/media'

const router = useRouter()
const keyword = ref('')
const list = ref([])
const loading = ref(false)
const searched = ref(false)
const currentSlide = ref({})

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

function goDetail(id) {
  router.push({ name: 'PostDetail', params: { id } })
}

function onMediaScroll(e) {
  const el = e.target
  const postId = el.dataset?.postId
  if (!postId || !el.scrollWidth) return
  const slideIndex = Math.round(el.scrollLeft / el.clientWidth)
  currentSlide.value = { ...currentSlide.value, [postId]: Math.min(slideIndex, el.children.length - 1) }
}

async function search() {
  searched.value = true
  loading.value = true
  try {
    const { data } = await postsApi.list({ keyword: keyword.value, limit: 20 })
    const items = data.list || []
    list.value = items
    currentSlide.value = Object.fromEntries(items.map((p) => [p.id, 0]))
  } catch (e) {
    list.value = []
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.search-page {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 20px;
  width: 100%;
  min-width: 0;
}
.search-form { display: flex; gap: 12px; align-items: flex-end; flex-wrap: wrap; }
.search-form .form-group { flex: 1; min-width: 200px; }
.search-form .btn { flex-shrink: 0; }
.feed-list { display: flex; flex-direction: column; gap: 16px; margin-top: var(--section-gap); }
.post-card { cursor: pointer; transition: transform var(--transition); padding: 16px; }
.post-card:hover { transform: translateY(-1px); }
.post-meta { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; color: var(--color-text-secondary); font-size: clamp(13px, 2vw, 14px); }
.author { font-weight: 500; color: var(--color-text); }
.post-content { margin: 0 0 12px; white-space: pre-wrap; word-break: break-word; color: var(--color-text); line-height: 1.6; font-size: clamp(0.9375rem, 2.5vw, 1rem); }

.post-media-wrap { position: relative; margin-bottom: 12px; border-radius: var(--radius-sm); overflow: hidden; background: #000; }
.post-media-scroll {
  display: flex;
  overflow-x: auto;
  overflow-y: hidden;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
}
.post-media-scroll::-webkit-scrollbar { display: none; }
.post-media-scroll.has-multi { cursor: grab; }
.post-media-scroll.has-multi:active { cursor: grabbing; }
.post-media-slide { flex: 0 0 100%; scroll-snap-align: start; scroll-snap-stop: always; }
.post-thumb { display: block; width: 100%; max-height: min(380px, 60vh); object-fit: contain; background: #111; }
.post-video { max-height: min(400px, 65vh); }
.media-dots {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 6px;
  z-index: 1;
}
.dot { width: 6px; height: 6px; border-radius: 50%; background: rgba(255,255,255,0.5); transition: background var(--transition); }
.dot.active { background: #fff; }

.post-actions { display: flex; gap: 24px; color: var(--color-text-muted); font-size: clamp(13px, 2vw, 14px); }
.action-item { display: inline-flex; align-items: center; gap: 6px; }

@media (max-width: 768px) {
  .search-page { padding: 0 16px; }
  .search-form { flex-direction: column; align-items: stretch; }
  .search-form .form-group { min-width: 0; }
  .search-form .btn { width: 100%; }
  .post-card { padding: 12px; }
  .post-thumb { max-height: min(280px, 50vh); }
  .post-video { max-height: min(300px, 55vh); }
}
@media (max-width: 375px) {
  .search-page { padding: 0 12px; }
  .post-card { padding: 10px; }
  .post-thumb { max-height: min(220px, 45vh); }
  .post-video { max-height: min(240px, 50vh); }
}
</style>
