// Format tampilan sesuai CLAUDE.md §9.

const rupiah = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  maximumFractionDigits: 0,
})

const tanggal = new Intl.DateTimeFormat('id-ID', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
  timeZone: 'Asia/Jakarta',
})

const jam = new Intl.DateTimeFormat('id-ID', { hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Jakarta' })

const angka = new Intl.NumberFormat('id-ID')

/** 125000 -> "Rp 125.000". Nilai DECIMAL dari API boleh berupa string. */
export function formatRupiah(nilai) {
  return rupiah.format(Number(nilai) || 0)
}

/** 4250000 -> "4.250.000" (tanpa "Rp"), mis. isian total belanja atau saldo poin. */
export function formatAngka(nilai) {
  return angka.format(Number(nilai) || 0)
}

/**
 * Isian uang -> rupiah utuh: "Rp 4.250.000", "4250000", "4.250.000,50", "4250000.00" -> 4250000.
 * Titik/koma diikuti 1–2 angka di akhir dianggap sen dan dibuang; pemisah ribuan selalu 3 angka.
 */
export function parseRupiah(teks) {
  const utuh = String(teks ?? '').trim().replace(/[.,]\d{1,2}$/, '')
  const digit = utuh.replace(/\D/g, '')
  return digit ? Number(digit) : 0
}

/** Harga selalu bersama satuan jual: (85000, 'sak') -> "Rp 85.000 per sak". */
export function formatHargaPerSatuan(nilai, satuan) {
  return satuan ? `${formatRupiah(nilai)} per ${satuan}` : formatRupiah(nilai)
}

/**
 * Kolom TIMESTAMP dari API tidak punya zona dan berisi waktu WIB, jadi dibaca sebagai +07:00
 * (bukan zona browser). String yang sudah punya zona (Z / +hh:mm) dipakai apa adanya.
 */
export function keWaktu(nilai) {
  if (nilai instanceof Date) return nilai
  const teks = String(nilai ?? '')
  const berzona = /(Z|[+-]\d{2}:?\d{2})$/.test(teks) || !teks.includes('T')
  return new Date(berzona ? teks : `${teks}+07:00`)
}

/** Date atau string ISO -> "6 Okt 2026" (zona Asia/Jakarta). */
export function formatTanggal(nilai) {
  if (!nilai) return ''
  const waktu = keWaktu(nilai)
  return Number.isNaN(waktu.getTime()) ? '' : tanggal.format(waktu)
}

/** "6 Okt 2026, 19.22" (zona Asia/Jakarta), mis. login terakhir di detail member. */
export function formatTanggalJam(nilai) {
  if (!nilai) return ''
  const waktu = keWaktu(nilai)
  return Number.isNaN(waktu.getTime()) ? '' : `${tanggal.format(waktu)}, ${jam.format(waktu)}`
}

/** "Budi Santoso" -> "BS", "Andi" -> "A": isi avatar bulat. */
export function inisialNama(nama) {
  return String(nama ?? '')
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((kata) => kata.charAt(0).toUpperCase())
    .join('')
}

/** "baru saja", "10 menit lalu", "3 jam lalu", "kemarin", lalu tanggal biasa ("6 Okt 2026"). */
export function formatWaktuRelatif(nilai, kini = Date.now()) {
  if (!nilai) return ''
  const waktu = keWaktu(nilai)
  if (Number.isNaN(waktu.getTime())) return ''
  const menit = Math.floor((kini - waktu.getTime()) / 60000)
  if (menit < 1) return 'baru saja'
  if (menit < 60) return `${menit} menit lalu`
  if (menit < 24 * 60) return `${Math.floor(menit / 60)} jam lalu`
  if (menit < 48 * 60) return 'kemarin'
  return formatTanggal(waktu)
}

/** Selisih hari kalender (WIB) dari hari ini ke tanggal itu: 0 = hari ini, 3 = tiga hari lagi. */
export function selisihHari(nilai, kini = Date.now()) {
  const hari = (ms) => Math.floor((ms + 7 * 3600 * 1000) / 86400000)
  return hari(keWaktu(nilai).getTime()) - hari(kini)
}

/** "6281234567890" (dari API) -> "0812 3456 7890" untuk ditampilkan. */
export function formatNomorHp(nilai) {
  const angka = normalisasiNomorHp(nilai)
  if (!angka.startsWith('62')) return String(nilai ?? '')
  return `0${angka.slice(2)}`.replace(/(\d{4})(?=\d)/g, '$1 ')
}

/** Input "0812-3456 789" atau "+62 812..." -> "628123456789" sebelum dikirim ke API. */
export function normalisasiNomorHp(input) {
  const angka = String(input ?? '').replace(/\D/g, '')
  if (angka.startsWith('62')) return angka
  if (angka.startsWith('0')) return `62${angka.slice(1)}`
  if (angka.startsWith('8')) return `62${angka}`
  return angka
}
