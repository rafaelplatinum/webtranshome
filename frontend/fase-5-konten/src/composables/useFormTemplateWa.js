import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { getWaTemplates, updateWaTemplates } from '@/services/adminKontenService'
import { errorField, pesanError } from '@/services/errors'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'
import { TEMPLATE_BAWAAN, cekTemplate } from '@/utils/templateWhatsApp'

// Kunci site_settings → jenis template (utils/templateWhatsApp.js).
export const KUNCI_TEMPLATE_WA = { wa_template_product: 'produk', wa_template_cart: 'keranjang', wa_template_cart_item: 'baris' }

/**
 * Form admin › Konten › Template WhatsApp (KEPUTUSAN #23). Ketiga template disimpan sekaligus.
 * Galat dicek di browser saat Simpan (lalu ikut diperbarui selagi diketik); server memakai aturan yang sama (422).
 */
export function useFormTemplateWa(formEl) {
  const settings = useSettingsStore()
  const ui = useUiStore()
  const status = ref('loading') // loading | ready | error
  const pesan = ref('')
  const galatUmum = ref('')
  const galatServer = ref({})
  const diperiksa = ref(false)
  const menyimpan = ref(false)
  const awal = ref('')
  const form = reactive(Object.fromEntries(Object.keys(KUNCI_TEMPLATE_WA).map((k) => [k, ''])))

  function isi(data) {
    for (const [kunci, jenis] of Object.entries(KUNCI_TEMPLATE_WA)) form[kunci] = data?.[kunci] || TEMPLATE_BAWAAN[jenis]
    awal.value = JSON.stringify(form)
  }

  async function muat() {
    status.value = 'loading'
    try {
      isi(await getWaTemplates())
      status.value = 'ready'
    } catch (error) {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }
  onMounted(muat)

  const kotor = computed(() => status.value === 'ready' && JSON.stringify(form) !== awal.value)
  const galatLokal = computed(() =>
    Object.fromEntries(
      Object.entries(KUNCI_TEMPLATE_WA)
        .map(([kunci, jenis]) => [kunci, cekTemplate(jenis, form[kunci])])
        .filter(([, galat]) => galat),
    ),
  )
  const galat = computed(() => ({ ...galatServer.value, ...(diperiksa.value ? galatLokal.value : {}) }))

  async function fokusGalat() {
    await nextTick()
    const el = formEl.value?.querySelector('[aria-invalid="true"]')
    el?.scrollIntoView({ block: 'center' })
    el?.focus({ preventScroll: true })
  }

  async function simpan() {
    if (menyimpan.value) return
    galatUmum.value = ''
    galatServer.value = {}
    diperiksa.value = true
    if (Object.keys(galatLokal.value).length) return fokusGalat()
    menyimpan.value = true
    try {
      isi(await updateWaTemplates({ ...form }))
      diperiksa.value = false
      // Tombol WhatsApp di website membaca template dari store ini.
      settings.load({ force: true })
      ui.tampilkanToast({ pesan: 'Template WhatsApp disimpan', jenis: 'success' })
    } catch (error) {
      const statusHttp = error?.response?.status
      if (statusHttp === 422) {
        galatServer.value = errorField(error)
        galatUmum.value = error.response.data?.message ?? 'Periksa kembali isian yang ditandai.'
        fokusGalat()
      } else if (statusHttp !== 401) {
        galatUmum.value = pesanError(error)
      }
    } finally {
      menyimpan.value = false
    }
  }

  function batal() {
    Object.assign(form, JSON.parse(awal.value))
    diperiksa.value = false
    galatServer.value = {}
    galatUmum.value = ''
  }

  return { status, pesan, muat, form, kotor, galat, galatUmum, menyimpan, simpan, batal }
}
