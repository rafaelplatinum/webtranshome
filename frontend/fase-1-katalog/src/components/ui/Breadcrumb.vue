<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  // [{ label, to }]; item terakhir adalah halaman yang sedang dibuka (tanpa tautan).
  items: { type: Array, required: true },
  variant: { type: String, default: 'default', validator: (v) => ['default', 'light'].includes(v) },
})

// HP: jejak panjang memakan 2–3 baris, jadi cukup satu tautan "‹ <induk>" ke tingkat di atasnya.
const induk = computed(() => [...props.items.slice(0, -1)].reverse().find((item) => item.to) ?? null)

const kelasTautan =
  'inline-flex min-h-11 items-center rounded-sm transition-colors duration-150 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-8'
const warnaTautan = computed(() => props.variant === 'light' ? 'text-white/85 hover:text-white focus-visible:ring-white' : 'text-muted hover:text-primary focus-visible:ring-ring')
const warnaAktif = computed(() => props.variant === 'light' ? 'text-white' : 'text-foreground')
const warnaPemisah = computed(() => props.variant === 'light' ? 'text-white/70' : 'text-faint')
</script>

<template>
  <nav aria-label="Breadcrumb" class="text-sm" :class="variant === 'light' ? 'text-white/85' : 'text-muted'">
    <RouterLink v-if="induk" :to="induk.to" :class="[kelasTautan, warnaTautan, 'gap-1 md:hidden']">
      <ChevronLeft class="size-4 shrink-0" aria-hidden="true" />
      {{ induk.label }}
    </RouterLink>
    <ol class="flex-wrap items-center gap-x-1" :class="induk ? 'hidden md:flex' : 'flex'">
      <li v-for="(item, i) in items" :key="`${i}-${item.label}`" class="flex min-w-0 items-center gap-1">
        <ChevronRight v-if="i > 0" class="size-4 shrink-0" :class="warnaPemisah" aria-hidden="true" />
        <RouterLink v-if="item.to && i < items.length - 1" :to="item.to" :class="[kelasTautan, warnaTautan]">{{ item.label }}</RouterLink>
        <span v-else aria-current="page" class="line-clamp-1 py-1 font-semibold" :class="warnaAktif">{{ item.label }}</span>
      </li>
    </ol>
  </nav>
</template>
