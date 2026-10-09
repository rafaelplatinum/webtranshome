<script setup>
defineProps({
  judul: { type: String, required: true },
  deskripsi: { type: String, default: '' },
})
</script>

<template>
  <!-- Kerangka halaman masuk, daftar, verifikasi, dan password (di dalam AuthLayout). -->
  <div class="flex flex-col gap-6 rounded-2xl border border-border bg-surface p-6 sm:p-8">
    <div class="flex flex-col gap-2">
      <span class="self-start rounded bg-accent px-2.5 py-1 text-xs font-bold text-accent-foreground">Trans Family</span>
      <h1 class="text-2xl font-extrabold">{{ judul }}</h1>
      <p v-if="deskripsi" class="text-sm text-muted">{{ deskripsi }}</p>
    </div>
    <slot />
  </div>
</template>
