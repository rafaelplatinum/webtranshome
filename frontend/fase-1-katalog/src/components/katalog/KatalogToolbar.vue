<script setup>
import { useId } from 'vue'
import { ArrowUpDown, ChevronDown, LayoutGrid, List, SlidersHorizontal } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

const urutan = defineModel('urutan', { type: String, default: 'terbaru' })
const tampilan = defineModel('tampilan', { type: String, default: 'grid' })

defineProps({
  // meta.total; null selama muatan pertama atau bila muatan pertama gagal.
  total: { type: Number, default: null },
  memuat: { type: Boolean, default: false },
  jumlahFilter: { type: Number, default: 0 },
})

const emit = defineEmits(['buka-filter'])

const idUrutan = useId()

const OPSI_URUTAN = [
  { value: 'terbaru', label: 'Terbaru' },
  { value: 'harga_terendah', label: 'Harga terendah' },
  { value: 'harga_tertinggi', label: 'Harga tertinggi' },
  { value: 'nama', label: 'Nama A–Z' },
]

const TAMPILAN = [
  { value: 'grid', label: 'Grid', ikon: LayoutGrid },
  { value: 'daftar', label: 'Daftar', ikon: List },
]
</script>

<template>
  <!-- HP: [jumlah | grid/daftar] lalu [Filter | Urutkan]. Desktop: satu baris, filter ada di sidebar. -->
  <div class="grid grid-cols-[auto_minmax(0,1fr)] items-center gap-2 lg:flex lg:gap-3">
    <p class="text-sm text-muted lg:mr-auto" aria-live="polite">
      <template v-if="total != null"><span class="font-bold text-foreground tabular-nums">{{ total }}</span> produk</template>
      <template v-else-if="memuat">Memuat produk…</template>
    </p>

    <div role="group" aria-label="Tampilan produk" class="flex justify-self-end overflow-hidden rounded-md border border-border-strong lg:order-last">
      <button
        v-for="t in TAMPILAN"
        :key="t.value"
        type="button"
        :aria-pressed="tampilan === t.value ? 'true' : 'false'"
        class="inline-flex h-10.5 w-11 items-center justify-center transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset"
        :class="tampilan === t.value ? 'bg-secondary text-secondary-foreground' : 'bg-surface text-muted hover:bg-subtle hover:text-foreground'"
        @click="tampilan = t.value"
      >
        <component :is="t.ikon" class="size-5" aria-hidden="true" />
        <span class="sr-only">{{ t.label }}</span>
      </button>
    </div>

    <!-- Fokus dipasang dulu (Safari tidak memfokuskan tombol saat diklik) agar setelah sheet ditutup fokus kembali ke sini. -->
    <Button variant="outline" class="lg:hidden" aria-haspopup="dialog" @click="$event.currentTarget.focus(), emit('buka-filter')">
      <template #icon><SlidersHorizontal class="size-4.5" aria-hidden="true" /></template>
      Filter<template v-if="jumlahFilter"> ({{ jumlahFilter }})</template>
    </Button>

    <!-- Label "Urutkan" terlihat di desktop; di HP diganti ikon agar muat, tetap terbaca pembaca layar. -->
    <div class="flex min-w-0 items-center gap-2">
      <label :for="idUrutan" class="sr-only text-sm font-semibold whitespace-nowrap lg:not-sr-only">Urutkan</label>
      <div class="relative min-w-0 flex-1 lg:w-48 lg:flex-none">
        <ArrowUpDown class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted lg:hidden" aria-hidden="true" />
        <select
          :id="idUrutan"
          v-model="urutan"
          class="h-11 w-full min-w-0 cursor-pointer appearance-none rounded-md border border-border-strong bg-surface pr-9 pl-9 text-base text-foreground transition-colors duration-150 focus-visible:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:pl-3 lg:text-sm"
        >
          <option v-for="o in OPSI_URUTAN" :key="o.value" :value="o.value">{{ o.label }}</option>
        </select>
        <ChevronDown class="pointer-events-none absolute top-1/2 right-3 size-4 -translate-y-1/2 text-muted" aria-hidden="true" />
      </div>
    </div>
  </div>
</template>
