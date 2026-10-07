import http from '@/services/http'

// Konten publik (Fase 5): banner Beranda, artikel/promo/event, dan pemetaan bagian rumah di hero.

/** Banner yang sedang tayang di satu posisi: HOME_SLIDER (slider utama) atau HOME_SIDE (samping slider). */
export async function listBanners(position) {
  const { data } = await http.get('/banners', { params: { position } })
  return data.data
}

/**
 * Artikel/promo/event yang sudah terbit, terbaru dulu. Parameter: type (mis. 'ARTIKEL' atau 'PROMO,EVENT'),
 * page, per_page. Hasil { data: [{ id, type, title, slug, thumbnail_url, excerpt, published_at }], meta }.
 */
export async function listArticles(params = {}, { signal } = {}) {
  const { data } = await http.get('/articles', { params, signal })
  return data
}

/** Detail artikel. Draf hanya bisa dibuka dengan token pratinjau (?preview=); selain itu 404. */
export async function getArticle(slug, preview = null, { signal } = {}) {
  const { data } = await http.get(`/articles/${encodeURIComponent(slug)}`, { params: preview ? { preview } : {}, signal })
  return data.data
}

/** Bagian rumah di hero Beranda → kategori katalog: [{ key, category: { slug, name } | null }]. */
export async function getHouseParts() {
  const { data } = await http.get('/home/house-parts')
  return data.data
}
