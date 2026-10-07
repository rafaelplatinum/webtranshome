import http from '@/services/http'

// Data milik member yang sedang masuk (/me/*). Pesan gagal ditampilkan di halaman, jadi tanpa toast.
const TANPA_TOAST = { tanpaToast: true }

/** { full_name, phone_number|null } → user terbaru. 409 bila nomor HP dipakai akun lain. */
export async function updateProfile(payload) {
  const { data } = await http.patch('/me', payload, TANPA_TOAST)
  return data.data.user
}

/** Persetujuan info promo & poin (communication_consent) → user terbaru. */
export async function updateConsent(setuju) {
  const { data } = await http.patch('/me/consent', { communication_consent: setuju }, TANPA_TOAST)
  return data.data.user
}

/** Ubah password ({ old_password, new_password }) atau buat password akun Google ({ new_password }). */
export async function changePassword(payload) {
  const { data } = await http.post('/me/password', payload, TANPA_TOAST)
  return data
}

/** Riwayat poin: type (EARN | REDEEM | kosong = semua), page, per_page. Hasil { data, meta: { balance, page, total, total_pages } }. */
export async function getPoints(params = {}, { signal } = {}) {
  const { data } = await http.get('/me/points', { params, signal })
  return data
}
