import { computed, ref } from 'vue'

// Kolom matriks. Aksi selain lihat/tambah/ubah (terbitkan, penyesuaian, rasio poin, coba ulang, hapus brand) masuk
// "Lainnya". Hanya brand yang bisa dihapus, jadi tidak perlu kolom Hapus tersendiri.
export const KOLOM_IZIN = [
  { kunci: 'view', label: 'Lihat' },
  { kunci: 'create', label: 'Tambah' },
  { kunci: 'update', label: 'Ubah' },
  { kunci: 'lainnya', label: 'Lainnya' },
]

const LABEL_KHUSUS = { publish: 'Terbitkan', adjust: 'Penyesuaian', point_ratio: 'Rasio poin', retry: 'Lihat & coba ulang', delete: 'Hapus' }

export const kolomDari = (izin) => (['view', 'create', 'update'].includes(izin.action) ? izin.action : 'lainnya')
export const labelKhusus = (izin) => LABEL_KHUSUS[izin.action] ?? izin.action

/**
 * Status centang izin satu role. `katalog` = GET /admin/permissions (izin per menu).
 * Aturan: izin lain di satu menu butuh `menu_permission` (izin yang membuat menunya terlihat):
 * mencentang Tambah/Ubah ikut mencentang Lihat; mencabut Lihat mencabut seluruh baris.
 */
export function useMatriksIzin(katalog) {
  const asal = ref([])
  const pilihan = ref(new Set())

  const menuDari = (kode) => katalog.value.find((g) => g.permissions.some((p) => p.code === kode))

  function setel(daftar) {
    asal.value = [...daftar]
    pilihan.value = new Set(daftar)
  }

  function ubah(kode, nyala, hasil = new Set(pilihan.value)) {
    const menu = menuDari(kode)
    if (nyala) {
      hasil.add(kode)
      if (menu?.menu_permission) hasil.add(menu.menu_permission)
    } else if (menu && kode === menu.menu_permission) {
      menu.permissions.forEach((p) => hasil.delete(p.code))
    } else {
      hasil.delete(kode)
    }
    return hasil
  }

  const toggle = (kode, nyala) => (pilihan.value = ubah(kode, nyala))

  function toggleBaris(menu, nyala) {
    const hasil = new Set(pilihan.value)
    menu.permissions.forEach((p) => (nyala ? hasil.add(p.code) : hasil.delete(p.code)))
    pilihan.value = hasil
  }

  const izinKolom = (kunci) => katalog.value.flatMap((g) => g.permissions.filter((p) => kolomDari(p) === kunci))

  function toggleKolom(kunci, nyala) {
    let hasil = new Set(pilihan.value)
    for (const p of izinKolom(kunci)) hasil = ubah(p.code, nyala, hasil)
    pilihan.value = hasil
  }

  /** 'semua' | 'sebagian' | 'kosong' untuk checkbox "pilih semua" (sebagian = indeterminate). */
  function keadaan(daftarIzin) {
    const n = daftarIzin.filter((p) => pilihan.value.has(p.code)).length
    return n === 0 ? 'kosong' : n === daftarIzin.length ? 'semua' : 'sebagian'
  }

  const ditambah = computed(() => [...pilihan.value].filter((k) => !asal.value.includes(k)))
  const dicabut = computed(() => asal.value.filter((k) => !pilihan.value.has(k)))
  const kotor = computed(() => ditambah.value.length + dicabut.value.length > 0)

  // Urutan mengikuti katalog supaya data yang dikirim rapi.
  const hasilAkhir = () => katalog.value.flatMap((g) => g.permissions.map((p) => p.code)).filter((k) => pilihan.value.has(k))

  return { pilihan, setel, toggle, toggleBaris, toggleKolom, izinKolom, keadaan, ditambah, dicabut, kotor, hasilAkhir }
}
