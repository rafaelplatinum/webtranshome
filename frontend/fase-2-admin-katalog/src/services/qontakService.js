import http from '@/services/http'

/** Ulangi sync kontak ke Qontak. Proses berjalan di background: status langsung PENDING. */
export async function retryContact(id) {
  const { data } = await http.post(`/admin/qontak/contacts/${id}/retry`)
  return data.data
}

/** Status sync terbaru: { id, name, sync_status, retry_count, status_code, last_synced_at }. */
export async function getContact(id) {
  const { data } = await http.get(`/admin/qontak/contacts/${id}`)
  return data.data
}

/**
 * Daftar kontak di halaman Sync Qontak (Fase 6). Parameter: status (PENDING | SYNCED | FAILED), q, page.
 * Item: { id, user_id, name, email, phone_number, member_code, sync_status, retry_count, status_code, last_synced_at, updated_at }.
 * `meta.counts` = jumlah per status { ALL, PENDING, SYNCED, FAILED } untuk pencarian yang sama.
 */
export async function listContacts(params = {}, { signal } = {}) {
  const { data } = await http.get('/admin/qontak/contacts', { params, signal })
  return data
}

/** Masukkan semua kontak FAILED ke antrean lagi. Hasil: jumlah kontak yang diulang. */
export async function retryFailed() {
  const { data } = await http.post('/admin/qontak/retry-failed')
  return data.data.count
}

/** Riwayat pengiriman (sync_logs) satu kontak, terbaru dulu: request, response, status code, retry_count. */
export async function getContactLogs(id) {
  const { data } = await http.get(`/admin/qontak/contacts/${id}/logs`)
  return data.data
}
