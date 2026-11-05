<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import DaftarArtikel from '@/components/konten/DaftarArtikel.vue'
import { useSeo } from '@/composables/useSeo'
import MediaPlaceholder from '@/components/public/MediaPlaceholder.vue'
import PageContactCta from '@/components/public/PageContactCta.vue'

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
  <div class="flex flex-col gap-12 sm:gap-16">
    <section class="bg-th-yellow py-10 sm:py-14 lg:py-16">
      <div class="th-container grid items-center gap-8 lg:grid-cols-[1.1fr_0.9fr] lg:gap-12">
        <div class="flex flex-col items-start gap-5">
          <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Promo dan event' }]" />
          <span class="th-eyebrow text-th-navy">Promo dan event</span>
          <h1 class="max-w-2xl text-4xl leading-[1.08] font-extrabold tracking-tight sm:text-5xl lg:text-[56px]">Promo untuk rumah yang lebih nyaman.</h1>
          <p class="max-w-xl text-base leading-relaxed text-th-navy/80 sm:text-lg">Temukan promo dan acara terbaru Transhome, lalu tanya tim kami untuk detailnya.</p>
        </div>
        <MediaPlaceholder tone="peach" class="h-56 rounded-[28px] sm:h-72 lg:h-[300px]" aria-hidden="true" />
      </div>
    </section>

    <section id="daftar-promo" class="th-container flex scroll-mt-28 flex-col gap-6">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div><h2 class="text-3xl leading-tight font-extrabold tracking-tight sm:text-4xl">Penawaran terbaru</h2><p class="mt-2 text-base text-th-muted">Lihat promo yang sedang berlangsung di Transhome.</p></div>
        <div role="group" aria-label="Saring promo dan event" class="flex flex-wrap gap-2">
        <button
          v-for="f in FILTER"
          :key="f.value"
          type="button"
          class="inline-flex min-h-11 items-center rounded-full border px-4 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          :class="filter.value === f.value ? 'border-th-navy bg-th-navy text-white' : 'border-th-border bg-white hover:border-th-navy'"
          :aria-pressed="filter.value === f.value ? 'true' : 'false'"
          @click="pilih(f)"
        >
          {{ f.label }}
        </button>
        </div>
      </div>
      <DaftarArtikel :tipe="filter.tipe" :gaya-promo="true" :judul="`Daftar ${filter.value ? filter.label.toLowerCase() : 'promo dan event'}`" judul-kosong="Belum ada promo saat ini" deskripsi-kosong="Promo dan event baru akan muncul di sini." />
    </section>
    <PageContactCta tone="navy" heading="Ada promo yang ingin kamu tanyakan?" description="Tim Transhome siap membantu menjelaskan penawaran yang sedang berlangsung." />
  </div>
</template>
