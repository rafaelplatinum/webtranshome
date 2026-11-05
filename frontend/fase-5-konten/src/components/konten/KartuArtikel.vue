<script setup>
import { RouterLink } from 'vue-router'
import { ArrowUpRight } from 'lucide-vue-next'
import { formatTanggal } from '@/utils/format'
import { TIPE_ARTIKEL, tautanArtikel } from '@/utils/artikel'
import MediaPlaceholder from '@/components/public/MediaPlaceholder.vue'

defineProps({
  // Item GET /articles: { type, title, slug, thumbnail_url, excerpt, published_at }
  artikel: { type: Object, required: true },
  // Ringkas: tanpa kutipan (mis. di Beranda).
  ringkas: { type: Boolean, default: false },
  gayaPromo: { type: Boolean, default: false },
  featured: { type: Boolean, default: false },
})
</script>

<template>
  <article v-if="featured" class="th-card relative grid overflow-hidden lg:grid-cols-[1.15fr_0.85fr]">
    <MediaPlaceholder :src="artikel.thumbnail_url" tone="peach" class="min-h-64 aspect-video lg:min-h-[380px] lg:aspect-auto" />
    <div class="flex flex-col justify-center gap-4 p-6 sm:p-8 lg:p-10">
      <span class="th-eyebrow">Artikel pilihan<span v-if="artikel.published_at"> · {{ formatTanggal(artikel.published_at) }}</span></span>
      <h2 class="text-2xl leading-tight font-extrabold tracking-tight sm:text-3xl">{{ artikel.title }}</h2>
      <p v-if="artikel.excerpt" class="text-base leading-relaxed text-th-muted">{{ artikel.excerpt }}</p>
      <RouterLink :to="tautanArtikel(artikel.slug)" class="relative z-10 inline-flex min-h-11 items-center gap-2 self-start font-bold text-th-orange hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Baca artikel <ArrowUpRight class="size-5" aria-hidden="true" /></RouterLink>
    </div>
  </article>

  <article v-else class="th-card relative flex h-full flex-col overflow-hidden">
    <span class="relative block aspect-video overflow-hidden">
      <MediaPlaceholder :src="artikel.thumbnail_url" tone="sand" class="size-full" aria-hidden="true" />
      <span v-if="gayaPromo" class="absolute top-3 left-3 rounded-full bg-th-yellow px-3 py-1 text-[13px] font-extrabold text-th-navy">{{ TIPE_ARTIKEL[artikel.type]?.label ?? artikel.type }}</span>
    </span>
    <div class="flex flex-1 flex-col gap-2 p-5 sm:p-6">
      <span v-if="!gayaPromo" class="self-start rounded px-2 py-0.5 text-[13px] font-bold" :class="TIPE_ARTIKEL[artikel.type]?.kelas">{{ TIPE_ARTIKEL[artikel.type]?.label ?? artikel.type }}</span>
      <h3 class="text-lg leading-snug font-extrabold">
        <RouterLink :to="tautanArtikel(artikel.slug)" class="rounded-sm after:absolute after:inset-0 after:content-[''] hover:text-th-orange focus-visible:outline-none focus-visible:after:rounded-[20px] focus-visible:after:ring-2 focus-visible:after:ring-ring">{{ artikel.title }}</RouterLink>
      </h3>
      <time v-if="artikel.published_at" :datetime="artikel.published_at" class="text-sm text-th-muted">{{ formatTanggal(artikel.published_at) }}</time>
      <p v-if="!ringkas && artikel.excerpt" class="line-clamp-3 text-base leading-relaxed text-th-muted">{{ artikel.excerpt }}</p>
      <span v-if="gayaPromo" class="mt-auto inline-flex items-center gap-1 pt-2 font-bold text-th-orange">Lihat detail <ArrowUpRight class="size-4" aria-hidden="true" /></span>
    </div>
  </article>
</template>
