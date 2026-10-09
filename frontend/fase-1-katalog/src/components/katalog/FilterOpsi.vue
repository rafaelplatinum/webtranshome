<script setup>
import { computed, ref, useId } from 'vue'

const props = defineProps({
  judul: { type: String, required: true },
  // [{ value, label, count }] dari meta.facets
  opsi: { type: Array, required: true },
  terpilih: { type: Array, default: () => [] },
  // Opsi di atas batas disembunyikan sampai "Lihat semua" ditekan.
  batas: { type: Number, default: 6 },
  // Lebih dari ini (mis. ratusan brand): muncul kotak cari, dan daftar lengkap bergulir di dalam panel.
  cariBila: { type: Number, default: 10 },
})

const emit = defineEmits(['toggle'])

const semua = ref(false)
const idDaftar = useId()

const cari = ref('')
const bisaCari = computed(() => props.opsi.length > props.cariBila)
const kataCari = computed(() => (bisaCari.value ? cari.value.trim().toLowerCase() : ''))

// Saat diringkas, opsi yang sudah dipilih tetap terlihat meski di luar batas. Saat mencari: semua yang cocok.
const tampil = computed(() => {
  if (kataCari.value) return props.opsi.filter((o) => o.label.toLowerCase().includes(kataCari.value))
  return semua.value ? props.opsi : props.opsi.filter((o, i) => i < props.batas || props.terpilih.includes(o.value))
})
const bergulir = computed(() => bisaCari.value && (semua.value || kataCari.value))

// Opsi tanpa hasil tidak bisa dipilih, kecuali sudah terpilih (agar tetap bisa dilepas).
const nonaktif = (o) => o.count === 0 && !props.terpilih.includes(o.value)
</script>

<template>
  <fieldset class="flex min-w-0 flex-col">
    <legend class="mb-1 text-sm font-bold">{{ judul }}</legend>
    <label v-if="bisaCari" class="mb-2 block">
      <span class="sr-only">Cari {{ judul.toLowerCase() }}</span>
      <input
        v-model="cari"
        type="search"
        :placeholder="`Cari ${judul.toLowerCase()}`"
        autocomplete="off"
        :aria-controls="idDaftar"
        class="h-10 w-full rounded-md border border-border-strong bg-surface px-3 text-sm placeholder:text-faint focus-visible:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      />
    </label>
    <div :id="idDaftar" class="flex flex-col" :class="bergulir && '-mr-2 max-h-72 overflow-y-auto pr-2'">
      <label
        v-for="o in tampil"
        :key="o.value"
        class="-mx-2 flex min-h-11 items-center gap-3 rounded-md px-2 transition-colors duration-150 md:min-h-9"
        :class="nonaktif(o) ? 'cursor-not-allowed text-faint' : 'cursor-pointer hover:bg-subtle'"
      >
        <input
          type="checkbox"
          class="size-5 shrink-0 cursor-pointer accent-primary disabled:cursor-not-allowed md:size-4.5"
          :checked="terpilih.includes(o.value)"
          :disabled="nonaktif(o)"
          @change="emit('toggle', o.value)"
        />
        <span class="min-w-0 flex-1 text-sm">{{ o.label }}</span>
        <span class="text-xs text-faint tabular-nums">{{ o.count }}<span class="sr-only"> produk</span></span>
      </label>
    </div>
    <p v-if="kataCari && !tampil.length" class="py-2 text-sm text-muted">Tidak ada {{ judul.toLowerCase() }} “{{ cari.trim() }}”.</p>
    <p v-if="kataCari" class="sr-only" aria-live="polite">{{ tampil.length }} {{ judul.toLowerCase() }} cocok</p>
    <button
      v-if="!kataCari && opsi.length > batas"
      type="button"
      class="mt-1 inline-flex min-h-11 items-center self-start rounded-sm text-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-8"
      :aria-expanded="semua ? 'true' : 'false'"
      :aria-controls="idDaftar"
      @click="semua = !semua"
    >
      {{ semua ? 'Tampilkan lebih sedikit' : `Lihat semua (${opsi.length})` }}
    </button>
  </fieldset>
</template>
