<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Search } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import { formatAngka } from '@/utils/format'

const props = defineProps({
  // Nilai filter dari URL: { status, q }
  nilai: { type: Object, required: true },
  // meta.counts dari GET /admin/qontak/contacts: { ALL, PENDING, SYNCED, FAILED }
  counts: { type: Object, default: null },
})

const emit = defineEmits(['ubah'])

const PILIHAN = [
  { value: '', label: 'Semua', kunci: 'ALL' },
  { value: 'FAILED', label: 'Gagal', kunci: 'FAILED' },
  { value: 'PENDING', label: 'Menunggu', kunci: 'PENDING' },
  { value: 'SYNCED', label: 'Tersinkron', kunci: 'SYNCED' },
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

const aktif = computed(() => props.nilai.status ?? '')
</script>

<template>
  <section aria-label="Filter kontak" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4 md:flex-row md:items-end md:justify-between">
    <div role="group" aria-label="Status sync" class="flex flex-wrap gap-2">
      <button
        v-for="p in PILIHAN"
        :key="p.kunci"
        type="button"
        :aria-pressed="aktif === p.value ? 'true' : 'false'"
        class="inline-flex min-h-11 cursor-pointer items-center gap-2 rounded-full border px-4 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-10"
        :class="aktif === p.value ? 'border-secondary bg-secondary text-secondary-foreground' : 'border-border-strong bg-surface hover:border-foreground'"
        @click="emit('ubah', { status: p.value || null })"
      >
        {{ p.label }}
        <span v-if="counts" class="tabular-nums" :class="aktif === p.value ? 'opacity-90' : 'text-muted'">{{ formatAngka(counts[p.kunci] ?? 0) }}</span>
      </button>
    </div>
    <div class="min-w-0 md:w-80">
      <Input
        v-model="teks"
        label="Cari member"
        type="search"
        placeholder="Nama, email, kode member, HP"
        autocomplete="off"
        enterkeyhint="search"
        @keydown.enter.prevent="cariSekarang"
      >
        <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
      </Input>
    </div>
  </section>
</template>
