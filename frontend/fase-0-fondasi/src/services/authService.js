import http from './http'

// Bentuk sesi dari API: { token, user: { id, email, phone_number, full_name, member_code, roles, ... }, permissions }.
// Semua aksi di sini memakai `tanpaToast`: pesan gagal (401/403/409/422/429) ditampilkan di form, bukan toast.
const TANPA_TOAST = { tanpaToast: true }
const PAKAI_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

/** Login pelanggan (email + password). 403 `code: EMAIL_BELUM_VERIFIKASI` bila email belum diverifikasi. */
export async function login({ email, password }) {
  const { data } = await http.post('/auth/login', { email, password }, TANPA_TOAST)
  return data.data
}

/** Login admin di /admin/masuk; akun tanpa role admin ditolak backend (403). */
export async function loginAdmin({ email, password }) {
  const { data } = await http.post('/admin/auth/login', { email, password }, TANPA_TOAST)
  return data.data
}

/** Keluar. Gagal jaringan tidak menghalangi keluar secara lokal (lihat useAuth). */
export async function logout() {
  await http.post('/auth/logout', null, TANPA_TOAST)
}

/** Profil dan izin user yang sedang masuk: { user, permissions }. */
export async function getMe() {
  const { data } = await http.get('/me')
  return data.data
}

/** Daftar member: { full_name, email, password, phone_number|null, communication_consent }. 409 email/HP dipakai. */
export async function register(payload) {
  const { data } = await http.post('/auth/register', payload, TANPA_TOAST)
  return data.data
}

/** Kirim ulang link verifikasi. Jawaban selalu sama, terdaftar atau tidak. */
export async function resendVerification(email) {
  await http.post('/auth/email/resend', { email }, TANPA_TOAST)
}

/** Token dari link email. 422 `code: TOKEN_TIDAK_VALID` bila salah atau sudah dipakai. */
export async function verifyEmail(token) {
  const { data } = await http.post('/auth/email/verify', { token }, TANPA_TOAST)
  return data.data
}

/** Minta link reset password. Jawaban selalu sama, terdaftar atau tidak. */
export async function forgotPassword(email) {
  await http.post('/auth/forgot-password', { email }, TANPA_TOAST)
}

/** 422 `code: TOKEN_TIDAK_VALID` bila link salah, sudah dipakai, atau kedaluwarsa. */
export async function resetPassword({ token, password }) {
  await http.post('/auth/reset-password', { token, password }, TANPA_TOAST)
}

/**
 * Alamat untuk memulai "Masuk/Daftar dengan Google" (pindah halaman penuh). OAuth dikerjakan backend,
 * lalu backend kembali ke `returnTo` dengan `#token=…` (berhasil) atau `?error=…` (batal/gagal).
 * Mode mock: halaman simulasi /dev/google (hanya saat development). null = tidak tersedia.
 */
export function urlGoogle(returnTo) {
  const tujuan = encodeURIComponent(returnTo)
  if (PAKAI_MOCK) return import.meta.env.DEV ? `/dev/google?return_to=${tujuan}` : null
  return `${import.meta.env.VITE_API_BASE_URL}/auth/google/redirect?return_to=${tujuan}`
}

/** Khusus mode mock: akun Google contoh → { token, new }. Tidak ada di API asli. */
export async function googleMock({ email, full_name }) {
  const { data } = await http.post('/auth/google/mock', { email, full_name }, TANPA_TOAST)
  return data.data
}
