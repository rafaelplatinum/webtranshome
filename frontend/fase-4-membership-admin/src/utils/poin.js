// Transaksi poin di panel admin (Fase 4). Label untuk member ada di components/akun/DaftarPoin.vue.

// Jenis dibedakan dengan label teks, bukan warna saja.
export const JENIS_TRANSAKSI = {
  EARN: { label: 'Belanja', kelas: 'bg-success-soft text-success-ink' },
  REDEEM: { label: 'Tukar poin', kelas: 'bg-primary-soft text-primary-hover' },
  ADJUST: { label: 'Penyesuaian', kelas: 'bg-info-soft text-info-ink' },
  EXPIRE: { label: 'Kedaluwarsa', kelas: 'bg-subtle text-muted' },
}

/** 4 -> "+4", -10 -> "−10" (tanda minus tipografis). */
export const tandaPoin = (n) => (n > 0 ? `+${n}` : `−${Math.abs(n)}`)

/** Poin dari belanja: total ÷ point_ratio_rupiah, dibulatkan ke bawah. Rasio dari pengaturan toko, bukan hardcode. */
export function hitungPoinBelanja(total, rasio) {
  return rasio > 0 ? Math.floor(total / rasio) : null
}

/** Isian poin "+3", "3", "-2", "−2" -> angka bulat; selain itu NaN. */
export function parsePoin(teks) {
  const isi = String(teks ?? '').trim().replace('−', '-').replace(/\s+/g, '')
  return /^[+-]?\d+$/.test(isi) ? Number(isi) : NaN
}

/** Nomor nota yang sama dicek tanpa beda huruf besar/kecil dan spasi di ujung. */
export const rapikanNota = (teks) => String(teks ?? '').trim().toUpperCase()
