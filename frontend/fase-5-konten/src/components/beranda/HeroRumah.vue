<script setup>
import { computed, ref, useId } from 'vue'
import { ArrowRight } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useMuat } from '@/composables/useMuat'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { getHouseParts } from '@/services/kontenService'
import { BAGIAN_RUMAH } from './bagianRumah'
import RumahTigaDimensi from './RumahTigaDimensi.vue'

// Hero Beranda (sketsa "Landing3D"): pilih bagian rumah → bagian terangkat + kartu kategori.
const SUDUT_AWAL = -40

const { generalLink } = useWhatsApp()
const idDaftar = useId()
const terpisah = ref(true)
const sudut = ref(SUDUT_AWAL)
const pilihan = ref(null)

// Kategori tujuan per bagian dari API. Bila gagal dimuat, "Lihat produk" mengarah ke katalog lengkap.
const peta = useMuat(getHouseParts)
const dipilih = computed(() => (pilihan.value === null ? null : BAGIAN_RUMAH[pilihan.value]))
const tautanDipilih = computed(() => {
  const kategori = dipilih.value && peta.data.value?.find((b) => b.key === dipilih.value.key)?.category
  return kategori ? `/katalog/${kategori.slug}` : '/katalog'
})

// Atap di atas, pondasi di bawah, seperti urutan rumah.
const daftar = [...BAGIAN_RUMAH].reverse()

function pilih(indeks) {
  pilihan.value = pilihan.value === indeks ? null : indeks
  terpisah.value = true
}

function rakitAtauPisah() {
  terpisah.value = !terpisah.value
  pilihan.value = null
}

function aturUlang() {
  terpisah.value = true
  sudut.value = SUDUT_AWAL
  pilihan.value = null
}

const kelasKontrol =
  'inline-flex h-11 items-center rounded-md border border-on-dark-muted px-4.5 text-sm font-semibold text-secondary-foreground transition-colors duration-150 hover:border-secondary-foreground hover:bg-secondary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <section aria-labelledby="judul-beranda" class="overflow-hidden bg-secondary text-secondary-foreground">
    <div class="mx-auto grid max-w-7xl items-center gap-10 px-4 py-10 sm:px-6 lg:grid-cols-2 lg:py-16">
      <div class="flex min-w-0 flex-col gap-5">
        <span class="text-sm font-semibold tracking-wide text-accent">Supermarket bahan bangunan dan elektronik · Sleman</span>
        <h1 id="judul-beranda" class="text-4xl leading-tight font-extrabold sm:text-5xl lg:text-6xl">Dari pondasi sampai atap, cukup di satu toko.</h1>
        <p class="max-w-lg text-lg text-on-dark">Cari bahan untuk tiap tahap rumahmu, cek harga dan stok, lalu tanya langsung ke tim kami lewat WhatsApp.</p>
        <div class="flex flex-wrap items-center gap-3">
          <Button to="/katalog" size="lg">
            Jelajahi katalog
            <ArrowRight class="size-5" aria-hidden="true" />
          </Button>
          <WhatsAppButton :href="generalLink()" size="lg" />
        </div>

        <div class="mt-2 flex flex-col gap-2">
          <span :id="idDaftar" class="text-xs font-bold tracking-wider text-on-dark-muted uppercase">Pilih bagian rumah</span>
          <ul :aria-labelledby="idDaftar" class="flex flex-col gap-1.5" data-daftar-bagian>
            <li v-for="b in daftar" :key="b.key">
              <button
                type="button"
                class="flex min-h-12 w-full items-center gap-3.5 rounded-lg border px-3.5 py-1.5 text-left transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
                :class="pilihan === b.nomor - 1 ? 'border-accent bg-secondary-hover' : 'border-on-dark-muted/40 hover:bg-secondary-hover'"
                :aria-pressed="pilihan === b.nomor - 1 ? 'true' : 'false'"
                @click="pilih(b.nomor - 1)"
              >
                <span
                  class="flex size-7.5 shrink-0 items-center justify-center rounded-md text-sm font-extrabold"
                  :style="{ background: `var(--color-${b.titik})`, color: `var(--color-${b.tinta})` }"
                  aria-hidden="true"
                >
                  {{ b.nomor }}
                </span>
                <span class="flex min-w-0 flex-col">
                  <strong class="font-semibold">{{ b.nama }}</strong>
                  <span class="text-sm text-on-dark-muted">{{ b.isi }}</span>
                </span>
              </button>
            </li>
          </ul>
        </div>
      </div>

      <div class="flex min-w-0 flex-col items-center gap-4">
        <RumahTigaDimensi :terpisah="terpisah" :sudut="sudut" :pilihan="pilihan" />
        <div class="flex flex-wrap items-center justify-center gap-x-4 gap-y-3">
          <button type="button" :class="kelasKontrol" data-tombol-rakit @click="rakitAtauPisah">{{ terpisah ? 'Rakit jadi rumah' : 'Pisahkan bagian' }}</button>
          <label class="flex items-center gap-2.5 text-sm text-on-dark">
            Putar rumah
            <input v-model.number="sudut" type="range" min="-78" max="-12" step="1" class="w-44 cursor-pointer accent-primary-on-dark" />
          </label>
          <button type="button" :class="kelasKontrol" @click="aturUlang">Atur ulang</button>
        </div>
        <div class="w-full max-w-md" aria-live="polite">
          <div v-if="dipilih" class="flex flex-col gap-2 rounded-xl border border-on-dark-muted bg-secondary-hover p-5" data-kartu-bagian>
            <span class="text-xs font-bold tracking-wider text-accent uppercase">Bagian {{ dipilih.nomor }}</span>
            <strong class="text-xl">{{ dipilih.nama }}</strong>
            <span class="text-sm text-on-dark">{{ dipilih.deskripsi }}</span>
            <Button :to="tautanDipilih" size="sm" class="mt-1 self-start">Lihat produk {{ dipilih.singkat }}</Button>
          </div>
          <p v-else class="text-center text-sm text-on-dark-muted">Pilih bagian rumah untuk melihat kategorinya.</p>
        </div>
      </div>
    </div>
  </section>
</template>
