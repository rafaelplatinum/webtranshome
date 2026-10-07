<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Modal from '@/components/ui/Modal.vue'
import FilterPanel from './FilterPanel.vue'
import { listProducts } from '@/services/productService'

// Filter di HP: pilihan disimpan sebagai draf dan baru berlaku setelah "Terapkan". Batal/Esc membuang draf.
const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  // Filter yang sedang berlaku (dari URL).
  nilai: { type: Object, required: true },
  facets: { type: Object, default: null },
  total: { type: Number, default: null },
  // Parameter yang dikunci halaman (kategori, q, brand, ruangan), ikut dikirim saat menghitung pratinjau.
  kunci: { type: Object, default: () => ({}) },
  sembunyikan: { type: Array, default: () => [] },
})

const emit = defineEmits(['terapkan'])

const draf = ref({})
const pratinjau = ref(null) // { total, facets } hasil draf; null = sama dengan yang berlaku
const isi = ref(null)
let pengendali = null
let jeda = null

const salin = (filter) => Object.fromEntries(Object.entries(filter).map(([k, v]) => [k, Array.isArray(v) ? [...v] : v]))

/** Draf → parameter API dengan urutan tetap, juga dipakai untuk membandingkan draf dengan filter yang berlaku. */
function keParams(filter) {
  const params = {}
  for (const k of Object.keys(filter).sort()) {
    const v = filter[k]
    const teks = Array.isArray(v) ? [...v].sort().join(',') : String(v ?? '')
    if (teks) params[k] = teks
  }
  return params
}

const sama = computed(() => JSON.stringify(keParams(draf.value)) === JSON.stringify(keParams(props.nilai)))
const adaDraf = computed(() => Object.keys(keParams(draf.value)).length > 0)
const facetsTampil = computed(() => (sama.value ? props.facets : (pratinjau.value?.facets ?? props.facets)))
const totalDraf = computed(() => (sama.value ? props.total : (pratinjau.value?.total ?? null)))

function hentikanPratinjau() {
  clearTimeout(jeda)
  pengendali?.abort()
}

async function muatPratinjau() {
  pengendali = new AbortController()
  try {
    const hasil = await listProducts({ ...keParams(draf.value), ...props.kunci, per_page: 1 }, { signal: pengendali.signal })
    pratinjau.value = { total: hasil.meta.total, facets: hasil.meta.facets }
  } catch (error) {
    // Pratinjau hanya bantuan: kalau gagal, tombol tetap "Terapkan" tanpa angka.
    if (error?.code !== 'ERR_CANCELED') pratinjau.value = null
  }
}

watch(
  open,
  (buka) => {
    hentikanPratinjau()
    if (buka) {
      draf.value = salin(props.nilai)
      pratinjau.value = null
    }
  },
  { immediate: true },
)

// Angka per opsi dan jumlah hasil mengikuti draf, dihitung ulang 300 ms setelah perubahan terakhir.
watch(draf, () => {
  hentikanPratinjau()
  if (sama.value) pratinjau.value = null
  else jeda = setTimeout(muatPratinjau, 300)
})

onBeforeUnmount(hentikanPratinjau)

function toggle(kunci, nilai) {
  const daftar = new Set(draf.value[kunci] ?? [])
  if (daftar.has(nilai)) daftar.delete(nilai)
  else daftar.add(nilai)
  draf.value = { ...draf.value, [kunci]: [...daftar] }
}

function ubahHarga({ min, max }) {
  draf.value = { ...draf.value, min, max }
}

function terapkan() {
  // Harga yang belum valid (min > maks) harus dibetulkan dulu.
  const salah = isi.value?.querySelector('[aria-invalid="true"]')
  if (salah) {
    salah.focus()
    return
  }
  emit('terapkan', salin(draf.value))
  open.value = false
}
</script>

<template>
  <Modal v-model:open="open" title="Filter" sheet>
    <div ref="isi" class="flex flex-col gap-2">
      <div class="flex justify-end">
        <Button variant="ghost" size="sm" :disabled="!adaDraf" @click="draf = {}">Reset</Button>
      </div>
      <FilterPanel :nilai="draf" :facets="facetsTampil" :sembunyikan="sembunyikan" @toggle="toggle" @harga="ubahHarga" />
    </div>
    <template #footer>
      <Button variant="outline" @click="open = false">Batal</Button>
      <Button @click="terapkan">
        Terapkan<template v-if="totalDraf != null"> ({{ totalDraf }} produk)</template>
      </Button>
    </template>
  </Modal>
</template>
