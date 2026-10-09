// Mengubah error axios menjadi pesan untuk pengguna dan error per field.
// Format error backend yang disepakati: { message, errors: { field: ['pesan'] } }.

/** Error per field untuk form, contoh { sku: 'SKU sudah dipakai produk lain' }. */
export function errorField(error) {
  const errors = error?.response?.data?.errors
  if (!errors || typeof errors !== 'object') return {}
  return Object.fromEntries(
    Object.entries(errors).map(([field, pesan]) => [field, Array.isArray(pesan) ? pesan[0] : String(pesan)]),
  )
}

/** Satu kalimat yang bisa langsung ditampilkan di ErrorState, toast, atau di atas form. */
export function pesanError(error) {
  if (!error?.response) return 'Tidak bisa terhubung ke server. Periksa koneksi, lalu coba lagi.'
  const { status, data } = error.response
  if (status === 401) return 'Sesi kamu sudah habis. Silakan masuk lagi.'
  if (status === 403) return data?.message || 'Kamu tidak punya akses untuk aksi ini.'
  if (status === 404) return data?.message || 'Data tidak ditemukan.'
  if (status === 429) return 'Terlalu banyak percobaan. Coba lagi beberapa saat.'
  if (status >= 500) return 'Terjadi gangguan di server. Coba lagi.'
  return data?.message || 'Permintaan tidak bisa diproses. Periksa isian, lalu coba lagi.'
}
