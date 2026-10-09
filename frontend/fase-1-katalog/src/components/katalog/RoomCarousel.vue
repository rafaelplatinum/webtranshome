<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import RoomCard from './RoomCard.vue'
import { perilakuGulir } from '@/utils/gerak'

const props = defineProps({
  ruangan: { type: Array, required: true },
  label: { type: String, default: 'Daftar ruangan' },
})

// Tidak berputar otomatis. Tombol ‹ › hanya muncul bila kartu tidak muat satu baris,
// dan nonaktif di ujung. Geser memakai animasi halus kecuali pengguna memilih reduced motion.
const jalur = ref(null)
const meluap = ref(false)
const bisaKiri = ref(false)
const bisaKanan = ref(false)
let pengamat = null

function perbarui() {
  const el = jalur.value
  if (!el) return
  meluap.value = el.scrollWidth > el.clientWidth + 1
  bisaKiri.value = el.scrollLeft > 1
  bisaKanan.value = el.scrollLeft + el.clientWidth < el.scrollWidth - 1
}

function geser(arah) {
  const el = jalur.value
  const kartu = el?.querySelector('li')
  if (!el || !kartu) return
  el.scrollBy({ left: arah * (kartu.getBoundingClientRect().width + 16), behavior: perilakuGulir() })
}

onMounted(() => {
  perbarui()
  pengamat = new ResizeObserver(perbarui)
  pengamat.observe(jalur.value)
})
onBeforeUnmount(() => pengamat?.disconnect())
watch(() => props.ruangan, () => nextTick(perbarui))

const kelasTombol =
  'inline-flex size-11 items-center justify-center rounded-full border border-border-strong bg-surface transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:border-border-strong'
</script>

<template>
  <div class="flex flex-col gap-3">
    <div class="flex items-end justify-between gap-4">
      <slot name="judul" />
      <div v-if="meluap" class="ml-auto flex shrink-0 gap-2">
        <button type="button" :class="kelasTombol" :disabled="!bisaKiri" aria-label="Geser ke ruangan sebelumnya" @click="geser(-1)">
          <ChevronLeft class="size-5" aria-hidden="true" />
        </button>
        <button type="button" :class="kelasTombol" :disabled="!bisaKanan" aria-label="Geser ke ruangan berikutnya" @click="geser(1)">
          <ChevronRight class="size-5" aria-hidden="true" />
        </button>
      </div>
    </div>
    <ul
      ref="jalur"
      :aria-label="label"
      class="-m-1 grid snap-x snap-mandatory scroll-px-1 auto-cols-[minmax(14rem,1fr)] grid-flow-col gap-4 overflow-x-auto overscroll-x-contain p-1 scrollbar-none"
      @scroll.passive="perbarui"
    >
      <li v-for="r in ruangan" :key="r.id" class="snap-start">
        <RoomCard :ruangan="r" />
      </li>
    </ul>
  </div>
</template>
