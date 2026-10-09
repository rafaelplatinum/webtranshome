<script setup>
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Pause, Play } from 'lucide-vue-next'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { listBrands } from '@/services/catalogService'
import { formatAngka } from '@/utils/format'
import PesanGagal from './PesanGagal.vue'

// "Brand pilihan": SEMUA brand aktif. Lebih dari 6 brand → logo berjalan pelan (CSS murni, tanpa library):
// berhenti saat disentuh/disorot, tombol Jeda, dan diam bila pengguna meminta animasi dikurangi.
// Saat dijeda atau ada elemen di dalamnya yang difokus, daftar berubah jadi baris yang bisa digulir sendiri.
const MIN_BERJALAN = 7
const DETIK_PER_BRAND = 3

const merek = useMuat(async () => [...(await listBrands())].sort((a, b) => a.name.localeCompare(b.name, 'id')))
const kurangiGerak = window.matchMedia?.('(prefers-reduced-motion: reduce)').matches === true
const berjalan = ref(!kurangiGerak)

const jumlah = computed(() => merek.data.value?.length ?? 0)
const bisaBerjalan = computed(() => jumlah.value >= MIN_BERJALAN && !kurangiGerak)
const bergerak = computed(() => bisaBerjalan.value && berjalan.value)
const durasi = computed(() => `${Math.max(20, jumlah.value * DETIK_PER_BRAND)}s`)

// Fokus keyboard masuk ke daftar: berhenti dan jadi baris biasa supaya tautan yang difokus terlihat.
function onFocusIn() {
  if (bergerak.value) berjalan.value = false
}

const kelasKartu =
  'flex h-20 w-40 items-center justify-center rounded-xl border border-border bg-surface px-3 text-center font-bold transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <section v-if="merek.status.value !== 'ready' || jumlah" aria-labelledby="judul-brand-pilihan" class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2">
      <h2 id="judul-brand-pilihan" class="text-2xl font-extrabold">Brand pilihan</h2>
      <div v-if="merek.status.value === 'ready'" class="flex items-center gap-2">
        <button
          v-if="bisaBerjalan"
          type="button"
          :aria-pressed="berjalan ? 'false' : 'true'"
          class="inline-flex min-h-11 items-center gap-1.5 rounded-md px-2 text-sm font-semibold text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-9"
          data-jeda-brand
          @click="berjalan = !berjalan"
        >
          <component :is="berjalan ? Pause : Play" class="size-4" aria-hidden="true" />
          {{ berjalan ? 'Jeda' : 'Putar' }}<span class="sr-only"> logo brand berjalan</span>
        </button>
        <RouterLink
          to="/brand"
          class="inline-flex min-h-11 items-center rounded-sm text-sm font-bold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
        >
          Semua brand ({{ formatAngka(jumlah) }})
        </RouterLink>
      </div>
    </div>

    <div v-if="merek.status.value === 'loading'" class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6" aria-busy="true">
      <span class="sr-only">Memuat brand…</span>
      <Skeleton v-for="n in 6" :key="n" class="h-20 rounded-xl" />
    </div>
    <PesanGagal v-else-if="merek.status.value === 'error'" pesan="Brand belum bisa dimuat." @coba-lagi="merek.muat" />

    <!-- Sedikit brand: kisi biasa. -->
    <ul v-else-if="!bisaBerjalan && jumlah < MIN_BERJALAN" class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6" data-brand-pilihan>
      <li v-for="b in merek.data.value" :key="b.id">
        <RouterLink :to="`/brand/${b.slug}`" :class="[kelasKartu, 'w-full']">
          <img v-if="b.logo_url" :src="b.logo_url" :alt="b.name" width="120" height="48" loading="lazy" class="max-h-12 w-auto object-contain" />
          <span v-else>{{ b.name }}</span>
        </RouterLink>
      </li>
    </ul>

    <!-- Banyak brand: berjalan (dua salinan berurutan, digeser -50% lalu berulang) atau baris yang bisa digulir. -->
    <div
      v-else
      class="relative"
      :class="bergerak ? 'overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_4%,black_96%,transparent)]' : '-mx-1 overflow-x-auto px-1 pb-2'"
      data-brand-berjalan
      :data-bergerak="bergerak ? 'ya' : 'tidak'"
      @focusin="onFocusIn"
    >
      <div class="flex w-max" :class="bergerak && 'brand-berjalan'" :style="{ '--durasi-brand': durasi }">
        <ul class="flex gap-3 pr-3" data-brand-pilihan>
          <li v-for="b in merek.data.value" :key="b.id">
            <RouterLink :to="`/brand/${b.slug}`" :class="kelasKartu">
              <img v-if="b.logo_url" :src="b.logo_url" :alt="b.name" width="120" height="48" loading="lazy" class="max-h-12 w-auto object-contain" />
              <span v-else>{{ b.name }}</span>
            </RouterLink>
          </li>
        </ul>
        <!-- Salinan kedua hanya untuk sambungan visual: disembunyikan dari pembaca layar dan keyboard. -->
        <ul v-if="bergerak" class="flex gap-3 pr-3" aria-hidden="true" inert>
          <li v-for="b in merek.data.value" :key="b.id">
            <span :class="kelasKartu">
              <img v-if="b.logo_url" :src="b.logo_url" alt="" width="120" height="48" loading="lazy" class="max-h-12 w-auto object-contain" />
              <span v-else>{{ b.name }}</span>
            </span>
          </li>
        </ul>
      </div>
    </div>
  </section>
</template>
