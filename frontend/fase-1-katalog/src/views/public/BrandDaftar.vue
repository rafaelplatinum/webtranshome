<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { ArrowRight, Search, SearchX } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useFilter } from '@/composables/useFilter'
import { remahSeo, urlSitus, useSeo } from '@/composables/useSeo'
import { listBrands } from '@/services/catalogService'
import { pesanError } from '@/services/errors'
import { formatAngka } from '@/utils/format'
import { perilakuGulir } from '@/utils/gerak'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'

// /brand mengikuti wireframe 02 (docs/wireframe-low-fidelity-transhome.xlsx, sheet "Brand"): judul + jumlah,
// cari dengan tombol Cari, filter kategori, A–Z, kartu 6/3/2 kolom, 24 per halaman. Filter tersimpan di URL.

const PER_HALAMAN = 24
const HURUF = [...'ABCDEFGHIJKLMNOPQRSTUVWXYZ', '#']
const { nilai, terapkan, reset } = useFilter()
const router = useRouter()

const status = ref('loading') // loading | ready | error
const brand = ref([])
const pesan = ref('')
const atasDaftar = ref(null)

async function muat() {
  status.value = 'loading'
  try {
    brand.value = [...(await listBrands())].sort((a, b) => a.name.localeCompare(b.name, 'id'))
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}
onMounted(muat)

// Cari diterapkan saat tombol Cari / Enter ditekan (wireframe), isian ikut nilai di URL.
const teks = ref(nilai.value.q ?? '')
watch(() => nilai.value.q, (q) => (teks.value = q ?? ''))
const cari = () => terapkan({ q: teks.value.trim() || null })

const hurufAwal = (nama) => {
  const h = String(nama).normalize('NFD').charAt(0).toUpperCase()
  return h >= 'A' && h <= 'Z' ? h : '#'
}
const kata = computed(() => (nilai.value.q ?? '').trim().toLowerCase())
const cocokCari = computed(() => brand.value.filter((b) => !kata.value || b.name.toLowerCase().includes(kata.value)))
// Kategori diambil dari brand yang ada (kategori utama produk tiap brand).
const kategori = computed(() => {
  const peta = new Map()
  for (const b of cocokCari.value) for (const k of b.categories ?? []) peta.set(k.slug, k.name)
  return [...peta].map(([slug, name]) => ({ slug, name }))
})
const perKategori = computed(() => cocokCari.value.filter((b) => !nilai.value.kategori || (b.categories ?? []).some((k) => k.slug === nilai.value.kategori)))
const hurufAda = computed(() => new Set(perKategori.value.map((b) => hurufAwal(b.name))))
const hasil = computed(() => perKategori.value.filter((b) => !nilai.value.huruf || hurufAwal(b.name) === nilai.value.huruf))

const totalHalaman = computed(() => Math.max(1, Math.ceil(hasil.value.length / PER_HALAMAN)))
const halaman = computed(() => Math.min(Math.max(Number(nilai.value.page) || 1, 1), totalHalaman.value))
const tampil = computed(() => hasil.value.slice((halaman.value - 1) * PER_HALAMAN, halaman.value * PER_HALAMAN))
const rentang = computed(() => {
  if (!hasil.value.length) return ''
  const awal = (halaman.value - 1) * PER_HALAMAN + 1
  return `Menampilkan ${awal}–${awal + tampil.value.length - 1} dari ${formatAngka(hasil.value.length)} brand`
})
const adaFilter = computed(() => Boolean(nilai.value.q || nilai.value.kategori || nilai.value.huruf))

// SEO: halaman bersih (/brand, /brand?page=N) boleh diindeks; hasil cari/filter = noindex + canonical ke /brand
// supaya Google tidak mengindeks ratusan variasi halaman yang isinya mirip. Tautan paginasi berupa <a href>.
const namaKategori = computed(() => kategori.value.find((k) => k.slug === nilai.value.kategori)?.name ?? nilai.value.kategori)
const judulFilter = computed(() => {
  if (nilai.value.q) return `Cari brand “${nilai.value.q}”`
  const huruf = nilai.value.huruf === '#' ? 'angka' : nilai.value.huruf
  if (nilai.value.kategori) return `Brand ${namaKategori.value}${huruf ? ` berawalan ${huruf}` : ''}`
  return huruf ? `Brand berawalan ${huruf}` : null
})
const tautanHalaman = (n) => router.resolve({ path: '/brand', query: { ...nilai.value, page: n > 1 ? n : undefined } }).href
useSeo({
  title: () => {
    if (status.value !== 'ready') return 'Brand partner'
    if (judulFilter.value) return judulFilter.value
    const n = formatAngka(brand.value.length)
    return halaman.value > 1 ? `${n} brand partner – halaman ${halaman.value}` : `${n} brand bahan bangunan & elektronik`
  },
  description: () => {
    if (status.value !== 'ready') return 'Daftar brand bahan bangunan dan elektronik yang tersedia di Transhome Sleman.'
    const contoh = (adaFilter.value ? hasil.value : brand.value).slice(0, 5).map((b) => b.name).join(', ')
    const jumlah = formatAngka(adaFilter.value ? hasil.value.length : brand.value.length)
    return `${jumlah} brand partner Transhome Sleman${contoh ? `, antara lain ${contoh}` : ''}. Lihat produk tiap brand, cek harga dan stok, lalu tanya lewat WhatsApp.`
  },
  canonical: () => (adaFilter.value || halaman.value === 1 ? '/brand' : `/brand?page=${halaman.value}`),
  robots: () => (adaFilter.value ? 'noindex, follow' : null),
  jsonLd: () =>
    status.value === 'ready' && [
      remahSeo([{ label: 'Beranda', to: '/' }, { label: 'Brand', to: '/brand' }]),
      {
        '@type': 'CollectionPage',
        name: 'Brand partner Transhome',
        url: urlSitus('/brand'),
        mainEntity: {
          '@type': 'ItemList',
          numberOfItems: hasil.value.length,
          itemListElement: tampil.value.map((b, i) => ({
            '@type': 'ListItem',
            position: (halaman.value - 1) * PER_HALAMAN + i + 1,
            url: urlSitus(`/brand/${b.slug}`),
            name: b.name,
          })),
        },
      },
    ],
})

async function gantiHalaman(n) {
  await terapkan({ page: n > 1 ? n : null })
  atasDaftar.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
}

const kelasChip = (aktif) => [
  'inline-flex min-h-11 cursor-pointer items-center rounded-full border px-4 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-10',
  aktif ? 'border-th-navy bg-th-navy text-white' : 'border-th-border bg-white hover:border-th-navy',
]
</script>

<template>
  <div class="flex flex-col gap-10 sm:gap-14">
    <section class="bg-th-navy py-10 text-white sm:py-14 lg:py-16">
      <div class="th-container flex flex-col gap-5">
        <Breadcrumb variant="light" :items="[{ label: 'Beranda', to: '/' }, { label: 'Brand' }]" />
        <span class="th-eyebrow text-orange-200">Brand partner</span>
        <h1 class="max-w-4xl text-4xl leading-[1.08] font-extrabold tracking-tight sm:text-5xl lg:text-[56px]">{{ status === 'ready' ? `${formatAngka(brand.length)} brand bahan bangunan dan elektronik pilihan.` : 'Brand bahan bangunan dan elektronik pilihan.' }}</h1>
        <p class="max-w-2xl text-base leading-relaxed text-white/80 sm:text-lg">Pilih kartu brand untuk melihat produknya.</p>
        <form role="search" aria-label="Cari brand" class="mt-2 flex w-full max-w-3xl items-start gap-3" @submit.prevent="cari">
          <div class="min-w-0 flex-1"><Input v-model="teks" label="Cari brand" hide-label type="search" placeholder="Cari nama brand, misalnya Alfa Keramik" autocomplete="off" enterkeyhint="search" /></div>
          <Button type="submit" class="shrink-0"><template #icon><Search class="size-4.5" aria-hidden="true" /></template>Cari</Button>
        </form>
      </div>
    </section>

    <div class="th-container flex flex-col gap-5">
      <div v-if="status === 'loading'" class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4" aria-busy="true">
        <span class="sr-only">Memuat brand…</span><Skeleton v-for="n in 8" :key="n" class="h-56 rounded-[20px]" />
      </div>
      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
      <EmptyState v-else-if="!brand.length" title="Belum ada brand" description="Daftar brand sedang disiapkan. Sementara itu, kamu bisa menjelajah katalog." />

      <template v-else>
        <div role="group" aria-label="Filter kategori brand" class="flex flex-wrap items-center gap-2">
          <span class="mr-1 font-bold">Kategori</span>
          <button type="button" :class="kelasChip(!nilai.kategori)" :aria-pressed="!nilai.kategori ? 'true' : 'false'" @click="terapkan({ kategori: null })">Semua</button>
          <button v-for="k in kategori" :key="k.slug" type="button" :class="kelasChip(nilai.kategori === k.slug)" :aria-pressed="nilai.kategori === k.slug ? 'true' : 'false'" @click="terapkan({ kategori: nilai.kategori === k.slug ? null : k.slug })">{{ k.name }}</button>
        </div>
        <div class="flex flex-col gap-3 rounded-2xl border border-th-border bg-white p-4">
          <div role="group" aria-label="Huruf awal brand" class="flex flex-wrap gap-1">
            <button v-for="h in HURUF" :key="h" type="button" :disabled="!hurufAda.has(h) && nilai.huruf !== h" :aria-pressed="nilai.huruf === h ? 'true' : 'false'" :aria-label="h === '#' ? 'Brand berawalan angka atau simbol' : `Brand berawalan ${h}`" class="grid size-10 place-items-center rounded-lg text-sm font-bold transition-colors disabled:cursor-not-allowed disabled:text-faint" :class="nilai.huruf === h ? 'bg-th-navy text-white' : 'hover:bg-th-cream'" @click="terapkan({ huruf: nilai.huruf === h ? null : h })">{{ h }}</button>
          </div>
          <p class="text-sm font-semibold text-th-muted">{{ formatAngka(perKategori.length) }} brand</p>
        </div>

        <section ref="atasDaftar" aria-label="Daftar brand" class="scroll-mt-24">
          <div class="mb-4 flex flex-wrap items-center justify-between gap-3"><p class="text-sm font-semibold text-th-muted" aria-live="polite">{{ rentang }}</p><RouterLink to="/" class="font-bold text-th-orange hover:underline">Kembali ke beranda</RouterLink></div>
          <EmptyState v-if="!hasil.length" title="Brand tidak ditemukan" description="Coba kata lain, kategori lain, atau hapus filter."><template #icon><SearchX class="size-10" aria-hidden="true" /></template><template v-if="adaFilter" #actions><Button variant="outline" @click="reset()">Hapus filter</Button></template></EmptyState>
          <ul v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 lg:gap-6">
            <li v-for="(b, i) in tampil" :key="b.id" v-reveal="i % 4">
              <RouterLink :to="`/brand/${b.slug}`" class="th-card group flex h-full min-h-56 flex-col items-center gap-3 p-5 text-center focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:p-6">
                <span class="grid size-[88px] shrink-0 place-items-center overflow-hidden rounded-full bg-th-sage p-3 transition-transform duration-200 group-hover:scale-105">
                  <img v-if="b.logo_url" :src="b.logo_url" :alt="`Logo ${b.name}`" width="120" height="80" loading="lazy" class="max-h-full max-w-full object-contain" />
                  <span v-else class="text-2xl font-extrabold text-th-navy" aria-hidden="true">{{ b.name.trim().split(/\s+/).map((bagian) => bagian[0]).join('').slice(0, 2).toUpperCase() }}</span>
                </span>
                <strong class="text-base leading-snug">{{ b.name }}</strong>
                <span class="text-sm text-th-muted">{{ b.product_count ? `${formatAngka(b.product_count)} produk` : 'Belum ada produk' }}</span>
                <span class="mt-auto inline-flex items-center gap-1 font-bold text-th-orange group-hover:underline">Lihat brand <ArrowRight class="size-4" aria-hidden="true" /></span>
              </RouterLink>
            </li>
          </ul>
          <PaginationBar v-if="totalHalaman > 1" class="mt-6" :page="halaman" :total-pages="totalHalaman" :tautan="tautanHalaman" @update:page="gantiHalaman" />
        </section>
      </template>
    </div>

    <section class="bg-th-yellow py-10 sm:py-14">
      <div class="th-container flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
        <div><h2 class="text-2xl leading-tight font-extrabold tracking-tight sm:text-3xl">Cari brand yang belum ada?</h2><p class="mt-2 text-base text-th-navy/80">Tanyakan ke tim kami, kami bantu carikan.</p></div>
        <WhatsAppButton class="shrink-0" size="lg" variant="secondary" :context="{ topic: 'Tanya produk', message: 'Saya mencari brand yang belum ada di daftar.' }" label="Chat via WhatsApp" />
      </div>
    </section>
  </div>
</template>
