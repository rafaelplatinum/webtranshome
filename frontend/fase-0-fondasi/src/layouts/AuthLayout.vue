<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<!-- Layout halaman masuk, daftar, verifikasi, dan password: fokus ke satu form. -->
<template>
  <div class="flex min-h-dvh flex-col">
    <header class="border-b border-border bg-surface">
      <div class="mx-auto flex h-16 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6">
        <RouterLink
          to="/"
          class="rounded-md text-2xl font-extrabold tracking-tight focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          aria-label="Transhome, ke beranda"
        >
          Trans<span class="text-primary">home</span>
        </RouterLink>
        <RouterLink
          to="/"
          class="rounded-sm text-sm font-semibold text-muted transition-colors duration-150 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Kembali ke beranda
        </RouterLink>
      </div>
    </header>
    <main id="konten" class="flex flex-1 justify-center px-4 py-10 sm:py-16">
      <div class="w-full max-w-md">
        <RouterView />
      </div>
    </main>
  </div>
</template>
