<script setup>
import { ShieldX } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
</script>

<template>
  <section class="mx-auto flex max-w-md flex-col items-center gap-4 px-4 py-20 text-center">
    <ShieldX class="size-10 text-danger" aria-hidden="true" />
    <h1 class="text-2xl font-bold">Tidak punya akses</h1>
    <p class="text-muted">Akun kamu tidak punya izin untuk membuka halaman ini. Hubungi Super Admin kalau kamu membutuhkan akses.</p>
    <Button to="/">Kembali ke beranda</Button>
  </section>
</template>
