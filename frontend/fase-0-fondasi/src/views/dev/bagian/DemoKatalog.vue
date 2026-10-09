<script setup>
import { onMounted, ref } from 'vue'
import { LayoutGrid, List } from 'lucide-vue-next'
import ProductCard from '@/components/katalog/ProductCard.vue'
import ProductListRow from '@/components/katalog/ProductListRow.vue'
import Button from '@/components/ui/Button.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import { listProducts } from '@/services/productService'
import { pesanError } from '@/services/errors'

const status = ref('loading') // loading | ready | error
const produk = ref([])
const pesan = ref('')
const tampilan = ref('grid')

async function muat() {
  status.value = 'loading'
  try {
    const hasil = await listProducts({ featured: 1, per_page: 4 })
    produk.value = hasil.data
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

onMounted(muat)
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Data dari <code class="font-mono">GET /products?featured=1</code> lewat service (mock atau API sesuai <code class="font-mono">VITE_USE_MOCK</code>).</p>
      <div class="flex gap-2">
        <Button variant="outline" size="sm" @click="muat">Muat ulang</Button>
        <div role="group" aria-label="Tampilan" class="flex overflow-hidden rounded-md border border-border-strong">
          <button
            type="button"
            :aria-pressed="tampilan === 'grid' ? 'true' : 'false'"
            class="inline-flex h-10 items-center gap-1.5 px-3 text-sm font-semibold transition-colors duration-150"
            :class="tampilan === 'grid' ? 'bg-secondary text-secondary-foreground' : 'bg-surface hover:bg-subtle'"
            @click="tampilan = 'grid'"
          >
            <LayoutGrid class="size-4" aria-hidden="true" />Grid
          </button>
          <button
            type="button"
            :aria-pressed="tampilan === 'list' ? 'true' : 'false'"
            class="inline-flex h-10 items-center gap-1.5 px-3 text-sm font-semibold transition-colors duration-150"
            :class="tampilan === 'list' ? 'bg-secondary text-secondary-foreground' : 'bg-surface hover:bg-subtle'"
            @click="tampilan = 'list'"
          >
            <List class="size-4" aria-hidden="true" />Daftar
          </button>
        </div>
      </div>
    </div>

    <div v-if="status === 'loading'" class="grid grid-cols-2 gap-4 md:grid-cols-4" aria-busy="true">
      <span class="sr-only">Memuat produk…</span>
      <SkeletonCard v-for="n in 4" :key="n" />
    </div>
    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
    <div v-else-if="tampilan === 'grid'" class="grid grid-cols-2 gap-4 md:grid-cols-4">
      <ProductCard v-for="p in produk" :key="p.id" :product="p" />
    </div>
    <div v-else class="flex flex-col gap-3">
      <ProductListRow v-for="p in produk" :key="p.id" :product="p" />
    </div>
  </div>
</template>
