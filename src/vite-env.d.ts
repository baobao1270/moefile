/// <reference types="vite/client" />
interface ImportMetaEnv {
  readonly NODE_ENV: 'development' | 'production'
  readonly APP_NAME: string
  readonly APP_VERSION: string
  readonly BUILD_TIMESTAMP: string
  readonly COPYRIGHT_HOLDER: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
