<script setup>
import { computed } from 'vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { listBrands } from '@/services/catalogService'
import { brandLingkaran } from '@/utils/brandBeranda'
import PesanGagal from './PesanGagal.vue'
import SusunanLingkaranBrand from './SusunanLingkaranBrand.vue'

// "Brand pilihan" Beranda (desain pemilik 2 Nov 2026, KEPUTUSAN #19): judul di tengah diapit lingkaran logo brand.
// Isi & urutan diatur di admin › Katalog › Brand di Beranda (KEPUTUSAN #21); belum ada yang dipilih → otomatis.
// Semua brand tetap bisa dilihat di /brand.
const merek = useMuat(listBrands)
const jumlah = computed(() => merek.data.value?.length ?? 0)
const susunan = computed(() => brandLingkaran(merek.data.value ?? []))
</script>

<template>
  <section v-if="merek.status.value !== 'ready' || jumlah" aria-labelledby="judul-brand-pilihan">
    <div v-if="merek.status.value === 'loading'" class="flex flex-col items-center gap-6 rounded-2xl border border-border bg-surface py-8" aria-busy="true">
      <span class="sr-only">Memuat brand…</span>
      <div class="flex gap-3"><Skeleton v-for="n in 4" :key="n" class="size-18 rounded-full" /></div>
      <Skeleton class="h-10 w-72 max-w-full" />
    </div>
    <PesanGagal v-else-if="merek.status.value === 'error'" pesan="Brand belum bisa dimuat." @coba-lagi="merek.muat" />
    <SusunanLingkaranBrand v-else :brands="susunan.brands" :total="jumlah" />
  </section>
</template>
