import { normalisasiNomorHp } from '@/utils/format'

// Validasi form di browser (backend tetap memvalidasi ulang).

export const PASSWORD_MIN = 8

export function emailValid(teks) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(teks ?? '').trim())
}

/** Nomor HP Indonesia: boleh ditulis 08…, +62…, atau 62…; dinormalisasi ke 62… sebelum dicek. */
export function nomorHpValid(teks) {
  return /^628\d{7,11}$/.test(normalisasiNomorHp(teks))
}

export const PESAN = {
  emailKosong: 'Email wajib diisi',
  emailFormat: 'Format email belum benar, contoh nama@email.com',
  passwordKosong: 'Password wajib diisi',
  passwordPendek: `Password minimal ${PASSWORD_MIN} karakter`,
  passwordBeda: 'Password tidak sama',
  hpFormat: 'Format nomor HP belum benar, contoh 0812 3456 7890',
  namaKosong: 'Nama lengkap wajib diisi',
}

/** Pesan galat email, atau '' bila valid. */
export function cekEmail(teks) {
  if (!String(teks ?? '').trim()) return PESAN.emailKosong
  return emailValid(teks) ? '' : PESAN.emailFormat
}
