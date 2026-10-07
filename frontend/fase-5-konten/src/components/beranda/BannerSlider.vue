<script setup>
import { computed, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { ChevronLeft, ChevronRight, Pause, Play } from 'lucide-vue-next'
import { formatTanggal } from '@/utils/format'

const props = defineProps({
  // Banner tayang posisi HOME_SLIDER: [{ id, title, image_url, link_url, end_at }]
  banner: { type: Array, required: true },
})

// Ganti otomatis tiap 6 detik. Berhenti saat kursor/fokus di slider, saat dijeda, saat tab tidak terlihat,
// dan tidak pernah berjalan bila pengguna memilih "kurangi animasi".
const JEDA_MS = 6000
const id = useId()
const indeks = ref(0)
const dijeda = ref(false)
const disentuh = ref(false)
const kurangiGerak = typeof window !== 'undefined' && window.matchMedia?.('(prefers-reduced-motion: reduce)').matches
let pengatur = null

const total = computed(() => props.banner.length)
const otomatis = computed(() => total.value > 1 && !kurangiGerak && !dijeda.value)

watch(total, () => {
  if (indeks.value >= total.value) indeks.value = 0
})

function ke(i) {
  indeks.value = (i + total.value) % total.value
}

onMounted(() => {
  pengatur = setInterval(() => {
    if (otomatis.value && !disentuh.value && document.visibilityState === 'visible') ke(indeks.value + 1)
  }, JEDA_MS)
})
onBeforeUnmount(() => clearInterval(pengatur))

const internal = (b) => b.link_url?.startsWith('/')

const kelasTombol =
  'inline-flex size-11 items-center justify-center rounded-full bg-surface/90 text-foreground transition-colors duration-150 hover:bg-surface focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <div
    class="relative overflow-hidden rounded-xl bg-secondary text-secondary-foreground"
    role="group"
    aria-roledescription="carousel"
    :aria-labelledby="id"
    data-slider-banner
    @mouseenter="disentuh = true"
    @mouseleave="disentuh = false"
    @focusin="disentuh = true"
    @focusout="disentuh = false"
  >
    <h3 :id="id" class="sr-only">Banner promo</h3>
    <Transition
      mode="out-in"
      enter-active-class="transition-opacity duration-200 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-150 ease-in"
      leave-to-class="opacity-0"
    >
      <div
        :key="banner[indeks]?.id"
        role="group"
        aria-roledescription="slide"
        :aria-label="`${indeks + 1} dari ${total}`"
        class="grid min-h-80 sm:min-h-90"
        data-slide
      >
        <img
          :src="banner[indeks].image_url"
          alt=""
          width="1600"
          height="640"
          class="col-start-1 row-start-1 size-full object-cover"
          :loading="indeks === 0 ? 'eager' : 'lazy'"
        />
        <div class="col-start-1 row-start-1 flex items-end p-4 sm:p-6">
          <div class="flex max-w-md flex-col gap-2 rounded-lg bg-secondary/90 p-4 sm:p-5">
            <span class="self-start rounded bg-accent px-2.5 py-0.5 text-xs font-bold text-accent-foreground">PROMO</span>
            <p class="text-2xl leading-tight font-extrabold sm:text-3xl">{{ banner[indeks].title }}</p>
            <p v-if="banner[indeks].end_at" class="text-sm text-on-dark">Sampai {{ formatTanggal(banner[indeks].end_at) }}</p>
            <template v-if="banner[indeks].link_url">
              <RouterLink
                v-if="internal(banner[indeks])"
                :to="banner[indeks].link_url"
                class="mt-1 inline-flex h-11 items-center self-start rounded-md bg-primary px-4.5 text-sm font-bold text-primary-foreground transition-colors duration-150 hover:bg-primary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
              >
                Lihat promo
              </RouterLink>
              <a
                v-else
                :href="banner[indeks].link_url"
                target="_blank"
                rel="noopener"
                class="mt-1 inline-flex h-11 items-center self-start rounded-md bg-primary px-4.5 text-sm font-bold text-primary-foreground transition-colors duration-150 hover:bg-primary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
              >
                Lihat promo<span class="sr-only"> (tab baru)</span>
              </a>
            </template>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Kontrol di atas supaya tidak menutupi panel judul di bawah. -->
    <div v-if="total > 1" class="absolute inset-x-3 top-3 flex items-center justify-between gap-2">
      <div class="flex rounded-full bg-secondary/80 px-1" data-titik-banner>
        <button
          v-for="(b, i) in banner"
          :key="b.id"
          type="button"
          class="flex size-11 items-center justify-center rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark md:size-8"
          :aria-label="`Banner ${i + 1}: ${b.title}`"
          :aria-current="i === indeks ? 'true' : undefined"
          @click="ke(i)"
        >
          <span class="block h-1.5 rounded-full transition-colors duration-200" :class="i === indeks ? 'w-6 bg-surface' : 'w-1.5 bg-on-dark-muted'" />
        </button>
      </div>
      <div class="flex items-center gap-2">
        <button type="button" :class="kelasTombol" :aria-label="dijeda ? 'Putar banner otomatis' : 'Jeda banner otomatis'" @click="dijeda = !dijeda">
          <component :is="dijeda ? Play : Pause" class="size-4.5" aria-hidden="true" />
        </button>
        <button type="button" :class="kelasTombol" aria-label="Banner sebelumnya" @click="ke(indeks - 1)">
          <ChevronLeft class="size-5" aria-hidden="true" />
        </button>
        <button type="button" :class="kelasTombol" aria-label="Banner berikutnya" @click="ke(indeks + 1)">
          <ChevronRight class="size-5" aria-hidden="true" />
        </button>
      </div>
    </div>
  </div>
</template>
