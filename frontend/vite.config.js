import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server : {
    proxy: {
      '/login': 'http://localhost:8080',
      '/upload': 'http://localhost:8080',
      '/members': 'http://localhost:8080',
      '/member': 'http://localhost:8080',
    }
  }
})
