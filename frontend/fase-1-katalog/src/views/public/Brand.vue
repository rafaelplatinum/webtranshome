<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import KatalogDaftar from '@/components/katalog/KatalogDaftar.vue'
import { useSeo } from '@/composables/useSeo'
import { getBrand } from '@/services/catalogService'
import { pesanError } from '@/services/errors'

// /brand/:slug : info brand + katalog dengan filter brand terkunci. Daftar produk dimuat bersamaan.
const route = useRoute()
const status = ref('loading') // loading | ready | tidak-ada | error
const brand = ref(null)
const pesan = ref('')

async function muat() {
  status.value = 'loading'
  try {
    brand.value = await getBrand(route.params.slug)
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

const kunci = computed(() => ({ brand: route.params.slug }))
const inisial = computed(() => brand.value?.name?.trim().charAt(0).toUpperCase() ?? '')

useSeo({
  title: () => (status.value === 'ready' ? `Produk ${brand.value.name}` : status.value === 'tidak-ada' ? 'Brand tidak ditemukan' : null),
  description: () =>
    status.value === 'ready' ? `Daftar produk ${brand.value.name} di Transhome Sleman: cek harga, stok, dan tanya lewat WhatsApp.` : null,
})
</script>

<template>
  <NotFoundState
    v-if="status === 'tidak-ada'"
    title="Brand tidak ditemukan"
    description="Brand yang kamu buka tidak ada atau sudah tidak aktif."
    action-label="Lihat katalog"
    action-to="/katalog"
  />

  <div v-else class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Katalog', to: '/katalog' }, { label: brand?.name ?? 'Brand' }]" />

    <div class="mt-3 flex items-center gap-4">
      <div class="flex size-16 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-border bg-surface">
        <img v-if="brand?.logo_url" :src="brand.logo_url" :alt="`Logo ${brand.name}`" width="64" height="64" class="size-full object-contain p-2" />
        <span v-else-if="brand" class="text-2xl font-extrabold text-muted" aria-hidden="true">{{ inisial }}</span>
      </div>
      <div class="flex min-w-0 flex-col gap-1">
        <h1 class="text-2xl font-extrabold sm:text-3xl">
          <template v-if="brand">{{ brand.name }}</template>
          <Skeleton v-else class="inline-block h-8 w-40 align-middle" />
        </h1>
        <p v-if="brand" class="text-sm text-muted">{{ brand.product_count }} produk di Transhome</p>
        <Skeleton v-else class="h-5 w-36" />
      </div>
    </div>

    <ErrorState v-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />
    <KatalogDaftar v-else class="mt-6" :kunci="kunci" :sembunyikan="['brand']" :topik="brand ? `produk ${brand.name}` : ''" />
  </div>
</template>
