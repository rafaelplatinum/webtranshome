<script setup>
import { nextTick, ref } from 'vue'
import { X } from 'lucide-vue-next'

defineProps({
  // [{ id, label, hapus }] dari useKatalog().chip
  chip: { type: Array, required: true },
})

// `kosong`: chip terakhir dihapus; induk memindahkan fokus agar tidak hilang ke <body>.
const emit = defineEmits(['hapus-semua', 'kosong'])

const wadah = ref(null)

async function hapus(item, indeks) {
  await item.hapus()
  await nextTick()
  const sisa = wadah.value?.querySelectorAll('[data-chip]') ?? []
  const tujuan = sisa[indeks] ?? sisa[indeks - 1]
  if (tujuan) tujuan.focus()
  else emit('kosong')
}
</script>

<template>
  <div v-if="chip.length" ref="wadah" class="flex flex-wrap items-center gap-2">
    <span class="text-sm font-semibold text-muted">Filter aktif:</span>
    <button
      v-for="(item, i) in chip"
      :key="item.id"
      type="button"
      data-chip
      class="inline-flex h-11 max-w-full items-center gap-1.5 rounded-full border border-primary bg-primary-soft pr-2.5 pl-3.5 text-sm font-semibold text-primary-hover transition-colors duration-150 hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 md:h-9"
      :aria-label="`Hapus filter ${item.label}`"
      @click="hapus(item, i)"
    >
      <span class="truncate">{{ item.label }}</span>
      <X class="size-4 shrink-0" aria-hidden="true" />
    </button>
    <button
      type="button"
      class="inline-flex h-11 items-center rounded-sm px-1 text-sm font-semibold text-foreground underline underline-offset-2 transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
      @click="emit('hapus-semua')"
    >
      Hapus semua
    </button>
  </div>
</template>
