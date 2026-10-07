<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Search } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'

const props = defineProps({
  // Nilai filter dari URL: { q, consent, hp, sync, active }
  nilai: { type: Object, required: true },
})

// `ubah`: perubahan sebagian (page kembali ke 1). `reset`: hapus semua filter.
const emit = defineEmits(['ubah', 'reset'])

const KUNCI = ['q', 'consent', 'hp', 'sync', 'active']
const OPSI = {
  consent: [
    { value: '', label: 'Semua' },
    { value: '1', label: 'Setuju' },
    { value: '0', label: 'Tidak setuju' },
  ],
  hp: [
    { value: '', label: 'Semua' },
    { value: '1', label: 'Ada nomor HP' },
    { value: '0', label: 'Belum ada' },
  ],
  sync: [
    { value: '', label: 'Semua' },
    { value: 'SYNCED', label: 'Tersinkron' },
    { value: 'PENDING', label: 'Menunggu' },
    { value: 'FAILED', label: 'Gagal' },
  ],
  active: [
    { value: '', label: 'Semua' },
    { value: '1', label: 'Aktif' },
    { value: '0', label: 'Nonaktif' },
  ],
}

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

const adaFilter = computed(() => KUNCI.some((k) => props.nilai[k]))

const pilih = (kunci) => ({
  modelValue: props.nilai[kunci] ?? '',
  'onUpdate:modelValue': (v) => emit('ubah', { [kunci]: v || null }),
})
</script>

<template>
  <section aria-label="Filter member" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
    <div class="flex flex-wrap items-end gap-2.5">
      <div class="min-w-0 flex-[2_1_16rem]">
        <Input
          v-model="teks"
          label="Cari"
          type="search"
          placeholder="Kode member, email, HP, nama"
          autocomplete="off"
          enterkeyhint="search"
          @keydown.enter.prevent="cariSekarang"
        >
          <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
        </Input>
      </div>
      <div class="min-w-0 flex-[1_1_9rem]"><Select v-bind="pilih('consent')" label="Persetujuan promo" :options="OPSI.consent" /></div>
      <div class="min-w-0 flex-[1_1_9rem]"><Select v-bind="pilih('hp')" label="Nomor HP" :options="OPSI.hp" /></div>
      <div class="min-w-0 flex-[1_1_9rem]"><Select v-bind="pilih('sync')" label="Sync Qontak" :options="OPSI.sync" /></div>
      <div class="min-w-0 flex-[1_1_8rem]"><Select v-bind="pilih('active')" label="Status akun" :options="OPSI.active" /></div>
    </div>
    <div v-if="adaFilter" class="flex text-sm">
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
