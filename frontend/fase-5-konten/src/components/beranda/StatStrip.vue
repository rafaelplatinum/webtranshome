<script setup>
import { computed } from 'vue'
import { useMuat } from '@/composables/useMuat'
import { listBrands } from '@/services/catalogService'
import { CONTACT } from '@/config/contact'

const brands = useMuat(listBrands)
const jumlahBrand = computed(() => brands.status.value === 'ready' ? brands.data.value.length : '—')
const statistik = computed(() => [
  { nilai: `${jumlahBrand.value} brand`, label: 'pilihan partner Transhome' },
  { nilai: 'Senin–Sabtu', label: CONTACT.openingHours.replace('Senin–Sabtu, ', '') },
  { nilai: 'Trans Family', label: 'belanja dan kumpulkan poin' },
])
</script>

<template>
  <section class="border-y border-white/10 bg-th-navy text-white" aria-label="Sekilas tentang Transhome">
    <div class="th-container grid gap-5 py-6 sm:grid-cols-3 sm:gap-8 sm:py-8">
      <div v-for="item in statistik" :key="item.nilai" class="flex flex-col gap-1">
        <strong class="text-2xl leading-tight font-extrabold tracking-tight text-th-yellow sm:text-3xl">{{ item.nilai }}</strong>
        <span class="text-sm text-white/75">{{ item.label }}</span>
      </div>
    </div>
  </section>
</template>
