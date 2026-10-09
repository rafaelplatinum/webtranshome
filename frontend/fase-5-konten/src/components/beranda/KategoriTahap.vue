<script setup>
import { RouterLink } from 'vue-router'
import { Boxes } from 'lucide-vue-next'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { getCategoryTree } from '@/services/categoryService'
import JudulSeksi from './JudulSeksi.vue'
import PesanGagal from './PesanGagal.vue'

// "Belanja per tahap proyek": kategori induk berurutan (urutan diatur admin), dengan contoh subkategorinya.
const pohon = useMuat(getCategoryTree)
const contohAnak = (k) => k.children.slice(0, 3).map((c) => c.name).join(', ')
</script>

<template>
  <section aria-labelledby="judul-kategori-tahap" class="flex flex-col gap-4">
    <JudulSeksi id="judul-kategori-tahap" judul="Belanja per tahap proyek" tautan="/katalog" label-tautan="Semua kategori" />
    <div v-if="pohon.status.value === 'loading'" class="grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-8" aria-busy="true">
      <span class="sr-only">Memuat kategori…</span>
      <Skeleton v-for="n in 8" :key="n" class="h-40 rounded-xl" />
    </div>
    <PesanGagal v-else-if="pohon.status.value === 'error'" pesan="Kategori belum bisa dimuat." @coba-lagi="pohon.muat" />
    <ul v-else class="grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-8" data-kategori-tahap>
      <li v-for="k in pohon.data.value" :key="k.id">
        <RouterLink
          :to="`/katalog/${k.slug}`"
          class="flex h-full flex-col gap-3 rounded-xl border border-border bg-surface p-3 transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <span class="flex aspect-4/3 items-center justify-center overflow-hidden rounded-lg bg-subtle text-faint">
            <img v-if="k.image_url" :src="k.image_url" alt="" width="160" height="120" loading="lazy" class="size-full object-cover" />
            <Boxes v-else class="size-7" aria-hidden="true" />
          </span>
          <span class="flex flex-col gap-0.5">
            <strong class="text-sm leading-snug font-bold">{{ k.name }}</strong>
            <span v-if="k.children.length" class="line-clamp-2 text-xs text-muted">{{ contohAnak(k) }}</span>
          </span>
        </RouterLink>
      </li>
    </ul>
  </section>
</template>
