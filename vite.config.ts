import fs from 'fs'
import { resolve } from 'path'
import { defineConfig, PluginOption } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { viteSingleFile } from 'vite-plugin-singlefile'

const ROOT_DIR = resolve(__dirname, '.')
const SOURCE_DIR = resolve(ROOT_DIR, 'src')
const OUTPUT_DIR = resolve(ROOT_DIR, 'dist')
const PUBLIC_DIR = resolve(ROOT_DIR, 'public')
const PAGE_NAME = process.env.PAGE_NAME || 'index'

const rewriteIndexPlugin: PluginOption = {
  name: 'vite-rewrite-index',
  apply: 'serve',
  enforce: 'post',
  configureServer: (server) => {
    server.middlewares.use(async (req, _res, next) => {
      if (!req.url) { return next() }
      const url = decodeURIComponent(req.url.split('?')[0])
      const file = (await server.moduleGraph.resolveUrl(url))[1]
      if (file.startsWith('/@')) { return next() }

      for (const search of [file, file.substring(1), resolve(PUBLIC_DIR, file.substring(1))]) {
        if (fs.existsSync(search) && fs.statSync(search).isFile()) { return next() }
      }

      console.log('REWRITE', url)
      req.url = `/${PAGE_NAME}.html`
      return next()
    })
  }
}

const ensureDocTypePlugin: PluginOption = {
  name: 'ensure-doctype',
  apply: 'build',
  enforce: 'post',
  generateBundle: async (_, bundle) => {
    for (const key in bundle) {
      const file = bundle[key]
      if (file.type == 'asset' && file.fileName.endsWith('.html')) {
        file.source = file.source
          .toString()
          .replaceAll("https://cdn.local", process.env.CDN_URL || "")
        if (PAGE_NAME === 'index') { continue }
        file.source = `<!DOCTYPE html>\n${file.source}`
      }
    }
  },
};

const ensureXmlCompatPlugin: PluginOption = {
  name: 'xml-tidy',
  apply: 'build',
  enforce: 'post',
  generateBundle: async (_, bundle) => {
    for (const key in bundle) {
      const file = bundle[key]
      if (PAGE_NAME !== 'index') { continue }
      if (file.type == 'asset' && file.fileName.endsWith('.html')) {
        file.source = file.source.toString()
          .replaceAll(/crossorigin /g, '')
          .replaceAll(/crossorigin>/g, '>')
      }
      if (file.type == 'asset' && file.fileName.endsWith('.css')) {
        file.source = `<![CDATA[\n${file.source}\n]]>`
      }

      if (file.type == 'chunk' && file.fileName.endsWith('.js')) {
        file.code = `\n<![CDATA[\n${file.code}\n]]>`
      }
    }
  },
};

function processEnv(defaultValues: Record<string, string>) {
  const result: Record<string, string> = {}
  for (const key in defaultValues) {
    result[`import.meta.env.${key}`] = JSON.stringify(process.env[key] || defaultValues[key])
  }
  return result
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [rewriteIndexPlugin, react(), ensureXmlCompatPlugin, viteSingleFile(), ensureDocTypePlugin],
  define: processEnv({
    APP_NAME: 'MoeFile',
    APP_VERSION: 'DEV',
    BUILD_TIMESTAMP: new Date().toISOString(),
    COPYRIGHT_HOLDER: 'MoeFile',
    NODE_ENV: 'development',
  }),
  resolve: {
    alias: {
      '@': SOURCE_DIR,
    },
  },
  publicDir: PUBLIC_DIR,
  root: resolve(SOURCE_DIR, PAGE_NAME),
  build: {
    rollupOptions: {
      input: resolve(SOURCE_DIR, PAGE_NAME, `${PAGE_NAME}.html`),
    },
    outDir: OUTPUT_DIR,
    emptyOutDir: true,
  },
})
