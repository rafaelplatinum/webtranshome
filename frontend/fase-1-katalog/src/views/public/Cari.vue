<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Search } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import KatalogDaftar from '@/components/katalog/KatalogDaftar.vue'
import { useSeo } from '@/composables/useSeo'

// /cari?q= : katalog dengan kata kunci terkunci. Kata kunci tetap ada saat filter direset.
const route = useRoute()
const q = computed(() => String(route.query.q ?? '').trim())
const kunci = computed(() => ({ q: q.value }))

useSeo({ title: () => (q.value ? `Hasil pencarian "${q.value}"` : 'Cari produk') })
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Pencarian' }]" />
    <h1 class="mt-2 text-2xl font-extrabold break-words sm:text-3xl">
      <template v-if="q">Hasil pencarian “{{ q }}”</template>
      <template v-else>Cari produk</template>
    </h1>

    <EmptyState v-if="!q" class="mt-6" title="Ketik kata kunci dulu" description="Cari nama produk, SKU, atau brand lewat kotak pencarian di bagian atas halaman.">
      <template #icon><Search class="size-10" aria-hidden="true" /></template>
    </EmptyState>
    <KatalogDaftar
      v-else
      class="mt-6"
      :kunci="kunci"
      :kunci-url="['q']"
      :topik="q"
      :judul-kosong="`Tidak ada hasil untuk “${q}”`"
      deskripsi-kosong="Periksa ejaan atau pakai kata yang lebih umum. Bisa juga tanya tim kami, mungkin barangnya ada di toko."
    />
  </div>
</template>
