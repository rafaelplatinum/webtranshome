<script setup>
import { SearchX } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
</script>

<template>
  <section class="mx-auto flex max-w-md flex-col items-center gap-4 px-4 py-16 text-center">
    <SearchX class="size-10 text-faint" aria-hidden="true" />
    <h2 class="text-2xl font-bold">Halaman admin tidak ditemukan</h2>
    <p class="text-muted">Alamat ini tidak ada di panel admin. Pilih menu di samping atau kembali ke dashboard.</p>
    <Button to="/admin">Kembali ke dashboard</Button>
  </section>
</template>
