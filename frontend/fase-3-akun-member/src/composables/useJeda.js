import { computed, onBeforeUnmount, ref } from 'vue'

/**
 * Jeda tombol "kirim ulang" (email verifikasi, link reset): setelah dipakai, tombol menunggu `detik` dulu.
 * Pakai: const jeda = useJeda(60); jeda.mulai(); label = jeda.aktif.value ? `Kirim ulang (${jeda.sisa.value})` : '…'
 */
export function useJeda(detik = 60) {
  const sisa = ref(0)
  let timer = null

  function mulai() {
    clearInterval(timer)
    sisa.value = detik
    timer = setInterval(() => {
      sisa.value -= 1
      if (sisa.value <= 0) clearInterval(timer)
    }, 1000)
  }

  onBeforeUnmount(() => clearInterval(timer))

  return { sisa, aktif: computed(() => sisa.value > 0), mulai }
}
