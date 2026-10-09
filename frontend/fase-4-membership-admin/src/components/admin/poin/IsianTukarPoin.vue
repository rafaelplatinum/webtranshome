<script setup>
import { computed } from 'vue'
import Input from '@/components/ui/Input.vue'
import { formatAngka } from '@/utils/format'
import RingkasanPoin from './RingkasanPoin.vue'

// Isian Tukar poin (REDEEM) dan Penyesuaian (ADJUST): jumlah poin + catatan/alasan wajib, dengan saldo setelahnya.
const poin = defineModel('poin', { type: String, required: true })
const catatan = defineModel('catatan', { type: String, required: true })

const props = defineProps({
  // REDEEM | ADJUST
  jenis: { type: String, required: true },
  errors: { type: Object, default: () => ({}) },
  saldo: { type: Number, default: 0 },
  perubahan: { type: Number, default: 0 },
  saldoSetelah: { type: Number, default: 0 },
  poinValid: { type: Boolean, default: false },
})

const tukar = computed(() => props.jenis === 'REDEEM')
</script>

<template>
  <div class="flex flex-col gap-4">
    <Input
      v-model="poin"
      :label="tukar ? 'Jumlah poin ditukar' : 'Poin (+/−)'"
      required
      :inputmode="tukar ? 'numeric' : 'text'"
      autocomplete="off"
      :placeholder="tukar ? '10' : '-2'"
      :hint="tukar ? `Saldo saat ini ${formatAngka(saldo)} poin.` : 'Pakai tanda minus untuk mengurangi, mis. -2.'"
      :error="errors.points ?? ''"
    />
    <RingkasanPoin :jenis="jenis" :perubahan="perubahan" :saldo="saldo" :saldo-setelah="saldoSetelah" :poin-valid="poinValid" />
    <Input
      v-model="catatan"
      :label="tukar ? 'Catatan' : 'Alasan'"
      required
      autocomplete="off"
      :placeholder="tukar ? 'Contoh: tukar potongan harga di kasir' : 'Contoh: koreksi nota ganda'"
      :error="errors.note ?? ''"
    />
  </div>
</template>
