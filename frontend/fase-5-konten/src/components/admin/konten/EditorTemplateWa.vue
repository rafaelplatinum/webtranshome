<script setup>
import { computed, nextTick, ref, useId } from 'vue'
import { RotateCcw } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Textarea from '@/components/ui/Textarea.vue'
import { BATAS_TEMPLATE, KODE_TEMPLATE, TEMPLATE_BAWAAN } from '@/utils/templateWhatsApp'

// Satu isian template WhatsApp (KEPUTUSAN #23): teks + tombol sisipkan {kode} di posisi kursor + kembalikan ke bawaan.
const model = defineModel({ type: String, default: '' })

const props = defineProps({
  // produk | keranjang | baris (lihat utils/templateWhatsApp.js)
  jenis: { type: String, required: true, validator: (v) => ['produk', 'keranjang', 'baris'].includes(v) },
  label: { type: String, required: true },
  rows: { type: Number, default: 6 },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
})

const id = useId()
const wadah = ref(null)
const kode = computed(() => KODE_TEMPLATE[props.jenis])
const bawaan = computed(() => TEMPLATE_BAWAAN[props.jenis])
// '{kode}' ditulis lewat fungsi: '}}' di dalam {{ }} akan menutup interpolasi Vue.
const tanda = (nama) => '{' + nama + '}'

async function sisipkan(nama) {
  const el = wadah.value?.querySelector('textarea')
  const teks = `{${nama}}`
  const nilai = model.value ?? ''
  const awal = el?.selectionStart ?? nilai.length
  const akhir = el?.selectionEnd ?? nilai.length
  model.value = nilai.slice(0, awal) + teks + nilai.slice(akhir)
  await nextTick()
  el?.focus()
  el?.setSelectionRange(awal + teks.length, awal + teks.length)
}
</script>

<template>
  <div class="flex flex-col gap-3" :data-editor-template="jenis">
    <div ref="wadah">
      <Textarea v-model="model" :label="label" :rows="rows" :error="error" :hint="hint" :maxlength="BATAS_TEMPLATE[jenis]" required spellcheck="false" class="font-mono" />
    </div>

    <div class="flex flex-col gap-1.5">
      <p :id="`${id}-kode`" class="text-sm font-semibold">Sisipkan kode di posisi kursor</p>
      <div class="flex flex-wrap gap-1.5" role="group" :aria-labelledby="`${id}-kode`">
        <button
          v-for="k in kode"
          :key="k.kode"
          type="button"
          class="inline-flex min-h-11 cursor-pointer items-center rounded-full border px-3 font-mono text-xs font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-8"
          :class="k.wajib ? 'border-primary bg-primary-soft text-primary-hover hover:bg-primary hover:text-primary-foreground' : 'border-border-strong bg-surface hover:border-foreground'"
          :title="k.arti"
          :aria-label="`Sisipkan {${k.kode}}${k.wajib ? ' (wajib)' : ''}: ${k.arti}`"
          :data-sisipkan-kode="k.kode"
          @click="sisipkan(k.kode)"
        >
          {{ tanda(k.kode) }}
        </button>
      </div>
      <p class="text-xs text-muted">Kode berwarna wajib ada. Baris yang semua kodenya kosong (mis. {pelanggan} saat pelanggan belum masuk) otomatis tidak ikut terkirim.</p>
    </div>

    <details class="rounded-lg border border-border text-sm">
      <summary class="flex min-h-11 cursor-pointer items-center rounded-lg px-3 font-semibold transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
        Arti setiap kode
      </summary>
      <dl class="flex flex-col gap-2 border-t border-border px-3 py-2.5">
        <div v-for="k in kode" :key="k.kode" class="flex flex-col gap-0.5 sm:flex-row sm:gap-3">
          <dt class="shrink-0 font-mono text-xs font-semibold sm:w-36">{{ tanda(k.kode) }}<span v-if="k.wajib" class="text-primary-hover"> · wajib</span></dt>
          <dd class="text-muted">{{ k.arti }}</dd>
        </div>
      </dl>
    </details>

    <Button variant="ghost" size="sm" class="self-start" :disabled="model === bawaan" @click="model = bawaan">
      <template #icon><RotateCcw class="size-4" aria-hidden="true" /></template>
      Kembalikan ke bawaan
    </Button>
  </div>
</template>
