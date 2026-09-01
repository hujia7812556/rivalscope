import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    // 开发期把 /api 代理到本地 Go 服务;生产环境前端嵌入 Go 二进制,天然同源
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:17317',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    chunkSizeWarningLimit: 1600,
  },
})
