import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1'

const client = axios.create({
  baseURL,
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('c_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const auth = {
  login: (phone, code) => client.post('/auth/login', { phone, code }),
}

export const posts = {
  list: (params) => client.get('/posts', { params }),
  get: (id) => client.get(`/posts/${id}`),
  create: (data) => client.post('/posts', data),
  myPosts: (params) => client.get('/users/me/posts', { params }),
}

export const comments = {
  list: (postId, params) => client.get(`/posts/${postId}/comments`, { params }),
  create: (postId, data) => client.post(`/posts/${postId}/comments`, data),
}

export const likes = {
  like: (postId) => client.post(`/posts/${postId}/like`),
  unlike: (postId) => client.delete(`/posts/${postId}/like`),
}

export const notifications = {
  list: (params) => client.get('/notifications', { params }),
  unreadCount: () => client.get('/notifications/unread-count'),
  markRead: (id) => client.patch(`/notifications/${id}/read`),
}

export const upload = (file) => {
  const form = new FormData()
  form.append('file', file)
  return client.post('/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export default client
