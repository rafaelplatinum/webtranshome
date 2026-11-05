<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight } from 'lucide-vue-next'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { listBanners } from '@/services/kontenService'
import { useAuthStore } from '@/stores/auth'
import BannerSlider from './BannerSlider.vue'
import PesanGagal from './PesanGagal.vue'
import MediaPlaceholder from '@/components/public/MediaPlaceholder.vue'

// Sketsa "Main": slider banner (HOME_SLIDER) + kolom samping (banner HOME_SIDE dan kartu Trans Family).
// Banner hanya yang sedang tayang (rentang tanggal & aktif), sudah difilter API.
const auth = useAuthStore()
// Hanya member yang punya kartu; admin yang sedang masuk diperlakukan seperti pengunjung.
const member = computed(() => auth.hasRole('CUSTOMER'))
const slider = useMuat(() => listBanners('HOME_SLIDER'))
const samping = useMuat(() => listBanners('HOME_SIDE'))

// Tanpa banner tayang, kolom slider disembunyikan dan kartu samping memakai lebar penuh.
const tampilSlider = computed(() => slider.status.value !== 'ready' || slider.data.value.length > 0)
const bannerSamping = computed(() => (samping.status.value === 'ready' ? (samping.data.value[0] ?? null) : null))
// Hanya atribut yang ada yang dipasang: `href: undefined` akan menimpa href buatan RouterLink.
const tautanSamping = computed(() => {
  const link = bannerSamping.value?.link_url
  if (!link) return { is: 'div', atribut: {} }
  return link.startsWith('/') ? { is: RouterLink, atribut: { to: link } } : { is: 'a', atribut: { href: link, target: '_blank', rel: 'noopener' } }
})
</script>

<template>
  <section aria-labelledby="judul-promo-pilihan" class="bg-th-yellow py-12 text-th-navy sm:py-16 lg:py-20">
    <div class="th-container grid items-center gap-8 lg:grid-cols-2 lg:gap-12">
      <div class="flex flex-col items-start gap-5">
        <span class="th-eyebrow text-th-navy">Promo bulan ini</span>
        <h2 id="judul-promo-pilihan" class="max-w-xl text-3xl leading-tight font-extrabold tracking-tight sm:text-4xl lg:text-[44px]">Promo untuk rumah yang lebih nyaman.</h2>
        <p class="max-w-xl text-base leading-relaxed text-th-navy/80 sm:text-lg">Temukan promo dan acara terbaru Transhome, lalu tanya tim kami untuk detailnya.</p>
        <RouterLink to="/promo" class="inline-flex min-h-12 items-center rounded-xl bg-th-navy px-6 font-bold text-white transition-colors hover:bg-th-navy/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-th-orange">Lihat promo</RouterLink>
      </div>

      <div class="flex flex-col gap-4">
        <Skeleton v-if="slider.status.value === 'loading' && tampilSlider" class="h-72 rounded-[28px]" />
        <PesanGagal v-else-if="slider.status.value === 'error'" pesan="Banner promo belum bisa dimuat." @coba-lagi="slider.muat" />
        <BannerSlider v-else-if="slider.data.value.length" :banner="slider.data.value" />
        <MediaPlaceholder v-else tone="peach" class="h-64 rounded-[28px] sm:h-72" aria-hidden="true" />

        <div class="grid gap-4 sm:grid-cols-2">
          <component
            :is="tautanSamping.is"
            v-if="bannerSamping"
            v-bind="tautanSamping.atribut"
            class="group relative flex min-h-36 items-end overflow-hidden rounded-2xl border border-white/70 bg-white/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-th-orange"
            data-banner-samping
          >
            <img :src="bannerSamping.image_url" alt="" width="800" height="520" loading="lazy" class="absolute inset-0 size-full object-cover" />
            <span class="relative flex p-3"><span class="rounded-xl bg-white px-4 py-3 text-base font-bold text-th-navy group-hover:underline">{{ bannerSamping.title }}<span v-if="tautanSamping.is === 'a'" class="sr-only"> (tab baru)</span></span></span>
          </component>
          <RouterLink
            :to="member ? '/akun' : '/membership'"
            class="flex min-h-36 flex-col justify-end gap-2 rounded-2xl bg-th-navy p-5 text-white transition-transform duration-200 hover:-translate-y-1 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-th-orange"
            data-kartu-trans-family
          >
            <span class="text-[13px] font-bold tracking-wide uppercase text-th-yellow">Trans Family</span>
            <strong class="text-lg leading-snug">Belanja di toko, kumpulkan poin</strong>
            <span class="inline-flex items-center gap-1 text-sm font-bold">{{ member ? 'Lihat kartu member' : 'Daftar gratis' }}<ArrowRight class="size-4" aria-hidden="true" /></span>
          </RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>
