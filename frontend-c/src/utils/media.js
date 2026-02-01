/**
 * 媒体 URL：开发环境 (5173/5174) 走 Nginx 80，避免 Vite 代理 /uploads 404
 */
export function mediaUrl(url) {
  if (!url) return ''
  const path = url.startsWith('http') ? url : (url.startsWith('/') ? url : '/' + url)
  if (path.startsWith('http')) return path
  if (typeof window === 'undefined') return path
  const origin = window.location.origin
  if (origin.includes(':5173') || origin.includes(':5174')) {
    return window.location.protocol + '//' + window.location.hostname + path
  }
  return origin + path
}

export function isVideo(media) {
  if (!media || !media.url) return false
  const url = (media.url || '').toLowerCase()
  return url.includes('video') || url.endsWith('.mp4') || url.endsWith('.webm')
}
