<script setup>
import { computed } from 'vue'
import { CircleCheck, CircleX, Clock, LoaderCircle } from 'lucide-vue-next'

const props = defineProps({
  // qontak_contacts.sync_status: PENDING | SYNCED | FAILED (null = belum ada kontak)
  status: { type: String, default: null },
  // true selama "Coba sync ulang" berjalan: ikon berputar.
  memproses: { type: Boolean, default: false },
})

// Status dibedakan dengan label + ikon, bukan warna saja.
const STATUS = {
  SYNCED: { label: 'Tersinkron', ikon: CircleCheck, kelas: 'bg-success-soft text-success-ink' },
  PENDING: { label: 'Menunggu', ikon: Clock, kelas: 'bg-info-soft text-info-ink' },
  FAILED: { label: 'Gagal', ikon: CircleX, kelas: 'bg-danger-soft text-danger-ink' },
}

const info = computed(() => STATUS[props.status] ?? null)
</script>

<template>
  <span v-if="info" class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs font-semibold whitespace-nowrap" :class="info.kelas">
    <LoaderCircle v-if="memproses" class="size-3.5 animate-spin" aria-hidden="true" />
    <component :is="info.ikon" v-else class="size-3.5" aria-hidden="true" />
    {{ info.label }}
  </span>
  <span v-else class="text-sm text-faint">Belum ada</span>
</template>
