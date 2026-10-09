<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import ChipNav from '@/components/katalog/ChipNav.vue'
import KatalogDaftar from '@/components/katalog/KatalogDaftar.vue'
import { useSeo } from '@/composables/useSeo'
import { getRoom, listRooms, ruanganTersimpan } from '@/services/catalogService'
import { pesanError } from '@/services/errors'

// /ruangan/:slug : katalog dengan filter ruangan terkunci (komponen katalog dipakai ulang).
const route = useRoute()
const status = ref('loading') // loading | ready | tidak-ada | error
const ruangan = ref(null)
// Daftar ruangan yang sudah dimuat dipakai langsung agar render pertama lengkap (posisi scroll saat Back tepat).
const semuaRuangan = ref(ruanganTersimpan() ?? [])
const pesan = ref('')

async function muat() {
  status.value = 'loading'
  // Judul langsung dari daftar ruangan bila ada; detail dari API tetap dimuat untuk memastikan masih aktif.
  ruangan.value = semuaRuangan.value.find((r) => r.slug === route.params.slug) ?? null
  try {
    ruangan.value = await getRoom(route.params.slug)
    status.value = 'ready'
  } catch (error) {
    if (error?.response?.status === 404) {
      status.value = 'tidak-ada'
      return
    }
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(
  () => route.params.slug,
  (slug) => {
    if (slug) muat()
  },
  { immediate: true },
)

// Navigasi antar-ruangan; kalau gagal dimuat cukup tidak ditampilkan.
onMounted(async () => {
  if (semuaRuangan.value.length) return
  try {
    semuaRuangan.value = await listRooms()
  } catch {
    semuaRuangan.value = []
  }
})

const kunci = computed(() => ({ room: route.params.slug }))
const navRuangan = computed(() =>
  semuaRuangan.value.map((r) => ({ key: r.slug, label: r.name, to: `/ruangan/${r.slug}`, aktif: r.slug === route.params.slug })),
)

useSeo({
  title: () => (status.value === 'ready' ? `Produk untuk ${ruangan.value.name}` : status.value === 'tidak-ada' ? 'Ruangan tidak ditemukan' : null),
  description: () =>
    status.value === 'ready'
      ? `Bahan bangunan dan perlengkapan untuk ${ruangan.value.name.toLowerCase()} di Transhome Sleman. Cek harga dan stok, tanya lewat WhatsApp.`
      : null,
})
</script>

<template>
  <NotFoundState
    v-if="status === 'tidak-ada'"
    title="Ruangan tidak ditemukan"
    description="Ruangan yang kamu buka tidak ada atau sudah tidak aktif."
    action-label="Lihat semua ruangan"
    action-to="/ruangan"
  />

  <div v-else class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Ruangan', to: '/ruangan' }, { label: ruangan?.name ?? 'Ruangan' }]" />

    <div class="mt-3 flex items-center gap-4">
      <div class="h-16 w-24 shrink-0 overflow-hidden rounded-xl border border-border bg-subtle sm:h-20 sm:w-32">
        <img v-if="ruangan?.image_cover" :src="ruangan.image_cover" alt="" width="128" height="80" class="size-full object-cover" />
      </div>
      <div class="flex min-w-0 flex-col gap-1">
        <h1 class="text-2xl font-extrabold sm:text-3xl">
          <template v-if="ruangan">{{ ruangan.name }}</template>
          <Skeleton v-else class="inline-block h-8 w-40 align-middle" />
        </h1>
        <p v-if="ruangan" class="text-sm text-muted">{{ ruangan.product_count }} produk cocok untuk {{ ruangan.name.toLowerCase() }}</p>
        <Skeleton v-else class="h-5 w-48" />
      </div>
    </div>

    <ChipNav v-if="navRuangan.length" class="mt-4" label="Pilih ruangan" :items="navRuangan" />

    <ErrorState v-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />
    <KatalogDaftar
      v-else
      class="mt-6"
      :kunci="kunci"
      :sembunyikan="['room']"
      :topik="ruangan ? `produk untuk ${ruangan.name.toLowerCase()}` : ''"
    />
  </div>
</template>
