<script setup>
import { computed, nextTick, ref, useId, watch } from 'vue'
import { FileText } from 'lucide-vue-next'
import { spesifikasiPublik } from '@/utils/produk'

const props = defineProps({
  produk: { type: Object, required: true },
})

const id = useId()
// Field internal (HPP, diskon, stok gudang, ...) dibuang di spesifikasiPublik, meski ikut terkirim API.
const spesifikasi = computed(() => spesifikasiPublik(props.produk.specifications))

const tab = computed(() => [
  { kunci: 'deskripsi', label: 'Deskripsi' },
  { kunci: 'spesifikasi', label: 'Spesifikasi' },
  ...(props.produk.datasheet_pdf_url ? [{ kunci: 'dokumen', label: 'Dokumen' }] : []),
])

const aktif = ref('spesifikasi')
const tombol = []

watch(tab, (daftar) => {
  if (!daftar.some((t) => t.kunci === aktif.value)) aktif.value = 'spesifikasi'
})

function pilih(indeks, fokus = false) {
  aktif.value = tab.value[indeks].kunci
  if (fokus) nextTick(() => tombol[indeks]?.focus())
}

// Pola tab WAI-ARIA: ← → pindah tab (berputar), Home/End ke ujung; tab langsung aktif saat difokus.
function onKeydown(event, indeks) {
  const jumlah = tab.value.length
  const tujuan = { ArrowRight: (indeks + 1) % jumlah, ArrowLeft: (indeks - 1 + jumlah) % jumlah, Home: 0, End: jumlah - 1 }[event.key]
  if (tujuan === undefined) return
  event.preventDefault()
  pilih(tujuan, true)
}
</script>

<template>
  <section class="flex flex-col">
    <div role="tablist" aria-label="Informasi produk" class="flex gap-1 overflow-x-auto border-b border-border scrollbar-none">
      <button
        v-for="(t, i) in tab"
        :id="`${id}-tab-${t.kunci}`"
        :key="t.kunci"
        :ref="(el) => (tombol[i] = el)"
        type="button"
        role="tab"
        :aria-selected="aktif === t.kunci ? 'true' : 'false'"
        :aria-controls="`${id}-panel-${t.kunci}`"
        :tabindex="aktif === t.kunci ? 0 : -1"
        class="-mb-px inline-flex h-12 items-center border-b-2 px-4 text-sm font-semibold whitespace-nowrap transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset"
        :class="aktif === t.kunci ? 'border-primary text-foreground' : 'border-transparent text-muted hover:text-foreground'"
        @click="pilih(i)"
        @keydown="onKeydown($event, i)"
      >
        {{ t.label }}
      </button>
    </div>

    <div
      v-for="t in tab"
      v-show="aktif === t.kunci"
      :id="`${id}-panel-${t.kunci}`"
      :key="t.kunci"
      role="tabpanel"
      :aria-labelledby="`${id}-tab-${t.kunci}`"
      tabindex="0"
      class="rounded-sm pt-5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
    >
      <template v-if="t.kunci === 'deskripsi'">
        <p v-if="produk.description" class="max-w-prose leading-relaxed whitespace-pre-line">{{ produk.description }}</p>
        <p v-else class="text-muted">Deskripsi belum tersedia.</p>
      </template>

      <template v-else-if="t.kunci === 'spesifikasi'">
        <dl v-if="spesifikasi.length" class="max-w-3xl divide-y divide-border overflow-hidden rounded-xl border border-border bg-surface">
          <div v-for="s in spesifikasi" :key="s.kunci" class="grid grid-cols-[minmax(0,2fr)_minmax(0,3fr)] gap-4 px-4 py-3 odd:bg-subtle sm:grid-cols-[14rem_minmax(0,1fr)]">
            <dt class="text-sm text-muted">{{ s.label }}</dt>
            <dd class="text-sm font-semibold break-words">{{ s.nilai }}</dd>
          </div>
        </dl>
        <p v-else class="text-muted">Spesifikasi belum tersedia.</p>
      </template>

      <ul v-else class="flex max-w-xl flex-col gap-2">
        <li>
          <a
            :href="produk.datasheet_pdf_url"
            target="_blank"
            rel="noopener"
            class="flex items-center gap-3 rounded-xl border border-border bg-surface p-4 transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <FileText class="size-6 shrink-0 text-primary" aria-hidden="true" />
            <span class="flex min-w-0 flex-1 flex-col">
              <span class="font-semibold">Datasheet produk</span>
              <span class="text-sm text-muted">PDF, dibuka di tab baru</span>
            </span>
            <span class="text-sm font-semibold text-primary">Unduh</span>
          </a>
        </li>
      </ul>
    </div>
  </section>
</template>
