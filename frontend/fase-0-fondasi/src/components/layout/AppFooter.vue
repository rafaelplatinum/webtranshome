<script setup>
// Footer (desain dari redesign Dylan, 6 Nov 2026, diminta pemilik): kolom Belanja, Informasi, toko, belanja online,
// sosial media. Tombol WhatsApp tetap langsung ke WhatsApp (pesan umum), sesuai tampilan lama.
import { computed } from 'vue'
import { Facebook, Instagram, MapPin, Phone, ShoppingBag, Store, Youtube } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'
import LogoTranshome from '@/components/layout/LogoTranshome.vue'
import WhatsAppButton from './WhatsAppButton.vue'
import { CONTACT } from '@/config/contact'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useSettingsStore } from '@/stores/settings'
import { formatNomorHp } from '@/utils/format'

const BELANJA = [
  { label: 'Semua produk', to: '/katalog' },
  { label: 'Kategori', to: '/katalog' },
  { label: 'Promo', to: '/promo' },
  { label: 'Keranjang belanja', to: '/keranjang' },
]
const INFORMASI = [
  { label: 'Tentang kami', to: '/info-toko' },
  { label: 'Membership', to: '/membership' },
  { label: 'Project', to: '/jasa' },
  { label: 'Fasilitas', to: '/info-toko#fasilitas' },
  { label: 'Tips & artikel', to: '/artikel' },
]
const SOSIAL = [
  { label: 'Instagram', url: CONTACT.social.instagram, icon: Instagram },
  { label: 'Facebook', url: CONTACT.social.facebook, icon: Facebook },
  { label: 'TikTok', url: CONTACT.social.tiktok },
  { label: 'YouTube', url: CONTACT.social.youtube, icon: Youtube },
]
const MARKETPLACE = [
  { label: 'Shopee', url: CONTACT.marketplace.shopee, icon: ShoppingBag },
  { label: 'Tokopedia', url: CONTACT.marketplace.tokopedia, icon: Store },
]
const settings = useSettingsStore()
const { generalLink } = useWhatsApp()
// Nomor yang tampil = nomor yang dipakai tombol (pengaturan toko); cadangan dari config kontak.
const nomorWa = computed(() => (settings.waNumber ? formatNomorHp(settings.waNumber) : CONTACT.whatsappDisplay))
const alamat = computed(() => settings.data?.address || CONTACT.address)
const jamBuka = computed(() => settings.data?.opening_hours || CONTACT.openingHours)
const peta = computed(() => settings.data?.maps_url)
const tahun = new Date().getFullYear()
const classesLink = 'rounded-sm text-on-dark/85 transition-colors hover:text-primary-on-dark hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <footer class="bg-secondary text-sm text-on-dark">
    <div class="mx-auto grid max-w-screen-2xl gap-x-8 gap-y-10 px-4 py-10 sm:grid-cols-2 sm:px-6 lg:grid-cols-[1.3fr_0.8fr_0.8fr_1.25fr_0.8fr] lg:gap-10 lg:px-8 lg:py-14">
      <section class="flex flex-col items-start gap-4">
        <RouterLink to="/" aria-label="Transhome, ke beranda" class="rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"><LogoTranshome gelap /></RouterLink>
        <p class="max-w-xs leading-relaxed">{{ CONTACT.description }}</p>
        <p class="font-bold text-secondary-foreground">Ikuti & hubungi kami</p>
        <div class="flex flex-wrap gap-2">
          <template v-for="item in SOSIAL" :key="item.label">
            <a v-if="item.url" :href="item.url" target="_blank" rel="noopener noreferrer" class="inline-flex min-h-10 items-center gap-2 rounded-md border border-on-dark-muted/50 px-3 transition-colors hover:border-on-dark hover:bg-secondary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark">
              <component :is="item.icon" v-if="item.icon" class="size-4" aria-hidden="true" />
              <svg v-else viewBox="0 0 24 24" fill="currentColor" class="size-4" aria-hidden="true"><path d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z" /></svg>
              {{ item.label }}<span class="sr-only">(tab baru)</span>
            </a>
            <button v-else type="button" disabled :aria-label="`${item.label}, tautan akan ditambahkan`" class="inline-flex min-h-10 items-center gap-2 rounded-md border border-on-dark-muted/50 px-3 text-on-dark/85">
              <component :is="item.icon" v-if="item.icon" class="size-4" aria-hidden="true" />
              <svg v-else viewBox="0 0 24 24" fill="currentColor" class="size-4" aria-hidden="true"><path d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z" /></svg>
              {{ item.label }}
            </button>
          </template>
        </div>
      </section>

      <nav aria-label="Belanja" class="flex flex-col items-start gap-3">
        <h2 class="mb-1 font-bold text-secondary-foreground">Belanja</h2>
        <RouterLink v-for="item in BELANJA" :key="item.label" :to="item.to" :class="classesLink">{{ item.label }}</RouterLink>
      </nav>

      <nav aria-label="Informasi" class="flex flex-col items-start gap-3">
        <h2 class="mb-1 font-bold text-secondary-foreground">Informasi</h2>
        <RouterLink v-for="item in INFORMASI" :key="item.label" :to="item.to" :class="classesLink">{{ item.label }}</RouterLink>
      </nav>

      <section class="flex flex-col items-start gap-3">
        <h2 class="mb-1 font-bold text-secondary-foreground">{{ CONTACT.storeName }}</h2>
        <a v-if="peta" :href="peta" target="_blank" rel="noopener noreferrer" class="flex items-start gap-2 rounded-sm hover:text-primary-on-dark focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark">
          <MapPin class="mt-0.5 size-4 shrink-0" aria-hidden="true" /><address class="not-italic">{{ alamat }}</address>
        </a>
        <p v-else class="flex items-start gap-2"><MapPin class="mt-0.5 size-4 shrink-0" aria-hidden="true" /><address class="not-italic">{{ alamat }}</address></p>
        <p>{{ jamBuka }}</p>
        <a :href="`tel:${CONTACT.telephone.replaceAll(' ', '')}`" class="inline-flex items-center gap-2 hover:text-primary-on-dark"><Phone class="size-4" aria-hidden="true" />{{ CONTACT.telephone }}</a>
        <div class="flex flex-col items-start gap-2">
          <p class="text-on-dark-muted">WhatsApp {{ nomorWa }}</p>
          <WhatsAppButton :href="generalLink()" :label="`WhatsApp ${CONTACT.storeName}`" />
        </div>
      </section>

      <section class="flex flex-col items-start gap-3">
        <h2 class="mb-1 font-bold text-secondary-foreground">Belanja online</h2>
        <template v-for="item in MARKETPLACE" :key="item.label">
          <a v-if="item.url" :href="item.url" target="_blank" rel="noopener noreferrer" class="inline-flex min-h-10 items-center gap-2 rounded-md border border-on-dark-muted/50 px-3 font-semibold transition-colors hover:border-on-dark hover:bg-secondary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark">
            <component :is="item.icon" class="size-4" aria-hidden="true" />{{ item.label }}<span class="sr-only">(tab baru)</span>
          </a>
          <button v-else type="button" disabled :aria-label="`${item.label}, tautan akan ditambahkan`" class="inline-flex min-h-10 items-center gap-2 rounded-md border border-on-dark-muted/50 px-3 font-semibold text-on-dark/85">
            <component :is="item.icon" class="size-4" aria-hidden="true" />{{ item.label }}
          </button>
        </template>
      </section>
    </div>

    <div class="border-t border-white/10">
      <p class="mx-auto max-w-screen-2xl px-4 py-5 text-xs text-on-dark-muted sm:px-6 lg:px-8">© {{ tahun }} CV Trans Home · Sleman, Yogyakarta</p>
    </div>
  </footer>
</template>
