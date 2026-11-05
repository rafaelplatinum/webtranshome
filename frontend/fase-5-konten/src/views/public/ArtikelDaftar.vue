<script setup>
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import DaftarArtikel from '@/components/konten/DaftarArtikel.vue'
import { useSeo } from '@/composables/useSeo'
import PageContactCta from '@/components/public/PageContactCta.vue'

useSeo({ title: 'Artikel dan tips', description: 'Tips memilih dan merawat bahan bangunan dari Transhome, Sleman.' })
</script>

<template>
  <div class="flex flex-col gap-12 sm:gap-16">
    <header class="th-container flex flex-col gap-4 pt-5 sm:gap-5 sm:pt-8 lg:pt-12">
      <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Artikel' }]" />
      <span class="th-eyebrow">Artikel dan tips</span>
      <h1 class="max-w-4xl text-4xl leading-[1.08] font-extrabold tracking-tight sm:text-5xl lg:text-[52px]">Panduan memilih, menghitung, dan merawat bahan bangunan.</h1>
      <p class="max-w-2xl text-base leading-relaxed text-th-muted sm:text-lg">Panduan singkat memilih, menghitung, dan merawat bahan bangunan untuk rumahmu.</p>
    </header>
    <div class="th-container"><DaftarArtikel tipe="ARTIKEL" judul="Daftar artikel" judul-kosong="Belum ada artikel" deskripsi-kosong="Artikel dan tips akan muncul di sini." featured /></div>
    <PageContactCta class="w-full" heading="Sudah punya daftar kebutuhan?" description="Kirimkan ke tim Transhome untuk cek stok, harga, dan pengiriman." />
  </div>
</template>
