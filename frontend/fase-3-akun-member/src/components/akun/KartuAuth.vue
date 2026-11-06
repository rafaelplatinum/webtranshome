<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const pendaftaran = computed(() => route.name === 'daftar')

defineProps({
  judul: { type: String, required: true },
  deskripsi: { type: String, default: '' },
})
</script>

<template>
  <!-- Kerangka halaman masuk, daftar, verifikasi, dan password (di dalam AuthLayout). -->
  <div class="flex flex-col gap-5" :class="pendaftaran ? '' : 'rounded-2xl border border-border bg-surface p-6'">
    <div class="flex flex-col gap-2">
      <span class="self-start rounded-md bg-th-yellow px-2.5 py-1 text-[13px] font-extrabold text-th-navy">Trans Family</span>
      <component :is="pendaftaran ? 'h2' : 'h1'" class="text-3xl leading-tight font-extrabold tracking-tight">{{ judul }}</component>
      <p v-if="deskripsi" class="text-base leading-relaxed text-th-muted">{{ deskripsi }}</p>
    </div>
    <slot />
  </div>
</template>
