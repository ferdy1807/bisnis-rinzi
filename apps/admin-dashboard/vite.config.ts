import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@frontend': fileURLToPath(new URL('../../packages/frontend', import.meta.url)),
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  envDir: '../../',
  server: {
    host: true,
    port: 5176 // Port khusus untuk admin-dashboard
  }
});