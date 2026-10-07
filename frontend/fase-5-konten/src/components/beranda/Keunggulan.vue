<script setup>
// "Kenapa Transhome" (teks dari sketsa "Landing3D").
const DAFTAR = [
  { nomor: '01', judul: 'Lengkap dalam satu atap', isi: 'Bahan struktur, finishing, sanitasi, sampai elektronik rumah.' },
  { nomor: '02', judul: 'Konsultasi via WhatsApp', isi: 'Kirim kebutuhan proyek, tim kami bantu cek stok dan jumlahnya.' },
  { nomor: '03', judul: 'Belanja jadi poin', isi: 'Setiap nota belanja di toko tercatat di akun Trans Family.' },
]
</script>

<template>
  <section aria-labelledby="judul-keunggulan" class="flex flex-col gap-4">
    <h2 id="judul-keunggulan" class="sr-only">Kenapa belanja di Transhome</h2>
    <ul class="grid gap-4 md:grid-cols-3">
      <li v-for="k in DAFTAR" :key="k.nomor" class="flex flex-col gap-2 rounded-xl border border-border bg-surface p-6">
        <span class="text-sm font-bold text-primary" aria-hidden="true">{{ k.nomor }}</span>
        <strong class="text-lg">{{ k.judul }}</strong>
        <span class="text-muted">{{ k.isi }}</span>
      </li>
    </ul>
  </section>
</template>
