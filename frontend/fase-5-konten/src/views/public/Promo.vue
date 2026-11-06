<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import DaftarArtikel from '@/components/konten/DaftarArtikel.vue'
import { useSeo } from '@/composables/useSeo'

// Filter disimpan di URL (?jenis=promo|event) supaya tetap setelah refresh dan bisa dibagikan.
const FILTER = [
  { value: '', label: 'Semua', tipe: 'PROMO,EVENT' },
  { value: 'promo', label: 'Promo', tipe: 'PROMO' },
  { value: 'event', label: 'Event', tipe: 'EVENT' },
]

const route = useRoute()
const router = useRouter()
const filter = computed(() => FILTER.find((f) => f.value === route.query.jenis) ?? FILTER[0])

function pilih(f) {
  router.replace({ query: f.value ? { jenis: f.value } : {} })
}

useSeo({ title: 'Promo dan event', description: 'Promo dan event terbaru di toko Transhome, Sleman.' })
</script>

<template>
  <div class="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-4 sm:px-6 lg:py-6">
    <div class="flex flex-col gap-2">
      <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Promo dan event' }]" />
      <h1 class="text-2xl font-extrabold sm:text-3xl">Promo dan event</h1>
      <p class="max-w-2xl text-muted">Promo yang berlaku di toko dan acara terbaru. Tanyakan detailnya ke tim kami lewat WhatsApp.</p>
    </div>
    <div role="group" aria-label="Saring promo dan event" class="flex flex-wrap gap-2">
      <button
        v-for="f in FILTER"
        :key="f.value"
        type="button"
        class="inline-flex h-11 items-center rounded-full border px-4 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
        :class="filter.value === f.value ? 'border-primary bg-primary-soft text-primary-hover' : 'border-border-strong bg-surface hover:border-foreground'"
        :aria-pressed="filter.value === f.value ? 'true' : 'false'"
        @click="pilih(f)"
      >
        {{ f.label }}
      </button>
    </div>
    <DaftarArtikel :tipe="filter.tipe" :judul="`Daftar ${filter.value ? filter.label.toLowerCase() : 'promo dan event'}`" judul-kosong="Belum ada promo saat ini" deskripsi-kosong="Promo dan event baru akan muncul di sini." />
  </div>
</template>
