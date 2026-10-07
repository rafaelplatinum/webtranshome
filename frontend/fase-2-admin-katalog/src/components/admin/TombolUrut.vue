<script setup>
import { ArrowDown, ArrowUp } from 'lucide-vue-next'

// Tombol ↑ ↓ untuk mengubah urutan baris tabel master (bisa dipakai keyboard).
defineProps({
  id: { type: [Number, String], required: true },
  nama: { type: String, required: true },
  pertama: { type: Boolean, default: false },
  terakhir: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['geser'])

const kelas =
  'inline-flex size-11 items-center justify-center rounded-md border border-border-strong bg-surface transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:border-border-strong md:size-9'
</script>

<template>
  <div class="inline-flex gap-1">
    <button type="button" :class="kelas" :data-urut="`${id}:-1`" :disabled="disabled || pertama" :aria-label="`Naikkan urutan ${nama}`" @click="emit('geser', -1)">
      <ArrowUp class="size-4" aria-hidden="true" />
    </button>
    <button type="button" :class="kelas" :data-urut="`${id}:1`" :disabled="disabled || terakhir" :aria-label="`Turunkan urutan ${nama}`" @click="emit('geser', 1)">
      <ArrowDown class="size-4" aria-hidden="true" />
    </button>
  </div>
</template>
