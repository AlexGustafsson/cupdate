/// <reference types="vite/client" />

interface ViteTypeOptions {
  strictImportMetaEnv: unknown
}

interface ImportMetaEnv {
  readonly VITE_API_ENDPOINT: string
  readonly VITE_DEMO_MODE?: string
  readonly VITE_CUPDATE_VERSION?: string
  readonly VITE_BASE_PATH?: string

  /** The product name (such as "Cupdate"). */
  readonly VITE_BRANDING_NAME: string
  /** The official OCI image reference (such as "github.com/alexgustafsson/cupdate"). */
  readonly VITE_BRANDING_OCI_REFERENCE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
