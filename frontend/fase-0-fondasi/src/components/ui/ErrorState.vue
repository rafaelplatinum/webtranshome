<script setup>
import { RefreshCw, TriangleAlert } from 'lucide-vue-next'
import Button from './Button.vue'

defineProps({
  title: { type: String, default: 'Data gagal dimuat' },
  // Pakai pesanError(error) dari services/errors.js
  message: { type: String, default: 'Periksa koneksi, lalu coba lagi.' },
  retrying: { type: Boolean, default: false },
})

const emit = defineEmits(['retry'])
</script>

<template>
  <div role="alert" class="flex flex-col items-center gap-3 rounded-xl border border-danger-soft bg-surface px-6 py-12 text-center">
    <TriangleAlert class="size-10 text-danger" aria-hidden="true" />
    <p class="text-lg font-bold">{{ title }}</p>
    <p class="max-w-md text-muted">{{ message }}</p>
    <Button variant="outline" class="mt-2" :loading="retrying" @click="emit('retry')">
      <template #icon><RefreshCw class="size-4" aria-hidden="true" /></template>
      Coba lagi
    </Button>
  </div>
</template>
