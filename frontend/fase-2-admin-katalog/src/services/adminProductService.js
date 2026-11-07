import http from '@/services/http'
import { bersihkanCachePublik } from './cache'

const MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const BASE = MOCK ? '/admin/products' : '/admin/catalog/products'
const snake = (obj = {}) => Object.fromEntries(Object.entries(obj).map(([k, v]) => [k.replace(/[A-Z]/g, s => '_' + s.toLowerCase()), v]))
const camel = (obj = {}) => Object.fromEntries(Object.entries(obj).map(([k, v]) => [k.replace(/_([a-z])/g, (_, c) => c.toUpperCase()), v]))
const normal = (p) => {
  const a = snake(p)
  try { if (typeof a.specifications === 'string') a.specifications = JSON.parse(a.specifications) } catch { a.specifications = {} }
  return { show_price: true, images: [], room_ids: [], ...a }
}
const payloadGo = (p) => {
  const result = camel(p)
  if (p.specifications && typeof p.specifications === 'object') result.specifications = JSON.stringify(p.specifications)
  // Legacy product mutation does not process images/show_price/room_ids; callers
  // must use image and room endpoints separately (see guide).
  delete result.showPrice
  delete result.images
  delete result.roomIds
  return result
}

export async function listAdminProducts(params = {}, { signal } = {}) {
  if (MOCK) { const { data } = await http.get(BASE, { params, signal }); return data }
  const { data } = await http.get(BASE, { signal })
  let items = (data.products ?? []).map(normal)
  const q = String(params.q ?? '').toLowerCase()
  if (q) items = items.filter(p => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q))
  if (params.category) items = items.filter(p => String(p.category_id) === String(params.category))
  if (params.brand) items = items.filter(p => String(p.brand_id) === String(params.brand))
  if (params.stock) items = items.filter(p => p.stock_status === params.stock)
  if (params.active === '1' || params.active === '0') items = items.filter(p => p.is_active === (params.active === '1'))
  const page = Math.max(1, Number(params.page) || 1)
  const perPage = Math.min(100, Math.max(1, Number(params.per_page) || 20))
  return { data: items.slice((page - 1) * perPage, page * perPage), meta: { page, per_page: perPage, total: items.length, total_pages: Math.ceil(items.length / perPage) } }
}
export async function getAdminProduct(id) {
  const { data } = await http.get(`${BASE}/${id}`)
  return MOCK ? data.data : normal(data)
}
export async function checkSku(sku, exclude = null) {
  const { data } = await http.get(`${BASE}/check-sku`, { params: { sku, exclude: exclude ?? undefined } })
  return data.data.available
}
export async function createProduct(payload) {
  if (!MOCK && (payload.images?.length || payload.show_price === false)) throw new Error('Backend Admin Produk belum mendukung penyimpanan foto/show_price. Jangan lanjut agar data tidak hilang.')
  const { data } = await http.post(BASE, MOCK ? payload : payloadGo(payload))
  bersihkanCachePublik()
  return MOCK ? data.data : normal(data)
}
export async function updateProduct(id, payload) {
  if (!MOCK && (payload.images?.length || payload.show_price === false)) throw new Error('Backend Admin Produk belum mendukung penyimpanan foto/show_price. Jangan lanjut agar data tidak hilang.')
  const { data } = await http[MOCK ? 'put' : 'patch'](`${BASE}/${id}`, MOCK ? payload : payloadGo(payload))
  if (!MOCK && payload.room_ids) await http.put(`${BASE}/${id}/rooms`, { roomIds: payload.room_ids })
  bersihkanCachePublik()
  return MOCK ? data.data : normal(data)
}
export async function patchProduct(id, patch) {
  if (!MOCK && ('show_price' in patch)) throw new Error('Backend belum menyimpan show_price.')
  const { data } = await http.patch(`${BASE}/${id}`, MOCK ? patch : camel(patch), { tanpaToast: true })
  bersihkanCachePublik()
  return MOCK ? data.data : normal(data)
}