import http from '@/services/http'

// Membership admin (Fase 4). Member dikenali dengan `id` = users.id.

/**
 * Daftar member. Parameter: q (kode member, email, HP, nama), consent (1/0), hp (1/0: punya nomor HP),
 * sync (PENDING | SYNCED | FAILED), active (1/0), page, per_page. Hasil: { data, meta }.
 */
export async function listMembers(params = {}, { signal } = {}) {
  const { data } = await http.get('/admin/members', { params, signal })
  return data
}

/** Cari member untuk input poin. Kode member / email / HP yang persis → 1 hasil; selain itu maks 10. */
export async function searchMembers(q) {
  const { data } = await http.get('/admin/members/search', { params: { q } })
  return data.data
}

/** Detail member: data akun, `balance`, `last_login_at`, dan `qontak` (status sync, untuk "Coba sync ulang"). */
export async function getMember(id) {
  const { data } = await http.get(`/admin/members/${id}`)
  return data.data
}

/** { is_active } → member terbaru. Nonaktif: member tidak bisa masuk. */
export async function patchMember(id, perubahan) {
  const { data } = await http.patch(`/admin/members/${id}`, perubahan)
  return data.data
}

/** Riwayat poin satu member, terbaru dulu: user_id, type (EARN | REDEEM | ADJUST), page, per_page. meta berisi `balance`. */
export async function listPoints(params, { signal } = {}) {
  const { data } = await http.get('/admin/points', { params, signal })
  return data
}

/** Nomor nota (EARN) sudah pernah diinput? → { available } atau { available: false, used_at, member_code }. */
export async function checkReference(ref) {
  const { data } = await http.get('/admin/points/check-reference', { params: { ref } })
  return data.data
}

/**
 * Catat transaksi poin → { transaction, balance }.
 * EARN { user_id, type, reference_no, purchase_amount, note }: poin dihitung backend dari point_ratio_rupiah.
 * REDEEM { user_id, type, points (positif), note }. ADJUST { user_id, type, points (+/−), note }, izin point.adjust.
 * 409 nomor nota sudah dipakai; 422 saldo kurang / isian salah (errors per field).
 */
export async function createPoint(payload) {
  const { data } = await http.post('/admin/points', payload)
  return data.data
}
