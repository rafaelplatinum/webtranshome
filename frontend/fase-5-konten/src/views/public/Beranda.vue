<script setup>
import AjakanMember from '@/components/beranda/AjakanMember.vue'
import BannerPromo from '@/components/beranda/BannerPromo.vue'
import StatStrip from '@/components/beranda/StatStrip.vue'
import BrandPilihan from '@/components/beranda/BrandPilihan.vue'
import HeroBanner from '@/components/beranda/HeroBanner.vue'
import InspirasiRuangan from '@/components/beranda/InspirasiRuangan.vue'
import JasaBeranda from '@/components/beranda/JasaBeranda.vue'
import KategoriTahap from '@/components/beranda/KategoriTahap.vue'
import Keunggulan from '@/components/beranda/Keunggulan.vue'
import ProdukUnggulan from '@/components/beranda/ProdukUnggulan.vue'
import PromoEvent from '@/components/beranda/PromoEvent.vue'
import { useSeo } from '@/composables/useSeo'
import PageContactCta from '@/components/public/PageContactCta.vue'

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
    <StatStrip />
    <!-- v-reveal: tiap seksi memudar dan naik sedikit saat masuk layar (composables/useAnimasi.js). -->
    <div class="th-container flex flex-col gap-14 py-14 sm:gap-16 sm:py-16 lg:gap-20 lg:py-20">
      <KategoriTahap v-reveal />
      <ProdukUnggulan v-reveal />
    </div>
    <BannerPromo v-reveal />
    <div class="th-container flex flex-col gap-14 pt-14 sm:gap-16 sm:pt-16 lg:gap-20 lg:pt-20">
      <!-- Jasa (KEPUTUSAN #20, 2 Nov 2026): setelah inspirasi ruangan, sebelum brand. -->
      <JasaBeranda v-reveal />
      <PromoEvent v-reveal />
      <InspirasiRuangan v-reveal />
      <BrandPilihan v-reveal />
      <Keunggulan v-reveal />
      <AjakanMember v-reveal />
    </div>
    <PageContactCta class="mt-14 sm:mt-16 lg:mt-20" v-reveal heading="Sudah punya daftar kebutuhan?" description="Kirimkan ke tim Transhome untuk cek stok, harga, dan pengiriman." />
  </div>
</template>
