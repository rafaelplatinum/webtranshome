<script setup>
import { computed } from 'vue'
import { ArrowRight } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { formatAngka } from '@/utils/format'
import { LINGKARAN_HP, MAKS_LINGKARAN_BRAND } from '@/utils/brandBeranda'
import LingkaranBrand from './LingkaranBrand.vue'

// Susunan lingkaran "Brand pilihan" (KEPUTUSAN #19 & #21): judul + tombol di tengah, semua lingkaran UTUH di dalam frame.
// Frame ≥ 58rem: tiap sisi pola sarang lebah 3–2–3 = 8 lingkaran (total 16). Lebih sempit (HP/tablet): 8 teratas,
// 4 di atas + 4 di bawah judul, masing-masing satu baris (di HP kecil diameter ikut mengecil). Ukuran mengikuti lebar FRAME (container query), bukan layar, jadi komponen ini juga
// dipakai untuk pratinjau di admin › Brand di Beranda. Lingkaran tanpa brand diisi hiasan dua warna.
const props = defineProps({
  // Brand berurutan (prioritas tertinggi dulu), maks 16. Item: { id, name, slug, logo_url }.
  brands: { type: Array, required: true },
  // Jumlah semua brand aktif untuk tombol "Telusuri N brand".
  total: { type: Number, required: true },
  idJudul: { type: String, default: 'judul-brand-pilihan' },
})

const WARNA = [
  'var(--color-primary-soft)',
  'var(--color-info-soft)',
  'var(--color-accent)',
  'var(--color-success-soft)',
  'var(--color-danger-soft)',
  'var(--color-on-dark)',
]
const HIASAN = [
  ['var(--color-primary)', 'var(--color-primary-hover)', 115],
  ['var(--color-accent)', 'var(--color-warning)', 65],
  ['var(--color-secondary)', 'var(--color-secondary-hover)', 150],
  ['var(--color-primary-on-dark)', 'var(--color-primary)', 30],
]
// Posisi sisi kiri [kolom awal dari 6 kolom setengah-langkah, baris], urut dari yang paling dekat judul.
// Baris 1 & 3: 3 lingkaran (kolom 1, 3, 5); baris 2: 2 lingkaran di celahnya (kolom 2, 4). Sisi kanan = cermin.
const POSISI_KIRI = [[4, 2], [5, 1], [5, 3], [2, 2], [3, 1], [3, 3], [1, 1], [1, 3]]

const slot = computed(() =>
  Array.from({ length: MAKS_LINGKARAN_BRAND }, (_, i) =>
    props.brands[i] ? { kunci: `b${props.brands[i].id}`, brand: props.brands[i] } : { kunci: `h${i}`, hiasan: HIASAN[i % HIASAN.length] },
  ),
)
// Warna menurut posisi: tetangga selalu berbeda.
const berwarna = (item, n) => (item.brand ? { ...item, warna: WARNA[n % WARNA.length] } : item)

// Peringkat genap ke kiri, ganjil ke kanan; tiap sisi diisi dari posisi terdekat judul.
function sisi(awal, kiri) {
  return POSISI_KIRI.map(([kolom, baris], k) => {
    const kol = kiri ? kolom : 6 - kolom
    return { ...berwarna(slot.value[awal + k * 2], kol + baris * 2 + (kiri ? 0 : 3)), gaya: { gridColumn: `${kol} / span 2`, gridRow: baris } }
  })
}
const kiri = computed(() => sisi(0, true))
const kanan = computed(() => sisi(1, false))
const atas = computed(() => slot.value.slice(0, LINGKARAN_HP / 2).map((item, i) => berwarna(item, i * 2)))
const bawah = computed(() => slot.value.slice(LINGKARAN_HP / 2, LINGKARAN_HP).map((item, i) => berwarna(item, i * 2 + 3)))
</script>

<template>
  <div class="@container overflow-hidden rounded-2xl border border-border bg-surface" data-brand-pilihan>
    <div
      class="grid grid-cols-[minmax(0,1fr)] items-center gap-6 px-4 py-8 [--d:min(4.5rem,calc((100cqw_-_4.5rem)/4))] @min-[36rem]:[--d:5.5rem] @min-[58rem]:grid-cols-[minmax(0,1fr)_minmax(16rem,20rem)_minmax(0,1fr)] @min-[58rem]:px-6 @min-[58rem]:py-7 @min-[58rem]:[--d:5rem] @min-[72rem]:[--d:7rem]"
    >
      <ul
        class="hidden justify-self-end gap-y-3 @min-[58rem]:grid @min-[58rem]:grid-cols-[repeat(6,calc((var(--d)+1rem)/2))]"
        data-lingkaran-kiri
      >
        <li v-for="item in kiri" :key="item.kunci" class="flex justify-center" :style="item.gaya"><LingkaranBrand :item="item" /></li>
      </ul>
      <ul class="flex flex-wrap justify-center gap-3 @min-[58rem]:hidden" data-lingkaran-atas>
        <li v-for="item in atas" :key="item.kunci"><LingkaranBrand :item="item" /></li>
      </ul>

      <div class="flex flex-col items-center gap-6 text-center">
        <h2 :id="idJudul" class="text-3xl leading-tight font-extrabold text-balance @min-[36rem]:text-4xl @min-[72rem]:text-5xl">
          Siap untuk menjelajah lebih banyak Brand?
        </h2>
        <Button to="/brand" size="lg">
          Telusuri {{ formatAngka(total) }} brand
          <ArrowRight class="size-5" aria-hidden="true" />
        </Button>
      </div>

      <ul class="flex flex-wrap justify-center gap-3 @min-[58rem]:hidden" data-lingkaran-bawah>
        <li v-for="item in bawah" :key="item.kunci"><LingkaranBrand :item="item" /></li>
      </ul>
      <ul
        class="hidden justify-self-start gap-y-3 @min-[58rem]:grid @min-[58rem]:grid-cols-[repeat(6,calc((var(--d)+1rem)/2))]"
        data-lingkaran-kanan
      >
        <li v-for="item in kanan" :key="item.kunci" class="flex justify-center" :style="item.gaya"><LingkaranBrand :item="item" /></li>
      </ul>
    </div>
  </div>
</template>
