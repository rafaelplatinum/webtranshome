<script setup>
import { computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import LogoTranshome from '@/components/layout/LogoTranshome.vue'

const route = useRoute()
const isDaftar = computed(() => route.name === 'daftar')
const CARA = [
  'Daftar pakai email atau akun Google.',
  'Dapat kode member digital di akunmu.',
  'Sebutkan kode ke kasir, poin belanja tercatat.',
]
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
          <LogoTranshome />
        </RouterLink>
        <RouterLink
          to="/"
          class="rounded-sm text-sm font-semibold text-muted transition-colors duration-150 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Kembali ke beranda
        </RouterLink>
      </div>
    </header>
    <main id="konten" class="flex flex-1" :class="isDaftar ? '' : 'justify-center px-4 py-10 sm:py-16'">
      <div v-if="isDaftar" class="grid w-full flex-1 lg:grid-cols-2">
        <section class="flex flex-col justify-center gap-8 overflow-hidden bg-th-navy px-6 py-10 text-white sm:px-10 lg:px-16 lg:py-16">
          <div class="relative mx-auto h-48 w-full max-w-[420px] sm:h-56" aria-hidden="true">
            <div class="absolute inset-x-4 top-3 h-40 rotate-[-7deg] rounded-[22px] bg-th-yellow sm:h-48"></div>
            <div class="absolute inset-x-0 top-0 flex h-40 rotate-[3deg] flex-col justify-between rounded-[22px] border border-white/10 bg-[#243b53] p-5 shadow-2xl sm:h-48 sm:p-6">
              <div class="flex items-center justify-between"><strong class="text-base sm:text-lg">Trans Family</strong><span class="rounded-full bg-th-yellow px-3 py-1 text-[13px] font-extrabold text-th-navy">Member</span></div>
              <div class="flex gap-3"><span class="h-3 w-3/5 rounded bg-white/15"></span><span class="h-3 w-1/4 rounded bg-white/10"></span></div>
              <span class="text-sm text-white/75">Kode member digital</span>
            </div>
          </div>
          <div class="mx-auto flex w-full max-w-xl flex-col gap-5">
            <span class="th-eyebrow text-orange-200">Trans Family · Gratis</span>
            <h1 class="text-3xl leading-tight font-extrabold tracking-tight sm:text-4xl">Jadi member Trans Family, gratis.</h1>
            <ol class="flex flex-col gap-4">
              <li v-for="(cara, i) in CARA" :key="cara" class="flex items-start gap-4 text-base leading-relaxed text-white/85">
                <strong class="text-xl font-extrabold text-th-yellow">{{ String(i + 1).padStart(2, '0') }}</strong><span>{{ cara }}</span>
              </li>
            </ol>
          </div>
        </section>
        <div class="flex items-center justify-center px-4 py-10 sm:px-8 lg:px-12 lg:py-14">
          <div class="w-full max-w-[460px]"><RouterView /></div>
        </div>
      </div>
      <div v-else class="flex w-full justify-center">
        <div class="w-full max-w-[336px]"><RouterView /></div>
      </div>
    </main>
  </div>
</template>
