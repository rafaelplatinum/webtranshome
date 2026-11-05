<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { MapPin, MessageCircle, ShoppingCart } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import MegaMenu from './MegaMenu.vue'
import SearchBox from './SearchBox.vue'
import UserMenu from './UserMenu.vue'
import LogoTranshome from '@/components/layout/LogoTranshome.vue'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { useCartStore } from '@/stores/cart'
import { useWhatsAppContactStore } from '@/stores/whatsappContact'
import { CONTACT } from '@/config/contact'

const NAVIGASI = [
  { label: 'Ruangan', to: '/ruangan' },
  { label: 'Brand', to: '/brand' },
  { label: 'Jasa', to: '/jasa' },
  { label: 'Promo', to: '/promo' },
  { label: 'Artikel', to: '/artikel' },
  { label: 'Membership', to: '/membership' },
  { label: 'Tentang Kami', to: '/info-toko' },
]

const route = useRoute()
const kelasNavigasiAktif = (tujuan) => route.path === tujuan || route.path.startsWith(`${tujuan}/`)
const auth = useAuthStore()
const settings = useSettingsStore()
const cart = useCartStore()
const whatsappForm = useWhatsAppContactStore()
const linkMasuk = computed(() => ({ path: '/masuk', query: { redirect: route.fullPath } }))

// Header menempel di atas saat halaman digulir. Di bawah 1024 px baris alamat ikut tergulir hilang
// (top negatif setinggi baris itu) supaya header yang menempel tidak memakan seperempat layar HP.
// Tinggi diukur ulang setiap baris turun / lebar layar berubah:
// --tinggi-bar-atas → top header; --tinggi-header (bagian yang menempel) → scroll-padding-top (main.css);
// --tinggi-header-penuh → batas tinggi panel mega menu.
const kepala = ref(null)
const barAtas = ref(null)
const VARIABEL = ['--tinggi-bar-atas', '--tinggi-header', '--tinggi-header-penuh']
let pengamat = null
let layarLebar = null
let ukur = null
onMounted(() => {
  const akar = document.documentElement
  layarLebar = window.matchMedia('(min-width: 1024px)')
  ukur = () => {
    const penuh = Math.ceil(kepala.value.offsetHeight)
    const bar = layarLebar.matches ? 0 : Math.ceil(barAtas.value.offsetHeight)
    akar.style.setProperty('--tinggi-bar-atas', `${bar}px`)
    akar.style.setProperty('--tinggi-header', `${penuh - bar}px`)
    akar.style.setProperty('--tinggi-header-penuh', `${penuh}px`)
  }
  pengamat = new ResizeObserver(ukur)
  pengamat.observe(kepala.value)
  pengamat.observe(barAtas.value)
  layarLebar.addEventListener('change', ukur)
})
onBeforeUnmount(() => {
  pengamat?.disconnect()
  layarLebar?.removeEventListener('change', ukur)
  VARIABEL.forEach((v) => document.documentElement.style.removeProperty(v))
})

const kelasTautanGelap =
  'inline-flex items-center gap-1.5 rounded-sm font-semibold text-secondary-foreground/90 transition-colors duration-150 hover:text-primary-on-dark focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <!-- view-transition-name: header tetap diam saat pindah halaman (hanya isi halaman yang crossfade). -->
  <header ref="kepala" class="[view-transition-name:header-publik] sticky top-[calc(-1*var(--tinggi-bar-atas,0px))] z-30 lg:top-0">
    <div ref="barAtas" class="border-b border-white/10 bg-secondary text-xs text-on-dark sm:text-sm">
      <div class="mx-auto flex max-w-screen-2xl flex-wrap items-center justify-between gap-x-6 gap-y-1 px-4 py-2 sm:px-6 lg:px-8">
        <p class="flex min-w-0 items-center gap-2 tracking-[0.01em]">
          <MapPin class="size-4 shrink-0 text-primary-on-dark" aria-hidden="true" />
          <span>{{ settings.data?.address || CONTACT.address }} · Buka {{ settings.data?.opening_hours || CONTACT.openingHours }}</span>
        </p>
        <div class="flex items-center gap-4 sm:gap-6">
          <RouterLink to="/info-toko" :class="kelasTautanGelap">Info toko</RouterLink>
          <button type="button" :class="kelasTautanGelap" @click="whatsappForm.buka()"><MessageCircle class="size-3.5 text-whatsapp" aria-hidden="true" />WhatsApp CS</button>
        </div>
      </div>
    </div>

    <div class="border-b border-border/80 bg-surface">
      <div class="mx-auto flex max-w-screen-2xl flex-wrap items-center gap-x-5 gap-y-3 px-4 py-3 sm:px-6 lg:gap-x-6 lg:px-8 lg:py-4">
        <RouterLink
          to="/"
          class="shrink-0 rounded-md text-2xl font-bold tracking-tight focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          aria-label="Transhome, ke beranda"
        >
          <LogoTranshome />
        </RouterLink>

        <div class="ml-auto flex items-center gap-2 lg:order-last lg:ml-0 lg:gap-3">
          <Button variant="ghost" size="sm" class="relative size-11 !px-0" to="/keranjang" :aria-label="`Keranjang belanja, ${cart.jumlah} produk`">
            <template #icon><ShoppingCart class="size-4.5" aria-hidden="true" /></template>
            <span v-if="cart.jumlah" class="absolute top-2 right-2 size-2 rounded-full bg-primary ring-2 ring-surface" aria-hidden="true"></span>
          </Button>
          <UserMenu v-if="auth.isLoggedIn" />
          <template v-else>
            <Button variant="outline" size="sm" class="sm:h-11" :to="linkMasuk">Masuk</Button>
            <Button size="sm" class="sm:h-11" to="/daftar"><span>Daftar<span class="hidden sm:inline"> member</span></span></Button>
          </template>
        </div>

        <div class="flex w-full min-w-0 items-center gap-3 lg:w-auto lg:flex-1">
          <MegaMenu />
          <SearchBox class="flex-1" placeholder="Cari SKU, nama produk, kategori, atau kebutuhan renovasi Anda..." />
        </div>
      </div>
    </div>

    <div class="hidden border-b border-white/10 bg-secondary text-on-dark lg:block">
      <nav aria-label="Navigasi utama" class="mx-auto flex max-w-screen-2xl items-center gap-1 overflow-x-auto px-6 lg:px-8">
          <RouterLink
            v-for="item in NAVIGASI.slice(0, 3)"
            :key="item.to"
            :to="item.to"
            class="relative inline-flex h-12 shrink-0 items-center px-4 text-sm font-semibold text-on-dark/85 transition-colors duration-150 after:absolute after:right-4 after:bottom-0 after:left-4 after:h-0.5 after:scale-x-0 after:bg-primary-on-dark after:transition-transform after:duration-200 hover:text-on-dark hover:after:scale-x-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
            :class="kelasNavigasiAktif(item.to) ? '!text-on-dark after:!scale-x-100' : ''"
            :aria-current="kelasNavigasiAktif(item.to) ? 'page' : undefined"
          >
            {{ item.label }}
          </RouterLink>
          <span class="mx-2 h-5 w-px shrink-0 bg-white/20" aria-hidden="true"></span>
          <RouterLink
            v-for="item in NAVIGASI.slice(3)"
            :key="item.to"
            :to="item.to"
            class="relative inline-flex h-12 shrink-0 items-center px-4 text-sm font-semibold text-on-dark/85 transition-colors duration-150 after:absolute after:right-4 after:bottom-0 after:left-4 after:h-0.5 after:scale-x-0 after:bg-primary-on-dark after:transition-transform after:duration-200 hover:text-on-dark hover:after:scale-x-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
            :class="kelasNavigasiAktif(item.to) ? '!text-on-dark after:!scale-x-100' : ''"
            :aria-current="kelasNavigasiAktif(item.to) ? 'page' : undefined"
          >
            {{ item.label }}
          </RouterLink>
      </nav>
      <div class="h-px bg-primary-on-dark/30" aria-hidden="true"></div>
    </div>
  </header>
</template>
