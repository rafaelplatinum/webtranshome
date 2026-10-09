<script setup>
import { RouterLink } from 'vue-router'
import { ArrowRight, ImageOff } from 'lucide-vue-next'

defineProps({
  // Item GET /rooms: { id, name, slug, image_cover, product_count }
  ruangan: { type: Object, required: true },
})
</script>

<template>
  <RouterLink
    :to="`/ruangan/${ruangan.slug}`"
    class="group flex h-full flex-col overflow-hidden rounded-xl border border-border bg-surface transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
  >
    <span class="block aspect-4/3 bg-subtle">
      <img v-if="ruangan.image_cover" :src="ruangan.image_cover" alt="" loading="lazy" width="400" height="300" class="size-full object-cover" />
      <span v-else class="flex size-full items-center justify-center text-faint"><ImageOff class="size-8" aria-hidden="true" /></span>
    </span>
    <span class="flex items-center justify-between gap-3 p-4">
      <span class="flex min-w-0 flex-col">
        <span class="font-bold">{{ ruangan.name }}</span>
        <span v-if="ruangan.product_count != null" class="text-sm text-muted">{{ ruangan.product_count }} produk</span>
      </span>
      <ArrowRight class="size-5 shrink-0 text-muted transition-colors duration-150 group-hover:text-primary" aria-hidden="true" />
    </span>
  </RouterLink>
</template>
