<script setup>
defineProps({
  // [{ key, label, align: 'left' | 'right' | 'center', class }]
  columns: { type: Array, required: true },
  rows: { type: Array, default: () => [] },
  rowKey: { type: String, default: 'id' },
  // Judul tabel untuk pembaca layar
  caption: { type: String, default: '' },
  // true: tanpa bingkai sendiri, untuk dipasang di dalam kartu yang punya footer (mis. paginasi).
  bare: { type: Boolean, default: false },
  // Lebar minimum tabel sebelum bergulir di wadahnya, mis. 'min-w-240' untuk tabel produk admin.
  minWidth: { type: String, default: 'min-w-160' },
})

const RATA = { left: 'text-left', right: 'text-right', center: 'text-center' }
</script>

<template>
  <!-- Tabel lebar bergulir di wadahnya sendiri, bukan halaman. Isi sel kustom: slot #cell-<key>.
       `relative`: teks sr-only (position absolute) ikut terpotong wadah ini, tidak memperlebar halaman. -->
  <div class="relative overflow-x-auto" :class="!bare && 'rounded-xl border border-border bg-surface'">
    <table class="w-full border-collapse text-sm" :class="minWidth">
      <caption v-if="caption" class="sr-only">{{ caption }}</caption>
      <thead class="bg-background text-xs font-semibold tracking-wide text-muted uppercase">
        <tr>
          <th v-for="kolom in columns" :key="kolom.key" scope="col" class="px-4 py-3" :class="[RATA[kolom.align ?? 'left'], kolom.class]">
            <span :class="kolom.srOnly && 'sr-only'">{{ kolom.label }}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="baris in rows" :key="baris[rowKey]" class="border-t border-border transition-colors duration-150 hover:bg-background">
          <td v-for="kolom in columns" :key="kolom.key" class="px-4 py-3" :class="[RATA[kolom.align ?? 'left'], kolom.class]">
            <slot :name="`cell-${kolom.key}`" :row="baris" :value="baris[kolom.key]">{{ baris[kolom.key] }}</slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
