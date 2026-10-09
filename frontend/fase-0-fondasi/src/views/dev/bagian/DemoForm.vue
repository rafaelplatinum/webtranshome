<script setup>
import { ref } from 'vue'
import Checkbox from '@/components/ui/Checkbox.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'

const nama = ref('')
const email = ref('budi@')
const hp = ref('')
const urutan = ref('terbaru')
const kategori = ref('')
const setuju = ref(false)
const aktif = ref(true)

const OPSI_URUT = [
  { value: 'terbaru', label: 'Terbaru' },
  { value: 'harga_terendah', label: 'Harga terendah' },
  { value: 'harga_tertinggi', label: 'Harga tertinggi' },
  { value: 'nama', label: 'Nama A–Z' },
]
</script>

<template>
  <div class="grid gap-6 md:grid-cols-2">
    <Input v-model="nama" label="Nama lengkap" required placeholder="Sesuai KTP" autocomplete="name" />
    <Input v-model="email" label="Email" type="email" required error="Format email belum benar." autocomplete="email" />
    <Input v-model="hp" label="Nomor HP (opsional)" type="tel" inputmode="tel" hint="Boleh 08… atau +62…" autocomplete="tel" />
    <Input model-value="TF-000123" label="Kode member" disabled />
    <Select v-model="urutan" label="Urutkan" :options="OPSI_URUT" />
    <Select v-model="kategori" label="Kategori" placeholder="Pilih kategori" :options="[{ value: 'keramik-lantai', label: 'Keramik lantai' }]" error="Kategori wajib dipilih." />
    <Checkbox v-model="setuju" label="Terima info promo dan poin" description="Bisa dimatikan kapan saja di halaman Keamanan." />
    <div class="flex flex-col gap-1">
      <Switch v-model="aktif" label="Tampil di website" />
      <p class="text-sm text-muted">Status: {{ aktif ? 'aktif' : 'nonaktif' }}</p>
    </div>
  </div>
</template>
