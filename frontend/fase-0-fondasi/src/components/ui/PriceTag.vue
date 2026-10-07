<script setup>
import { formatRupiah } from '@/utils/format'

defineProps({
  // products.price_general (DECIMAL dari API berupa string)
  price: { type: [Number, String], required: true },
  // products.unit_sale, contoh "dus", "sak"
  unit: { type: String, default: '' },
  // `card`: kecil di kartu 2 kolom HP, normal mulai 640px.
  size: { type: String, default: 'md', validator: (v) => ['sm', 'md', 'lg', 'card'].includes(v) },
})

const UKURAN = { sm: 'text-lg', md: 'text-xl', lg: 'text-3xl', card: 'text-base sm:text-xl' }
</script>

<template>
  <span class="inline-flex flex-wrap items-baseline gap-x-1">
    <span class="font-extrabold tabular-nums" :class="UKURAN[size]">{{ formatRupiah(price) }}</span>
    <span v-if="unit" class="text-sm text-muted"><span class="sr-only">per </span><span aria-hidden="true">/</span>{{ unit }}</span>
  </span>
</template>
