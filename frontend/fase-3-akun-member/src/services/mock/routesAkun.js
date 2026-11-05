// Rute API tiruan untuk akun (Fase 3): masuk, daftar, verifikasi email, lupa/reset password, Google, dan /me.
// TODO: ganti ke API begitu endpoint /auth/* dan /me/* tersedia di backend.
import { normalisasiNomorHp } from '@/utils/format'
import { MOCK_PASSWORD, emailKeluar, otpKeluar, pointTransactions, tokenReset, tokenVerifikasi, users } from '@/services/mock/data/akun'
import { qontakContacts } from '@/services/mock/data/admin'
import { idBaru, simpanDb } from '@/services/mock/db'
import { antrekan } from '@/services/mock/antreanQontak'
import { ROLE_ADMIN, akhiriSimulasiSesiHabis, belumMasuk, izinPengguna, penggunaDariToken, profilPengguna, sekarang } from '@/services/mock/sesi'

const POLA_EMAIL = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const POLA_HP = /^628\d{7,11}$/
const huruf = (teks) => String(teks ?? '').trim().toLowerCase()
const galatIsian = (errors) => ({ status: 422, data: { message: 'Periksa kembali isian yang ditandai.', errors } })
const tokenTidakValid = (pesan) => ({ status: 422, data: { message: pesan, code: 'TOKEN_TIDAK_VALID' } })
const cariEmail = (email) => users.find((u) => huruf(u.email) === huruf(email))
const tokenAcak = () => [...crypto.getRandomValues(new Uint8Array(16))].map((b) => b.toString(16).padStart(2, '0')).join('')
const tokenKuat = () => [...crypto.getRandomValues(new Uint8Array(32))].map((b) => b.toString(16).padStart(2, '0')).join('')
const otpPassword = []
const tokenResetOtp = []
const hashTeks = async (value) => [...new Uint8Array(await crypto.subtle.digest('SHA-256', new TextEncoder().encode(value)))].map((b) => b.toString(16).padStart(2, '0')).join('')

// Password akun baru disimpan sebagai hash (bukan teks asli), meski hanya mock.
export async function hashPassword(teks) {
  const hasil = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(`transhome-mock:${teks}`))
  return [...new Uint8Array(hasil)].map((b) => b.toString(16).padStart(2, '0')).join('')
}

async function cocokPassword(u, teks) {
  if (u.password_hash === null) return false
  if (u.password_hash === undefined) return teks === MOCK_PASSWORD
  return u.password_hash === (await hashPassword(teks))
}

function kirimEmail(kepada, jenis, link) {
  const subjek = jenis === 'verifikasi' ? 'Verifikasi email akun Trans Family' : 'Atur ulang password Transhome'
  emailKeluar.unshift({ id: idBaru(emailKeluar), to: kepada, jenis, subject: subjek, link, created_at: sekarang() })
}

function kirimVerifikasi(u) {
  const token = tokenAcak()
  tokenVerifikasi.push({ token, user_id: u.id, used_at: null, created_at: sekarang() })
  kirimEmail(u.email, 'verifikasi', `/verifikasi?token=${token}`)
}

function buatMember({ email, fullName, phone = null, passwordHash, consent = false, verified = false, googleId = null }) {
  const nomorMember = users.reduce((maks, u) => Math.max(maks, Number(u.profile?.member_code?.replace('TF-', '')) || 0), 0) + 1
  const idProfil = users.reduce((maks, u) => Math.max(maks, u.profile?.id ?? 0), 0) + 1
  const u = {
    id: idBaru(users), email, phone_number: phone, is_active: true, roles: ['CUSTOMER'], password_hash: passwordHash,
    email_verified_at: verified ? sekarang() : null, google_id: googleId,
    profile: {
      id: idProfil, member_code: `TF-${String(nomorMember).padStart(6, '0')}`, full_name: fullName, tier: 'TAHAP_1',
      communication_consent: consent, consent_at: consent ? sekarang() : null, created_at: sekarang(),
    },
  }
  users.push(u)
  // Sync kontak ke Qontak berjalan di background; registrasi tidak menunggu dan tetap berhasil walau Qontak mati.
  const kontak = { id: idBaru(qontakContacts), user_id: u.id, qontak_contact_id: null, sync_status: 'PENDING', last_synced_at: null, created_at: sekarang(), updated_at: sekarang() }
  antrekan(kontak)
  qontakContacts.push(kontak)
  return u
}

const sesi = (u) => ({ token: `mock-token-${u.id}`, user: profilPengguna(u), permissions: izinPengguna(u) })

// ---------- masuk & daftar ----------

async function masuk({ body }, { admin = false } = {}) {
  const u = cariEmail(body.email)
  if (!u || !(await cocokPassword(u, String(body.password ?? '')))) return { status: 401, data: { message: 'Email atau password salah.' } }
  if (!u.is_active) return { status: 403, data: { message: 'Akun ini dinonaktifkan.' } }
  if (admin && !u.roles.some((r) => ROLE_ADMIN.includes(r))) return { status: 403, data: { message: 'Akun ini tidak punya akses admin.' } }
  if (!admin && !u.email_verified_at) return { status: 403, data: { message: 'Email kamu belum diverifikasi.', code: 'EMAIL_BELUM_VERIFIKASI' } }
  akhiriSimulasiSesiHabis()
  u.last_login_at = sekarang()
  simpanDb()
  return { data: { data: sesi(u) } }
}

async function daftar({ body }) {
  const nama = String(body.full_name ?? '').trim()
  const email = huruf(body.email)
  const password = String(body.password ?? '')
  const hp = body.phone_number ? normalisasiNomorHp(body.phone_number) : null
  const errors = {}
  if (!nama) errors.full_name = ['Nama lengkap wajib diisi']
  if (!email) errors.email = ['Email wajib diisi']
  else if (!POLA_EMAIL.test(email)) errors.email = ['Format email belum benar']
  if (password.length < 8) errors.password = ['Password minimal 8 karakter']
  if (hp && !POLA_HP.test(hp)) errors.phone_number = ['Format nomor HP belum benar']
  if (Object.keys(errors).length) return galatIsian(errors)
  if (cariEmail(email)) return { status: 409, data: { message: 'Email sudah terdaftar.', errors: { email: ['Email sudah terdaftar'] } } }
  if (hp && users.some((u) => u.phone_number === hp)) {
    return { status: 409, data: { message: 'Nomor HP sudah dipakai akun lain.', errors: { phone_number: ['Nomor HP sudah dipakai akun lain'] } } }
  }
  const u = buatMember({ email, fullName: nama, phone: hp, passwordHash: await hashPassword(password), consent: Boolean(body.communication_consent) })
  kirimVerifikasi(u)
  simpanDb()
  return { status: 201, data: { data: { email: u.email } } }
}

// Jawaban selalu sama, terdaftar atau tidak, supaya email orang lain tidak bisa ditebak.
function kirimUlangVerifikasi({ body }) {
  const u = cariEmail(body.email)
  if (u && !u.email_verified_at) {
    kirimVerifikasi(u)
    simpanDb()
  }
  return { data: { message: 'Jika email terdaftar dan belum diverifikasi, link baru sudah dikirim.' } }
}

function verifikasiEmail({ body }) {
  const t = tokenVerifikasi.find((x) => x.token === body.token && !x.used_at)
  const u = t && users.find((x) => x.id === t.user_id)
  if (!u) return tokenTidakValid('Link verifikasi tidak valid atau sudah dipakai.')
  t.used_at = sekarang()
  u.email_verified_at ??= sekarang()
  simpanDb()
  return { data: { data: { email: u.email } } }
}

// ---------- lupa & reset password ----------

/** Link reset sekali pakai, berlaku 1 jam. Dipakai juga "Reset password" di /admin/pengguna. */
export function kirimLinkReset(u) {
  const token = tokenAcak()
  tokenReset.push({ token, user_id: u.id, used_at: null, expires_at: Date.now() + 60 * 60 * 1000 })
  kirimEmail(u.email, 'reset', `/reset-password?token=${token}`)
}

function lupaPassword({ body }) {
  const u = cariEmail(body.email)
  if (u?.is_active) {
    kirimLinkReset(u)
    simpanDb()
  }
  return { data: { message: 'Jika email terdaftar, link reset sudah dikirim.' } }
}

async function resetPassword({ body }) {
  const t = tokenReset.find((x) => x.token === body.token)
  if (!t || t.used_at || t.expires_at < Date.now()) return tokenTidakValid('Link reset tidak valid, sudah dipakai, atau sudah kedaluwarsa.')
  if (String(body.password ?? '').length < 8) return galatIsian({ password: ['Password minimal 8 karakter'] })
  const u = users.find((x) => x.id === t.user_id)
  u.password_hash = await hashPassword(body.password)
  // Bisa membuka email berarti pemilik email: akun yang belum terverifikasi ikut terverifikasi.
  u.email_verified_at ??= sekarang()
  t.used_at = sekarang()
  simpanDb()
  return { data: { message: 'Password diperbarui.' } }
}

function normalTujuanOtp(channel, raw) {
  if (channel === 'email') return huruf(raw)
  if (channel === 'whatsapp') return normalisasiNomorHp(raw)
  return ''
}

function lupaPasswordOtp({ body }) {
  const channel = body.channel
  const destination = normalTujuanOtp(channel, String(body.destination ?? ''))
  if (!destination) return galatIsian({ destination: ['Alamat email atau nomor WhatsApp belum benar'] })
  const u = users.find((x) => x.is_active && (channel === 'email' ? huruf(x.email) === destination && x.email_verified_at : x.phone_number === destination && x.phone_verified_at))
  const now = Date.now()
  otpPassword.filter((x) => x.destination === destination && x.channel === channel && !x.used_at).forEach((x) => { x.used_at = sekarang() })
  if (u) {
    const code = String(crypto.getRandomValues(new Uint32Array(1))[0] % 1000000).padStart(6, '0')
    otpPassword.push({ channel, destination, user_id: u.id, code, attempts: 0, created_at: now, expires_at: now + 5 * 60 * 1000, used_at: null })
    otpKeluar.unshift({ id: idBaru(otpKeluar), channel, to: destination, code, created_at: sekarang(), expires_at: now + 5 * 60 * 1000 })
  } else otpPassword.push({ channel, destination, user_id: null, code: '', attempts: 0, created_at: now, expires_at: now + 5 * 60 * 1000, used_at: null })
  return { data: { message: 'Jika data terdaftar, kode verifikasi sudah dikirim.' } }
}

async function verifikasiPasswordOtp({ body }) {
  const channel = body.channel
  const destination = normalTujuanOtp(channel, String(body.destination ?? ''))
  const code = String(body.code ?? '')
  const challenge = otpPassword.filter((x) => x.destination === destination && x.channel === channel && !x.purpose).at(-1)
  if (!challenge || !/^\d{6}$/.test(code) || !challenge.user_id || challenge.used_at || challenge.expires_at <= Date.now() || challenge.attempts >= 5) {
    return { status: 422, data: { message: 'Kode tidak valid atau sudah kedaluwarsa.' } }
  }
  if (code !== challenge.code) {
    challenge.attempts += 1
    if (challenge.attempts >= 5) challenge.used_at = sekarang()
    return { status: 422, data: { message: 'Kode tidak valid atau sudah kedaluwarsa.' } }
  }
  challenge.used_at = sekarang()
  const token = tokenKuat()
  tokenResetOtp.push({ token_hash: await hashTeks(token), user_id: challenge.user_id, expires_at: Date.now() + 10 * 60 * 1000, used_at: null })
  return { data: { data: { reset_token: token } } }
}

async function resetPasswordOtp({ body }) {
  const tokenHash = await hashTeks(String(body.token ?? ''))
  const record = tokenResetOtp.find((x) => x.token_hash === tokenHash && !x.used_at && x.expires_at > Date.now())
  if (!record) return tokenTidakValid('TOKEN_TIDAK_VALID')
  const password = String(body.password ?? '')
  if (password.length < 8 || password.length > 72) return galatIsian({ password: ['Password minimal 8 karakter'] })
  const u = users.find((x) => x.id === record.user_id)
  if (!u || !u.is_active) return tokenTidakValid('TOKEN_TIDAK_VALID')
  u.password_hash = await hashPassword(password)
  u.email_verified_at ??= sekarang()
  u.auth_version = (u.auth_version ?? 0) + 1
  record.used_at = sekarang()
  simpanDb()
  return { data: { message: 'Password diperbarui.' } }
}

function mintaOtpHp({ headers }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  if (!u.phone_number) return { status: 422, data: { message: 'Tambahkan nomor WhatsApp valid terlebih dahulu.' } }
  otpPassword.filter((x) => x.destination === u.phone_number && x.channel === 'whatsapp' && x.purpose === 'phone_verify' && !x.used_at).forEach((x) => { x.used_at = sekarang() })
  const code = String(crypto.getRandomValues(new Uint32Array(1))[0] % 1000000).padStart(6, '0')
  otpPassword.push({ channel: 'whatsapp', destination: u.phone_number, user_id: u.id, purpose: 'phone_verify', code, attempts: 0, created_at: Date.now(), expires_at: Date.now() + 5 * 60 * 1000, used_at: null })
  otpKeluar.unshift({ id: idBaru(otpKeluar), channel: 'whatsapp', to: u.phone_number, code, created_at: sekarang(), expires_at: Date.now() + 5 * 60 * 1000 })
  return { data: { message: 'Jika nomor valid, kode verifikasi sudah dikirim.' } }
}

function verifikasiOtpHp({ headers, body }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  const challenge = otpPassword.filter((x) => x.destination === u.phone_number && x.channel === 'whatsapp' && x.purpose === 'phone_verify').at(-1)
  if (!challenge || challenge.used_at || challenge.expires_at <= Date.now() || challenge.attempts >= 5) return { status: 422, data: { message: 'Kode tidak valid atau sudah kedaluwarsa.' } }
  if (String(body.code ?? '') !== challenge.code) {
    challenge.attempts += 1
    if (challenge.attempts >= 5) challenge.used_at = sekarang()
    return { status: 422, data: { message: 'Kode tidak valid atau sudah kedaluwarsa.' } }
  }
  challenge.used_at = sekarang()
  u.phone_verified_at = sekarang()
  simpanDb()
  return { data: { message: 'Nomor WhatsApp terverifikasi.' } }
}

// ---------- Google (khusus mode mock) ----------
// API asli: GET /auth/google/redirect → OAuth → kembali ke return_to#token=… (lihat lampiran rencana build).

function googleMock({ body }) {
  const email = huruf(body.email)
  if (!POLA_EMAIL.test(email)) return galatIsian({ email: ['Format email belum benar'] })
  let u = cariEmail(email)
  const baru = !u
  if (!u) u = buatMember({ email, fullName: String(body.full_name ?? '').trim() || email, passwordHash: null, verified: true, googleId: `google-${tokenAcak().slice(0, 12)}` })
  else {
    u.google_id ??= `google-${tokenAcak().slice(0, 12)}`
    u.email_verified_at ??= sekarang()
  }
  if (!u.is_active) return { status: 403, data: { message: 'Akun ini dinonaktifkan.' } }
  akhiriSimulasiSesiHabis()
  u.last_login_at = sekarang()
  simpanDb()
  return { data: { data: { token: `mock-token-${u.id}`, new: baru } } }
}

// ---------- /me ----------

function profilSaya({ headers }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  return { data: { data: { user: profilPengguna(u), permissions: izinPengguna(u) } } }
}

function ubahProfil({ headers, body }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  const nama = String(body.full_name ?? '').trim()
  const hp = body.phone_number ? normalisasiNomorHp(body.phone_number) : null
  const errors = {}
  if (!nama) errors.full_name = ['Nama lengkap wajib diisi']
  if (hp && !POLA_HP.test(hp)) errors.phone_number = ['Format nomor HP belum benar']
  if (Object.keys(errors).length) return galatIsian(errors)
  if (hp && users.some((x) => x.phone_number === hp && x.id !== u.id)) {
    return { status: 409, data: { message: 'Nomor HP sudah dipakai akun lain.', errors: { phone_number: ['Nomor HP sudah dipakai akun lain'] } } }
  }
  u.profile = { ...u.profile, full_name: nama }
  if (u.phone_number !== hp) u.phone_verified_at = null
  u.phone_number = hp
  simpanDb()
  return { data: { data: { user: profilPengguna(u) } } }
}

function ubahPersetujuan({ headers, body }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  const setuju = Boolean(body.communication_consent)
  u.profile = { ...u.profile, communication_consent: setuju, consent_at: setuju ? sekarang() : null }
  simpanDb()
  return { data: { data: { user: profilPengguna(u) } } }
}

async function ubahPassword({ headers, body }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  const sudahPunya = u.password_hash !== null
  const errors = {}
  if (sudahPunya && !(await cocokPassword(u, String(body.old_password ?? '')))) errors.old_password = ['Password lama salah']
  if (String(body.new_password ?? '').length < 8) errors.new_password = ['Password baru minimal 8 karakter']
  if (Object.keys(errors).length) return galatIsian(errors)
  u.password_hash = await hashPassword(body.new_password)
  simpanDb()
  // Backend asli juga mengeluarkan sesi lain milik akun ini.
  return { data: { message: sudahPunya ? 'Password diperbarui.' : 'Password dibuat.' } }
}

function poinSaya({ headers, query }) {
  const u = penggunaDariToken(headers)
  if (!u) return belumMasuk
  const milik = pointTransactions
    .filter((t) => t.member_id === u.profile?.id)
    .sort((a, b) => b.created_at.localeCompare(a.created_at) || b.id - a.id)
  const saldo = milik.reduce((jumlah, t) => jumlah + t.points, 0)
  const riwayat = query.type ? milik.filter((t) => t.type === query.type) : milik
  const perPage = Math.min(Number(query.per_page) || 10, 50)
  const totalPages = Math.max(1, Math.ceil(riwayat.length / perPage))
  const page = Math.min(Math.max(Number(query.page) || 1, 1), totalPages)
  return {
    data: {
      data: riwayat.slice((page - 1) * perPage, page * perPage),
      meta: { balance: saldo, page, per_page: perPage, total: riwayat.length, total_pages: totalPages },
    },
  }
}

// Kotak email tiruan (/dev/email): tidak ada di API asli.
function emailTiruan() {
  return { data: { data: emailKeluar.slice(0, 30) } }
}
function otpTiruan() { return { data: { data: otpKeluar.slice(0, 30) } } }

export const ruteAkun = [
  ['post', '/auth/login', (req) => masuk(req)],
  ['post', '/admin/auth/login', (req) => masuk(req, { admin: true })],
  ['post', '/auth/logout', () => ({ status: 204, data: null })],
  ['post', '/auth/register', daftar],
  ['post', '/auth/email/resend', kirimUlangVerifikasi],
  ['post', '/auth/email/verify', verifikasiEmail],
  ['post', '/auth/forgot-password', lupaPassword],
  ['post', '/auth/password/forgot', lupaPasswordOtp],
  ['post', '/auth/password/verify-otp', verifikasiPasswordOtp],
  ['post', '/auth/password/reset', resetPasswordOtp],
  ['post', '/auth/reset-password', resetPassword],
  ['post', '/auth/google/mock', googleMock],
  ['get', '/me', profilSaya],
  ['patch', '/me', ubahProfil],
  ['patch', '/me/consent', ubahPersetujuan],
  ['post', '/me/password', ubahPassword],
  ['post', '/member/phone/request-otp', mintaOtpHp],
  ['post', '/member/phone/verify', verifikasiOtpHp],
  ['get', '/me/points', poinSaya],
  ['get', '/dev/emails', emailTiruan],
  ['get', '/dev/otps', otpTiruan],
]
