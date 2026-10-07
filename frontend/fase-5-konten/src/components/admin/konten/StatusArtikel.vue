<script setup>
import { CircleCheck, FilePen } from 'lucide-vue-next'

// Draf vs terbit di daftar dan editor artikel admin.
defineProps({
  status: { type: String, required: true },
})

const STATUS = {
  PUBLISHED: { label: 'Terbit', ikon: CircleCheck, kelas: 'bg-success-soft text-success-ink' },
  DRAFT: { label: 'Draf', ikon: FilePen, kelas: 'bg-subtle text-foreground' },
}
</script>

<template>
  <span
    class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs font-semibold whitespace-nowrap"
    :class="STATUS[status]?.kelas ?? 'bg-subtle text-foreground'"
    :data-status-artikel="status"
  >
    <component :is="STATUS[status].ikon" v-if="STATUS[status]" class="size-3.5" aria-hidden="true" />
    {{ STATUS[status]?.label ?? status }}
  </span>
</template>
