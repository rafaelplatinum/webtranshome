import { existsSync, mkdirSync, readFileSync, statSync, writeFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv, normalizePath } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// Kode dipisah per fase: setiap <folder fase>/src memakai susunan yang sama (components/, views/, services/, ...).
// Import `@/...` dicari di src semua folder fase, berurutan. Satu path (mis. components/ui/Button.vue) hanya boleh ada di satu fase.
const FOLDER_FASE = [
  'fase-0-fondasi',
  'fase-1-katalog',
  'fase-2-admin-katalog',
  'fase-3-akun-member',
  'fase-4-membership-admin',
  'fase-5-konten',
  'fase-6-sistem',
]
const DASAR = fileURLToPath(new URL('.', import.meta.url))
const AKHIRAN = ['', '.js', '.vue', '.json', '/index.js']

function aliasFase() {
  return {
    name: 'transhome-alias-fase',
    enforce: 'pre',
    resolveId(sumber) {
      if (!sumber.startsWith('@/')) return null
      const [jalur, query] = sumber.slice(2).split('?')
      for (const folder of FOLDER_FASE) {
        for (const akhiran of AKHIRAN) {
          const calon = resolve(DASAR, folder, 'src', jalur + akhiran)
          // Path dinormalisasi seperti hasil resolve Vite, supaya satu file tidak termuat sebagai dua modul.
          if (existsSync(calon) && statSync(calon).isFile()) return normalizePath(calon) + (query ? `?${query}` : '')
        }
      }
      return null
    },
  }
}

// SEO di HTML mentah untuk halaman publik beralamat tetap (KEPUTUSAN #22). Tanpa ini, sumber halaman (Ctrl+U), alat cek
// SEO, dan pratinjau tautan WhatsApp/Facebook hanya melihat judul umum "Transhome", karena judul per halaman baru dipasang
// JavaScript (useSeo). Teks di sini harus sejalan dengan useSeo di halaman terkait. Halaman berparameter (/brand/:slug,
// /produk/:slug, /jasa/:slug, /artikel/:slug) butuh server yang menyisipkan meta per URL (catatan backend di
// docs/FRONTEND_BUILD_PLAN.md). Tag ditandai data-seo="statis" dan diambil alih / dibersihkan useSeo & router guard.
const SEO_STATIS = {
  '/': {
    judul: 'Bahan bangunan dan elektronik di Sleman',
    deskripsi: 'Katalog bahan bangunan dan elektronik Transhome di Sleman, Yogyakarta. Cek harga dan stok, lalu tanya langsung lewat WhatsApp.',
  },
  '/katalog': {
    judul: 'Katalog produk',
    deskripsi: 'Katalog bahan bangunan dan elektronik Transhome di Sleman: pilih kategori, brand, ruangan, dan harga, lalu tanya stok lewat WhatsApp.',
    remah: 'Katalog',
    jenis: 'CollectionPage',
  },
  '/brand': {
    judul: 'Brand partner bahan bangunan & elektronik',
    deskripsi: 'Daftar brand partner Transhome Sleman untuk bahan bangunan dan elektronik. Lihat produk tiap brand, cek harga dan stok, lalu tanya lewat WhatsApp.',
    remah: 'Brand',
    jenis: 'CollectionPage',
  },
  '/ruangan': {
    judul: 'Belanja per ruangan',
    deskripsi: 'Pilih ruangan seperti dapur, kamar mandi, atau teras, lalu lihat bahan bangunan yang cocok di Transhome Sleman.',
    remah: 'Ruangan',
    jenis: 'CollectionPage',
  },
  '/jasa': {
    judul: 'Jasa pasang & renovasi rumah',
    deskripsi: 'Jasa Transhome di Sleman: pasang keramik, pengecatan, perlengkapan kamar mandi, AC & elektronik, service, pintu & jendela, rakit perabot. Konsultasi lewat WhatsApp.',
    remah: 'Jasa',
    jenis: 'CollectionPage',
  },
  '/artikel': { judul: 'Artikel dan tips', deskripsi: 'Tips memilih dan merawat bahan bangunan dari Transhome, Sleman.', remah: 'Artikel', jenis: 'CollectionPage' },
  '/promo': { judul: 'Promo dan event', deskripsi: 'Promo dan event terbaru di toko Transhome, Sleman.', remah: 'Promo dan event', jenis: 'CollectionPage' },
  '/info-toko': {
    judul: 'Info toko',
    deskripsi: 'Alamat, jam buka, dan peta toko Transhome di Berbah, Sleman, Yogyakarta. Tanya langsung lewat WhatsApp.',
    remah: 'Info toko',
  },
}

const escHtml = (teks) => String(teks).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
const rapikanPath = (url = '/') => (url.split(/[?#]/)[0].replace(/\/+$/, '') || '/')

function pasangSeoStatis(html, path, situs) {
  const meta = SEO_STATIS[path]
  if (!meta) return html
  const judul = `${meta.judul} | Transhome`
  const url = situs ? `${situs}${path === '/' ? '/' : path}` : null
  const graf = [
    { '@type': meta.jenis ?? 'WebPage', name: judul, description: meta.deskripsi, ...(url ? { url } : {}) },
    ...(meta.remah && situs
      ? [{
          '@type': 'BreadcrumbList',
          itemListElement: [
            { '@type': 'ListItem', position: 1, name: 'Beranda', item: `${situs}/` },
            { '@type': 'ListItem', position: 2, name: meta.remah, item: url },
          ],
        }]
      : []),
  ]
  const tag = [
    url && `<link rel="canonical" href="${escHtml(url)}" data-seo="statis" />`,
    `<meta property="og:title" content="${escHtml(judul)}" data-seo="statis" />`,
    `<meta property="og:description" content="${escHtml(meta.deskripsi)}" data-seo="statis" />`,
    url && `<meta property="og:url" content="${escHtml(url)}" data-seo="statis" />`,
    `<script type="application/ld+json" data-seo="statis">${JSON.stringify({ '@context': 'https://schema.org', '@graph': graf }).replace(/</g, '\\u003c')}</script>`,
  ].filter(Boolean)
  return html
    .replace(/\s*<!-- seo-statis -->[\s\S]*?<!-- \/seo-statis -->/, '')
    .replace(/<title>[\s\S]*?<\/title>/, `<title>${escHtml(judul)}</title>`)
    .replace(/(<meta name="description" content=")[^"]*(")/, `$1${escHtml(meta.deskripsi)}$2`)
    .replace('</head>', `  <!-- seo-statis -->\n    ${tag.join('\n    ')}\n    <!-- /seo-statis -->\n  </head>`)
}

function seoStatis(situsProduksi) {
  let dist = null
  return {
    name: 'transhome-seo-statis',
    configResolved(config) {
      dist = resolve(config.root, config.build.outDir)
    },
    transformIndexHtml: {
      order: 'post',
      handler(html, ctx) {
        // Dev: ctx.originalUrl = alamat yang dibuka (SPA fallback). Build: index.html = Beranda.
        const path = ctx.server ? rapikanPath(ctx.originalUrl) : '/'
        const situs = situsProduksi || ctx.server?.resolvedUrls?.local?.[0]?.replace(/\/+$/, '') || null
        return pasangSeoStatis(html, path, situs)
      },
    },
    // Build: salin index.html ke <halaman>/index.html dengan meta halaman itu, mis. dist/brand/index.html.
    // Server/hosting perlu menyajikan file ini untuk /brand (try_files $uri $uri/ /index.html).
    closeBundle() {
      const sumber = dist && resolve(dist, 'index.html')
      if (!sumber || !existsSync(sumber)) return
      const html = readFileSync(sumber, 'utf8')
      for (const path of Object.keys(SEO_STATIS).filter((p) => p !== '/')) {
        const folder = resolve(dist, `.${path}`)
        mkdirSync(folder, { recursive: true })
        writeFileSync(resolve(folder, 'index.html'), pasangSeoStatis(html, path, situsProduksi || null))
      }
    },
  }
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, DASAR, 'VITE_')
  const situs = (env.VITE_SITE_URL || '').replace(/\/+$/, '') || null
  return {
    plugins: [aliasFase(), vue(), tailwindcss(), seoStatis(situs)],
  }
})
