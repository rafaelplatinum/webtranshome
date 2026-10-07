<script setup>
import { computed, ref, useId, watch } from 'vue'

const props = defineProps({
  min: { type: [String, Number], default: '' },
  max: { type: [String, Number], default: '' },
  // meta.facets.price { min, max }: harga terendah & tertinggi yang ada, untuk placeholder.
  batas: { type: Object, default: null },
})

// Dipanggil saat Enter atau keluar dari field, hanya bila nilainya berubah dan valid.
const emit = defineEmits(['ubah'])

const id = useId()
const minLokal = ref('')
const maksLokal = ref('')
const error = ref('')

const angka = (teks) => String(teks ?? '').replace(/\D/g, '')
const formatAngka = (teks) => (teks ? new Intl.NumberFormat('id-ID').format(Number(teks)) : '')

// Disinkronkan per field, agar mengetik di field kedua tidak tertimpa saat field pertama diterapkan.
watch(
  () => props.min,
  (v) => {
    minLokal.value = formatAngka(angka(v))
    error.value = ''
  },
  { immediate: true },
)
watch(
  () => props.max,
  (v) => {
    maksLokal.value = formatAngka(angka(v))
    error.value = ''
  },
  { immediate: true },
)

const placeholderMin = computed(() => (props.batas ? formatAngka(String(Math.floor(props.batas.min))) : '0'))
const placeholderMaks = computed(() => (props.batas ? formatAngka(String(Math.ceil(props.batas.max))) : ''))

function terapkan() {
  const min = angka(minLokal.value)
  const max = angka(maksLokal.value)
  minLokal.value = formatAngka(min)
  maksLokal.value = formatAngka(max)
  if (min && max && Number(min) > Number(max)) {
    error.value = 'Harga minimum melebihi maksimum'
    return
  }
  error.value = ''
  if (min === angka(props.min) && max === angka(props.max)) return
  emit('ubah', { min: min || null, max: max || null })
}

// HP: 16px agar iOS tidak memperbesar halaman saat field disentuh; desktop: 14px agar angka jutaan muat.
const kelasInput =
  'h-11 w-full min-w-0 rounded-md border bg-surface pr-2 pl-8 text-base tabular-nums text-foreground transition-colors duration-150 placeholder:text-faint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:text-sm'
</script>

<template>
  <fieldset class="flex min-w-0 flex-col gap-2">
    <legend class="mb-1 text-sm font-bold">Harga</legend>
    <div class="grid grid-cols-2 gap-2">
      <div class="flex min-w-0 flex-col gap-1">
        <label :for="`${id}-min`" class="text-xs font-semibold text-muted">Minimum</label>
        <div class="relative">
          <span class="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-sm text-muted" aria-hidden="true">Rp</span>
          <input
            :id="`${id}-min`"
            v-model="minLokal"
            type="text"
            inputmode="numeric"
            autocomplete="off"
            enterkeyhint="done"
            :placeholder="placeholderMin"
            :aria-invalid="error ? 'true' : undefined"
            :aria-describedby="error ? `${id}-error` : undefined"
            :class="[kelasInput, error ? 'border-danger' : 'border-border-strong focus-visible:border-primary']"
            @keydown.enter.prevent="terapkan"
            @blur="terapkan"
          />
        </div>
      </div>
      <div class="flex min-w-0 flex-col gap-1">
        <label :for="`${id}-max`" class="text-xs font-semibold text-muted">Maksimum</label>
        <div class="relative">
          <span class="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-sm text-muted" aria-hidden="true">Rp</span>
          <input
            :id="`${id}-max`"
            v-model="maksLokal"
            type="text"
            inputmode="numeric"
            autocomplete="off"
            enterkeyhint="done"
            :placeholder="placeholderMaks"
            :aria-describedby="error ? `${id}-error` : undefined"
            :class="[kelasInput, 'border-border-strong focus-visible:border-primary']"
            @keydown.enter.prevent="terapkan"
            @blur="terapkan"
          />
        </div>
      </div>
    </div>
    <p v-if="error" :id="`${id}-error`" role="alert" class="text-sm font-medium text-danger-ink">{{ error }}</p>
  </fieldset>
</template>
