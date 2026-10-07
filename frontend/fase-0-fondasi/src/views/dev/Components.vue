<script setup>
import DemoAkun from './bagian/DemoAkun.vue'
import DemoData from './bagian/DemoData.vue'
import DemoForm from './bagian/DemoForm.vue'
import DemoKatalog from './bagian/DemoKatalog.vue'
import DemoKeadaan from './bagian/DemoKeadaan.vue'
import DemoStatus from './bagian/DemoStatus.vue'
import DemoTombol from './bagian/DemoTombol.vue'
import DemoUmpanBalik from './bagian/DemoUmpanBalik.vue'

// Galeri komponen Fase 0 (hanya mode development). Dipakai untuk checklist QA di 375/768/1024/1440 px.
const BAGIAN = [
  { id: 'tombol', judul: 'Tombol', komponen: DemoTombol },
  { id: 'form', judul: 'Form', komponen: DemoForm },
  { id: 'status', judul: 'Badge, stok, harga', komponen: DemoStatus },
  { id: 'umpan-balik', judul: 'Modal dan toast', komponen: DemoUmpanBalik },
  { id: 'keadaan', judul: 'Kosong, error, loading', komponen: DemoKeadaan },
  { id: 'data', judul: 'Tabel dan paginasi', komponen: DemoData },
  { id: 'katalog', judul: 'Kartu produk', komponen: DemoKatalog },
  { id: 'akun', judul: 'Akun contoh dan route guard', komponen: DemoAkun },
]
</script>

<template>
  <div class="mx-auto flex max-w-7xl flex-col gap-10 px-4 py-10 sm:px-6">
    <div class="flex flex-col gap-3">
      <h1 class="text-3xl font-extrabold">Galeri komponen</h1>
      <p class="max-w-2xl text-muted">
        Semua komponen dasar Fase 0 dalam setiap state-nya. Token warna dan font mengikuti
        <code class="font-mono text-sm">design-system/transhome/MASTER.md</code>.
      </p>
      <nav aria-label="Bagian galeri" class="flex flex-wrap gap-2">
        <a
          v-for="b in BAGIAN"
          :key="b.id"
          :href="`#${b.id}`"
          class="inline-flex h-10 items-center rounded-full bg-subtle px-4 text-sm font-semibold transition-colors duration-150 hover:bg-border focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          {{ b.judul }}
        </a>
      </nav>
    </div>

    <section
      v-for="b in BAGIAN"
      :id="b.id"
      :key="b.id"
      :aria-labelledby="`judul-${b.id}`"
      class="flex scroll-mt-6 flex-col gap-5 rounded-2xl border border-border bg-surface p-5 sm:p-7"
    >
      <h2 :id="`judul-${b.id}`" class="text-xl font-bold">{{ b.judul }}</h2>
      <component :is="b.komponen" />
    </section>
  </div>
</template>
