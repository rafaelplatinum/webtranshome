<script setup>
import Badge from '@/components/ui/Badge.vue'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'

const VARIAN = ['neutral', 'primary', 'success', 'warning', 'info', 'danger', 'accent']
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Badge</p>
      <div class="flex flex-wrap gap-2">
        <Badge v-for="v in VARIAN" :key="v" :variant="v">{{ v === 'accent' ? 'PROMO' : v }}</Badge>
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Status stok (label + ikon, bukan warna saja)</p>
      <div class="flex flex-wrap gap-2">
        <StockBadge status="TERSEDIA" />
        <StockBadge status="SISA_STOK" keterangan="3 dus" />
        <StockBadge status="PRE_ORDER" />
        <StockBadge status="HABIS" />
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Harga per satuan jual (sm / md / lg)</p>
      <div class="flex flex-wrap items-baseline gap-6">
        <PriceTag price="68000.00" unit="sak" size="sm" />
        <PriceTag price="89500.00" unit="dus" />
        <PriceTag price="3250000.00" unit="unit" size="lg" />
      </div>
    </div>
  </div>
</template>
