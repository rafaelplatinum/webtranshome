// Pembantu sesi untuk rute mock: token, profil, izin, dan pemeriksaan akses admin.
import { rolePermissions, users } from '@/services/mock/data/akun'

export const ROLE_ADMIN = ['SUPER_ADMIN', 'ADMIN_KATALOG', 'ADMIN_MEMBERSHIP', 'ADMIN_KONTEN']

// Tombol "Simulasikan sesi habis" di /dev/components memasang penanda ini; login berikutnya menghapusnya.
export const KUNCI_SESI_HABIS = 'transhome.mockSesiHabis'

// Kolom waktu di database tanpa zona (TIMESTAMP), diisi waktu WIB (UTC+7).
export const sekarang = () => new Date(Date.now() + 7 * 3600 * 1000).toISOString().slice(0, 19)

export function profilPengguna(u) {
  return {
    id: u.id, email: u.email, phone_number: u.phone_number,
    full_name: u.profile?.full_name ?? null, member_code: u.profile?.member_code ?? null,
    tier: u.profile?.tier ?? null, communication_consent: u.profile?.communication_consent ?? false,
    member_since: u.profile?.created_at ?? null,
    email_verified: Boolean(u.email_verified_at),
    // Akun yang dibuat lewat Google belum punya password sampai member membuatnya di /akun/keamanan.
    has_password: u.password_hash !== null,
    google_connected: Boolean(u.google_id),
    roles: u.roles,
  }
}

export const izinPengguna = (u) => [...new Set(u.roles.flatMap((r) => rolePermissions[r] ?? []))]

function sesiHabisDisimulasikan() {
  try {
    return localStorage.getItem(KUNCI_SESI_HABIS) === '1'
  } catch {
    return false
  }
}

export function akhiriSimulasiSesiHabis() {
  try {
    localStorage.removeItem(KUNCI_SESI_HABIS)
  } catch {
    // diabaikan
  }
}

export function penggunaDariToken(headers) {
  if (sesiHabisDisimulasikan()) return null
  const auth = headers?.get?.('Authorization') ?? headers?.Authorization ?? ''
  const id = Number(String(auth).replace('Bearer mock-token-', ''))
  return users.find((u) => u.id === id && u.is_active) ?? null
}

export const belumMasuk = { status: 401, data: { message: 'Sesi kamu sudah habis. Silakan masuk lagi.' } }
export const ditolak = (pesan = 'Kamu tidak punya akses untuk aksi ini.') => ({ status: 403, data: { message: pesan } })

/**
 * Pemeriksaan endpoint /admin/*: 401 tanpa sesi, 403 bukan admin atau tidak punya izin.
 * `izin` boleh satu kode atau daftar kode (cukup punya salah satu).
 */
export function cekAdmin(headers, izin = null) {
  const u = penggunaDariToken(headers)
  if (!u) return { gagal: belumMasuk }
  if (!u.roles.some((r) => ROLE_ADMIN.includes(r))) return { gagal: ditolak('Akun ini tidak punya akses admin.') }
  const syarat = [izin].flat().filter(Boolean)
  if (syarat.length && !syarat.some((k) => izinPengguna(u).includes(k))) return { gagal: ditolak() }
  return { u }
}

export const namaPengguna = (id) => {
  const u = users.find((x) => x.id === id)
  return u?.profile?.full_name ?? u?.email ?? 'Admin'
}
