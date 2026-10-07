<script setup>
import { CircleCheck, LoaderCircle } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import RingkasanPoin from './RingkasanPoin.vue'

// Isian transaksi Belanja (EARN): nomor nota, total belanja, catatan, dan poin yang dihitung otomatis.
const nota = defineModel('nota', { type: String, required: true })
const total = defineModel('total', { type: String, required: true })
const catatan = defineModel('catatan', { type: String, required: true })

defineProps({
  errors: { type: Object, default: () => ({}) },
  // '' | memeriksa | tersedia | dipakai
  statusNota: { type: String, default: '' },
  // Poin dari total belanja; null bila rasio belum dimuat.
  poin: { type: Number, default: null },
  totalAngka: { type: Number, default: 0 },
  rasio: { type: Number, default: null },
})

// `periksa-nota`: keluar field nota (dicek ke server). `rapikan-total`: keluar field total.
const emit = defineEmits(['periksa-nota', 'rapikan-total'])
</script>

<template>
  <div class="flex flex-col gap-4">
    <Input
      v-model="nota"
      label="Nomor nota"
      required
      autocomplete="off"
      class="font-mono uppercase"
      placeholder="NT-2026-10-0457"
      :error="errors.reference_no ?? ''"
      :hint="statusNota === 'tersedia' ? 'Nomor nota belum pernah diinput' : 'Sesuai nota kasir. Satu nota hanya bisa diinput sekali.'"
      @blur="emit('periksa-nota')"
    >
      <template #akhir>
        <LoaderCircle v-if="statusNota === 'memeriksa'" class="mr-3 size-4.5 animate-spin text-muted" aria-hidden="true" />
        <CircleCheck v-else-if="statusNota === 'tersedia'" class="mr-3 size-4.5 text-success" aria-hidden="true" />
      </template>
    </Input>
    <Input
      v-model="total"
      label="Total belanja"
      required
      prefix="Rp"
      inputmode="numeric"
      autocomplete="off"
      placeholder="4.250.000"
      :error="errors.purchase_amount ?? ''"
      @blur="emit('rapikan-total')"
    />
    <RingkasanPoin jenis="EARN" :perubahan="poin ?? 0" :total="totalAngka" :rasio="rasio" />
    <Input v-model="catatan" label="Catatan (opsional)" autocomplete="off" placeholder="Contoh: belanja keramik proyek rumah" />
  </div>
</template>
