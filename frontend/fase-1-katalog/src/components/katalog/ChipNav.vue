<script setup>
import { RouterLink } from 'vue-router'

defineProps({
  // Nama navigasi untuk pembaca layar, mis. "Subkategori".
  label: { type: String, required: true },
  // [{ key, label, to, aktif }]
  items: { type: Array, required: true },
})
</script>

<template>
  <!-- HP: satu baris yang bisa digeser; mulai 768px membungkus ke baris berikutnya. -->
  <nav :aria-label="label">
    <ul class="-mx-4 flex gap-2 overflow-x-auto px-4 py-1 scrollbar-none sm:-mx-6 sm:px-6 md:mx-0 md:flex-wrap md:overflow-visible md:px-0">
      <li v-for="item in items" :key="item.key" class="shrink-0">
        <RouterLink
          :to="item.to"
          :aria-current="item.aktif ? 'page' : undefined"
          class="inline-flex h-11 items-center rounded-full border px-4 text-sm font-semibold whitespace-nowrap transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 md:h-9"
          :class="item.aktif ? 'border-secondary bg-secondary text-secondary-foreground' : 'border-border-strong bg-surface hover:border-foreground'"
        >
          {{ item.label }}
        </RouterLink>
      </li>
    </ul>
  </nav>
</template>
