import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { getPublicSettings } from '@/services/settingsService'

// Pengaturan toko (site_settings) yang dipakai banyak komponen: nomor WhatsApp, alamat, jam buka.
export const useSettingsStore = defineStore('settings', () => {
  const data = ref(null)
  const status = ref('idle') // idle | loading | ready | error

  const waNumber = computed(() => data.value?.wa_number ?? null)

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

  return { data, status, waNumber, load }
})
