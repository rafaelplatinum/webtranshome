<script setup>
import { computed, defineAsyncComponent, ref, watch } from 'vue'
import { Expand, ImageOff } from 'lucide-vue-next'

const Lightbox = defineAsyncComponent(() => import('./Lightbox.vue'))

const props = defineProps({
  // product.images: [{ id, image_url, is_primary, sort_order }]
  images: { type: Array, default: () => [] },
  nama: { type: String, required: true },
})

// Foto utama selalu pertama, sisanya menurut sort_order.
const foto = computed(() =>
  [...props.images].sort((a, b) => Number(b.is_primary) - Number(a.is_primary) || a.sort_order - b.sort_order),
)

const aktif = ref(0)
const lightboxBuka = ref(false)
const lightboxDimuat = ref(false)
const tombolPerbesar = ref(null)

watch(() => props.images, () => (aktif.value = 0))

const teksAlt = (i) => (i === 0 ? props.nama : `${props.nama}, foto ${i + 1} dari ${foto.value.length}`)

function perbesar() {
  // Fokus dipindah ke tombol Perbesar dulu, supaya saat lightbox ditutup fokus kembali ke sana.
  tombolPerbesar.value?.focus({ preventScroll: true })
  lightboxDimuat.value = true
  lightboxBuka.value = true
}
</script>

<template>
  <div class="flex min-w-0 flex-col gap-3">
    <div class="relative overflow-hidden rounded-xl border border-border bg-surface">
      <!-- Klik foto = jalan pintas mouse; jalur keyboard lewat tombol Perbesar. -->
      <div v-if="foto.length" class="aspect-square cursor-zoom-in" @click="perbesar">
        <img :src="foto[aktif].image_url" :alt="teksAlt(aktif)" width="800" height="800" class="size-full object-contain" />
      </div>
      <div v-else class="flex aspect-square flex-col items-center justify-center gap-2 bg-subtle text-faint">
        <ImageOff class="size-10" aria-hidden="true" />
        <span class="text-sm">Foto belum ada</span>
      </div>
      <button
        v-if="foto.length"
        ref="tombolPerbesar"
        type="button"
        class="absolute right-3 bottom-3 inline-flex h-11 items-center gap-2 rounded-md border border-border-strong bg-surface px-3.5 text-sm font-semibold shadow-sm transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        @click="perbesar"
      >
        <Expand class="size-4" aria-hidden="true" />
        Perbesar<span class="sr-only"> foto {{ aktif + 1 }} dari {{ foto.length }}</span>
      </button>
    </div>

    <div v-if="foto.length > 1" role="group" aria-label="Pilih foto" class="-m-1 flex gap-2 overflow-x-auto p-1 scrollbar-none">
      <button
        v-for="(f, i) in foto"
        :key="f.id"
        type="button"
        class="size-16 shrink-0 overflow-hidden rounded-lg border-2 bg-subtle transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 sm:size-20"
        :class="i === aktif ? 'border-primary' : 'border-border hover:border-border-strong'"
        :aria-pressed="i === aktif ? 'true' : 'false'"
        :aria-label="`Foto ${i + 1} dari ${foto.length}`"
        @click="aktif = i"
      >
        <img :src="f.image_url" alt="" width="80" height="80" loading="lazy" class="size-full object-cover" />
      </button>
    </div>

    <Lightbox v-if="lightboxDimuat" v-model:open="lightboxBuka" v-model:index="aktif" :foto="foto" :nama="nama" />
  </div>
</template>
