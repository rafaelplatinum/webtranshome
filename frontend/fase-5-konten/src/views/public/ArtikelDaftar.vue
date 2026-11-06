<script setup>
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import DaftarArtikel from '@/components/konten/DaftarArtikel.vue'
import { useSeo } from '@/composables/useSeo'

useSeo({ title: 'Artikel dan tips', description: 'Tips memilih dan merawat bahan bangunan dari Transhome, Sleman.' })
</script>

<template>
  <div class="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-4 sm:px-6 lg:py-6">
    <div class="flex flex-col gap-2">
      <Breadcrumb :items="[{ label: 'Beranda', to: '/' }, { label: 'Artikel' }]" />
      <h1 class="text-2xl font-extrabold sm:text-3xl">Artikel dan tips</h1>
      <p class="max-w-2xl text-muted">Panduan singkat memilih, menghitung, dan merawat bahan bangunan untuk rumahmu.</p>
    </div>
    <DaftarArtikel tipe="ARTIKEL" judul="Daftar artikel" judul-kosong="Belum ada artikel" deskripsi-kosong="Artikel dan tips akan muncul di sini." />
  </div>
</template>
