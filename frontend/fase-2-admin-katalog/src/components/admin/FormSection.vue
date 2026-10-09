<script setup>
import { useId } from 'vue'

defineProps({
  judul: { type: String, required: true },
  deskripsi: { type: String, default: '' },
})

const idJudul = useId()
</script>

<template>
  <!-- `id` dari luar (mis. "foto", "stok") dipakai tautan dashboard /admin/produk/:id#foto. -->
  <section :aria-labelledby="idJudul" class="flex scroll-mt-24 flex-col gap-4 rounded-xl border border-border bg-surface p-5">
    <div class="flex flex-col gap-0.5">
      <h2 :id="idJudul" class="text-base font-bold">{{ judul }}</h2>
      <p v-if="deskripsi" class="text-sm text-muted">{{ deskripsi }}</p>
    </div>
    <slot />
  </section>
</template>
