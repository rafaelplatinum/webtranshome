<script setup>
import { onMounted, ref } from 'vue'
import { LayoutDashboard } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import RoomCarousel from '@/components/katalog/RoomCarousel.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { listRooms } from '@/services/catalogService'
import { pesanError } from '@/services/errors'

const status = ref('loading') // loading | ready | error
const ruangan = ref([])
const pesan = ref('')
const { generalLink } = useWhatsApp()

async function muat() {
  status.value = 'loading'
  try {
    ruangan.value = await listRooms()
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

onMounted(muat)

useSeo({
  title: 'Belanja per ruangan',
  description: 'Pilih ruangan seperti dapur, kamar mandi, atau teras, lalu lihat bahan bangunan yang cocok di Transhome Sleman.',
})
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Ruangan' }]" />
    <h1 class="mt-2 text-2xl font-extrabold sm:text-3xl">Belanja per ruangan</h1>
    <p class="mt-2 max-w-2xl text-muted">Pilih ruangan yang sedang kamu bangun atau renovasi. Kami kumpulkan produk yang cocok untuk tiap ruangan.</p>

    <div v-if="status === 'loading'" class="mt-6 flex flex-col gap-3" aria-busy="true">
      <span class="sr-only">Memuat ruangan…</span>
      <Skeleton class="h-6 w-40" />
      <div class="grid auto-cols-[minmax(14rem,1fr)] grid-flow-col gap-4 overflow-hidden">
        <Skeleton v-for="n in 5" :key="n" class="aspect-4/3 rounded-xl" />
      </div>
    </div>

    <ErrorState v-else-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />

    <EmptyState v-else-if="ruangan.length === 0" class="mt-6" title="Belum ada ruangan" description="Daftar ruangan sedang disiapkan. Sementara itu, kamu bisa menjelajah katalog.">
      <template #icon><LayoutDashboard class="size-10" aria-hidden="true" /></template>
    </EmptyState>

    <RoomCarousel v-else class="mt-6" :ruangan="ruangan">
      <template #judul>
        <h2 class="text-lg font-bold">Pilih ruangan</h2>
      </template>
    </RoomCarousel>

    <div v-if="status === 'ready'" class="mt-10 flex flex-col items-start gap-3 rounded-xl border border-border bg-surface p-5 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-col gap-1">
        <p class="font-bold">Bingung memilih bahan untuk ruanganmu?</p>
        <p class="text-sm text-muted">Ceritakan kebutuhanmu, tim kami bantu carikan produknya.</p>
      </div>
      <WhatsAppButton class="shrink-0" :href="generalLink()" label="Tanya via WhatsApp" />
    </div>
  </div>
</template>
