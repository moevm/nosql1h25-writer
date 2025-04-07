import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { reactRouter } from '@react-router/dev/vite'

import path from 'node:path'
import { packageDirectorySync } from 'pkg-dir'

const packageRoot = packageDirectorySync() || '/'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), reactRouter()],
  resolve: {
    alias: {
      'src': path.resolve(packageRoot, './src'),
    },
  },
})
