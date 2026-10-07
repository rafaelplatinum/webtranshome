// Rute API tiruan halaman Sync Qontak (Fase 6). Retry satu kontak & status per kontak ada di routesAdmin.js (Fase 2).
// TODO: ganti ke API begitu endpoint /admin/qontak/* tersedia.
import { qontakContacts, syncLogs } from '@/services/mock/data/admin'
import { users } from '@/services/mock/data/akun'
import { simpanDb } from '@/services/mock/db'
import { antrekan, permintaanKontak, prosesAntrean, responsKontak } from '@/services/mock/antreanQontak'
import { bentukKontak, catat, halaman } from '@/services/mock/routesAdmin'
import { cekAdmin } from '@/services/mock/sesi'

const STATUS = ['PENDING', 'SYNCED', 'FAILED']
const huruf = (teks) => String(teks ?? '').trim().toLowerCase()

function bentukLengkap(c) {
  const u = users.find((x) => x.id === c.user_id)
  return {
    ...bentukKontak(c), qontak_contact_id: c.qontak_contact_id, email: u?.email ?? null, phone_number: u?.phone_number ?? null,
    member_code: u?.profile?.member_code ?? null, created_at: c.created_at, updated_at: c.updated_at,
  }
}

function cocok(c, q) {
  if (!q) return true
  const u = users.find((x) => x.id === c.user_id)
  return [u?.email, u?.profile?.full_name, u?.profile?.member_code, u?.phone_number].some((v) => huruf(v).includes(q))
}

function daftarKontak({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'qontak.retry')
  if (gagal) return gagal
  prosesAntrean()
  const q = huruf(query.q)
  const cocokCari = qontakContacts.filter((c) => cocok(c, q))
  const hasil = cocokCari
    .filter((c) => !STATUS.includes(query.status) || c.sync_status === query.status)
    // Urut terdaftar terbaru: stabil walau status berubah saat di-retry (baris tidak melompat).
    .sort((a, b) => String(b.created_at).localeCompare(String(a.created_at)) || b.id - a.id)
  const { data, meta } = halaman(hasil, query)
  // Jumlah per status mengikuti pencarian, tidak mengikuti filter status (untuk angka di tab).
  const counts = Object.fromEntries(STATUS.map((s) => [s, cocokCari.filter((c) => c.sync_status === s).length]))
  return { data: { data: data.map(bentukLengkap), meta: { ...meta, counts: { ALL: cocokCari.length, ...counts } } } }
}

function ulangSemuaGagal({ headers }) {
  const { u, gagal } = cekAdmin(headers, 'qontak.retry')
  if (gagal) return gagal
  prosesAntrean()
  const daftar = qontakContacts.filter((c) => c.sync_status === 'FAILED')
  daftar.forEach(antrekan)
  if (daftar.length) {
    catat(u, 'qontak.retry_failed', 'Sync Qontak', { sesudah: { count: daftar.length, sync_status: 'PENDING' } })
    simpanDb()
  }
  return { status: daftar.length ? 202 : 200, data: { data: { count: daftar.length } } }
}

function logKontak({ headers, params }) {
  const { gagal } = cekAdmin(headers, 'qontak.retry')
  if (gagal) return gagal
  const c = qontakContacts.find((x) => String(x.id) === String(params.id))
  if (!c) return { status: 404, data: { message: 'Kontak tidak ditemukan.' } }
  prosesAntrean()
  // Log contoh lama belum menyimpan isi permintaan/respons: dibentuk dari data member dan kode status.
  const data = syncLogs
    .filter((l) => l.user_id === c.user_id)
    .sort((a, b) => b.created_at.localeCompare(a.created_at) || b.id - a.id)
    .map((l) => ({
      id: l.id, event_type: l.event_type, status_code: l.status_code, retry_count: l.retry_count, created_at: l.created_at,
      request_payload: l.request_payload ?? permintaanKontak(l.user_id),
      response_payload: l.response_payload ?? responsKontak(l.status_code, c.qontak_contact_id),
    }))
  return { data: { data, meta: { contact: bentukLengkap(c) } } }
}

export const ruteQontak = [
  ['get', '/admin/qontak/contacts', daftarKontak],
  ['post', '/admin/qontak/retry-failed', ulangSemuaGagal],
  ['get', '/admin/qontak/contacts/:id/logs', logKontak],
]
