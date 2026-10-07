import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

/**
 * Filter yang disimpan di URL query: bisa dibagikan, tetap ada setelah refresh, aman untuk tombol Back.
 * Kunci multi berisi banyak nilai, disimpan dipisah koma: ?brand=merek-a,merek-b
 * `multi` berupa daftar kunci atau fungsi (kunci) => boolean, contoh untuk kunci dinamis attr[...].
 * Setiap perubahan filter mengembalikan `page` ke 1 (kecuali `page` sendiri yang diubah).
 */
export function useFilter({ multi = [] } = {}) {
  const route = useRoute()
  const router = useRouter()
  const apakahMulti = typeof multi === 'function' ? multi : (kunci) => multi.includes(kunci)

  const nilai = computed(() => {
    const hasil = {}
    for (const [kunci, isi] of Object.entries(route.query)) {
      const teks = Array.isArray(isi) ? isi.join(',') : String(isi ?? '')
      hasil[kunci] = apakahMulti(kunci) ? teks.split(',').filter(Boolean) : teks
    }
    if (Array.isArray(multi)) for (const kunci of multi) hasil[kunci] ??= []
    return hasil
  })

  function terapkan(perubahan) {
    const query = { ...route.query }
    for (const [kunci, isi] of Object.entries(perubahan)) {
      const kosong = isi == null || isi === '' || (Array.isArray(isi) && isi.length === 0)
      if (kosong) delete query[kunci]
      else query[kunci] = Array.isArray(isi) ? isi.join(',') : String(isi)
    }
    if (!('page' in perubahan)) delete query.page
    return router.replace({ query })
  }

  function toggle(kunci, item) {
    const daftar = new Set(nilai.value[kunci] ?? [])
    if (daftar.has(item)) daftar.delete(item)
    else daftar.add(item)
    return terapkan({ [kunci]: [...daftar] })
  }

  /** Hapus semua filter kecuali kunci yang disebut, misalnya `q` di halaman pencarian. */
  function reset(kecuali = []) {
    const query = Object.fromEntries(Object.entries(route.query).filter(([kunci]) => kecuali.includes(kunci)))
    return router.replace({ query })
  }

  return { nilai, terapkan, toggle, reset }
}
