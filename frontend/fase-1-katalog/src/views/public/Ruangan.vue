<script setup>
import { onMounted, ref } from 'vue'
import { LayoutDashboard } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import RoomCarousel from '@/components/katalog/RoomCarousel.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import PageContactCta from '@/components/public/PageContactCta.vue'
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
  <div class="flex flex-col gap-12 sm:gap-16">
    <section class="bg-th-navy py-10 text-white sm:py-14 lg:py-16">
      <div class="th-container flex flex-col gap-5">
        <Breadcrumb variant="light" :items="[{ label: 'Beranda', to: '/' }, { label: 'Ruangan' }]" />
        <span class="th-eyebrow text-orange-200">Belanja per ruangan</span>
        <h1 class="max-w-4xl text-4xl leading-[1.08] font-extrabold tracking-tight sm:text-5xl lg:text-[56px]">Pilih ruangan untuk mulai merencanakan.</h1>
        <p class="max-w-2xl text-base leading-relaxed text-white/80 sm:text-lg">Pilih ruangan yang sedang kamu bangun atau renovasi. Kami kumpulkan produk yang cocok untuk tiap ruangan.</p>
        <p v-if="status === 'ready'" class="font-semibold text-th-yellow">{{ ruangan.length }} pilihan ruangan</p>
      </div>
    </section>

    <div class="th-container flex flex-col gap-8">
      <div v-if="status === 'loading'" class="flex flex-col gap-4" aria-busy="true">
        <span class="sr-only">Memuat ruangan…</span><Skeleton class="h-9 w-52" /><div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4"><Skeleton v-for="n in 4" :key="n" class="aspect-4/3 rounded-[20px]" /></div>
      </div>
      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
      <EmptyState v-else-if="ruangan.length === 0" title="Belum ada ruangan" description="Daftar ruangan sedang disiapkan. Sementara itu, kamu bisa menjelajah katalog."><template #icon><LayoutDashboard class="size-10" aria-hidden="true" /></template></EmptyState>
      <RoomCarousel v-else :ruangan="ruangan"><template #judul><h2 class="text-2xl font-extrabold tracking-tight sm:text-3xl">Pilih ruangan</h2></template></RoomCarousel>
    </div>

    <PageContactCta v-if="status === 'ready'" heading="Bingung memilih bahan untuk ruanganmu?" description="Ceritakan kebutuhanmu, tim kami bantu carikan produknya." button-label="Tanya via WhatsApp" :context="generalLink()" />
  </div>
</template>
