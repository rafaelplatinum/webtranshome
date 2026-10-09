import http from './http'

/**
 * Daftar produk aktif. Parameter: q, category (slug), brand / stock / room (dipisah koma),
 * attr[kunci] (dipisah koma), min, max, exclude, featured,
 * sort (terbaru | harga_terendah | harga_tertinggi | nama), page, per_page.
 * Hasil: { data: [...], meta: { page, per_page, total, total_pages, facets } }.
 */
export async function listProducts(params = {}, { signal } = {}) {
  const { data } = await http.get('/products', { params, signal })
  return data
}

/** Detail produk lengkap dengan images, rooms, brand, category (+ parent). 404 bila tidak ada atau nonaktif. */
export async function getProduct(slug, { signal } = {}) {
  const { data } = await http.get(`/products/${encodeURIComponent(slug)}`, { signal })
  return data.data
}

/** Saran pencarian: { products, categories, brands }. Dipanggil mulai 2 huruf. */
export async function suggestProducts(q, { signal } = {}) {
  const { data } = await http.get('/products/suggest', { params: { q }, signal })
  return data.data
}
