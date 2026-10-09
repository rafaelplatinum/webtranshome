<script setup>
import { computed } from 'vue'
import { blokIsi } from '@/utils/artikel'

// Isi artikel dirender per blok (paragraf, subjudul, daftar) tanpa v-html, jadi teks dari CMS aman ditampilkan.
const props = defineProps({
  isi: { type: String, required: true },
})

const blok = computed(() => blokIsi(props.isi))
</script>

<template>
  <div class="flex flex-col gap-4 text-lg leading-relaxed" data-isi-artikel>
    <template v-for="(b, i) in blok" :key="i">
      <h2 v-if="b.jenis === 'h2'" class="mt-2 text-xl font-extrabold">{{ b.teks }}</h2>
      <ul v-else-if="b.jenis === 'ul'" class="flex list-disc flex-col gap-1 pl-6 marker:text-primary">
        <li v-for="(butir, j) in b.butir" :key="j">{{ butir }}</li>
      </ul>
      <p v-else>{{ b.teks }}</p>
    </template>
  </div>
</template>
