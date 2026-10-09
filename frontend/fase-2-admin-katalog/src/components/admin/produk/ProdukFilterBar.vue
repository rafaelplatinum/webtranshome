<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Search, X } from 'lucide-vue-next'
import Combobox from '@/components/ui/Combobox.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import { LABEL_STOK } from '@/utils/produk'

const props = defineProps({
  // Nilai filter dari URL: { q, category, brand, stock, active, foto, room }
  nilai: { type: Object, required: true },
  kategori: { type: Array, default: () => [] },
  brand: { type: Array, default: () => [] },
  ruangan: { type: Array, default: () => [] },
})

// `ubah`: perubahan sebagian (router.replace di halaman, page kembali ke 1). `reset`: hapus semua filter.
const emit = defineEmits(['ubah', 'reset'])

const teks = ref(props.nilai.q ?? '')
let jeda = null

// Cari diterapkan 400 ms setelah berhenti mengetik, atau langsung saat Enter.
watch(teks, (isi) => {
  clearTimeout(jeda)
  if (isi.trim() === (props.nilai.q ?? '')) return
  jeda = setTimeout(() => emit('ubah', { q: isi.trim() || null }), 400)
})
watch(
  () => props.nilai.q,
  (q) => {
    if ((q ?? '') !== teks.value.trim()) teks.value = q ?? ''
  },
)
onBeforeUnmount(() => clearTimeout(jeda))

function cariSekarang() {
  clearTimeout(jeda)
  if (teks.value.trim() !== (props.nilai.q ?? '')) emit('ubah', { q: teks.value.trim() || null })
}

const opsiKategori = computed(() => {
  const induk = props.kategori.filter((k) => k.parent_id == null)
  return induk.map((p) => ({
    label: p.name,
    options: [
      { value: String(p.id), label: `Semua ${p.name}` },
      ...props.kategori.filter((k) => k.parent_id === p.id).map((k) => ({ value: String(k.id), label: k.name })),
    ],
  }))
})
const opsiBrand = computed(() => [
  { value: '', label: 'Semua brand' },
  ...props.brand.map((b) => ({ value: String(b.id), label: b.is_active ? b.name : `${b.name} (nonaktif)` })),
])
const OPSI_STOK = [{ value: '', label: 'Semua' }, ...Object.entries(LABEL_STOK).map(([value, label]) => ({ value, label }))]
const OPSI_AKTIF = [
  { value: '', label: 'Semua' },
  { value: '1', label: 'Aktif' },
  { value: '0', label: 'Nonaktif' },
]
const OPSI_FOTO = [
  { value: '', label: 'Semua' },
  { value: '1', label: 'Ada foto utama' },
  { value: '0', label: 'Belum ada foto' },
]

const namaRuangan = computed(() => props.ruangan.find((r) => String(r.id) === String(props.nilai.room))?.name ?? `#${props.nilai.room}`)
const adaFilter = computed(() => ['q', 'category', 'brand', 'stock', 'active', 'foto', 'room'].some((k) => props.nilai[k]))

const pilih = (kunci) => ({
  modelValue: props.nilai[kunci] ?? '',
  'onUpdate:modelValue': (v) => emit('ubah', { [kunci]: v || null }),
})
</script>

<template>
  <section aria-label="Filter produk" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
    <div class="flex flex-wrap items-end gap-2.5">
      <div class="min-w-0 flex-[2_1_15rem]">
        <Input v-model="teks" label="Cari" type="search" placeholder="SKU atau nama produk" autocomplete="off" enterkeyhint="search" @keydown.enter.prevent="cariSekarang">
          <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
        </Input>
      </div>
      <div class="min-w-0 flex-[1_1_10rem]">
        <Select v-bind="pilih('category')" label="Kategori" :options="[{ value: '', label: 'Semua kategori' }]" :groups="opsiKategori" />
      </div>
      <!-- Brand bisa ratusan: isian yang bisa diketik untuk mencari, bukan dropdown panjang. -->
      <div class="min-w-0 flex-[1_1_11rem]"><Combobox v-bind="pilih('brand')" label="Brand" :options="opsiBrand" kosong="Brand tidak ditemukan" /></div>
      <div class="min-w-0 flex-[1_1_8rem]"><Select v-bind="pilih('stock')" label="Status stok" :options="OPSI_STOK" /></div>
      <div class="min-w-0 flex-[1_1_8rem]"><Select v-bind="pilih('active')" label="Tampil di web" :options="OPSI_AKTIF" /></div>
      <div class="min-w-0 flex-[1_1_8rem]"><Select v-bind="pilih('foto')" label="Foto" :options="OPSI_FOTO" /></div>
    </div>
    <div v-if="adaFilter" class="flex flex-wrap items-center gap-2 text-sm">
      <button
        v-if="nilai.room"
        type="button"
        class="inline-flex h-11 items-center gap-1.5 rounded-full border border-primary bg-primary-soft pr-2.5 pl-3.5 font-semibold text-primary-hover transition-colors duration-150 hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
        :aria-label="`Hapus filter ruangan ${namaRuangan}`"
        @click="emit('ubah', { room: null })"
      >
        Ruangan: {{ namaRuangan }}
        <X class="size-4" aria-hidden="true" />
      </button>
      <button
        type="button"
        class="inline-flex h-11 items-center rounded-sm px-1 font-semibold underline underline-offset-2 transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
        @click="emit('reset')"
      >
        Reset filter
      </button>
    </div>
  </section>
</template>
