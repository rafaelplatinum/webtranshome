import { computed, onBeforeUnmount, ref, toValue, watch } from 'vue'
import { useFilter } from '@/composables/useFilter'
import { listProducts } from '@/services/productService'
import { daftarkanPembersihCache } from '@/services/cache'
import { pesanError } from '@/services/errors'
import { formatRupiah } from '@/utils/format'
import { LABEL_STOK } from '@/utils/produk'

export const PER_HALAMAN = 24
const KUNCI_MULTI = ['brand', 'stock', 'room']
const apakahMulti = (kunci) => KUNCI_MULTI.includes(kunci) || /^attr\[.+\]$/.test(kunci)
const apakahFilter = (kunci) => apakahMulti(kunci) || kunci === 'min' || kunci === 'max'

// Cache singkat hasil per query: tombol Back langsung menampilkan daftar, posisi scroll bisa dipulihkan.
const cache = new Map()
const UMUR_CACHE_MS = 60_000
daftarkanPembersihCache(() => cache.clear())

/**
 * Data katalog dari URL query (filter, urutan, halaman) + parameter terkunci dari halaman
 * (kategori, brand, ruangan, kata kunci). Dipakai /katalog, /cari, /brand, /ruangan/:slug.
 */
export function useKatalog({ kunci = {}, kunciUrl = [] } = {}) {
  const { nilai, terapkan, toggle, reset } = useFilter({ multi: apakahMulti })

  const status = ref('loading') // loading | ready | error
  const produk = ref([])
  const meta = ref(null)
  const pesan = ref('')
  let pengendali = null

  // Kunci yang dikunci halaman (mis. brand di /brand/:slug) bukan filter pilihan pengguna.
  const terkunci = computed(() => Object.keys(toValue(kunci) ?? {}))
  const filterUrl = computed(() =>
    Object.fromEntries(Object.entries(nilai.value).filter(([k]) => apakahFilter(k) && !terkunci.value.includes(k))),
  )

  const parameterApi = computed(() => {
    const params = { per_page: PER_HALAMAN, sort: nilai.value.sort || 'terbaru', page: nilai.value.page || 1 }
    for (const [k, v] of Object.entries(filterUrl.value)) params[k] = Array.isArray(v) ? v.join(',') : v
    return { ...params, ...toValue(kunci) }
  })
  const kunciCache = computed(() => JSON.stringify(parameterApi.value))

  function ambilCache() {
    const isi = cache.get(kunciCache.value)
    if (!isi || Date.now() - isi.waktu > UMUR_CACHE_MS) return false
    produk.value = isi.data.data
    meta.value = isi.data.meta
    status.value = 'ready'
    return true
  }

  async function muat({ paksa = false } = {}) {
    if (!paksa && ambilCache()) return
    pengendali?.abort()
    pengendali = new AbortController()
    status.value = 'loading'
    try {
      const hasil = await listProducts(parameterApi.value, { signal: pengendali.signal })
      cache.set(kunciCache.value, { waktu: Date.now(), data: hasil })
      if (cache.size > 30) cache.delete(cache.keys().next().value)
      produk.value = hasil.data
      meta.value = hasil.meta
      status.value = 'ready'
    } catch (error) {
      if (error?.code === 'ERR_CANCELED') return
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }

  // Cache dibaca langsung saat setup agar daftar sudah ada pada render pertama (untuk tombol Back).
  if (!ambilCache()) muat()
  watch(kunciCache, () => muat())
  onBeforeUnmount(() => pengendali?.abort())

  const facets = computed(() => meta.value?.facets ?? null)

  /** Chip untuk setiap nilai filter aktif (yang terkunci halaman tidak ikut). */
  const chip = computed(() => {
    const f = facets.value
    const aktif = filterUrl.value
    const nama = (daftar, slug) => daftar?.find((x) => x.slug === slug)?.name ?? slug
    const hasil = []
    for (const slug of aktif.brand ?? []) hasil.push({ id: `brand-${slug}`, label: nama(f?.brands, slug), hapus: () => toggle('brand', slug) })
    for (const s of aktif.stock ?? []) hasil.push({ id: `stock-${s}`, label: LABEL_STOK[s] ?? s, hapus: () => toggle('stock', s) })
    for (const slug of aktif.room ?? []) hasil.push({ id: `room-${slug}`, label: nama(f?.rooms, slug), hapus: () => toggle('room', slug) })
    for (const [k, daftar] of Object.entries(aktif)) {
      if (!/^attr\[.+\]$/.test(k)) continue
      for (const v of daftar) hasil.push({ id: `${k}-${v}`, label: v, hapus: () => toggle(k, v) })
    }
    const { min, max } = aktif
    if (min || max) {
      const label = min && max ? `${formatRupiah(min)} – ${formatRupiah(max)}` : min ? `Mulai ${formatRupiah(min)}` : `Sampai ${formatRupiah(max)}`
      hasil.push({ id: 'harga', label, hapus: () => terapkan({ min: null, max: null }) })
    }
    return hasil
  })

  /** Ganti seluruh filter sekaligus (bottom sheet di HP). Kunci yang tidak ada di `baru` ikut dihapus. */
  function gantiFilter(baru) {
    const perubahan = {}
    for (const k of new Set([...Object.keys(filterUrl.value), ...Object.keys(baru)])) perubahan[k] = baru[k] ?? null
    return terapkan(perubahan)
  }

  return {
    status,
    produk,
    meta,
    facets,
    pesan,
    chip,
    filter: filterUrl,
    halaman: computed(() => Number(nilai.value.page) || 1),
    urutan: computed(() => nilai.value.sort || 'terbaru'),
    toggle,
    terapkan,
    gantiFilter,
    resetFilter: () => reset(kunciUrl.concat(['sort'])),
    muatUlang: () => muat({ paksa: true }),
  }
}
