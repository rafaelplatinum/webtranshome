<script setup>
import { RefreshCw, TriangleAlert } from 'lucide-vue-next'

// Pesan gagal ringkas untuk satu seksi Beranda: seksi lain tetap tampil.
defineProps({
  pesan: { type: String, required: true },
})

const emit = defineEmits(['coba-lagi'])
</script>

<template>
  <div role="alert" class="flex flex-wrap items-center gap-3 rounded-xl border border-danger-soft bg-surface px-4 py-3 text-sm">
    <TriangleAlert class="size-5 shrink-0 text-danger" aria-hidden="true" />
    <span class="text-muted">{{ pesan }}</span>
    <button
      type="button"
      class="ml-auto inline-flex min-h-11 items-center gap-1.5 rounded-md px-2 font-semibold text-primary transition-colors duration-150 hover:text-primary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-9"
      @click="emit('coba-lagi')"
    >
      <RefreshCw class="size-4" aria-hidden="true" />
      Coba lagi
    </button>
  </div>
</template>
