<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { ChevronLeft, ChevronRight, Pause, Play } from 'lucide-vue-next'
import { useMuat } from '@/composables/useMuat'
import { listBanners } from '@/services/kontenService'

// Hero Beranda hanya menampilkan foto dari banner HOME_HERO. Beberapa banner aktif
// menjadi slide; judul admin hanya dipakai sebagai teks alternatif untuk pembaca layar.
const muat = useMuat(async () => (await listBanners('HOME_HERO')).filter((banner) => banner.image_url))
const banners = computed(() => muat.data.value ?? [])
const indeks = ref(0)
const dijeda = ref(false)
const disentuh = ref(false)
const kurangiGerak = window.matchMedia?.('(prefers-reduced-motion: reduce)').matches === true
const otomatis = computed(() => banners.value.length > 1 && !kurangiGerak && !dijeda.value)
let pengatur = null

watch(banners, (daftar) => {
  if (indeks.value >= daftar.length) indeks.value = 0
})

function ke(i) {
  const total = banners.value.length
  if (!total) return
  indeks.value = (i + total) % total
}

onMounted(() => {
  pengatur = window.setInterval(() => {
    if (otomatis.value && !disentuh.value && document.visibilityState === 'visible') ke(indeks.value + 1)
  }, 6000)
})
onBeforeUnmount(() => window.clearInterval(pengatur))

const internal = (banner) => banner.link_url?.startsWith('/')
</script>

<template>
  <section
    class="relative w-full overflow-hidden bg-secondary"
    :class="banners.length ? 'aspect-[12/5] min-h-[17rem] sm:min-h-[24rem] lg:min-h-[30rem]' : 'h-[min(62svh,800px)] min-h-[17rem]'"
    role="group"
    aria-roledescription="carousel"
    aria-label="Banner utama Transhome"
    data-hero-banner
    @mouseenter="disentuh = true"
    @mouseleave="disentuh = false"
    @focusin="disentuh = true"
    @focusout="disentuh = false"
  >
    <span class="sr-only" aria-live="polite">{{ banners.length ? `Banner ${indeks + 1} dari ${banners.length}` : 'Banner utama' }}</span>
    <Transition
      mode="out-in"
      enter-active-class="transition-opacity duration-500 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-300 ease-in"
      leave-to-class="opacity-0"
    >
      <div v-if="banners[indeks]" :key="banners[indeks].id" class="absolute inset-0" role="group" aria-roledescription="slide">
        <RouterLink v-if="internal(banners[indeks])" :to="banners[indeks].link_url" class="absolute inset-0 block focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-inset focus-visible:ring-primary">
          <img :src="banners[indeks].image_url" :alt="banners[indeks].title || 'Banner Transhome'" :fetchpriority="indeks === 0 ? 'high' : undefined" :loading="indeks === 0 ? 'eager' : 'lazy'" width="1920" height="800" class="size-full object-cover" />
        </RouterLink>
        <a v-else-if="banners[indeks].link_url" :href="banners[indeks].link_url" target="_blank" rel="noopener" class="absolute inset-0 block focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-inset focus-visible:ring-primary">
          <img :src="banners[indeks].image_url" :alt="banners[indeks].title || 'Banner Transhome'" :fetchpriority="indeks === 0 ? 'high' : undefined" :loading="indeks === 0 ? 'eager' : 'lazy'" width="1920" height="800" class="size-full object-cover" />
        </a>
        <img v-else :src="banners[indeks].image_url" :alt="banners[indeks].title || 'Banner Transhome'" :fetchpriority="indeks === 0 ? 'high' : undefined" :loading="indeks === 0 ? 'eager' : 'lazy'" width="1920" height="800" class="size-full object-cover" />
      </div>
      <div v-else key="empty" class="absolute inset-0" aria-hidden="true" />
    </Transition>

    <div v-if="banners.length > 1" class="absolute inset-x-0 bottom-0 flex items-center justify-center gap-2 bg-gradient-to-t from-black/45 to-transparent px-4 pt-12 pb-4 sm:justify-between sm:px-6 sm:pb-5 lg:px-10">
      <div class="hidden items-center gap-1.5 sm:flex" aria-label="Pilih banner">
        <button
          v-for="(banner, i) in banners"
          :key="banner.id"
          type="button"
          class="flex size-8 items-center justify-center rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white"
          :aria-label="`Tampilkan banner ${i + 1}`"
          :aria-current="i === indeks ? 'true' : undefined"
          @click="ke(i)"
        >
          <span class="h-1.5 rounded-full transition-all" :class="i === indeks ? 'w-6 bg-white' : 'w-1.5 bg-white/65'" />
        </button>
      </div>
      <div class="flex items-center gap-2">
        <button type="button" class="inline-flex size-10 items-center justify-center rounded-full border border-white/45 bg-black/25 text-white backdrop-blur-sm transition hover:bg-black/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white" :aria-label="dijeda ? 'Putar slide otomatis' : 'Jeda slide otomatis'" @click="dijeda = !dijeda">
          <component :is="dijeda ? Play : Pause" class="size-4" aria-hidden="true" />
        </button>
        <button type="button" class="inline-flex size-10 items-center justify-center rounded-full border border-white/45 bg-black/25 text-white backdrop-blur-sm transition hover:bg-black/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white" aria-label="Banner sebelumnya" @click="ke(indeks - 1)">
          <ChevronLeft class="size-5" aria-hidden="true" />
        </button>
        <button type="button" class="inline-flex size-10 items-center justify-center rounded-full border border-white/45 bg-black/25 text-white backdrop-blur-sm transition hover:bg-black/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white" aria-label="Banner berikutnya" @click="ke(indeks + 1)">
          <ChevronRight class="size-5" aria-hidden="true" />
        </button>
      </div>
    </div>
  </section>
</template>
