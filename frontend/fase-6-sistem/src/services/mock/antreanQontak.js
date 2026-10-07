// Simulasi proses background sync kontak ke Mekari Qontak (mode mock). TODO: ganti ke API.
// Backend mengirim kontak lewat antrean; frontend hanya membaca status PENDING / SYNCED / FAILED.
import { qontakContacts, syncLogs } from '@/services/mock/data/admin'
import { users } from '@/services/mock/data/akun'
import { idBaru, simpanDb } from '@/services/mock/db'
import { sekarang } from '@/services/mock/sesi'

// Tombol "Matikan Qontak" di /dev/components memasang penanda ini: semua pengiriman berikutnya gagal (503).
export const KUNCI_QONTAK_MATI = 'transhome.mockQontakMati'
const LAMA_PROSES = 1500

export function qontakMati() {
  try {
    return localStorage.getItem(KUNCI_QONTAK_MATI) === '1'
  } catch {
    return false
  }
}

/** Masukkan kontak ke antrean: langsung PENDING, selesai ±1,5 detik kemudian (saat statusnya dibaca). */
export function antrekan(c) {
  c.sync_status = 'PENDING'
  c.mulai_ulang = Date.now()
  c.updated_at = sekarang()
}

/** Isi yang dikirim ke Qontak. Member tanpa HP tetap dikirim, dengan email saja. */
export function permintaanKontak(userId) {
  const u = users.find((x) => x.id === userId)
  if (!u) return null
  return {
    name: u.profile?.full_name ?? null,
    email: u.email,
    ...(u.phone_number ? { phone_number: u.phone_number } : {}),
    member_code: u.profile?.member_code ?? null,
    communication_consent: Boolean(u.profile?.communication_consent),
  }
}

const GALAT = { 429: 'Too many requests', 500: 'Internal server error', 503: 'Service unavailable' }

export function responsKontak(kode, idQontak) {
  if (kode == null) return null
  if (kode === 200) return { id: idQontak, status: 'success' }
  return { status: 'error', message: GALAT[kode] ?? `HTTP ${kode}` }
}

/**
 * Selesaikan kontak yang sudah lewat waktu prosesnya. Hasil true bila status berubah.
 * Qontak dimatikan → 503; Andi (user 103) selalu 500 sebagai contoh gagal; selain itu berhasil.
 */
export function prosesKontak(c) {
  if (c.sync_status !== 'PENDING' || !c.mulai_ulang || Date.now() - c.mulai_ulang < LAMA_PROSES) return false
  const kode = qontakMati() ? 503 : c.user_id === 103 ? 500 : 200
  const sebelumnya = syncLogs.filter((l) => l.user_id === c.user_id).at(-1)
  if (kode === 200) {
    c.qontak_contact_id ??= `QC-${String(c.user_id).padStart(4, '0')}`
    c.last_synced_at = sekarang()
  }
  c.sync_status = kode === 200 ? 'SYNCED' : 'FAILED'
  c.updated_at = sekarang()
  delete c.mulai_ulang
  syncLogs.push({
    id: idBaru(syncLogs), user_id: c.user_id, event_type: 'CONTACT_UPSERT',
    request_payload: permintaanKontak(c.user_id), response_payload: responsKontak(kode, c.qontak_contact_id),
    status_code: kode, retry_count: sebelumnya ? sebelumnya.retry_count + 1 : 0, created_at: sekarang(),
  })
  return true
}

/** Proses semua antrean yang sudah waktunya; dipanggil sebelum status kontak dibaca. */
export function prosesAntrean() {
  let berubah = false
  for (const c of qontakContacts) berubah = prosesKontak(c) || berubah
  if (berubah) simpanDb()
}
