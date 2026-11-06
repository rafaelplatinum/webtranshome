import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { getPublicSettings } from '@/services/settingsService'

// Pengaturan toko (site_settings) yang dipakai banyak komponen: nomor WhatsApp, template pesan, alamat, jam buka.
export const useSettingsStore = defineStore('settings', () => {
  const data = ref(null)
  const status = ref('idle') // idle | loading | ready | error

  const waNumber = computed(() => data.value?.wa_number ?? null)
  // Template pesan WhatsApp (KEPUTUSAN #23). null/kosong → bawaan di utils/templateWhatsApp.js.
  const templateWa = computed(() => ({
    produk: data.value?.wa_template_product ?? null,
    keranjang: data.value?.wa_template_cart ?? null,
    baris: data.value?.wa_template_cart_item ?? null,
  }))

  async function load({ force = false } = {}) {
    if (status.value === 'loading' || (status.value === 'ready' && !force)) return
    status.value = 'loading'
    try {
      data.value = await getPublicSettings()
      status.value = 'ready'
    } catch {
      status.value = 'error'
    }
  }

  return { data, status, waNumber, templateWa, load }
})
