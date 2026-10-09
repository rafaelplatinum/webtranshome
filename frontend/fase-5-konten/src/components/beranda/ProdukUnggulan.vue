<script setup>
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import ProductCard from '@/components/katalog/ProductCard.vue'
import { useMuat } from '@/composables/useMuat'
import { listProducts } from '@/services/productService'
import JudulSeksi from './JudulSeksi.vue'
import PesanGagal from './PesanGagal.vue'

// Produk bertanda unggulan di admin (★), paling banyak 8. Seksi disembunyikan bila belum ada.
const unggulan = useMuat(async () => (await listProducts({ featured: 1, per_page: 8 })).data)
</script>

<template>
  <section
    v-if="unggulan.status.value !== 'ready' || unggulan.data.value.length"
    aria-labelledby="judul-produk-unggulan"
    class="flex flex-col gap-4"
  >
    <JudulSeksi id="judul-produk-unggulan" judul="Produk unggulan" tautan="/katalog" label-tautan="Lihat katalog" />
    <div v-if="unggulan.status.value === 'loading'" class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4" aria-busy="true">
      <span class="sr-only">Memuat produk unggulan…</span>
      <SkeletonCard v-for="n in 4" :key="n" />
    </div>
    <PesanGagal v-else-if="unggulan.status.value === 'error'" pesan="Produk unggulan belum bisa dimuat." @coba-lagi="unggulan.muat" />
    <ul v-else class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4" data-produk-unggulan>
      <li v-for="p in unggulan.data.value" :key="p.id" class="flex">
        <ProductCard :product="p" class="w-full" />
      </li>
    </ul>
  </section>
</template>
