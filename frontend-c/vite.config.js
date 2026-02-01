import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  base: '/c/',
  server: {
    port: 5173,
    host: true,
    proxy: {
      '/api': { target: process.env.VITE_PROXY_TARGET || 'http://localhost:8080', changeOrigin: true },
      '/uploads': { target: process.env.VITE_PROXY_TARGET || 'http://localhost:8080', changeOrigin: true },
    },
  },
  build: {
    outDir: 'dist',
  },
})
