<script setup>
import { ref } from 'vue'
import { Download } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'

const VARIAN = ['primary', 'secondary', 'outline', 'ghost', 'whatsapp', 'danger']
const memuat = ref(false)
const { generalLink } = useWhatsApp()

function simulasiMuat() {
  memuat.value = true
  setTimeout(() => (memuat.value = false), 1500)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Varian (ukuran md, 44px)</p>
      <div class="flex flex-wrap gap-3">
        <Button v-for="v in VARIAN" :key="v" :variant="v">{{ v }}</Button>
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Ukuran sm / md / lg, dengan ikon</p>
      <div class="flex flex-wrap items-center gap-3">
        <Button size="sm">Kecil</Button>
        <Button>
          <template #icon><Download class="size-4" aria-hidden="true" /></template>
          Unduh datasheet
        </Button>
        <Button size="lg">Jelajahi katalog</Button>
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Loading (klik), nonaktif, tautan</p>
      <div class="flex flex-wrap items-center gap-3">
        <Button :loading="memuat" @click="simulasiMuat">{{ memuat ? 'Menyimpan…' : 'Simpan' }}</Button>
        <Button variant="outline" disabled>Nonaktif</Button>
        <Button variant="ghost" to="/">RouterLink ke beranda</Button>
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">WhatsApp (nomor dari pengaturan toko) dan salin</p>
      <div class="flex flex-wrap items-center gap-3">
        <WhatsAppButton :href="generalLink()" />
        <WhatsAppButton :href="null" label="Nomor belum tersedia" />
        <CopyButton text="KRM-4040-01" label="Salin SKU" copied-label="SKU disalin" />
      </div>
    </div>
  </div>
</template>
