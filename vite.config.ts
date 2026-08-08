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

const generateManifest: () => Plugin = () => {
  let env: Record<string, any>
  return {
    name: 'substitute manifest values',
    configResolved(config) {
      env = config.env
    },
    generateBundle() {
      this.emitFile({
        type: 'asset',
        fileName: 'manifest.json',
        source: JSON.stringify(
          {
            short_name: env.VITE_BRANDING_NAME,
            name: env.VITE_BRANDING_NAME,
            icons: [
              {
                src: '/assets/icon.png',
                purpose: 'any',
                sizes: '512x512',
                type: 'image/png',
              },
              {
                src: '/assets/maskable-icon.png',
                purpose: 'maskable',
                sizes: '512x512',
                type: 'image/png',
              },
            ],
            start_url: '/',
            display: 'standalone',
          },
          null,
          2
        ),
      })
    },
  }
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

  if (!process.env.VITE_BRANDING_NAME) {
    process.env.VITE_BRANDING_NAME = 'Cupdate'
  }

  if (!process.env.VITE_BRANDING_OCI_REFERENCE) {
    process.env.VITE_BRANDING_OCI_REFERENCE = 'ghcr.io/alexgustafsson/cupdate'
  }

  return defineConfig({
    plugins: [react(), tailwindcss(), removeUselessRule, generateManifest()],
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
