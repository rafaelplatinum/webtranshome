import { ref } from 'vue'
import { defineStore } from 'pinia'

const topikAwal = 'Lainnya'

export const useWhatsAppContactStore = defineStore('whatsappContact', () => {
  const terbuka = ref(false)
  const konteks = ref({ topic: topikAwal, product: '', message: '' })

  function buka(prefill = {}) {
    konteks.value = {
      topic: prefill.topic || topikAwal,
      product: prefill.product || '',
      message: prefill.message || '',
    }
    terbuka.value = true
  }

  function tutup() { terbuka.value = false }

  return { terbuka, konteks, buka, tutup }
})
