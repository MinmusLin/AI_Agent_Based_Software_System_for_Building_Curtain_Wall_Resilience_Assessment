import { fileURLToPath, URL } from 'node:url';

import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// https://vite.dev/config
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) {
            return undefined;
          }
          if (id.includes('/react') || id.includes('/react-dom') || id.includes('/react-router-dom')) {
            return 'vendor-react';
          }
          if (id.includes('/antd') || id.includes('/@ant-design')) {
            return 'vendor-antd';
          }
          return 'vendor';
        },
      },
    },
  },
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: [{ find: '@', replacement: fileURLToPath(new URL('./src', import.meta.url)) }],
  },
});
