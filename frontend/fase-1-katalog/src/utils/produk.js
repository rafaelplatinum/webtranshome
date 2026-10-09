// Aturan tampilan data produk (CLAUDE.md §7).

export const LABEL_STOK = {
  TERSEDIA: 'Tersedia',
  SISA_STOK: 'Sisa stok',
  PRE_ORDER: 'Pre-order',
  HABIS: 'Habis',
}

// Field internal Master Barang yang TIDAK BOLEH tampil ke pelanggan, meski ikut terkirim dari API:
// HPP, hnet, diskon (disk1/disk2/proyek), stok per lokasi/gudang, konsinyasi, karantina.
const POLA_INTERNAL = /^(hpp|hnet|harga_(beli|pokok|net|proyek)|disk\d*|diskon.*|stok_.*|lokasi_(stok|gudang).*|gudang.*|konsinyasi.*|karantina.*)$/i

const LABEL_KHUSUS = { sni: 'SNI' }

/** "isi_per_dus" → "Isi per dus". */
export function labelSpesifikasi(kunci) {
  if (LABEL_KHUSUS[kunci]) return LABEL_KHUSUS[kunci]
  const teks = String(kunci).replace(/_/g, ' ').trim()
  return teks.charAt(0).toUpperCase() + teks.slice(1)
}

export function kunciInternal(kunci) {
  return POLA_INTERNAL.test(String(kunci))
}

/** Spesifikasi yang aman ditampilkan: [{ kunci, label, nilai }], tanpa field internal dan nilai kosong. */
export function spesifikasiPublik(spesifikasi) {
  return Object.entries(spesifikasi ?? {})
    .filter(([kunci, nilai]) => !kunciInternal(kunci) && nilai != null && String(nilai).trim() !== '')
    .map(([kunci, nilai]) => ({ kunci, label: labelSpesifikasi(kunci), nilai: String(nilai) }))
}
