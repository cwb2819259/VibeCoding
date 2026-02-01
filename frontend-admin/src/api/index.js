import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1'

const client = axios.create({
  baseURL,
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const auth = {
  login: (username, password) => client.post('/admin/auth/login', { username, password }),
}

/** 媒体 URL：与 C 端一致，相对路径补全为当前 origin */
export function mediaUrl(url) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  const path = url.startsWith('/') ? url : '/' + url
  if (typeof window === 'undefined') return path
  return window.location.origin + path
}

export const posts = {
  list: (params) => client.get('/admin/posts', { params }),
  get: (id) => client.get(`/admin/posts/${id}`),
  updateStatus: (id, status) => client.patch(`/admin/posts/${id}`, { status }),
  delete: (id) => client.delete(`/admin/posts/${id}`),
}

export const users = {
  list: (params) => client.get('/admin/users', { params }),
}

export const stats = {
  get: () => client.get('/admin/stats'),
}

export default client
