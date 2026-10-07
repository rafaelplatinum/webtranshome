<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Eye, EyeOff } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import ProductGallery from '@/components/produk/ProductGallery.vue'
import ProductSummary from '@/components/produk/ProductSummary.vue'
import ProductTabs from '@/components/produk/ProductTabs.vue'
import { usePermission } from '@/composables/usePermission'
import { getAdminProduct } from '@/services/adminProductService'
import { pesanError } from '@/services/errors'

// Pratinjau tampilan detail produk untuk admin, termasuk produk nonaktif (halaman publiknya 404).
const route = useRoute()
const { can } = usePermission()
const status = ref('loading') // loading | ready | tidak-ada | error
const produk = ref(null)
const pesan = ref('')

async function muat() {
  status.value = 'loading'
  try {
    produk.value = await getAdminProduct(route.params.id)
    status.value = 'ready'
  } catch (error) {
    if (error?.response?.status === 404) status.value = 'tidak-ada'
    else {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }
}

onMounted(muat)
</script>

<template>
  <div v-if="status === 'loading'" class="grid gap-6 lg:grid-cols-2" aria-busy="true">
    <span class="sr-only">Memuat pratinjau…</span>
    <Skeleton class="aspect-square rounded-xl" />
    <Skeleton class="h-80 rounded-xl" />
  </div>

  <section v-else-if="status === 'tidak-ada'" class="mx-auto flex max-w-md flex-col items-center gap-4 py-16 text-center">
    <h2 class="text-2xl font-bold">Produk tidak ditemukan</h2>
    <Button to="/admin/produk">Kembali ke daftar produk</Button>
  </section>

  <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

  <div v-else class="flex flex-col gap-6">
    <div
      role="note"
      class="flex flex-wrap items-center gap-3 rounded-xl px-4 py-3"
      :class="produk.is_active ? 'bg-success-soft text-success-ink' : 'bg-warning-soft text-warning-ink'"
    >
      <component :is="produk.is_active ? Eye : EyeOff" class="size-5 shrink-0" aria-hidden="true" />
      <p class="min-w-0 flex-1 text-sm font-semibold">
        {{ produk.is_active ? 'Produk ini tampil di website.' : 'Tidak tampil di website. Pratinjau ini hanya terlihat oleh admin.' }}
      </p>
      <Button v-if="can('product.update')" variant="outline" size="sm" :to="`/admin/produk/${produk.id}`">Ubah produk</Button>
      <Button v-if="produk.is_active" variant="ghost" size="sm" :href="`/produk/${produk.slug}`">Buka halaman publik</Button>
    </div>

    <div class="rounded-xl border border-border bg-background p-4 sm:p-6">
      <div class="grid items-start gap-8 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)] lg:gap-12">
        <ProductGallery :images="produk.images ?? []" :nama="produk.name" />
        <ProductSummary :produk="produk" tag-judul="h2" />
      </div>
      <ProductTabs class="mt-10" :produk="produk" />
    </div>
  </div>
</template>
