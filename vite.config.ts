import { resolve } from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig, type Plugin, type UserConfigExport } from 'vite'

// Markdown can have imgs using the "height" attribute. This is overridden by
// tailwind's defaults as it has a higher specificity. The rule is unused by us
// and can be introduced manually when needed instead.
// SEE: https://github.com/tailwindlabs/tailwindcss/pull/7742#issuecomment-1061332148
const uselessRules = [
  `  img, video {
    max-width: 100%;
    height: auto;
  }
 `,
  'img,video{max-width:100%;height:auto}',
]

const removeUselessRule: Plugin = {
  name: 'remove useless rule',
  transform(code, id) {
    if (id.endsWith('.css')) {
      for (const rule of uselessRules) {
        code = code.replace(rule, '')
      }
      return code
    }
  },
}

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
    plugins: [react(), tailwindcss(), removeUselessRule],
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
