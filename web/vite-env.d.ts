/// <reference types="vite/client" />

interface ViteTypeOptions {
  strictImportMetaEnv: unknown
}

interface ImportMetaEnv {
  readonly VITE_API_ENDPOINT: string
  readonly VITE_DEMO_MODE?: string
  readonly VITE_CUPDATE_VERSION?: string
  readonly VITE_BASE_PATH?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
