<script setup>
import { RouterLink } from 'vue-router'

// Judul seksi Beranda + tautan "lihat semua" di kanan.
defineProps({
  id: { type: String, required: true },
  judul: { type: String, required: true },
  tautan: { type: String, default: '' },
  labelTautan: { type: String, default: '' },
})
</script>

<template>
  <div class="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
    <h2 :id="id" class="text-2xl font-extrabold">{{ judul }}</h2>
    <RouterLink
      v-if="tautan"
      :to="tautan"
      class="inline-flex min-h-11 items-center rounded-sm text-sm font-bold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
    >
      {{ labelTautan }}
    </RouterLink>
  </div>
</template>
