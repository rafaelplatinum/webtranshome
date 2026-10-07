<script setup>
import { computed, nextTick, ref, useId, watch } from 'vue'
import { ChevronDown } from 'lucide-vue-next'

// Pilihan yang bisa diketik untuk mencari (pola ARIA combobox + listbox), untuk daftar panjang seperti brand.
// Pilihan hanya berubah saat opsi dipilih (klik / Enter); keluar tanpa memilih mengembalikan label semula.
const model = defineModel({ type: [String, Number], default: '' })

const props = defineProps({
  label: { type: String, required: true },
  // [{ value, label }]; opsi bernilai '' (mis. "Semua brand") selalu di atas dan ikut dicari.
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Ketik untuk mencari' },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  // Teks bila tidak ada yang cocok, mis. "Brand tidak ditemukan".
  kosong: { type: String, default: 'Tidak ada yang cocok' },
})

const id = useId()
const idDaftar = `${id}-daftar`
const input = ref(null)
const daftar = ref(null)
const buka = ref(false)
const teks = ref('')
const aktif = ref(-1)

const terpilih = computed(() => props.options.find((o) => String(o.value) === String(model.value ?? '')) ?? null)
const normal = (s) => String(s ?? '').toLowerCase().trim()
// Saat teks sama dengan pilihan sekarang (baru dibuka), tampilkan semua opsi.
const hasil = computed(() => {
  const q = normal(teks.value)
  if (!q || q === normal(terpilih.value?.label)) return props.options
  return props.options.filter((o) => normal(o.label).includes(q))
})
const keterangan = computed(() => (props.error ? `${id}-error` : props.hint ? `${id}-hint` : undefined))

watch(terpilih, (o) => !buka.value && (teks.value = o?.label ?? ''), { immediate: true })
watch(hasil, () => (aktif.value = hasil.value.length ? 0 : -1))

async function gulirKeAktif() {
  await nextTick()
  daftar.value?.querySelector(`#${CSS.escape(`${id}-opsi-${aktif.value}`)}`)?.scrollIntoView({ block: 'nearest' })
}

function bukaDaftar() {
  if (buka.value) return
  buka.value = true
  aktif.value = Math.max(0, hasil.value.findIndex((o) => String(o.value) === String(model.value ?? '')))
  gulirKeAktif()
}

function tutup({ kembalikan = true } = {}) {
  buka.value = false
  if (kembalikan) teks.value = terpilih.value?.label ?? ''
}

function pilih(o) {
  model.value = o.value
  teks.value = o.label
  tutup({ kembalikan: false })
}

function onInput() {
  buka.value = true
}

function onKeydown(e) {
  if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
    e.preventDefault()
    if (!buka.value) return bukaDaftar()
    const n = hasil.value.length
    if (!n) return
    aktif.value = (aktif.value + (e.key === 'ArrowDown' ? 1 : -1) + n) % n
    gulirKeAktif()
  } else if (e.key === 'Enter' && buka.value) {
    e.preventDefault()
    if (hasil.value[aktif.value]) pilih(hasil.value[aktif.value])
  } else if (e.key === 'Escape' && buka.value) {
    e.stopPropagation()
    tutup()
  }
}
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label :for="id" class="text-sm font-semibold">{{ label }}</label>
    <div class="relative" @focusout="(e) => !e.currentTarget.contains(e.relatedTarget) && tutup()">
      <input
        :id="id"
        ref="input"
        v-model="teks"
        type="text"
        role="combobox"
        autocomplete="off"
        :placeholder="placeholder"
        aria-autocomplete="list"
        :aria-expanded="buka ? 'true' : 'false'"
        :aria-controls="idDaftar"
        :aria-activedescendant="buka && aktif >= 0 ? `${id}-opsi-${aktif}` : undefined"
        :aria-invalid="error ? 'true' : undefined"
        :aria-describedby="keterangan"
        class="h-11 w-full min-w-0 rounded-md border bg-surface pr-10 pl-3 text-base text-foreground transition-colors duration-150 placeholder:text-faint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        :class="error ? 'border-danger' : 'border-border-strong focus-visible:border-primary'"
        @focus="$event.target.select()"
        @click="bukaDaftar"
        @input="onInput"
        @keydown="onKeydown"
      />
      <button
        type="button"
        tabindex="-1"
        aria-hidden="true"
        class="absolute inset-y-0 right-0 flex w-10 items-center justify-center text-muted"
        @mousedown.prevent="buka ? tutup() : (input.focus(), bukaDaftar())"
      >
        <ChevronDown class="size-4.5 transition-transform duration-150" :class="buka && 'rotate-180'" />
      </button>
      <ul
        v-show="buka"
        :id="idDaftar"
        ref="daftar"
        role="listbox"
        :aria-label="label"
        class="absolute inset-x-0 top-full z-40 mt-1 max-h-64 overflow-y-auto rounded-md border border-border bg-surface py-1 shadow-lg"
      >
        <li
          v-for="(o, i) in hasil"
          :id="`${id}-opsi-${i}`"
          :key="o.value"
          role="option"
          :aria-selected="String(o.value) === String(model ?? '') ? 'true' : 'false'"
          class="flex min-h-10 cursor-pointer items-center px-3 text-sm"
          :class="[i === aktif && 'bg-subtle', String(o.value) === String(model ?? '') && 'font-semibold']"
          @mousedown.prevent="pilih(o)"
          @mousemove="aktif = i"
        >
          {{ o.label }}
        </li>
        <li v-if="!hasil.length" class="px-3 py-2 text-sm text-muted" role="presentation">{{ kosong }}</li>
      </ul>
    </div>
    <p v-if="error" :id="`${id}-error`" class="text-sm font-medium text-danger-ink">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-hint`" class="text-sm text-muted">{{ hint }}</p>
  </div>
</template>
