<script setup>
import { computed } from 'vue'

// Pratinjau pesan WhatsApp (gelembung chat keluar). Format WhatsApp ikut ditampilkan: *tebal*, _miring_, ~coret~.
const props = defineProps({
  pesan: { type: String, default: '' },
  label: { type: String, default: 'Pratinjau pesan' },
})

const POLA = /(?<![\p{L}\p{N}])([*_~])(?=\S)([^\n]*?\S)\1(?![\p{L}\p{N}])/gu
const GAYA = { '*': 'tebal', _: 'miring', '~': 'coret' }

// Dipecah jadi potongan teks biasa / berformat, tanpa v-html.
const potongan = computed(() => {
  const hasil = []
  let posisi = 0
  for (const m of props.pesan.matchAll(POLA)) {
    if (m.index > posisi) hasil.push({ teks: props.pesan.slice(posisi, m.index) })
    hasil.push({ teks: m[2], gaya: GAYA[m[1]] })
    posisi = m.index + m[0].length
  }
  if (posisi < props.pesan.length) hasil.push({ teks: props.pesan.slice(posisi) })
  return hasil
})
</script>

<template>
  <figure class="flex min-w-0 flex-col gap-2">
    <figcaption class="text-sm font-semibold">{{ label }}</figcaption>
    <div class="rounded-xl bg-subtle p-3 sm:p-4">
      <div v-if="pesan" class="ml-auto w-fit max-w-full rounded-lg rounded-tr-sm bg-success-soft px-3 py-2 text-sm leading-relaxed break-words whitespace-pre-wrap text-foreground" data-pratinjau-wa>
        <template v-for="(p, i) in potongan" :key="i">
          <strong v-if="p.gaya === 'tebal'">{{ p.teks }}</strong>
          <em v-else-if="p.gaya === 'miring'">{{ p.teks }}</em>
          <s v-else-if="p.gaya === 'coret'">{{ p.teks }}</s>
          <template v-else>{{ p.teks }}</template>
        </template>
      </div>
      <p v-else class="text-sm text-muted">Template masih kosong.</p>
      <p class="mt-1.5 text-right text-xs text-muted tabular-nums">{{ pesan.length }} karakter</p>
    </div>
  </figure>
</template>
