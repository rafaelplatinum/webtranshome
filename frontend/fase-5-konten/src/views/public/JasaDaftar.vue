<script setup>
import { Wrench } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import KartuJasa from '@/components/jasa/KartuJasa.vue'
import PageContactCta from '@/components/public/PageContactCta.vue'
import { useMuat } from '@/composables/useMuat'
import { remahSeo, urlSitus, useSeo } from '@/composables/useSeo'
import { listServices } from '@/services/kontenService'

// /jasa (KEPUTUSAN #20): daftar jasa aktif dari admin › Jasa. Konsultasi hanya lewat WhatsApp (KEPUTUSAN #12).
const CARA = [
  { judul: 'Pilih jasa', isi: 'Sesuai pekerjaan di rumahmu.' },
  { judul: 'Isi form singkat', isi: 'Nama, kebutuhan, dan lokasi.' },
  { judul: 'Atur jadwal dan biaya', isi: 'Tim kami balas lewat WhatsApp.' },
]
const TONE = ['sand', 'sage', 'blue', 'yellow', 'stone', 'peach']

const jasa = useMuat(() => listServices())
const REMAH = [{ label: 'Beranda', to: '/' }, { label: 'Jasa', to: '/jasa' }]

useSeo({
  title: 'Jasa pasang & renovasi rumah',
  description: 'Jasa Transhome di Sleman: pasang keramik, pengecatan, perlengkapan kamar mandi, AC & elektronik, service, pintu & jendela, rakit perabot. Konsultasi lewat WhatsApp.',
  canonical: '/jasa',
  jsonLd: () =>
    jasa.status.value === 'ready' && [
      remahSeo(REMAH),
      {
        '@type': 'ItemList',
        name: 'Jasa Transhome',
        itemListElement: jasa.data.value.map((j, i) => ({ '@type': 'ListItem', position: i + 1, url: urlSitus(`/jasa/${j.slug}`), name: j.name })),
      },
    ],
})
</script>

<template>
  <div class="flex flex-col gap-12 sm:gap-16">
    <section class="bg-th-cream pt-5">
      <div class="th-container">
        <Breadcrumb :items="REMAH.slice(0, 1).concat({ label: 'Jasa' })" />
      </div>
      <div class="mt-5 bg-th-yellow py-12 sm:py-16 lg:py-20">
        <div class="th-container grid items-center gap-8 lg:grid-cols-[1.15fr_0.85fr] lg:gap-14">
          <div class="flex flex-col gap-5">
            <span class="th-eyebrow text-th-navy">Jasa Transhome</span>
            <h1 class="max-w-3xl text-4xl leading-[1.08] font-extrabold tracking-tight text-th-navy sm:text-5xl lg:text-[52px]">Beli bahannya di sini, pasangnya juga bisa.</h1>
            <p class="max-w-2xl text-base leading-relaxed text-th-navy/80 sm:text-lg">Pilih jasa yang kamu butuhkan, isi form singkat, dan tim kami membalas lewat WhatsApp dengan biaya dan jadwal.</p>
          </div>
          <ol aria-label="Cara pesan jasa" class="grid gap-3">
            <li v-for="(c, i) in CARA" :key="c.judul" class="flex items-center gap-4 rounded-2xl border border-th-border bg-white p-4 sm:p-5">
              <span class="grid size-11 shrink-0 place-items-center rounded-full bg-th-orange text-xl font-extrabold text-white" aria-hidden="true">{{ i + 1 }}</span>
              <span class="flex flex-col gap-1"><strong>{{ c.judul }}</strong><span class="text-sm text-th-muted">{{ c.isi }}</span></span>
            </li>
          </ol>
        </div>
      </div>
    </section>

    <div v-if="jasa.status.value === 'loading'" class="th-container grid gap-5 sm:grid-cols-2 lg:grid-cols-3" aria-busy="true">
      <span class="sr-only">Memuat jasa…</span>
      <SkeletonCard v-for="n in 6" :key="n" />
    </div>
    <ErrorState v-else-if="jasa.status.value === 'error'" class="th-container" :message="jasa.pesan.value" @retry="jasa.muat" />
    <EmptyState v-else-if="!jasa.data.value.length" class="th-container" title="Daftar jasa sedang disiapkan" description="Tanyakan jasa yang kamu butuhkan langsung ke tim kami.">
      <template #icon><Wrench class="size-10" aria-hidden="true" /></template>
      <template #actions><WhatsAppButton :context="{ topic: 'Jasa', message: 'Saya ingin bertanya tentang jasa Transhome.' }" label="Tanya via WhatsApp" /></template>
    </EmptyState>
    <ul v-else class="th-container grid gap-5 sm:grid-cols-2 lg:grid-cols-3 lg:gap-6" data-daftar-jasa>
      <li v-for="(j, i) in jasa.data.value" :key="j.id" v-reveal="i % 3"><KartuJasa :jasa="j" tingkat="h2" :tone="TONE[i % TONE.length]" /></li>
    </ul>

    <PageContactCta v-if="jasa.status.value === 'ready' && jasa.data.value.length" tone="navy" heading="Jasa yang kamu cari belum ada?" description="Tanyakan saja, tim kami bantu carikan." button-label="Tanya via WhatsApp" :context="{ topic: 'Jasa', message: 'Saya mencari jasa yang belum ada di daftar.' }" />
  </div>
</template>
