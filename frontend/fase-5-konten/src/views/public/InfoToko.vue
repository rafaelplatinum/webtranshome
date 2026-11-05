<script setup>
import { computed, ref } from 'vue'
import { Clock, Coins, ExternalLink, MapPin, MessageCircle, Phone, Store } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import Button from '@/components/ui/Button.vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useSeo } from '@/composables/useSeo'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'
import { CONTACT } from '@/config/contact'
import PageContactCta from '@/components/public/PageContactCta.vue'
import { useMuat } from '@/composables/useMuat'
import { listBanners } from '@/services/kontenService'

const LAYANAN = [
  { ikon: Store, judul: 'Toko fisik di Sleman', isi: 'Lihat langsung barangnya dan ambil kebutuhan rumah di toko.' },
  { ikon: MessageCircle, judul: 'Konsultasi kebutuhan', isi: 'Ceritakan kebutuhan produk atau proyek kepada tim Transhome.' },
  { ikon: Coins, judul: 'Poin Trans Family', isi: 'Setiap belanja di toko dapat tercatat sebagai poin member.' },
]

const settings = useSettingsStore()
const ui = useUiStore()
settings.load()

const alamat = computed(() => CONTACT.address)
const jamBuka = computed(() => CONTACT.openingHours)
const mapsUrl = computed(() => settings.data?.maps_url)
const memuat = computed(() => settings.status === 'idle' || settings.status === 'loading')
const teksAlamat = ref(null)
const bannerTentang = useMuat(async () => (await listBanners('ABOUT_HERO'))[0] ?? null)

function alamatDisalin() { ui.tampilkanToast({ pesan: 'Alamat toko disalin.', jenis: 'success' }) }
function alamatGagalDisalin() {
  if (teksAlamat.value) window.getSelection()?.selectAllChildren(teksAlamat.value)
  ui.tampilkanToast({ pesan: 'Alamat sudah dipilih. Salin secara manual.', jenis: 'info' })
}

useSeo({ title: 'Info toko', description: () => `Kunjungi ${CONTACT.storeName} di ${alamat.value}. Buka ${jamBuka.value}.` })
</script>

<template>
  <div class="flex flex-col gap-14 sm:gap-20">
    <section class="relative isolate flex min-h-[32rem] overflow-hidden bg-secondary text-on-dark sm:min-h-[36rem]">
      <img v-if="bannerTentang.data.value?.image_url" :src="bannerTentang.data.value.image_url" alt="" width="1920" height="800" fetchpriority="high" class="absolute inset-0 z-0 size-full object-cover" aria-hidden="true" />
      <div class="absolute inset-0 z-10 bg-gradient-to-r from-secondary via-secondary/90 to-secondary/35" aria-hidden="true"></div>
      <div class="th-container relative z-20 flex w-full flex-col items-start justify-center gap-5 py-12 sm:py-16 lg:py-20">
        <Breadcrumb variant="light" :items="[{ label: 'Beranda', to: '/' }, { label: 'Info toko' }]" />
        <span class="th-eyebrow text-orange-200">Trans Home Store · Yogyakarta</span>
        <h1 class="max-w-3xl text-4xl leading-[1.08] font-extrabold tracking-tight text-secondary-foreground sm:text-5xl lg:text-[56px]">Ruang yang nyaman dimulai dari kebutuhan yang tepat.</h1>
        <p class="max-w-2xl text-base leading-relaxed text-on-dark sm:text-lg">{{ CONTACT.description }}</p>
        <div class="flex flex-wrap gap-3 pt-1">
          <Button v-if="mapsUrl" :href="mapsUrl" variant="outline" class="!border-th-yellow !bg-th-yellow !text-th-navy hover:!border-amber-300 hover:!bg-amber-300"><template #icon><MapPin class="size-4" aria-hidden="true" /></template>Buka lokasi</Button>
          <WhatsAppButton :context="{ topic: 'Lainnya', message: 'Saya ingin bertanya tentang toko Transhome.' }" label="Tanya ke tim kami" />
        </div>
      </div>
    </section>

    <section class="th-container flex flex-col gap-8" aria-labelledby="judul-kunjungi">
      <div class="grid items-center gap-8 lg:grid-cols-[0.8fr_1.2fr] lg:gap-16">
        <div class="flex flex-col gap-4"><span class="th-eyebrow">Kunjungi kami</span><h2 id="judul-kunjungi" class="text-3xl leading-tight font-extrabold tracking-tight sm:text-4xl">Kami siap menyambutmu di toko.</h2><p class="max-w-lg text-base leading-relaxed text-th-muted">Datang untuk melihat pilihan produk secara langsung, atau hubungi tim kami untuk membicarakan kebutuhan rumah dan proyekmu.</p></div>
        <div v-if="memuat" class="h-64 animate-pulse rounded-[20px] bg-th-blue" aria-hidden="true"></div>
        <ErrorState v-else-if="settings.status === 'error' && !settings.data" title="Info toko belum bisa dimuat" @retry="settings.load({ force: true })" />
        <a v-else-if="mapsUrl" :href="mapsUrl" target="_blank" rel="noopener noreferrer" class="group relative flex min-h-64 items-end overflow-hidden rounded-[20px] border border-th-border bg-th-blue p-6 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:p-8">
          <div class="absolute inset-0 opacity-70" aria-hidden="true"><span class="absolute inset-0 bg-[linear-gradient(34deg,transparent_47%,rgba(255,255,255,.75)_48%,rgba(255,255,255,.75)_51%,transparent_52%),linear-gradient(145deg,transparent_44%,rgba(255,255,255,.55)_45%,rgba(255,255,255,.55)_48%,transparent_49%)]"></span><span class="absolute top-1/4 left-1/3 size-28 rounded-full border-2 border-white/70"></span></div>
          <div class="relative flex w-full items-end justify-between gap-4"><div><span class="th-eyebrow text-th-navy">Lokasi toko</span><p class="mt-2 max-w-md text-lg leading-snug font-extrabold text-th-navy">{{ alamat }}</p></div><span class="grid size-12 shrink-0 place-items-center rounded-full bg-th-orange text-white transition-transform group-hover:-translate-y-1"><MapPin class="size-6" aria-hidden="true" /></span></div>
        </a>
        <div v-else class="relative flex min-h-64 items-end overflow-hidden rounded-[20px] border border-th-border bg-th-blue p-6 sm:p-8"><div class="absolute inset-0 opacity-70" aria-hidden="true"><span class="absolute inset-0 bg-[linear-gradient(34deg,transparent_47%,rgba(255,255,255,.75)_48%,rgba(255,255,255,.75)_51%,transparent_52%),linear-gradient(145deg,transparent_44%,rgba(255,255,255,.55)_45%,rgba(255,255,255,.55)_48%,transparent_49%)]"></span></div><div class="relative"><span class="th-eyebrow text-th-navy">Lokasi toko</span><p class="mt-2 max-w-md text-lg leading-snug font-extrabold text-th-navy">{{ alamat }}</p></div></div>
      </div>

      <div v-if="memuat" class="grid gap-4 md:grid-cols-3" aria-busy="true"><span class="sr-only">Memuat info toko…</span><Skeleton v-for="n in 3" :key="n" class="h-52 rounded-[20px]" /></div>
      <div v-else class="grid gap-4 md:grid-cols-3 lg:gap-6">
        <section aria-labelledby="judul-alamat" class="th-card flex min-h-56 flex-col gap-4 p-5 sm:p-6">
          <MapPin class="size-6 text-th-orange" aria-hidden="true" /><div><h3 id="judul-alamat" class="text-lg font-extrabold">Alamat toko</h3><address ref="teksAlamat" class="mt-2 text-base leading-relaxed not-italic" data-alamat>{{ alamat }}</address></div>
          <div class="mt-auto flex flex-wrap gap-2"><Button v-if="mapsUrl" variant="outline" size="sm" :href="mapsUrl" data-tombol-peta><template #icon><ExternalLink class="size-4" aria-hidden="true" /></template>Buka peta</Button><CopyButton :text="alamat" label="Salin alamat" copied-label="Alamat disalin" @copied="alamatDisalin" @failed="alamatGagalDisalin" /></div>
        </section>
        <section class="th-card flex min-h-56 flex-col gap-4 bg-th-blue p-5 sm:p-6"><Clock class="size-6 text-th-navy" aria-hidden="true" /><div><h3 class="text-lg font-extrabold">Jam operasional</h3><p class="mt-2 text-base leading-relaxed">{{ jamBuka }}</p></div><a :href="`tel:${CONTACT.telephone.replaceAll(' ', '')}`" class="mt-auto inline-flex min-h-11 items-center gap-2 font-bold text-th-navy hover:underline"><Phone class="size-4" aria-hidden="true" />{{ CONTACT.telephone }}</a></section>
        <section class="th-card flex min-h-56 flex-col items-start gap-4 bg-th-peach p-5 sm:p-6"><MessageCircle class="size-6 text-th-orange" aria-hidden="true" /><div><h3 class="text-lg font-extrabold">Hubungi tim kami</h3><p class="mt-2 text-base leading-relaxed">WhatsApp {{ CONTACT.whatsappDisplay }}</p></div><WhatsAppButton class="mt-auto" variant="whatsapp-outline" :context="{ topic: 'Lainnya', message: 'Saya ingin menghubungi tim Transhome.' }" label="WhatsApp CS" /></section>
      </div>
    </section>

    <section id="fasilitas" class="th-container scroll-mt-28 flex flex-col gap-6" aria-labelledby="judul-layanan">
      <div class="max-w-2xl"><span class="th-eyebrow">Layanan Transhome</span><h2 id="judul-layanan" class="mt-2 text-3xl font-extrabold tracking-tight sm:text-4xl">Lebih dari sekadar tempat belanja.</h2></div>
      <ul class="grid gap-4 md:grid-cols-3 lg:gap-6">
        <li v-for="(layanan, i) in LAYANAN" :key="layanan.judul" class="th-card flex min-h-56 flex-col gap-5 p-6 sm:p-8" :class="['bg-th-yellow', 'bg-th-blue', 'bg-th-peach'][i]">
          <component :is="layanan.ikon" class="size-6 text-primary" aria-hidden="true" />
          <div class="mt-auto"><h3 class="font-bold">{{ layanan.judul }}</h3><p class="mt-1 text-sm leading-relaxed text-muted">{{ layanan.isi }}</p></div>
        </li>
      </ul>
    </section>

    <PageContactCta heading="Ada yang ingin ditanyakan?" description="Ceritakan kebutuhan produk atau proyekmu kepada tim kami." />
  </div>
</template>
