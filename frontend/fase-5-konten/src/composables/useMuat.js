import { onMounted, ref, shallowRef } from 'vue'
import { pesanError } from '@/services/errors'

/**
 * Muat satu sumber data saat komponen tampil, dengan status untuk skeleton / kosong / gagal + Coba lagi.
 * Dipakai seksi-seksi Beranda yang saling terpisah: satu seksi gagal tidak menggagalkan yang lain.
 */
export function useMuat(ambil, { segera = true } = {}) {
  const status = ref('loading') // loading | ready | error
  const data = shallowRef(null)
  const pesan = ref('')

  async function muat() {
    status.value = 'loading'
    try {
      data.value = await ambil()
      status.value = 'ready'
    } catch (error) {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }

  if (segera) onMounted(muat)
  return { status, data, pesan, muat }
}
