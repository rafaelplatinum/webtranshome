<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { MapPin, ShoppingCart } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import MegaMenu from './MegaMenu.vue'
import SearchBox from './SearchBox.vue'
import UserMenu from './UserMenu.vue'
import LogoTranshome from '@/components/layout/LogoTranshome.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useAuthStore } from '@/stores/auth'
import { useKeranjangStore } from '@/stores/keranjang'
import { useSettingsStore } from '@/stores/settings'

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
const auth = useAuthStore()
const settings = useSettingsStore()
const keranjang = useKeranjangStore()
const { generalLink } = useWhatsApp()

const linkWhatsApp = computed(() => generalLink())
const linkMasuk = computed(() => ({ path: '/masuk', query: { redirect: route.fullPath } }))
// Keranjang (KEPUTUSAN #23): angka = banyaknya produk berbeda.
const labelKeranjang = computed(() => (keranjang.jumlahProduk ? `Keranjang belanja, ${keranjang.jumlahProduk} produk` : 'Keranjang belanja, masih kosong'))

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
  'rounded-sm font-semibold text-secondary-foreground hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <!-- view-transition-name: header tetap diam saat pindah halaman (hanya isi halaman yang crossfade). -->
  <header ref="kepala" class="[view-transition-name:header-publik] sticky top-[calc(-1*var(--tinggi-bar-atas,0px))] z-30 lg:top-0">
    <div ref="barAtas" class="bg-secondary text-sm text-on-dark">
      <div class="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-x-6 gap-y-1 px-4 py-2 sm:px-6">
        <p class="flex min-w-0 items-center gap-2">
          <MapPin class="size-4 shrink-0 text-primary-on-dark" aria-hidden="true" />
          <span v-if="settings.data">{{ settings.data.address }} · Buka {{ settings.data.opening_hours }}</span>
          <span v-else-if="settings.status === 'error'" class="text-on-dark-muted">Info toko belum bisa dimuat.</span>
          <span v-else class="text-on-dark-muted">Memuat info toko…</span>
        </p>
        <div class="flex gap-5">
          <RouterLink to="/info-toko" :class="kelasTautanGelap">Info toko</RouterLink>
          <a v-if="linkWhatsApp" :href="linkWhatsApp" target="_blank" rel="noopener" :class="kelasTautanGelap">WhatsApp CS</a>
        </div>
      </div>
    </div>

    <div class="border-b border-border bg-surface">
      <div class="mx-auto flex max-w-7xl flex-wrap items-center gap-x-6 gap-y-3 px-4 py-3 sm:px-6">
        <RouterLink
          to="/"
          class="rounded-md text-2xl font-extrabold tracking-tight focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          aria-label="Transhome, ke beranda"
        >
          <LogoTranshome />
        </RouterLink>

        <div class="ml-auto flex items-center gap-2 lg:order-last lg:ml-0">
          <!-- Menu keranjang (permintaan pemilik 6 Nov 2026): ikon + tulisan "Keranjang" di layar lebar, angka = banyaknya produk. -->
          <RouterLink
            to="/keranjang"
            class="inline-flex h-11 min-w-11 shrink-0 items-center justify-center gap-2 rounded-md px-2 text-foreground transition-colors duration-150 hover:bg-subtle hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:px-3"
            active-class="text-primary"
            :aria-label="labelKeranjang"
            data-ikon-keranjang
          >
            <span class="relative">
              <ShoppingCart class="size-6" aria-hidden="true" />
              <span
                v-if="keranjang.jumlahProduk"
                class="absolute -top-2 -right-2.5 flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1 text-xs font-bold text-primary-foreground tabular-nums ring-2 ring-surface"
                aria-hidden="true"
                data-angka-keranjang
              >
                {{ keranjang.jumlahProduk > 99 ? '99+' : keranjang.jumlahProduk }}
              </span>
            </span>
            <span class="hidden text-sm font-semibold lg:inline" aria-hidden="true">Keranjang</span>
          </RouterLink>
          <UserMenu v-if="auth.isLoggedIn" />
          <template v-else>
            <Button variant="outline" size="sm" class="sm:h-11" :to="linkMasuk">Masuk</Button>
            <Button size="sm" class="sm:h-11" to="/daftar"><span>Daftar<span class="hidden sm:inline"> member</span></span></Button>
          </template>
        </div>

        <div class="flex w-full min-w-0 items-center gap-3 lg:w-auto lg:flex-1">
          <MegaMenu />
          <SearchBox class="flex-1" />
        </div>

        <!-- Baris menu sendiri di bawah pencarian (6 Nov 2026): 7 menu + Keranjang tidak muat satu baris dengan kotak cari. -->
        <nav aria-label="Navigasi utama" class="hidden w-full items-center gap-x-7 border-t border-border pt-3 text-sm font-semibold lg:order-[10000] lg:flex">
          <RouterLink
            v-for="item in NAVIGASI"
            :key="item.to"
            :to="item.to"
            class="rounded-sm transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            active-class="text-primary"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
      </div>
    </div>
  </header>
</template>
