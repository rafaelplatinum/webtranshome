import http from '@/services/http'
import { bersihkanCachePublik } from './cache'

// Canonical backend URL: /api/v1/admin/catalog/* (never /admin/* for catalog).
// Mock adapter retains its original /admin/* contract until removed.
const MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const URL = (jenis) => MOCK ? `/admin/${jenis}` : `/admin/catalog/${jenis}`

const snake = (obj = {}) => Object.fromEntries(Object.entries(obj).map(([k, v]) => [k.replace(/[A-Z]/g, x => '_' + x.toLowerCase()), v]))
const camel = (obj = {}) => Object.fromEntries(Object.entries(obj).map(([k, v]) => [k.replace(/_([a-z])/g, (_, c) => c.toUpperCase()), v]))
const daftar = (res, jenis) => (MOCK ? res.data.data : (res.data[jenis] ?? []).map(snake))
const satu = (res) => MOCK ? res.data.data : snake(res.data)

export async function listMaster(jenis) {
  const response = await http.get(URL(jenis))
  let items = daftar(response, jenis)
  // Brand home flags come from the public brand projection; admin legacy list lacks them.
  if (!MOCK && jenis === 'brands') {
    try {
      const { data } = await http.get('/brands')
      const home = new Map((data.data ?? []).map(x => [x.id, x]))
      items = items.map(x => ({ ...x, show_on_home: home.get(x.id)?.show_on_home ?? false, home_order: home.get(x.id)?.home_order ?? null }))
    } catch { /* List remains usable even if the public endpoint is temporarily unavailable. */ }
  }
  return items
}
export async function createMaster(jenis, payload) {
  const { data } = await http.post(URL(jenis), MOCK ? payload : camel(payload))
  bersihkanCachePublik()
  return MOCK ? data.data : snake(data)
}
export async function updateMaster(jenis, id, payload) {
  const { data } = await http[MOCK ? 'put' : 'patch'](`${URL(jenis)}/${id}`, MOCK ? payload : camel(payload))
  bersihkanCachePublik()
  return MOCK ? data.data : snake(data)
}
export async function patchMaster(jenis, id, payload) {
  // The Go backend accepts PATCH with only the changed fields.
  const { data } = await http.patch(`${URL(jenis)}/${id}`, MOCK ? payload : camel(payload), { tanpaToast: true })
  bersihkanCachePublik()
  return MOCK ? data.data : snake(data)
}
export async function deleteBrand(id) {
  await http.delete(MOCK ? `/admin/brands/${id}` : `/admin/catalog/brands/${id}`, { tanpaToast: true })
  bersihkanCachePublik()
}
export async function saveBrandHome(brandIds) {
  const { data } = await http.put(MOCK ? '/admin/brands/home' : '/admin/catalog/brands/home', { brand_ids: brandIds }, { tanpaToast: true })
  bersihkanCachePublik()
  return MOCK ? data.data : data.data
}