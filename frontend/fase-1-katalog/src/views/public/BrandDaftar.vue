<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, SearchX } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useSeo } from '@/composables/useSeo'
import { listBrands } from '@/services/catalogService'
import { pesanError } from '@/services/errors'
import { formatAngka } from '@/utils/format'

// /brand: semua brand aktif urut A–Z, dengan cari dan lompat per huruf (untuk ratusan brand).
useSeo({ title: 'Semua brand', description: 'Daftar brand bahan bangunan dan elektronik yang tersedia di Transhome Sleman.' })

const HURUF = [...'ABCDEFGHIJKLMNOPQRSTUVWXYZ', '#']
const status = ref('loading') // loading | ready | error
const brand = ref([])
const pesan = ref('')
const cari = ref('')

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

// Huruf pertama tanpa aksen; angka/simbol masuk kelompok "#".
const hurufAwal = (nama) => {
  const h = String(nama).normalize('NFD').charAt(0).toUpperCase()
  return h >= 'A' && h <= 'Z' ? h : '#'
}
const cocok = computed(() => {
  const q = cari.value.trim().toLowerCase()
  return q ? brand.value.filter((b) => b.name.toLowerCase().includes(q)) : brand.value
})
const kelompok = computed(() =>
  HURUF.map((h) => ({ huruf: h, brand: cocok.value.filter((b) => hurufAwal(b.name) === h) })).filter((k) => k.brand.length),
)
const adaHuruf = computed(() => new Set(kelompok.value.map((k) => k.huruf)))
const idHuruf = (h) => (h === '#' ? 'huruf-lain' : `huruf-${h}`)

// Bar cari + huruf menempel di bawah header dan tingginya berubah (huruf bisa turun baris):
// jarak berhenti tiap kelompok huruf mengikuti tinggi bar supaya judulnya tidak tertutup.
const bar = ref(null)
const tinggiBar = ref(0)
let pengamat = null
watch(bar, (el) => {
  pengamat?.disconnect()
  if (!el) return
  pengamat = new ResizeObserver(() => (tinggiBar.value = el.offsetHeight))
  pengamat.observe(el)
})
onBeforeUnmount(() => pengamat?.disconnect())
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Brand' }]" />
    <h1 class="mt-2 text-2xl font-extrabold sm:text-3xl">Semua brand</h1>
    <p class="mt-2 max-w-2xl text-muted">
      <template v-if="status === 'ready'">{{ formatAngka(brand.length) }} brand tersedia di toko.</template>
      Pilih brand untuk melihat semua produknya.
    </p>

    <div v-if="status === 'loading'" class="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6" aria-busy="true">
      <span class="sr-only">Memuat brand…</span>
      <Skeleton v-for="n in 12" :key="n" class="h-24 rounded-xl" />
    </div>

    <ErrorState v-else-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />

    <EmptyState v-else-if="!brand.length" class="mt-6" title="Belum ada brand" description="Daftar brand sedang disiapkan. Sementara itu, kamu bisa menjelajah katalog." />

    <template v-else>
      <!-- Cari + lompat huruf menempel di bawah header saat digulir. -->
      <div ref="bar" class="sticky top-(--tinggi-header,0px) z-20 -mx-4 mt-6 flex flex-col gap-3 border-b border-border bg-background px-4 py-3 sm:-mx-6 sm:px-6 lg:flex-row lg:items-center">
        <div class="lg:w-80">
          <Input v-model="cari" label="Cari brand" hide-label type="search" placeholder="Cari brand" autocomplete="off" enterkeyhint="search">
            <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
          </Input>
        </div>
        <nav aria-label="Lompat ke huruf" class="-mx-1 flex gap-0.5 overflow-x-auto pb-1 lg:flex-wrap lg:overflow-visible lg:pb-0">
          <template v-for="h in HURUF" :key="h">
            <a
              v-if="adaHuruf.has(h)"
              :href="`#${idHuruf(h)}`"
              class="inline-flex size-9 shrink-0 items-center justify-center rounded-md text-sm font-bold transition-colors duration-150 hover:bg-subtle hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              :aria-label="h === '#' ? 'Brand berawalan angka atau simbol' : `Brand berawalan ${h}`"
            >
              {{ h }}
            </a>
            <span v-else class="inline-flex size-9 shrink-0 items-center justify-center text-sm text-faint" aria-hidden="true">{{ h }}</span>
          </template>
        </nav>
      </div>

      <EmptyState v-if="!kelompok.length" class="mt-6" :title="`Brand “${cari.trim()}” tidak ditemukan`" description="Periksa ejaan, atau tanyakan ke tim kami lewat WhatsApp.">
        <template #icon><SearchX class="size-10" aria-hidden="true" /></template>
        <template #actions><Button variant="outline" @click="cari = ''">Tampilkan semua brand</Button></template>
      </EmptyState>

      <p class="sr-only" aria-live="polite">{{ cari.trim() ? `${cocok.length} brand cocok` : '' }}</p>

      <section v-for="k in kelompok" :id="idHuruf(k.huruf)" :key="k.huruf" class="mt-8" :style="{ scrollMarginTop: `${tinggiBar + 16}px` }" :aria-labelledby="`judul-${idHuruf(k.huruf)}`">
        <h2 :id="`judul-${idHuruf(k.huruf)}`" class="border-b border-border pb-2 text-xl font-extrabold">{{ k.huruf }}</h2>
        <ul class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
          <li v-for="(b, i) in k.brand" :key="b.id" v-reveal="i % 6">
            <RouterLink
              :to="`/brand/${b.slug}`"
              class="flex h-24 flex-col items-center justify-center gap-1 rounded-xl border border-border bg-surface px-3 text-center transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <img v-if="b.logo_url" :src="b.logo_url" :alt="b.name" width="120" height="40" loading="lazy" class="max-h-10 w-auto object-contain" />
              <span v-else class="font-bold">{{ b.name }}</span>
              <span class="text-xs text-muted">{{ formatAngka(b.product_count ?? 0) }} produk</span>
            </RouterLink>
          </li>
        </ul>
      </section>
    </template>
  </div>
</template>
