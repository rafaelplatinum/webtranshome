import { ref } from 'vue'
import { defineStore } from 'pinia'

// State tampilan global. Toast ditampilkan oleh components/ui/Toast (belum dibuat).
export const useUiStore = defineStore('ui', () => {
  // jenis: 'info' | 'success' | 'warning' | 'danger'
  const toasts = ref([])
  let idBerikutnya = 1

  function tampilkanToast({ pesan, jenis = 'info', durasi = 4000 }) {
    const id = idBerikutnya++
    toasts.value.push({ id, pesan, jenis })
    if (durasi > 0) setTimeout(() => hapusToast(id), durasi)
    return id
  }

  function hapusToast(id) {
    toasts.value = toasts.value.filter((toast) => toast.id !== id)
  }

  return { toasts, tampilkanToast, hapusToast }
})
