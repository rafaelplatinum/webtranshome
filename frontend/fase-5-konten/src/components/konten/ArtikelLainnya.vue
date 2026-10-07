<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import Skeleton from '@/components/ui/Skeleton.vue'
import PesanGagal from '@/components/beranda/PesanGagal.vue'
import { useMuat } from '@/composables/useMuat'
import { listArticles } from '@/services/kontenService'
import { daftarAsal } from '@/utils/artikel'
import KartuArtikel from './KartuArtikel.vue'

// Di bawah detail: 3 terbaru dari kelompok yang sama (promo & event, atau artikel), tanpa yang sedang dibuka.
const props = defineProps({
  artikel: { type: Object, required: true },
})

const promo = props.artikel.type === 'PROMO' || props.artikel.type === 'EVENT'
const asal = daftarAsal(props.artikel.type)
const judul = promo ? 'Promo dan event lainnya' : 'Artikel lainnya'

const lainnya = useMuat(async () => {
  const { data } = await listArticles({ type: promo ? 'PROMO,EVENT' : 'ARTIKEL', per_page: 4 })
  return data.filter((a) => a.id !== props.artikel.id).slice(0, 3)
})

const tampil = computed(() => lainnya.status.value !== 'ready' || lainnya.data.value.length > 0)
</script>

<template>
  <section v-if="tampil" aria-labelledby="judul-lainnya" class="flex flex-col gap-4">
    <div class="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
      <h2 id="judul-lainnya" class="text-xl font-extrabold">{{ judul }}</h2>
      <RouterLink
        :to="asal.to"
        class="inline-flex min-h-11 items-center rounded-sm text-sm font-bold text-primary hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
      >
        Lihat semua
      </RouterLink>
    </div>
    <div v-if="lainnya.status.value === 'loading'" class="grid gap-4 md:grid-cols-3" aria-busy="true">
      <span class="sr-only">Memuat…</span>
      <Skeleton v-for="n in 3" :key="n" class="h-56 rounded-xl" />
    </div>
    <PesanGagal v-else-if="lainnya.status.value === 'error'" :pesan="`${judul} belum bisa dimuat.`" @coba-lagi="lainnya.muat" />
    <ul v-else class="grid gap-4 md:grid-cols-3" data-artikel-lainnya>
      <li v-for="a in lainnya.data.value" :key="a.id">
        <KartuArtikel :artikel="a" ringkas />
      </li>
    </ul>
  </section>
</template>
