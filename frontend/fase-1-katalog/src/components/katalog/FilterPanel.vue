<script setup>
import { computed } from 'vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import FilterOpsi from './FilterOpsi.vue'
import RentangHarga from './RentangHarga.vue'
import { LABEL_STOK, labelSpesifikasi } from '@/utils/produk'

const props = defineProps({
  // Nilai filter: { brand: [], stock: [], room: [], 'attr[ukuran]': [], min, max }
  nilai: { type: Object, required: true },
  // meta.facets dari GET /products; null selama muatan pertama.
  facets: { type: Object, default: null },
  // Grup yang dikunci halaman, mis. ['brand'] di /brand/:slug atau ['room'] di /ruangan/:slug.
  sembunyikan: { type: Array, default: () => [] },
})

// Panel ini tidak mengubah URL sendiri: sidebar desktop langsung menerapkan, bottom sheet HP menyimpan ke draf.
const emit = defineEmits(['toggle', 'harga'])

const grup = computed(() => {
  const f = props.facets
  if (!f) return []
  const tampil = (kunci) => !props.sembunyikan.includes(kunci)
  const hasil = []
  if (tampil('brand') && f.brands?.length) {
    hasil.push({ kunci: 'brand', judul: 'Brand', opsi: f.brands.map((b) => ({ value: b.slug, label: b.name, count: b.count })) })
  }
  hasil.push({ kunci: 'harga' })
  if (tampil('stock') && f.stock?.length) {
    hasil.push({ kunci: 'stock', judul: 'Status stok', opsi: f.stock.map((s) => ({ value: s.value, label: LABEL_STOK[s.value] ?? s.value, count: s.count })) })
  }
  if (tampil('room') && f.rooms?.length) {
    hasil.push({ kunci: 'room', judul: 'Ruangan', opsi: f.rooms.map((r) => ({ value: r.slug, label: r.name, count: r.count })) })
  }
  for (const a of f.attributes ?? []) {
    hasil.push({ kunci: `attr[${a.key}]`, judul: labelSpesifikasi(a.key), opsi: a.values.map((v) => ({ value: v.value, label: v.value, count: v.count })) })
  }
  return hasil
})
</script>

<template>
  <div v-if="!facets" class="flex flex-col gap-6" aria-hidden="true">
    <div v-for="n in 3" :key="n" class="flex flex-col gap-3">
      <Skeleton class="h-4 w-24" />
      <Skeleton class="h-4 w-full" />
      <Skeleton class="h-4 w-4/5" />
      <Skeleton class="h-4 w-3/5" />
    </div>
  </div>
  <div v-else class="flex flex-col divide-y divide-border">
    <div v-for="g in grup" :key="g.kunci" class="py-4 first:pt-0 last:pb-0">
      <RentangHarga v-if="g.kunci === 'harga'" :min="nilai.min ?? ''" :max="nilai.max ?? ''" :batas="facets.price" @ubah="emit('harga', $event)" />
      <FilterOpsi v-else :judul="g.judul" :opsi="g.opsi" :terpilih="nilai[g.kunci] ?? []" @toggle="emit('toggle', g.kunci, $event)" />
    </div>
  </div>
</template>
