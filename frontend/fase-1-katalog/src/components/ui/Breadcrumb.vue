<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  // [{ label, to }]; item terakhir adalah halaman yang sedang dibuka (tanpa tautan).
  items: { type: Array, required: true },
})

// HP: jejak panjang memakan 2–3 baris, jadi cukup satu tautan "‹ <induk>" ke tingkat di atasnya.
const induk = computed(() => [...props.items.slice(0, -1)].reverse().find((item) => item.to) ?? null)

const kelasTautan =
  'inline-flex min-h-11 items-center rounded-sm transition-colors duration-150 hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-8'
</script>

<template>
  <nav aria-label="Breadcrumb" class="text-sm text-muted">
    <RouterLink v-if="induk" :to="induk.to" :class="[kelasTautan, 'gap-1 md:hidden']">
      <ChevronLeft class="size-4 shrink-0" aria-hidden="true" />
      {{ induk.label }}
    </RouterLink>
    <ol class="flex-wrap items-center gap-x-1" :class="induk ? 'hidden md:flex' : 'flex'">
      <li v-for="(item, i) in items" :key="`${i}-${item.label}`" class="flex min-w-0 items-center gap-1">
        <ChevronRight v-if="i > 0" class="size-4 shrink-0 text-faint" aria-hidden="true" />
        <RouterLink v-if="item.to && i < items.length - 1" :to="item.to" :class="kelasTautan">{{ item.label }}</RouterLink>
        <span v-else aria-current="page" class="line-clamp-1 py-1 font-semibold text-foreground">{{ item.label }}</span>
      </li>
    </ol>
  </nav>
</template>
