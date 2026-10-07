// Rute API tiruan membership admin (Fase 4): daftar, detail & cari member, nonaktifkan akun, transaksi poin.
// Member dikenali dengan `id` = users.id; point_transactions.member_id tetap merujuk member_profiles.id.
// TODO: ganti ke API begitu endpoint /admin/members dan /admin/points tersedia di backend.
import { formatRupiah, formatTanggal, normalisasiNomorHp } from '@/utils/format'
import { pointTransactions, publicSettings, users } from '@/services/mock/data/akun'
import { qontakContacts } from '@/services/mock/data/admin'
import { idBaru, simpanDb } from '@/services/mock/db'
import { prosesAntrean } from '@/services/mock/antreanQontak'
import { bentukKontak, catat, halaman } from '@/services/mock/routesAdmin'
import { cekAdmin, ditolak, izinPengguna, namaPengguna, sekarang } from '@/services/mock/sesi'

const JENIS = ['EARN', 'REDEEM', 'ADJUST']
const huruf = (teks) => String(teks ?? '').trim().toLowerCase()
const sama = (a, b) => String(a) === String(b)
const galatIsian = (errors) => ({ status: 422, data: { message: 'Periksa kembali isian yang ditandai.', errors } })
const tidakDitemukan = { status: 404, data: { message: 'Member tidak ditemukan.' } }

const semuaMember = () => users.filter((u) => u.roles.includes('CUSTOMER') && u.profile?.member_code)
const cariMember = (id) => semuaMember().find((u) => sama(u.id, id))
const transaksiMember = (u) => pointTransactions.filter((t) => t.member_id === u.profile.id)
const saldo = (u) => transaksiMember(u).reduce((jumlah, t) => jumlah + t.points, 0)
// Antrean sync diproses dulu supaya status yang dibaca sudah yang terbaru.
const kontak = (u) => {
  prosesAntrean()
  return qontakContacts.find((c) => c.user_id === u.id)
}
const terbaru = (a, b) => b.created_at.localeCompare(a.created_at) || b.id - a.id

function ringkas(u) {
  return {
    id: u.id, member_code: u.profile.member_code, full_name: u.profile.full_name, email: u.email, phone_number: u.phone_number,
    tier: u.profile.tier, is_active: u.is_active, email_verified: Boolean(u.email_verified_at),
    communication_consent: u.profile.communication_consent, balance: saldo(u), sync_status: kontak(u)?.sync_status ?? null,
    created_at: u.profile.created_at,
  }
}

function detail(u) {
  const c = kontak(u)
  return {
    ...ringkas(u), consent_at: u.profile.consent_at, google_connected: Boolean(u.google_id), last_login_at: u.last_login_at ?? null,
    qontak: c ? bentukKontak(c) : null,
  }
}

function bentukTransaksi(t) {
  const u = users.find((x) => x.profile?.id === t.member_id)
  return {
    id: t.id, user_id: u?.id ?? null, member_code: u?.profile.member_code ?? null, type: t.type, points: t.points,
    reference_no: t.reference_no, purchase_amount: t.purchase_amount, note: t.note,
    created_by_name: t.created_by ? namaPengguna(t.created_by) : null, created_at: t.created_at,
  }
}

// Cocok sebagian di kode member, email, nama, atau nomor HP (08…, +62…, atau potongan angka).
function cocok(u, q) {
  if ([u.profile.member_code, u.email, u.profile.full_name].some((teks) => huruf(teks).includes(q))) return true
  if (/[a-z]/.test(q) || !u.phone_number) return false
  const angka = q.replace(/\D/g, '')
  return angka.length >= 4 && (u.phone_number.includes(angka) || u.phone_number.includes(normalisasiNomorHp(q)))
}

// ---------- member ----------

function daftar({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'member.view')
  if (gagal) return gagal
  let hasil = semuaMember()
  const q = huruf(query.q)
  if (q) hasil = hasil.filter((u) => cocok(u, q))
  const ya = (kunci) => query[kunci] === '1'
  if (['1', '0'].includes(query.consent)) hasil = hasil.filter((u) => u.profile.communication_consent === ya('consent'))
  if (['1', '0'].includes(query.hp)) hasil = hasil.filter((u) => Boolean(u.phone_number) === ya('hp'))
  if (query.sync) hasil = hasil.filter((u) => kontak(u)?.sync_status === query.sync)
  if (['1', '0'].includes(query.active)) hasil = hasil.filter((u) => u.is_active === ya('active'))
  hasil.sort((a, b) => b.profile.created_at.localeCompare(a.profile.created_at) || b.id - a.id)
  const { data, meta } = halaman(hasil, query)
  return { data: { data: data.map(ringkas), meta } }
}

// Input poin: kode member / email / HP yang persis sama → hanya member itu; selain itu cocok sebagian (maks 10).
function cari({ headers, query }) {
  const { gagal } = cekAdmin(headers, ['member.view', 'point.create'])
  if (gagal) return gagal
  const q = huruf(query.q)
  if (q.length < 3) return galatIsian({ q: ['Isi minimal 3 karakter.'] })
  const semua = semuaMember()
  const hp = /[a-z]/.test(q) ? null : normalisasiNomorHp(q)
  const persis = semua.filter((u) => huruf(u.profile.member_code) === q || huruf(u.email) === q || (hp && u.phone_number === hp))
  const hasil = persis.length ? persis : semua.filter((u) => cocok(u, q)).sort((a, b) => a.profile.full_name.localeCompare(b.profile.full_name, 'id'))
  return { data: { data: hasil.slice(0, 10).map(ringkas), meta: { total: hasil.length } } }
}

function ambil({ headers, params }) {
  const { gagal } = cekAdmin(headers, ['member.view', 'point.create'])
  if (gagal) return gagal
  const u = cariMember(params.id)
  return u ? { data: { data: detail(u) } } : tidakDitemukan
}

// Nonaktif: member tidak bisa masuk (sesi yang berjalan ikut berakhir). Poin & riwayat tetap tersimpan.
function ubah({ headers, params, body }) {
  const { u: admin, gagal } = cekAdmin(headers, 'member.update')
  if (gagal) return gagal
  const u = cariMember(params.id)
  if (!u) return tidakDitemukan
  if (typeof body.is_active !== 'boolean') return galatIsian({ is_active: ['Nilai tidak valid.'] })
  const sebelum = { member_code: u.profile.member_code, is_active: u.is_active }
  u.is_active = body.is_active
  catat(admin, 'member.update', 'Member', { sebelum, sesudah: { member_code: u.profile.member_code, is_active: u.is_active } })
  simpanDb()
  return { data: { data: detail(u) } }
}

// ---------- poin ----------

const cariNota = (ref) => {
  const r = huruf(ref)
  return r ? pointTransactions.find((t) => t.type === 'EARN' && huruf(t.reference_no) === r) : null
}

function pesanNota(t) {
  const kode = bentukTransaksi(t).member_code
  return `Nomor nota ini sudah pernah diinput pada ${formatTanggal(t.created_at)}${kode ? ` (${kode})` : ''}.`
}

function riwayat({ headers, query }) {
  const { gagal } = cekAdmin(headers, ['member.view', 'point.create'])
  if (gagal) return gagal
  const u = cariMember(query.user_id)
  if (!u) return tidakDitemukan
  let hasil = transaksiMember(u).sort(terbaru)
  if (JENIS.includes(query.type)) hasil = hasil.filter((t) => t.type === query.type)
  const { data, meta } = halaman(hasil, query, 10)
  return { data: { data: data.map(bentukTransaksi), meta: { ...meta, balance: saldo(u) } } }
}

function cekNota({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'point.create')
  if (gagal) return gagal
  const t = cariNota(query.ref)
  return { data: { data: t ? { available: false, used_at: t.created_at, member_code: bentukTransaksi(t).member_code } : { available: true } } }
}

// EARN: poin dihitung di server dari total belanja ÷ point_ratio_rupiah (dibulatkan ke bawah).
// REDEEM: `points` positif, disimpan negatif. ADJUST (izin point.adjust): plus/minus. Saldo tidak boleh minus.
function catatPoin({ headers, body }) {
  const { u: admin, gagal } = cekAdmin(headers, 'point.create')
  if (gagal) return gagal
  const jenis = body.type
  if (!JENIS.includes(jenis)) return galatIsian({ type: ['Pilih jenis transaksi.'] })
  if (jenis === 'ADJUST' && !izinPengguna(admin).includes('point.adjust')) return ditolak()
  const u = cariMember(body.user_id)
  if (!u) return tidakDitemukan
  if (!u.is_active) return galatIsian({ user_id: ['Akun member ini nonaktif. Aktifkan dulu sebelum mencatat poin.'] })

  const saldoKini = saldo(u)
  const catatan = String(body.note ?? '').trim()
  const errors = {}
  let poin = 0
  let nota = null
  let belanja = null
  if (jenis === 'EARN') {
    const rasio = Number(publicSettings.point_ratio_rupiah)
    nota = String(body.reference_no ?? '').trim().toUpperCase()
    belanja = Number(body.purchase_amount)
    if (!nota) errors.reference_no = ['Nomor nota wajib diisi.']
    if (!(belanja > 0)) errors.purchase_amount = ['Total belanja wajib diisi.']
    else {
      poin = Math.floor(belanja / rasio)
      if (poin < 1) errors.purchase_amount = [`Total belanja di bawah ${formatRupiah(rasio)}, jadi belum dapat poin.`]
    }
  } else {
    const n = Number(body.points)
    if (!Number.isInteger(n) || n === 0 || (jenis === 'REDEEM' && n < 0)) {
      errors.points = [jenis === 'REDEEM' ? 'Isi jumlah poin yang ditukar.' : 'Isi jumlah poin, boleh plus atau minus.']
    }
    if (!catatan) errors.note = [jenis === 'REDEEM' ? 'Catatan wajib diisi.' : 'Alasan wajib diisi.']
    poin = jenis === 'REDEEM' ? -n : n
    if (!errors.points && saldoKini + poin < 0) {
      errors.points = [jenis === 'REDEEM' ? `Saldo poin tidak cukup. Saldo saat ini ${saldoKini} poin.` : `Saldo tidak boleh kurang dari 0. Saldo saat ini ${saldoKini} poin.`]
    }
  }
  if (Object.keys(errors).length) return galatIsian(errors)

  // Nomor nota unik (UNIQUE (type, reference_no)): ditolak walau pengecekan di browser terlewat.
  const lama = jenis === 'EARN' && cariNota(nota)
  if (lama) return { status: 409, data: { message: pesanNota(lama), errors: { reference_no: [pesanNota(lama)] } } }

  const t = {
    id: idBaru(pointTransactions), member_id: u.profile.id, type: jenis, points: poin, reference_no: nota,
    purchase_amount: jenis === 'EARN' ? belanja.toFixed(2) : null, note: catatan || null, created_by: admin.id, created_at: sekarang(),
  }
  pointTransactions.push(t)
  catat(admin, jenis === 'ADJUST' ? 'point.adjust' : 'point.create', 'Poin', {
    sesudah: { member_code: u.profile.member_code, type: jenis, points: poin, reference_no: nota, purchase_amount: t.purchase_amount, note: t.note },
  })
  simpanDb()
  return { status: 201, data: { data: { transaction: bentukTransaksi(t), balance: saldo(u) } } }
}

// Urutan penting: path tetap (search, check-reference) sebelum path ber-parameter.
export const ruteMember = [
  ['get', '/admin/members/search', cari],
  ['get', '/admin/members', daftar],
  ['get', '/admin/members/:id', ambil],
  ['patch', '/admin/members/:id', ubah],
  ['get', '/admin/points/check-reference', cekNota],
  ['get', '/admin/points', riwayat],
  ['post', '/admin/points', catatPoin],
]
