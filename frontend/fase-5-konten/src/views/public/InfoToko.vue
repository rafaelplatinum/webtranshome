<script setup>
import { computed, ref } from 'vue'
import { Clock, Coins, ExternalLink, MapPin, MessageCircle, Store } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import Button from '@/components/ui/Button.vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'

// Teks layanan dari sketsa "Main". Alamat, jam buka, dan link peta dari pengaturan toko (FR-P).
const LAYANAN = [
  { ikon: Store, judul: 'Toko fisik di Sleman', isi: 'Lihat langsung barangnya, ambil di toko.' },
  { ikon: MessageCircle, judul: 'Konsultasi proyek via WhatsApp', isi: 'Kirim kebutuhan proyekmu, tim kami bantu hitung.' },
  { ikon: Coins, judul: 'Poin Trans Family', isi: 'Setiap belanja di toko tercatat jadi poin.' },
]

const settings = useSettingsStore()
const ui = useUiStore()
const { generalLink } = useWhatsApp()
settings.load()

const toko = computed(() => settings.data)
const memuat = computed(() => settings.status === 'idle' || settings.status === 'loading')
const teksAlamat = ref(null)

function alamatDisalin() {
  ui.tampilkanToast({ pesan: 'Alamat toko disalin.', jenis: 'success' })
}

function alamatGagalDisalin() {
  if (teksAlamat.value) window.getSelection()?.selectAllChildren(teksAlamat.value)
  ui.tampilkanToast({ pesan: 'Alamat sudah dipilih. Salin secara manual.', jenis: 'info' })
}

useSeo({
  title: 'Info toko',
  description: () =>
    toko.value
      ? `Toko Transhome di ${toko.value.address}. Buka ${toko.value.opening_hours}.`
      : 'Alamat, jam buka, dan kontak WhatsApp toko Transhome di Sleman.',
})
</script>

<template>
  <div class="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-4 sm:px-6 lg:py-6">
    <div class="flex flex-col gap-2">
      <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Info toko' }]" />
      <h1 class="text-2xl font-extrabold sm:text-3xl">Info toko</h1>
      <p class="max-w-2xl text-muted">Datang langsung untuk lihat barangnya, atau tanya dulu lewat WhatsApp.</p>
    </div>

    <div v-if="memuat" class="grid gap-4 lg:grid-cols-2" aria-busy="true">
      <span class="sr-only">Memuat info toko…</span>
      <Skeleton class="h-64 rounded-xl" />
      <Skeleton class="h-64 rounded-xl" />
    </div>

    <ErrorState v-else-if="settings.status === 'error' || !toko" title="Info toko belum bisa dimuat" @retry="settings.load()" />

    <div v-else class="grid items-start gap-4 lg:grid-cols-[minmax(0,1.2fr)_minmax(0,1fr)]">
      <section aria-labelledby="judul-alamat" class="flex flex-col gap-5 rounded-xl border border-border bg-surface p-5 sm:p-6">
        <div class="flex items-start gap-3">
          <MapPin class="mt-1 size-5 shrink-0 text-primary" aria-hidden="true" />
          <div class="flex min-w-0 flex-col gap-1">
            <h2 id="judul-alamat" class="text-lg font-bold">Alamat</h2>
            <address ref="teksAlamat" class="text-lg not-italic" data-alamat>{{ toko.address }}</address>
          </div>
        </div>
        <div class="flex flex-wrap gap-2 pl-8">
          <Button v-if="toko.maps_url" variant="secondary" :href="toko.maps_url" data-tombol-peta>
            <template #icon><ExternalLink class="size-4" aria-hidden="true" /></template>
            Buka di Google Maps
            <span class="sr-only">(tab baru)</span>
          </Button>
          <CopyButton :text="toko.address" label="Salin alamat" copied-label="Alamat disalin" @copied="alamatDisalin" @failed="alamatGagalDisalin" />
        </div>
        <div class="flex items-start gap-3 border-t border-border pt-5">
          <Clock class="mt-1 size-5 shrink-0 text-primary" aria-hidden="true" />
          <div class="flex flex-col gap-1">
            <h2 class="text-lg font-bold">Jam buka</h2>
            <p class="text-lg" data-jam-buka>{{ toko.opening_hours }}</p>
          </div>
        </div>
      </section>

      <section aria-labelledby="judul-layanan" class="flex flex-col gap-5 rounded-xl bg-secondary p-5 text-on-dark sm:p-6">
        <h2 id="judul-layanan" class="text-lg font-bold text-secondary-foreground">Layanan</h2>
        <ul class="flex flex-col gap-4">
          <li v-for="l in LAYANAN" :key="l.judul" class="flex items-start gap-3">
            <component :is="l.ikon" class="mt-0.5 size-5 shrink-0 text-primary-on-dark" aria-hidden="true" />
            <div class="flex flex-col gap-0.5">
              <strong class="text-secondary-foreground">{{ l.judul }}</strong>
              <span>{{ l.isi }}</span>
            </div>
          </li>
        </ul>
        <WhatsAppButton class="self-start" :href="generalLink()" label="Chat WhatsApp CS" />
      </section>
    </div>
  </div>
</template>
