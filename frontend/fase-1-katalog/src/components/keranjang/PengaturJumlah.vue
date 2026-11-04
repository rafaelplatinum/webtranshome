<script setup>
import { computed, ref, useId, watch } from 'vue'
import { Minus, Plus } from 'lucide-vue-next'

// Isian jumlah dengan tombol − dan + (detail produk & baris keranjang). Ketikan di luar batas dirapikan saat
// keluar dari isian / Enter; panah atas/bawah menambah/mengurangi.
const model = defineModel({ type: Number, default: 1 })

const props = defineProps({
  min: { type: Number, default: 1 },
  max: { type: Number, default: 9999 },
  // Satuan jual di samping isian, mis. "sak".
  satuan: { type: String, default: '' },
  label: { type: String, default: 'Jumlah' },
  hideLabel: { type: Boolean, default: false },
  // Nama produk untuk pembaca layar di keranjang ("Jumlah Semen PCC 50 kg").
  nama: { type: String, default: '' },
})

const id = useId()
const teks = ref(String(model.value))
watch(model, (n) => {
  teks.value = String(n)
})

const labelLengkap = computed(() => [props.label, props.nama].filter(Boolean).join(' '))
// Label tombol: 'Kurangi jumlah Semen PCC 50 kg' (nama produk tidak diubah huruf kecil).
const sasaran = computed(() => [props.label.toLowerCase(), props.nama].filter(Boolean).join(' '))

function atur(n) {
  const nilai = Math.min(props.max, Math.max(props.min, Math.floor(Number(n) || 0)))
  teks.value = String(nilai)
  if (nilai !== model.value) model.value = nilai
}

function selesaiKetik() {
  const angka = Number.parseInt(teks.value.replace(/\D/g, ''), 10)
  atur(Number.isNaN(angka) ? model.value : angka)
}

const kelasTombol =
  'inline-flex size-11 shrink-0 items-center justify-center text-foreground transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent'
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label :for="`${id}-isian`" :class="hideLabel ? 'sr-only' : 'text-sm font-semibold'">{{ hideLabel ? labelLengkap : label }}</label>
    <div class="flex items-center gap-2">
      <div class="inline-flex items-center overflow-hidden rounded-md border border-border-strong bg-surface">
        <button type="button" :class="kelasTombol" :disabled="model <= min" :aria-label="`Kurangi ${sasaran}`" @click="atur(model - 1)">
          <Minus class="size-4" aria-hidden="true" />
        </button>
        <input
          :id="`${id}-isian`"
          v-model="teks"
          type="text"
          inputmode="numeric"
          autocomplete="off"
          :aria-describedby="satuan ? `${id}-satuan` : undefined"
          class="h-11 w-16 min-w-0 border-x border-border bg-surface text-center text-base font-semibold tabular-nums focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset"
          data-isian-jumlah
          @blur="selesaiKetik"
          @keydown.enter.prevent="selesaiKetik"
          @keydown.up.prevent="atur(model + 1)"
          @keydown.down.prevent="atur(model - 1)"
        />
        <button type="button" :class="kelasTombol" :disabled="model >= max" :aria-label="`Tambah ${sasaran}`" @click="atur(model + 1)">
          <Plus class="size-4" aria-hidden="true" />
        </button>
      </div>
      <span v-if="satuan" :id="`${id}-satuan`" class="text-sm text-muted">{{ satuan }}</span>
    </div>
  </div>
</template>
