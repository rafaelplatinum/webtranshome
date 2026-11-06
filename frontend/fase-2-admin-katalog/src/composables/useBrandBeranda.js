import { computed, nextTick, onMounted, ref } from 'vue'
import { listMaster, saveBrandHome } from '@/services/adminMasterService'
import { errorField, pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'
import { MAKS_LINGKARAN_BRAND } from '@/utils/brandBeranda'

const urutNama = (a, b) => a.name.localeCompare(b.name, 'id')

/**
 * Halaman admin › Brand di Beranda (KEPUTUSAN #21): pilih & urutkan brand lingkaran Beranda (maks 16), lalu simpan
 * sekaligus (PUT /admin/brands/home). Perubahan baru tersimpan setelah Simpan; keluar halaman → "Buang perubahan?".
 */
export function useBrandBeranda() {
  const ui = useUiStore()
  const status = ref('loading') // loading | ready | error
  const pesan = ref('')
  const semua = ref([])
  const dipilih = ref([])
  const awal = ref('[]')
  const cari = ref('')
  const menyimpan = ref(false)
  const galat = ref('')

  const kunci = () => JSON.stringify(dipilih.value.map((b) => b.id))
  const urutTersimpan = (daftar) =>
    daftar.filter((b) => b.show_on_home && b.is_active).sort((a, b) => (a.home_order ?? Infinity) - (b.home_order ?? Infinity) || urutNama(a, b))

  function isi(daftar) {
    semua.value = daftar
    dipilih.value = urutTersimpan(daftar)
    awal.value = kunci()
  }

  async function muat() {
    status.value = 'loading'
    try {
      isi(await listMaster('brands'))
      status.value = 'ready'
    } catch (error) {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }
  onMounted(muat)

  const aktif = computed(() => semua.value.filter((b) => b.is_active))
  const penuh = computed(() => dipilih.value.length >= MAKS_LINGKARAN_BRAND)
  const kotor = computed(() => status.value === 'ready' && kunci() !== awal.value)
  // Dicentang tapi nonaktif: tidak tampil di Beranda, dan ikut dilepas saat disimpan.
  const nonaktifDicentang = computed(() => semua.value.filter((b) => b.show_on_home && !b.is_active))
  const tersedia = computed(() => {
    const ids = new Set(dipilih.value.map((b) => b.id))
    const q = cari.value.trim().toLowerCase()
    return aktif.value.filter((b) => !ids.has(b.id) && (!q || b.name.toLowerCase().includes(q) || b.slug.includes(q))).sort(urutNama)
  })

  // Fokus pindah ke elemen sebelah setelah tombol yang diklik hilang dari daftar.
  async function fokus(selektor, cadangan) {
    await nextTick()
    ;(document.querySelector(selektor) ?? document.querySelector(cadangan))?.focus()
  }

  function tambah(brand) {
    if (penuh.value) return
    const i = tersedia.value.indexOf(brand)
    const berikut = tersedia.value[i + 1] ?? tersedia.value[i - 1]
    dipilih.value = [...dipilih.value, brand]
    fokus(berikut ? `[data-tambah-brand="${berikut.id}"]` : null, '[data-cari-brand-beranda] input')
  }

  function lepas(brand) {
    const i = dipilih.value.indexOf(brand)
    const berikut = dipilih.value[i + 1] ?? dipilih.value[i - 1]
    dipilih.value = dipilih.value.filter((b) => b !== brand)
    fokus(berikut ? `[data-lepas-brand="${berikut.id}"]` : null, '[data-cari-brand-beranda] input')
  }

  async function geser(brand, arah) {
    const daftar = [...dipilih.value]
    const i = daftar.indexOf(brand)
    if (!daftar[i + arah]) return
    ;[daftar[i], daftar[i + arah]] = [daftar[i + arah], daftar[i]]
    dipilih.value = daftar
    await nextTick()
    const tombol = (a) => document.querySelector(`[data-urut="${brand.id}:${a}"]:not(:disabled)`)
    ;(tombol(arah) ?? tombol(-arah))?.focus()
  }

  function batal() {
    const ids = JSON.parse(awal.value)
    dipilih.value = ids.map((id) => semua.value.find((b) => b.id === id)).filter(Boolean)
    galat.value = ''
  }

  async function simpan() {
    if (menyimpan.value) return false
    galat.value = ''
    menyimpan.value = true
    try {
      isi(await saveBrandHome(dipilih.value.map((b) => b.id)))
      ui.tampilkanToast({ pesan: 'Brand di Beranda disimpan', jenis: 'success' })
      return true
    } catch (error) {
      if (error?.response?.status === 422) galat.value = errorField(error).brand_ids ?? 'Periksa kembali pilihan brand.'
      else if (error?.response?.status !== 401) galat.value = pesanError(error)
      return false
    } finally {
      menyimpan.value = false
    }
  }

  return { status, pesan, muat, semua, aktif, dipilih, tersedia, cari, penuh, kotor, nonaktifDicentang, menyimpan, galat, tambah, lepas, geser, batal, simpan }
}
