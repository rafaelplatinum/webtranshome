<script setup>
import { RouterLink } from 'vue-router'
import Skeleton from '@/components/ui/Skeleton.vue'
import KartuArtikel from '@/components/konten/KartuArtikel.vue'
import { useMuat } from '@/composables/useMuat'
import { listArticles } from '@/services/kontenService'
import PesanGagal from './PesanGagal.vue'

// "Promo dan event" (latar kuning, sketsa "Main"): 3 terbaru yang sudah terbit. Disembunyikan bila belum ada.
const promo = useMuat(async () => (await listArticles({ type: 'PROMO,EVENT', per_page: 3 })).data)
</script>

<template>
  <section
    v-if="promo.status.value !== 'ready' || promo.data.value.length"
    aria-labelledby="judul-promo-event"
    class="flex flex-col gap-5 rounded-xl bg-accent p-5 text-accent-foreground sm:p-8"
  >
    <div class="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
      <h2 id="judul-promo-event" class="text-2xl font-extrabold">Promo dan event</h2>
      <RouterLink
        to="/promo"
        class="inline-flex min-h-11 items-center rounded-sm text-sm font-bold underline-offset-2 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-foreground md:min-h-0"
      >
        Semua promo
      </RouterLink>
    </div>
    <div v-if="promo.status.value === 'loading'" class="grid gap-4 sm:grid-cols-3" aria-busy="true">
      <span class="sr-only">Memuat promo…</span>
      <Skeleton v-for="n in 3" :key="n" class="h-56 rounded-xl" />
    </div>
    <PesanGagal v-else-if="promo.status.value === 'error'" pesan="Promo belum bisa dimuat." @coba-lagi="promo.muat" />
    <ul v-else class="grid gap-4 sm:grid-cols-3" data-promo-event>
      <li v-for="a in promo.data.value" :key="a.id">
        <KartuArtikel :artikel="a" ringkas />
      </li>
    </ul>
  </section>
</template>
