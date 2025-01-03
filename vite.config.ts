import { resolve } from 'node:path'
import react from '@vitejs/plugin-react'
import { type UserConfigExport, defineConfig } from 'vite'

export default ({ mode }: { mode: string }): UserConfigExport => {
  // The dev server listens on port 8080, use it during development with vite
  if (!process.env.VITE_API_ENDPOINT) {
    if (mode === 'development') {
      process.env.VITE_API_ENDPOINT = 'http://localhost:8080/api/v1'
    } else {
      process.env.VITE_API_ENDPOINT = '/api/v1'
    }
  }

  return defineConfig({
    plugins: [react()],
    root: 'web',
    base: process.env.VITE_BASE_PATH,
    build: {
      outDir: '../internal/web/public',
      emptyOutDir: true,
      sourcemap: true,
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, '/web'),
      },
    },
  })
}
