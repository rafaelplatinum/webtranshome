<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Search } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import { TIPE_ARTIKEL } from '@/utils/artikel'

const props = defineProps({
  // Nilai filter dari URL: { q, type, status }
  nilai: { type: Object, required: true },
})

// `ubah`: perubahan sebagian (halaman kembali ke 1). `reset`: hapus semua filter.
const emit = defineEmits(['ubah', 'reset'])

const OPSI_TIPE = [{ value: '', label: 'Semua tipe' }, ...Object.entries(TIPE_ARTIKEL).map(([value, t]) => ({ value, label: t.label }))]
const OPSI_STATUS = [
  { value: '', label: 'Semua status' },
  { value: 'DRAFT', label: 'Draf' },
  { value: 'PUBLISHED', label: 'Terbit' },
]

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

const adaFilter = computed(() => ['q', 'type', 'status'].some((k) => props.nilai[k]))
const pilih = (kunci) => ({
  modelValue: props.nilai[kunci] ?? '',
  'onUpdate:modelValue': (v) => emit('ubah', { [kunci]: v || null }),
})
</script>

<template>
  <section aria-label="Filter artikel" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
    <div class="flex flex-wrap items-end gap-2.5">
      <div class="min-w-0 flex-[2_1_15rem]">
        <Input v-model="teks" label="Cari" type="search" placeholder="Judul atau slug" autocomplete="off" enterkeyhint="search" @keydown.enter.prevent="cariSekarang">
          <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
        </Input>
      </div>
      <div class="min-w-0 flex-[1_1_9rem]"><Select v-bind="pilih('type')" label="Tipe" :options="OPSI_TIPE" /></div>
      <div class="min-w-0 flex-[1_1_9rem]"><Select v-bind="pilih('status')" label="Status" :options="OPSI_STATUS" /></div>
    </div>
    <button
      v-if="adaFilter"
      type="button"
      class="inline-flex h-11 items-center self-start rounded-sm px-1 text-sm font-semibold underline underline-offset-2 transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
      @click="emit('reset')"
    >
      Reset filter
    </button>
  </section>
</template>
