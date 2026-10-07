<script setup>
import { CalendarClock, CalendarX, CircleCheck, EyeOff } from 'lucide-vue-next'

// Label status banner dari backend: warna + ikon + teks (tidak hanya warna).
defineProps({
  status: { type: String, required: true },
})

const STATUS = {
  TAYANG: { label: 'Tayang', ikon: CircleCheck, kelas: 'bg-success-soft text-success-ink' },
  TERJADWAL: { label: 'Terjadwal', ikon: CalendarClock, kelas: 'bg-info-soft text-info-ink' },
  BERAKHIR: { label: 'Berakhir', ikon: CalendarX, kelas: 'bg-warning-soft text-warning-ink' },
  NONAKTIF: { label: 'Nonaktif', ikon: EyeOff, kelas: 'bg-subtle text-muted' },
}
</script>

<template>
  <span
    class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs font-semibold whitespace-nowrap"
    :class="STATUS[status]?.kelas ?? 'bg-subtle text-foreground'"
    :data-status-banner="status"
  >
    <component :is="STATUS[status].ikon" v-if="STATUS[status]" class="size-3.5" aria-hidden="true" />
    {{ STATUS[status]?.label ?? status }}
  </span>
</template>
