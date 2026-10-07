<script setup>
import { nextTick, ref, useId } from 'vue'
import { ArrowDown, ArrowUp, Plus, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

// Baris spesifikasi: [{ id, label, nilai }] → disimpan sebagai JSON specifications (lihat useFormProduk).
const baris = defineModel({ type: Array, default: () => [] })

defineProps({
  // Galat per baris dari validasi: [{ label, nilai }] sejajar dengan `baris`.
  errors: { type: Array, default: () => [] },
})

// Nama atribut yang sering dipakai; ukuran, finishing, warna, bahan, daya juga jadi filter di katalog.
const SARAN = ['Ukuran', 'Finishing', 'Warna', 'Bahan', 'Daya', 'Berat', 'Isi per dus', 'Area aplikasi', 'Tipe', 'Kapasitas', 'Garansi', 'SNI']

const idSaran = useId()
const wadah = ref(null)
let nomor = 0

async function tambah() {
  baris.value = [...baris.value, { id: `spek-baru-${++nomor}`, label: '', nilai: '' }]
  await nextTick()
  const semua = wadah.value?.querySelectorAll('[data-baris]') ?? []
  semua[semua.length - 1]?.querySelector('input')?.focus()
}

async function geser(i, arah) {
  const j = i + arah
  const baru = [...baris.value]
  ;[baru[i], baru[j]] = [baru[j], baru[i]]
  baris.value = baru
  await nextTick()
  wadah.value?.querySelector(`[data-geser="${j}:${arah}"]:not(:disabled)`)?.focus()
}

async function hapus(i) {
  baris.value = baris.value.filter((_, j) => j !== i)
  await nextTick()
  const sisa = wadah.value?.querySelectorAll('[data-baris]')
  ;(sisa?.[i] ?? sisa?.[i - 1])?.querySelector('input')?.focus()
}

const kelasIkon =
  'inline-flex size-11 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-40'
</script>

<template>
  <div ref="wadah" class="flex flex-col gap-3">
    <p v-if="!baris.length" class="text-sm text-muted">Belum ada spesifikasi. Tambahkan misalnya ukuran, warna, atau isi per dus.</p>
    <div v-else class="hidden grid-cols-[minmax(0,1fr)_minmax(0,1.4fr)_8.5rem] gap-2 text-xs font-semibold tracking-wide text-muted uppercase sm:grid">
      <span>Nama atribut</span>
      <span>Isi</span>
    </div>
    <div
      v-for="(b, i) in baris"
      :key="b.id"
      data-baris
      class="grid grid-cols-1 gap-2 rounded-lg border border-border p-2 sm:grid-cols-[minmax(0,1fr)_minmax(0,1.4fr)_8.5rem] sm:items-start sm:border-0 sm:p-0"
    >
      <Input v-model="b.label" :label="`Nama atribut baris ${i + 1}`" hide-label placeholder="Mis. Ukuran" :list="idSaran" autocomplete="off" :error="errors[i]?.label ?? ''" />
      <Input v-model="b.nilai" :label="`Isi atribut baris ${i + 1}`" hide-label placeholder="Mis. 40 × 40 cm" autocomplete="off" :error="errors[i]?.nilai ?? ''" />
      <div class="flex justify-end">
        <button type="button" :class="kelasIkon" :data-geser="`${i}:-1`" :disabled="i === 0" :aria-label="`Naikkan baris ${i + 1}`" @click="geser(i, -1)">
          <ArrowUp class="size-4" aria-hidden="true" />
        </button>
        <button type="button" :class="kelasIkon" :data-geser="`${i}:1`" :disabled="i === baris.length - 1" :aria-label="`Turunkan baris ${i + 1}`" @click="geser(i, 1)">
          <ArrowDown class="size-4" aria-hidden="true" />
        </button>
        <button type="button" :class="[kelasIkon, 'hover:text-danger']" :aria-label="`Hapus baris ${i + 1}`" @click="hapus(i)">
          <X class="size-4" aria-hidden="true" />
        </button>
      </div>
    </div>
    <datalist :id="idSaran">
      <option v-for="s in SARAN" :key="s" :value="s" />
    </datalist>
    <Button variant="outline" size="sm" class="self-start" @click="tambah">
      <template #icon><Plus class="size-4" aria-hidden="true" /></template>
      Tambah baris
    </Button>
  </div>
</template>
