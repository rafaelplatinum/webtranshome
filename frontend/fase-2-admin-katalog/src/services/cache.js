// Cache data publik (daftar katalog, pohon kategori, ruangan) didaftarkan di sini dan dibersihkan
// setelah admin mengubah data, supaya halaman publik di tab yang sama langsung menampilkan perubahan.
const pembersih = new Set()

export function daftarkanPembersihCache(fungsi) {
  pembersih.add(fungsi)
}

export function bersihkanCachePublik() {
  for (const fungsi of pembersih) fungsi()
}
