<script setup>
import { computed } from 'vue'
import { formatAngka, formatRupiah } from '@/utils/format'
import { tandaPoin } from '@/utils/poin'

const props = defineProps({
  jenis: { type: String, required: true },
  // EARN: poin dari total belanja (null bila rasio belum dimuat). REDEEM/ADJUST: perubahan saldo bertanda.
  perubahan: { type: Number, default: 0 },
  total: { type: Number, default: 0 },
  rasio: { type: Number, default: null },
  saldo: { type: Number, default: 0 },
  saldoSetelah: { type: Number, default: 0 },
  // REDEEM/ADJUST: isian poin sudah berupa angka yang valid.
  poinValid: { type: Boolean, default: false },
})

// Kotak hasil hitung, berubah langsung saat mengetik (sketsa "Poin yang didapat").
const isi = computed(() => {
  if (props.jenis === 'EARN') {
    if (!props.rasio) return { judul: 'Poin yang didapat', ket: 'Poin dihitung otomatis dari total belanja saat disimpan.', angka: '—', kelas: 'bg-info-soft text-info-ink' }
    const ketRasio = `Otomatis: 1 poin per ${formatRupiah(props.rasio)}`
    if (!props.total) return { judul: 'Poin yang didapat', ket: ketRasio, angka: '—', kelas: 'bg-subtle text-foreground' }
    if (props.perubahan < 1) return { judul: 'Belum dapat poin', ket: `Total belanja di bawah ${formatRupiah(props.rasio)}.`, angka: '0', kelas: 'bg-warning-soft text-warning-ink' }
    return { judul: 'Poin yang didapat', ket: ketRasio, angka: tandaPoin(props.perubahan), kelas: 'bg-success-soft text-success-ink' }
  }
  const judul = props.jenis === 'REDEEM' ? 'Saldo setelah penukaran' : 'Saldo setelah penyesuaian'
  const ket = `Saldo sekarang ${formatAngka(props.saldo)} poin`
  if (!props.poinValid) return { judul, ket, angka: '—', kelas: 'bg-subtle text-foreground' }
  if (props.saldoSetelah < 0) return { judul: 'Saldo tidak cukup', ket, angka: formatAngka(props.saldoSetelah), kelas: 'bg-danger-soft text-danger-ink' }
  return { judul, ket: `${ket} · ${tandaPoin(props.perubahan)} poin`, angka: formatAngka(props.saldoSetelah), kelas: 'bg-primary-soft text-primary-hover' }
})
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-3 rounded-lg px-4 py-3.5 transition-colors duration-150" :class="isi.kelas">
    <span class="flex flex-col">
      <strong>{{ isi.judul }}</strong>
      <span class="text-xs">{{ isi.ket }}</span>
    </span>
    <strong class="text-3xl leading-none tabular-nums" aria-live="polite" data-poin-hasil>{{ isi.angka }}</strong>
  </div>
</template>
