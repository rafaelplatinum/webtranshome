<script setup>
import AjakanMember from '@/components/beranda/AjakanMember.vue'
import BannerPromo from '@/components/beranda/BannerPromo.vue'
import BrandPilihan from '@/components/beranda/BrandPilihan.vue'
import HeroBanner from '@/components/beranda/HeroBanner.vue'
import InspirasiRuangan from '@/components/beranda/InspirasiRuangan.vue'
import KategoriTahap from '@/components/beranda/KategoriTahap.vue'
import Keunggulan from '@/components/beranda/Keunggulan.vue'
import ProdukUnggulan from '@/components/beranda/ProdukUnggulan.vue'
import PromoEvent from '@/components/beranda/PromoEvent.vue'
import { useSeo } from '@/composables/useSeo'

// Beranda = hero banner (gambar/video dari admin, wireframe 01; menggantikan rumah 3D sejak 30 Okt 2026) + seksi
// Beranda (sketsa "Main"). Setiap seksi memuat datanya sendiri, jadi satu seksi yang gagal tidak menggagalkan halaman.
useSeo({
  title: 'Bahan bangunan dan elektronik di Sleman',
  description: 'Katalog bahan bangunan dan elektronik Transhome di Sleman, Yogyakarta. Cek harga dan stok, lalu tanya langsung lewat WhatsApp.',
})
</script>

<template>
  <div>
    <HeroBanner />
    <!-- v-reveal: tiap seksi memudar dan naik sedikit saat masuk layar (composables/useAnimasi.js). -->
    <div class="mx-auto flex max-w-7xl flex-col gap-12 px-4 py-10 sm:px-6 lg:gap-16 lg:py-14">
      <BannerPromo v-reveal />
      <KategoriTahap v-reveal />
      <ProdukUnggulan v-reveal />
      <PromoEvent v-reveal />
      <InspirasiRuangan v-reveal />
      <BrandPilihan v-reveal />
      <Keunggulan v-reveal />
      <AjakanMember v-reveal />
    </div>
  </div>
</template>
