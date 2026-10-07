import http from './http'
import { daftarkanPembersihCache } from '@/services/cache'

let cachePohon = null
let pohonSiap = null
daftarkanPembersihCache(() => {
  cachePohon = null
  pohonSiap = null
})

/**
 * Pohon kategori aktif (induk → anak) untuk mega menu, filter, dan breadcrumb.
 * Disimpan di memori selama tab terbuka; kalau gagal, cache dikosongkan agar bisa dicoba lagi.
 */
export function getCategoryTree() {
  if (!cachePohon) {
    cachePohon = http
      .get('/categories', { params: { tree: 1 } })
      .then(({ data }) => (pohonSiap = data.data))
      .catch((error) => {
        cachePohon = null
        throw error
      })
  }
  return cachePohon
}

/** Pohon yang sudah pernah dimuat (atau null): halaman bisa langsung render tanpa menunggu, penting untuk tombol Back. */
export function pohonTersimpan() {
  return pohonSiap
}

/** Cari kategori berdasarkan slug di pohon: { kategori, induk } (induk null untuk kategori induk), atau null. */
export function cariKategori(pohon, slug) {
  for (const induk of pohon ?? []) {
    if (induk.slug === slug) return { kategori: induk, induk: null }
    const anak = induk.children?.find((c) => c.slug === slug)
    if (anak) return { kategori: anak, induk }
  }
  return null
}
