import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  server: {
    port: 3000,
    proxy: {
      '/health': 'http://localhost:8080',
      '/accounts': 'http://localhost:8080',
      '/transfers': 'http://localhost:8080',
    },
  },
})
