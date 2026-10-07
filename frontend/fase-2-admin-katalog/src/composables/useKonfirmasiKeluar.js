import { onBeforeUnmount, onMounted, ref } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

/**
 * "Buang perubahan?" saat meninggalkan form yang belum disimpan: pindah menu, Batal, atau tutup tab.
 * Sesi habis (401 → diarahkan ke halaman masuk): keluar tanpa bertanya, dan `saatSesiHabis`
 * dipanggil supaya isian bisa diselamatkan lalu dipulihkan setelah masuk lagi.
 */
export function useKonfirmasiKeluar(kotor, { saatSesiHabis } = {}) {
  const auth = useAuthStore()
  const terbuka = ref(false)
  let putuskan = null
  let bebas = false

  onBeforeRouteLeave(() => {
    if (bebas || !kotor.value) return true
    if (!auth.isLoggedIn) {
      saatSesiHabis?.()
      return true
    }
    terbuka.value = true
    return new Promise((selesai) => {
      putuskan = selesai
    })
  })

  function jawab(keluar) {
    terbuka.value = false
    putuskan?.(keluar)
    putuskan = null
  }

  function peringatanTab(event) {
    if (bebas || !kotor.value) return
    event.preventDefault()
    event.returnValue = ''
  }

  onMounted(() => window.addEventListener('beforeunload', peringatanTab))
  onBeforeUnmount(() => window.removeEventListener('beforeunload', peringatanTab))

  return {
    terbuka,
    buang: () => jawab(true),
    tetap: () => jawab(false),
    /** Setelah berhasil disimpan: pindah halaman tanpa bertanya. */
    izinkanKeluar: () => {
      bebas = true
    },
  }
}
