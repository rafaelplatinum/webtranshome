import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { pesanError } from '@/services/errors'

/**
 * Daftar ber-paginasi dari API: dimuat ulang setiap `params` berubah (permintaan lama dibatalkan).
 * `params` computed; bernilai null = jangan memuat (mis. rentang tanggal salah).
 * `muat({ diam: true })`: muat ulang tanpa skeleton dan tanpa mengganti tampilan bila gagal (polling).
 */
export function useDaftar(ambil, params, { satuan = 'data' } = {}) {
  const status = ref('loading') // loading | ready | error
  const baris = ref([])
  const meta = ref(null)
  const pesan = ref('')
  let pengendali = null

  async function muat({ diam = false } = {}) {
    if (params.value === null) return
    pengendali?.abort()
    pengendali = new AbortController()
    if (!diam) status.value = 'loading'
    try {
      const hasil = await ambil(params.value, { signal: pengendali.signal })
      baris.value = hasil.data
      meta.value = hasil.meta ?? null
      status.value = 'ready'
    } catch (error) {
      if (error?.code === 'ERR_CANCELED' || diam) return
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }

  watch(() => JSON.stringify(params.value), () => muat(), { immediate: true })
  onBeforeUnmount(() => pengendali?.abort())

  const rentang = computed(() => {
    if (!meta.value?.total) return ''
    const awal = (meta.value.page - 1) * meta.value.per_page + 1
    return `Menampilkan ${awal}–${awal + baris.value.length - 1} dari ${meta.value.total} ${satuan}`
  })

  return { status, baris, meta, pesan, muat, rentang }
}
