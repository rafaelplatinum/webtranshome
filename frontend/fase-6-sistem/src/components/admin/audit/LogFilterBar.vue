<script setup>
import { computed } from 'vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'

const props = defineProps({
  // Nilai filter dari URL: { admin_id, module, from, to }
  nilai: { type: Object, required: true },
  // Dari GET /admin/activity-logs/options: { admins: [{ id, name }], modules: [] }
  opsi: { type: Object, default: () => ({ admins: [], modules: [] }) },
  // Pesan bila tanggal akhir sebelum tanggal awal (daftar tidak dimuat).
  errorTanggal: { type: String, default: '' },
})

// `ubah`: perubahan sebagian (page kembali ke 1). `reset`: hapus semua filter.
const emit = defineEmits(['ubah', 'reset'])

const KUNCI = ['admin_id', 'module', 'from', 'to']
const adaFilter = computed(() => KUNCI.some((k) => props.nilai[k]))

const opsiAdmin = computed(() => [{ value: '', label: 'Semua admin' }, ...props.opsi.admins.map((a) => ({ value: String(a.id), label: a.name }))])
const opsiModul = computed(() => [{ value: '', label: 'Semua modul' }, ...props.opsi.modules.map((m) => ({ value: m, label: m }))])

const pilih = (kunci) => ({
  modelValue: props.nilai[kunci] ?? '',
  'onUpdate:modelValue': (v) => emit('ubah', { [kunci]: v || null }),
})
</script>

<template>
  <section aria-label="Filter log aktivitas" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
    <div class="flex flex-wrap items-start gap-2.5">
      <div class="min-w-0 flex-[1_1_12rem]"><Select v-bind="pilih('admin_id')" label="Admin" :options="opsiAdmin" /></div>
      <div class="min-w-0 flex-[1_1_11rem]"><Select v-bind="pilih('module')" label="Modul" :options="opsiModul" /></div>
      <div class="min-w-0 flex-[1_1_9.5rem]">
        <Input v-bind="pilih('from')" label="Dari tanggal" type="date" :max="nilai.to || undefined" />
      </div>
      <div class="min-w-0 flex-[1_1_9.5rem]">
        <Input v-bind="pilih('to')" label="Sampai tanggal" type="date" :min="nilai.from || undefined" :error="errorTanggal" />
      </div>
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
