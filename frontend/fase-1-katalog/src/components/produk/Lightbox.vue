<script setup>
import { computed, onBeforeUnmount, watch } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import Modal from '@/components/ui/Modal.vue'

// Foto produk ukuran besar. Fokus terkunci & Esc dari Modal; ← → ganti foto.
const open = defineModel('open', { type: Boolean, default: false })
const index = defineModel('index', { type: Number, default: 0 })

const props = defineProps({
  // [{ id, image_url }] urutan sama dengan galeri
  foto: { type: Array, required: true },
  nama: { type: String, required: true },
})

const total = computed(() => props.foto.length)
const sekarang = computed(() => props.foto[index.value] ?? props.foto[0])

function geser(arah) {
  index.value = (index.value + arah + total.value) % total.value
}

function onKeydown(event) {
  if (total.value < 2) return
  if (event.key === 'ArrowLeft') {
    event.preventDefault()
    geser(-1)
  } else if (event.key === 'ArrowRight') {
    event.preventDefault()
    geser(1)
  }
}

watch(
  open,
  (buka) => {
    if (buka) document.addEventListener('keydown', onKeydown)
    else document.removeEventListener('keydown', onKeydown)
  },
  { immediate: true },
)
onBeforeUnmount(() => document.removeEventListener('keydown', onKeydown))

const kelasTombol =
  'absolute top-1/2 inline-flex size-11 -translate-y-1/2 items-center justify-center rounded-full border border-border-strong bg-surface shadow-sm transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <Modal v-model:open="open" :title="nama" :description="`Foto ${index + 1} dari ${total}`" size="xl">
    <div class="relative flex items-center justify-center rounded-lg bg-subtle">
      <img
        v-if="sekarang"
        :src="sekarang.image_url"
        :alt="`${nama}, foto ${index + 1} dari ${total}`"
        width="800"
        height="800"
        class="max-h-[65dvh] w-auto max-w-full object-contain"
      />
      <template v-if="total > 1">
        <button type="button" :class="[kelasTombol, 'left-2']" aria-label="Foto sebelumnya" @click="geser(-1)">
          <ChevronLeft class="size-5" aria-hidden="true" />
        </button>
        <button type="button" :class="[kelasTombol, 'right-2']" aria-label="Foto berikutnya" @click="geser(1)">
          <ChevronRight class="size-5" aria-hidden="true" />
        </button>
      </template>
    </div>
    <p class="sr-only" aria-live="polite">Foto {{ index + 1 }} dari {{ total }}</p>
  </Modal>
</template>
