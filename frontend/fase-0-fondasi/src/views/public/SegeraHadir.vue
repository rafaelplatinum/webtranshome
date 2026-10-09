<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Hammer } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

// Pengganti sementara untuk halaman yang dibangun di fase berikutnya (meta.fase di route).
// Di area admin judul halaman sudah ada di header (h1), jadi di sini memakai h2.
const route = useRoute()
const judul = computed(() => route.meta.title ?? 'Halaman ini')
const diAdmin = computed(() => route.meta.area === 'admin')
</script>

<template>
  <section class="mx-auto flex max-w-md flex-col items-center gap-4 px-4 py-16 text-center">
    <Hammer class="size-10 text-primary" aria-hidden="true" />
    <component :is="diAdmin ? 'h2' : 'h1'" class="text-2xl font-bold">{{ judul }} sedang dibangun</component>
    <p class="text-muted">
      Halaman ini dikerjakan di Fase {{ route.meta.fase ?? '-' }} rencana build. Sementara itu, kamu bisa kembali ke
      {{ diAdmin ? 'dashboard' : 'beranda' }}.
    </p>
    <Button :to="diAdmin ? '/admin' : '/'">{{ diAdmin ? 'Kembali ke dashboard' : 'Kembali ke beranda' }}</Button>
  </section>
</template>
