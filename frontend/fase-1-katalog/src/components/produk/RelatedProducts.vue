<script setup>
import { onBeforeUnmount, ref, useId, watch } from 'vue'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import ProductCard from '@/components/katalog/ProductCard.vue'
import { listProducts } from '@/services/productService'

const props = defineProps({
  produk: { type: Object, required: true },
})

const id = useId()
const status = ref('loading') // loading | ready | error
const daftar = ref([])
let pengendali = null

async function muat() {
  if (!props.produk.category) {
    daftar.value = []
    status.value = 'ready'
    return
  }
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listProducts(
      { category: props.produk.category.slug, exclude: props.produk.id, per_page: 4 },
      { signal: pengendali.signal },
    )
    daftar.value = hasil.data
    status.value = 'ready'
  } catch (error) {
    if (error?.code !== 'ERR_CANCELED') status.value = 'error'
  }
}

watch(() => props.produk.id, muat, { immediate: true })
onBeforeUnmount(() => pengendali?.abort())
</script>

<template>
  <!-- Tidak ada produk serupa → bagian ini tidak ditampilkan. -->
  <section v-if="status !== 'ready' || daftar.length" :aria-labelledby="id" class="flex flex-col gap-4">
    <div class="flex items-end justify-between gap-4">
      <h2 :id="id" class="text-xl font-bold">Produk serupa</h2>
      <RouterLink
        v-if="produk.category"
        :to="`/katalog/${produk.category.slug}`"
        class="inline-flex min-h-11 items-center rounded-sm text-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
      >
        Lihat semua <span class="sr-only">{{ produk.category.name }}</span>
      </RouterLink>
    </div>

    <div v-if="status === 'loading'" class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4" aria-busy="true">
      <span class="sr-only">Memuat produk serupa…</span>
      <SkeletonCard v-for="n in 4" :key="n" />
    </div>
    <div v-else-if="status === 'error'" class="flex flex-wrap items-center gap-3 rounded-xl border border-border bg-surface p-4">
      <p class="flex-1 text-sm text-muted">Produk serupa gagal dimuat.</p>
      <Button variant="outline" size="sm" @click="muat">Coba lagi</Button>
    </div>
    <ul v-else class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-4">
      <li v-for="p in daftar" :key="p.id" class="flex">
        <ProductCard class="w-full" :product="p" />
      </li>
    </ul>
  </section>
</template>
