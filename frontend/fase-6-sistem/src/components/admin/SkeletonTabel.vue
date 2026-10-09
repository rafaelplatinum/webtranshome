<script setup>
import Skeleton from '@/components/ui/Skeleton.vue'

defineProps({
  // Teks untuk pembaca layar, mis. "Memuat log aktivitas…"
  label: { type: String, required: true },
  baris: { type: Number, default: 8 },
})
</script>

<template>
  <div class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
    <span class="sr-only">{{ label }}</span>
    <div v-for="n in baris" :key="n" class="flex items-center gap-3">
      <Skeleton class="h-4 w-32" />
      <Skeleton class="h-4 flex-1" />
      <Skeleton class="h-4 w-24" />
      <Skeleton class="h-4 w-16" />
    </div>
  </div>
</template>
