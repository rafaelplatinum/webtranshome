import http from '@/services/http'

/**
 * Log aktivitas semua admin (izin log.view). Parameter: admin_id, module, from, to (YYYY-MM-DD, WIB), page.
 * Item: { id, admin_id, admin_name, action, module, before_data, after_data, ip_address, created_at }. Hasil: { data, meta }.
 */
export async function listActivityLogs(params = {}, { signal } = {}) {
  const { data } = await http.get('/admin/activity-logs', { params, signal })
  return data
}

/** Pilihan filter: { admins: [{ id, name }], modules: ['Produk', …] }. */
export async function getActivityLogOptions() {
  const { data } = await http.get('/admin/activity-logs/options')
  return data.data
}
