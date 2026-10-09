import http from '@/services/http'
import { bersihkanCachePublik } from './cache'

// Master katalog: jenis = 'categories' | 'brands' | 'rooms'. Daftar berisi yang nonaktif juga + product_count.

export async function listMaster(jenis) {
  const { data } = await http.get(`/admin/${jenis}`)
  return data.data
}

export async function createMaster(jenis, payload) {
  const { data } = await http.post(`/admin/${jenis}`, payload)
  bersihkanCachePublik()
  return data.data
}

export async function updateMaster(jenis, id, payload) {
  const { data } = await http.put(`/admin/${jenis}/${id}`, payload)
  bersihkanCachePublik()
  return data.data
}

/** Hapus permanen brand nonaktif yang tidak dipakai produk. 409 bila masih dipakai, 422 bila masih aktif. */
export async function deleteBrand(id) {
  await http.delete(`/admin/brands/${id}`, { tanpaToast: true })
  bersihkanCachePublik()
}

/** { is_active } atau { sort_order }. */
export async function patchMaster(jenis, id, perubahan) {
  const { data } = await http.patch(`/admin/${jenis}/${id}`, perubahan, { tanpaToast: true })
  bersihkanCachePublik()
  return data.data
}
