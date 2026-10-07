<script setup>
import { RouterLink } from 'vue-router'
import Skeleton from '@/components/ui/Skeleton.vue'
import RoomCarousel from '@/components/katalog/RoomCarousel.vue'
import { useMuat } from '@/composables/useMuat'
import { listRooms } from '@/services/catalogService'
import PesanGagal from './PesanGagal.vue'

// "Inspirasi ruangan": memakai carousel ruangan dari halaman /ruangan (‹ › hanya bila tidak muat).
const ruangan = useMuat(listRooms)
</script>

<template>
  <section
    v-if="ruangan.status.value !== 'ready' || ruangan.data.value.length"
    aria-labelledby="judul-inspirasi-ruangan"
    class="flex flex-col gap-4"
  >
    <div v-if="ruangan.status.value === 'loading'" class="flex flex-col gap-3" aria-busy="true">
      <span class="sr-only">Memuat ruangan…</span>
      <Skeleton class="h-7 w-48" />
      <div class="grid auto-cols-[minmax(14rem,1fr)] grid-flow-col gap-4 overflow-hidden">
        <Skeleton v-for="n in 4" :key="n" class="aspect-4/3 rounded-xl" />
      </div>
    </div>
    <template v-else-if="ruangan.status.value === 'error'">
      <h2 id="judul-inspirasi-ruangan" class="text-2xl font-extrabold">Inspirasi ruangan</h2>
      <PesanGagal pesan="Ruangan belum bisa dimuat." @coba-lagi="ruangan.muat" />
    </template>
    <RoomCarousel v-else :ruangan="ruangan.data.value" label="Inspirasi ruangan">
      <template #judul>
        <div class="flex flex-wrap items-baseline gap-x-4 gap-y-1">
          <h2 id="judul-inspirasi-ruangan" class="text-2xl font-extrabold">Inspirasi ruangan</h2>
          <RouterLink
            to="/ruangan"
            class="inline-flex min-h-11 items-center rounded-sm text-sm font-bold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
          >
            Semua ruangan
          </RouterLink>
        </div>
      </template>
    </RoomCarousel>
  </section>
</template>
