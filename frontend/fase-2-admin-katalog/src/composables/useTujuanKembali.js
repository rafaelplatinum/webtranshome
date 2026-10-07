import { watch } from 'vue'
import { useRoute } from 'vue-router'

// Alamat terakhir halaman daftar (lengkap dengan filter & halaman), supaya Batal / Simpan di form
// kembali ke daftar yang sama, bukan daftar kosong.
const terakhir = new Map()

/** Dipanggil di halaman daftar: catat alamatnya setiap filter berubah. */
export function ingatDaftar(kunci) {
  const route = useRoute()
  watch(() => route.fullPath, (path) => terakhir.set(kunci, path), { immediate: true })
}

export function tujuanKembali(kunci, bawaan) {
  return terakhir.get(kunci) ?? bawaan
}
