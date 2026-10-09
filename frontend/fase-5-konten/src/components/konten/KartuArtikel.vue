<script setup>
import { RouterLink } from 'vue-router'
import { Newspaper } from 'lucide-vue-next'
import { formatTanggal } from '@/utils/format'
import { TIPE_ARTIKEL, tautanArtikel } from '@/utils/artikel'

defineProps({
  // Item GET /articles: { type, title, slug, thumbnail_url, excerpt, published_at }
  artikel: { type: Object, required: true },
  // Ringkas: tanpa kutipan (mis. di Beranda).
  ringkas: { type: Boolean, default: false },
})
</script>

<template>
  <article class="relative flex h-full flex-col overflow-hidden rounded-xl border border-border bg-surface transition-colors duration-150 hover:border-border-strong">
    <span class="flex aspect-video items-center justify-center bg-subtle text-faint">
      <img v-if="artikel.thumbnail_url" :src="artikel.thumbnail_url" alt="" width="800" height="450" loading="lazy" class="size-full object-cover" />
      <Newspaper v-else class="size-8" aria-hidden="true" />
    </span>
    <div class="flex flex-1 flex-col gap-1.5 p-4">
      <span class="self-start rounded px-2 py-0.5 text-xs font-bold" :class="TIPE_ARTIKEL[artikel.type]?.kelas">{{ TIPE_ARTIKEL[artikel.type]?.label ?? artikel.type }}</span>
      <h3 class="text-base leading-snug font-bold">
        <!-- Tautan meliputi seluruh kartu (::after), satu tab stop per kartu. -->
        <RouterLink
          :to="tautanArtikel(artikel.slug)"
          class="rounded-sm after:absolute after:inset-0 after:content-[''] hover:text-primary focus-visible:outline-none focus-visible:after:rounded-xl focus-visible:after:ring-2 focus-visible:after:ring-ring"
        >
          {{ artikel.title }}
        </RouterLink>
      </h3>
      <time v-if="artikel.published_at" :datetime="artikel.published_at" class="text-sm text-muted">{{ formatTanggal(artikel.published_at) }}</time>
      <p v-if="!ringkas && artikel.excerpt" class="line-clamp-3 text-sm text-muted">{{ artikel.excerpt }}</p>
    </div>
  </article>
</template>
