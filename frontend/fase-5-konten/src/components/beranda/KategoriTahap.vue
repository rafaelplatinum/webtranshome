<script setup>
import { RouterLink } from 'vue-router'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { getCategoryTree } from '@/services/categoryService'
import JudulSeksi from './JudulSeksi.vue'
import PesanGagal from './PesanGagal.vue'

const TILE = ['bg-th-sand', 'bg-th-blue', 'bg-th-yellow', 'bg-th-sage', 'bg-th-peach']

// "Belanja per tahap proyek": kategori induk berurutan (urutan diatur admin), dengan contoh subkategorinya.
const pohon = useMuat(getCategoryTree)
const contohAnak = (k) => k.children.slice(0, 3).map((c) => c.name).join(', ')
</script>

<template>
  <section aria-labelledby="judul-kategori-tahap" class="flex flex-col gap-4">
    <JudulSeksi id="judul-kategori-tahap" judul="Belanja per tahap proyek" tautan="/katalog" label-tautan="Semua kategori" />
    <div v-if="pohon.status.value === 'loading'" class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5" aria-busy="true">
      <span class="sr-only">Memuat kategori…</span>
      <Skeleton v-for="n in 5" :key="n" class="h-52 rounded-[20px]" />
    </div>
    <PesanGagal v-else-if="pohon.status.value === 'error'" pesan="Kategori belum bisa dimuat." @coba-lagi="pohon.muat" />
    <ul v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5 lg:gap-5" data-kategori-tahap>
      <li v-for="(k, i) in pohon.data.value" :key="k.id" v-reveal="i % 5">
        <RouterLink
          :to="`/katalog/${k.slug}`"
          class="group relative flex min-h-48 h-full flex-col justify-end overflow-hidden rounded-[20px] border border-th-border p-4 transition-all duration-200 hover:-translate-y-1 hover:shadow-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:min-h-56 sm:p-5"
          :class="[TILE[i % TILE.length], k.image_url && 'text-white']"
        >
          <span v-if="k.image_url" class="absolute inset-0">
            <img :src="k.image_url" alt="" width="400" height="320" loading="lazy" class="size-full object-cover transition-transform duration-200 group-hover:scale-105" />
            <span class="absolute inset-0 bg-gradient-to-t from-th-navy/85 via-th-navy/10 to-transparent"></span>
          </span>
          <span v-else aria-hidden="true" class="absolute -top-10 -right-8 size-36 rounded-full border border-white/50"></span>
          <span class="relative flex flex-col gap-1" :class="k.image_url ? 'text-white' : 'text-th-navy'">
            <strong class="text-base leading-snug font-extrabold sm:text-lg">{{ k.name }}</strong>
            <span v-if="k.children.length" class="line-clamp-2 text-sm opacity-80">{{ contohAnak(k) }}</span>
          </span>
        </RouterLink>
      </li>
    </ul>
  </section>
</template>
