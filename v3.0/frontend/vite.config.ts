import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(),tailwindcss()],
  server: {
    proxy: {
      '/api/engine/readyz': { target: 'http://127.0.0.1:6081', changeOrigin: true, rewrite: path => path.replace(/^\/api\/engine\/readyz$/, '/readyz') },
      '/api/engine': { target: 'http://127.0.0.1:6081', changeOrigin: true, rewrite: path => path.replace(/^\/api\/engine/, '/api/v1') },
      '/api/gateway': { target: 'http://127.0.0.1:6080', changeOrigin: true, rewrite: path => path.replace(/^\/api\/gateway/, '/api/v1') },
      '/ws': { target: 'ws://127.0.0.1:6080', ws: true, changeOrigin: true },
    },
  },
})
