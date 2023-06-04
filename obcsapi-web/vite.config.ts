import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import viteCompression from 'vite-plugin-compression'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), 
    viteCompression({
      threshold: 1024,
    }),VitePWA({
    registerType: 'autoUpdate',
    manifest: {
      name: 'Obcsapi',
      short_name: 'Note',
      description: 'Obcsapi',
      theme_color: '#ffffff',
      icons: [
        {
          src: 'pwa-192x192.png',
          sizes: '192x192',
          type: 'image/png'
        }
      ]
    }
  })],
  base: "/web/",
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
})
