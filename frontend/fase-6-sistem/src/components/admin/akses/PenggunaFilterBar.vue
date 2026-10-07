<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Search } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'

const props = defineProps({
  // Nilai filter dari URL: { q, role, active }
  nilai: { type: Object, required: true },
  // meta.roles dari GET /admin/users: [{ code, name }]
  roles: { type: Array, default: () => [] },
})

const emit = defineEmits(['ubah', 'reset'])

const OPSI_STATUS = [
  { value: '', label: 'Semua' },
  { value: '1', label: 'Aktif' },
  { value: '0', label: 'Nonaktif' },
]
const opsiRole = computed(() => [{ value: '', label: 'Semua role' }, ...props.roles.map((r) => ({ value: r.code, label: r.name }))])

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

const adaFilter = computed(() => ['q', 'role', 'active'].some((k) => props.nilai[k]))
const pilih = (kunci) => ({
  modelValue: props.nilai[kunci] ?? '',
  'onUpdate:modelValue': (v) => emit('ubah', { [kunci]: v || null }),
})
</script>

<template>
  <section aria-label="Filter pengguna admin" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
    <div class="flex flex-wrap items-end gap-2.5">
      <div class="min-w-0 flex-[2_1_16rem]">
        <Input v-model="teks" label="Cari" type="search" placeholder="Nama atau email" autocomplete="off" enterkeyhint="search" @keydown.enter.prevent="cariSekarang">
          <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
        </Input>
      </div>
      <div class="min-w-0 flex-[1_1_11rem]"><Select v-bind="pilih('role')" label="Role" :options="opsiRole" /></div>
      <div class="min-w-0 flex-[1_1_9rem]"><Select v-bind="pilih('active')" label="Status akun" :options="OPSI_STATUS" /></div>
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
