<script setup>
import { SearchX } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

// Isi halaman 404: rute tidak dikenal, atau produk/kategori/brand/ruangan yang tidak ada atau nonaktif.
defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
  actionLabel: { type: String, default: 'Kembali ke beranda' },
  actionTo: { type: [String, Object], default: '/' },
})
</script>

<template>
  <section class="mx-auto flex max-w-md flex-col items-center gap-4 px-4 py-20 text-center">
    <SearchX class="size-10 text-faint" aria-hidden="true" />
    <h1 class="text-2xl font-bold">{{ title }}</h1>
    <p v-if="description" class="text-muted">{{ description }}</p>
    <Button :to="actionTo">{{ actionLabel }}</Button>
  </section>
</template>
