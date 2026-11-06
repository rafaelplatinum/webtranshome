<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppFooter from './AppFooter.vue'
import PageContactCta from '@/components/public/PageContactCta.vue'

// Footer situs (permintaan pemilik 6 Nov 2026): pita ajakan WhatsApp + footer tampil di SEMUA halaman
// (publik, masuk/daftar, akun member). Halaman boleh mengganti teks pita lewat `meta.cta` di router.
const CTA_BAWAAN = {
  heading: 'Sudah punya daftar kebutuhan?',
  description: 'Kirimkan ke tim Transhome untuk cek stok, harga, dan pengiriman.',
}

const route = useRoute()
const cta = computed(() => route.meta.cta ?? CTA_BAWAAN)
</script>

<template>
  <div data-footer-situs>
    <PageContactCta v-bind="cta" />
    <AppFooter />
  </div>
</template>
