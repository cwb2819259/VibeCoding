<template>
  <div class="publish-page">
    <h1 class="page-title">发帖</h1>
    <div class="card">
      <form @submit.prevent="submit" class="publish-form">
        <div class="form-group">
          <label for="publish-content">内容</label>
          <textarea
            id="publish-content"
            v-model="content"
            rows="4"
            placeholder="写点什么..."
            class="publish-textarea"
          ></textarea>
        </div>
        <div class="form-group">
          <label class="media-label">图片 / 视频</label>
          <p class="media-hint">最多 9 张图片或 1 个视频</p>

          <!-- 上传区：拖拽或点击 -->
          <div
            v-show="canAddMore"
            class="upload-zone"
            :class="{ 'upload-zone--dragging': isDragging, 'upload-zone--disabled': uploading }"
            role="button"
            tabindex="0"
            aria-label="选择图片或视频"
            @click="triggerFileInput"
            @keydown.enter.space.prevent="triggerFileInput"
            @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="onDrop"
          >
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*,video/*"
              multiple
              class="upload-input"
              aria-hidden="true"
              @change="onFileChange"
            />
            <span class="upload-zone-icon">
              <Icons name="upload" size="lg" />
            </span>
            <span class="upload-zone-text">拖拽到此处，或点击选择</span>
            <span class="upload-zone-sub">支持图片、视频</span>
          </div>

          <!-- 预览列表 -->
          <div v-if="mediaItems.length" class="preview-list">
            <div
              v-for="(item, i) in mediaItems"
              :key="i"
              class="preview-item"
            >
              <div class="preview-inner">
                <img
                  v-if="!isVideo(item) && item.previewUrl"
                  :src="item.previewUrl"
                  :alt="`图片 ${i + 1}`"
                  class="preview-img"
                />
                <div v-else-if="isVideo(item)" class="preview-video">
                  <Icons name="video" size="lg" />
                  <span class="preview-video-label">视频</span>
                </div>
                <div v-else class="preview-loading">
                  <span class="preview-spinner" aria-hidden="true" />
                  <span class="preview-loading-text">上传中</span>
                </div>
              </div>
              <div v-if="item.uploading" class="preview-overlay" aria-hidden="true">
                <span class="preview-spinner" />
              </div>
              <button
                type="button"
                class="preview-remove"
                :disabled="item.uploading"
                :aria-label="`移除第 ${i + 1} 项`"
                @click.stop="removeMedia(i)"
              >
                <Icons name="close" size="sm" />
              </button>
            </div>
          </div>
        </div>
        <button type="submit" class="btn btn-primary btn-submit" :disabled="loading">
          {{ loading ? '发布中...' : '发布' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { posts as postsApi, upload as uploadApi } from '../api'
import Icons from '../components/Icons.vue'

const router = useRouter()
const content = ref('')
const fileInputRef = ref(null)
const mediaItems = ref([]) // { serverUrl?, previewUrl, uploading? }
const loading = ref(false)
const isDragging = ref(false)
const uploading = ref(false)

const hasVideo = computed(() => mediaItems.value.some((item) => isVideo(item)))
const canAddMore = computed(() => {
  if (hasVideo.value) return false
  return mediaItems.value.length < 9
})

function isVideo(item) {
  const url = item.serverUrl || ''
  return url.includes('video') || url.endsWith('.mp4') || url.endsWith('.webm')
}

function triggerFileInput() {
  if (!canAddMore.value || uploading.value) return
  fileInputRef.value?.click()
}

function removeMedia(i) {
  const item = mediaItems.value[i]
  if (item.uploading) return
  if (item.previewUrl && item.previewUrl.startsWith('blob:')) {
    URL.revokeObjectURL(item.previewUrl)
  }
  mediaItems.value.splice(i, 1)
}

function onDrop(e) {
  isDragging.value = false
  const files = Array.from(e.dataTransfer?.files || []).filter(
    (f) => f.type.startsWith('image/') || f.type.startsWith('video/')
  )
  if (!files.length) return
  processFiles(files)
}

function onFileChange(e) {
  const files = Array.from(e.target.files || [])
  e.target.value = ''
  processFiles(files)
}

async function processFiles(files) {
  if (!canAddMore.value || uploading.value) return
  const remain = hasVideo.value ? 0 : 9 - mediaItems.value.length
  if (remain <= 0) return

  uploading.value = true
  const toAdd = files.slice(0, remain)
  const hasVideoFile = toAdd.some((f) => f.type.startsWith('video/'))
  const finalList = hasVideoFile ? toAdd.slice(0, 1) : toAdd

  for (const file of finalList) {
    const previewUrl = file.type.startsWith('image/') ? URL.createObjectURL(file) : ''
    mediaItems.value.push({ previewUrl, uploading: true })
    try {
      const { data } = await uploadApi(file)
      const last = mediaItems.value[mediaItems.value.length - 1]
      last.serverUrl = data.url
      last.uploading = false
    } catch (err) {
      mediaItems.value.pop()
      if (previewUrl) URL.revokeObjectURL(previewUrl)
      console.error(err)
    }
  }
  uploading.value = false
}

async function submit() {
  const done = mediaItems.value.filter((i) => i.serverUrl && !i.uploading)
  if (!content.value.trim() && done.length === 0) {
    alert('请填写内容或上传图片/视频')
    return
  }
  loading.value = true
  try {
    const serverUrls = done.map((i) => i.serverUrl)
    const type = serverUrls.some((url) => isVideo({ serverUrl: url })) ? 'video' : 'image'
    await postsApi.create({
      content: content.value.trim(),
      type,
      media_urls: serverUrls,
      topic_names: [],
    })
    mediaItems.value.forEach((i) => {
      if (i.previewUrl && i.previewUrl.startsWith('blob:')) URL.revokeObjectURL(i.previewUrl)
    })
    router.push('/posts')
  } catch (e) {
    alert(e.response?.data?.error || '发布失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.publish-page {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 20px;
  width: 100%;
  min-width: 0;
}
.publish-form { margin: 0; }
.media-label { display: block; margin-bottom: 4px; color: var(--color-text-secondary); font-size: 14px; }
.media-hint { margin: 0 0 12px; color: var(--color-text-muted); font-size: 13px; }

.upload-input { position: absolute; width: 0; height: 0; opacity: 0; pointer-events: none; }
.upload-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 140px;
  padding: 24px 20px;
  border: 2px dashed var(--color-border);
  border-radius: var(--radius);
  background: var(--color-bg);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: border-color var(--transition), background var(--transition), color var(--transition);
}
.upload-zone:hover:not(.upload-zone--disabled) {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  color: var(--color-primary-dark);
}
.upload-zone:focus-visible { outline: 2px solid var(--color-primary); outline-offset: 2px; }
.upload-zone--dragging {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  color: var(--color-primary-dark);
}
.upload-zone--disabled { opacity: 0.7; cursor: not-allowed; }
.upload-zone-icon { color: var(--color-primary); }
.upload-zone:hover:not(.upload-zone--disabled) .upload-zone-icon { color: var(--color-primary-dark); }
.upload-zone-text { font-size: 14px; font-weight: 500; color: var(--color-text-secondary); }
.upload-zone:hover:not(.upload-zone--disabled) .upload-zone-text { color: var(--color-primary-dark); }
.upload-zone-sub { font-size: 12px; }

.preview-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
  margin-top: 16px;
}
.preview-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: var(--radius-sm);
  overflow: hidden;
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  transition: box-shadow var(--transition);
}
.preview-item:hover { box-shadow: var(--shadow-hover); }
.preview-inner {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.preview-video {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: var(--color-primary);
}
.preview-video-label { font-size: 12px; color: var(--color-text-muted); }
.preview-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--color-text-muted);
}
.preview-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.preview-loading-text { font-size: 12px; }
.preview-remove {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.55);
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background var(--transition), transform var(--transition);
}
.preview-remove:hover:not(:disabled) { background: rgba(0, 0, 0, 0.75); transform: scale(1.05); }
.preview-remove:disabled { opacity: 0.6; cursor: not-allowed; }
.preview-remove:focus-visible { outline: 2px solid var(--color-primary); outline-offset: 2px; }
.preview-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.75);
  border-radius: var(--radius-sm);
}
.preview-overlay .preview-spinner { position: static; }

.publish-textarea { resize: vertical; min-height: 100px; width: 100%; box-sizing: border-box; }
.btn-submit { width: 100%; max-width: 100%; margin-top: 8px; }

@media (max-width: 768px) {
  .publish-page { padding: 0 16px; }
  .publish-page .card { padding: 16px; }
  .upload-zone { min-height: 120px; padding: 20px 16px; }
  .preview-list { grid-template-columns: repeat(auto-fill, minmax(88px, 1fr)); gap: 10px; margin-top: 12px; }
  .preview-remove { width: 24px; height: 24px; top: 4px; right: 4px; }
}
@media (max-width: 375px) {
  .publish-page { padding: 0 12px; }
  .publish-page .card { padding: 12px; }
  .upload-zone { min-height: 110px; padding: 16px 12px; }
  .preview-list { grid-template-columns: repeat(auto-fill, minmax(76px, 1fr)); gap: 8px; }
  .preview-remove { width: 22px; height: 22px; top: 4px; right: 4px; }
}
@media (prefers-reduced-motion: reduce) {
  .preview-spinner { animation: none; }
  .preview-remove:hover:not(:disabled) { transform: none; }
}
</style>
