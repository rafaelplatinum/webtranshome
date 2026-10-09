import http from '@/services/http'
import { bersihkanCachePublik } from './cache'

/**
 * Daftar produk admin (aktif dan nonaktif). Parameter: q, category (id), brand (id), stock, active (1/0),
 * foto (1/0: punya foto utama), room (id), page, per_page. Hasil: { data, meta }.
 */
export async function listAdminProducts(params = {}, { signal } = {}) {
  const { data } = await http.get('/admin/products', { params, signal })
  return data
}

/** Produk lengkap untuk form: kolom products + images + room_ids. */
export async function getAdminProduct(id) {
  const { data } = await http.get(`/admin/products/${id}`)
  return data.data
}

/** true bila SKU belum dipakai produk lain (`exclude` = id produk yang sedang diubah). */
export async function checkSku(sku, exclude = null) {
  const { data } = await http.get('/admin/products/check-sku', { params: { sku, exclude: exclude ?? undefined } })
  return data.data.available
}

export async function createProduct(payload) {
  const { data } = await http.post('/admin/products', payload)
  bersihkanCachePublik()
  return data.data
}

export async function updateProduct(id, payload) {
  const { data } = await http.put(`/admin/products/${id}`, payload)
  bersihkanCachePublik()
  return data.data
}

/** Ubah sebagian: { is_featured } / { is_active } (422 bila diaktifkan tanpa foto utama). */
export async function patchProduct(id, perubahan) {
  const { data } = await http.patch(`/admin/products/${id}`, perubahan, { tanpaToast: true })
  bersihkanCachePublik()
  return data.data
}
