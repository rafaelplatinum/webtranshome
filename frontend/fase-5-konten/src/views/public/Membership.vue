<script setup>
import { computed } from 'vue'
import { Check, CircleHelp, Coins, CreditCard, UserRound } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import TombolGoogle from '@/components/akun/TombolGoogle.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useSettingsStore } from '@/stores/settings'
import { formatRupiah } from '@/utils/format'

const settings = useSettingsStore()
const { generalLink } = useWhatsApp()
const rasio = computed(() => {
  const nilai = settings.data?.point_ratio_rupiah
  return nilai ? `Rasio poin saat ini: 1 poin tiap belanja ${formatRupiah(nilai)} di toko.` : 'Poin tercatat berdasarkan total belanja di toko.'
})
const LANGKAH = [
  { angka: '01', judul: 'Daftar gratis', isi: 'Pakai email atau akun Google. Tidak ada biaya.' },
  { angka: '02', judul: 'Dapat kode member', isi: 'Kode digital langsung tersedia di akunmu.' },
  { angka: '03', judul: 'Sebutkan ke kasir', isi: 'Setiap nota belanja di toko tercatat sebagai poin.' },
]
const MANFAAT = computed(() => [
  { ikon: Coins, judul: 'Poin tiap belanja', isi: rasio.value },
  { ikon: CreditCard, judul: 'Kode digital', isi: 'Tidak perlu kartu fisik. Cukup sebutkan kode member di kasir.' },
  { ikon: UserRound, judul: 'Pantau di akunmu', isi: 'Lihat jumlah poin dan riwayat nota belanja kapan saja.' },
])
const FAQ = [
  { tanya: 'Apakah membership ini gratis?', jawab: 'Ya, pendaftaran Trans Family gratis.' },
  { tanya: 'Bagaimana cara mendaftar?', jawab: 'Daftar memakai email atau akun Google, lalu kode member digital muncul di akunmu.' },
  { tanya: 'Bagaimana poin belanja tercatat?', jawab: 'Sebutkan kode member ke kasir saat membayar agar poin dari nota tercatat di akunmu.' },
]

useSeo({ title: 'Trans Family', description: 'Kenali program Trans Family dan daftar member gratis untuk mengumpulkan poin belanja di Transhome.' })
</script>

<template>
  <main class="flex flex-col gap-14 pb-16 sm:gap-20 sm:pb-20">
    <section class="overflow-hidden bg-th-navy text-white">
      <div class="th-container grid min-h-[32rem] items-center gap-10 py-12 sm:py-16 lg:grid-cols-[1.1fr_0.9fr] lg:gap-16 lg:py-20">
        <div class="flex flex-col items-start gap-6">
          <span class="th-eyebrow text-orange-200">Trans Family · Gratis</span>
          <h1 class="max-w-2xl text-4xl leading-[1.08] font-extrabold tracking-tight sm:text-5xl lg:text-[56px]">Belanja di toko, kumpulkan poin.</h1>
          <p class="max-w-xl text-base leading-relaxed text-white/80 sm:text-lg">Daftar pakai email, dapatkan kode member digital, lalu sebutkan kodenya ke kasir. Poin belanja langsung tercatat.</p>
          <div class="flex w-full flex-col items-start gap-3 sm:w-auto">
            <Button to="/daftar" size="lg">Daftar sekarang</Button>
          </div>
        </div>
        <div class="relative mx-auto flex h-64 w-full max-w-lg items-center justify-center sm:h-80" aria-hidden="true">
          <div class="absolute h-[72%] w-[90%] rotate-[-8deg] rounded-[24px] bg-th-yellow"></div>
          <div class="relative flex h-[72%] w-[90%] rotate-[3deg] flex-col justify-between rounded-[24px] border border-white/10 bg-secondary-hover p-6 shadow-2xl sm:p-8">
            <div class="flex items-center justify-between"><strong class="text-lg">Trans Family</strong><span class="rounded-full bg-th-yellow px-3 py-1 text-[13px] font-extrabold text-th-navy">Member</span></div>
            <div class="flex items-end justify-between gap-4"><span class="h-3 w-2/3 rounded bg-white/15"></span><span class="h-3 w-1/4 rounded bg-white/10"></span></div>
            <div class="flex items-center gap-3"><span class="grid size-10 place-items-center rounded-full bg-th-yellow text-th-navy"><Check class="size-5" /></span><span class="font-semibold text-white/80">Kode member digital</span></div>
          </div>
        </div>
      </div>
    </section>

    <section class="th-container flex flex-col gap-8" aria-labelledby="judul-cara-kerja">
      <div class="max-w-2xl"><span class="th-eyebrow">Cara kerja</span><h2 id="judul-cara-kerja" class="mt-2 text-3xl leading-tight font-extrabold tracking-tight sm:text-4xl">Tiga langkah, selesai di kasir.</h2></div>
      <ol class="grid gap-4 md:grid-cols-3 lg:gap-6">
        <li v-for="langkah in LANGKAH" :key="langkah.angka" class="th-card flex min-h-44 gap-4 p-5 sm:p-7">
          <strong class="text-3xl leading-none font-extrabold text-th-orange">{{ langkah.angka }}</strong>
          <div><h3 class="text-lg font-extrabold">{{ langkah.judul }}</h3><p class="mt-2 text-base leading-relaxed text-th-muted">{{ langkah.isi }}</p></div>
        </li>
      </ol>
    </section>

    <section class="th-container grid gap-4 md:grid-cols-3 lg:gap-6" aria-label="Manfaat Trans Family">
      <article v-for="(manfaat, i) in MANFAAT" :key="manfaat.judul" class="flex min-h-56 flex-col justify-end gap-4 rounded-[24px] p-6 sm:p-8" :class="['bg-th-yellow', 'bg-th-blue', 'bg-th-peach'][i]">
        <component :is="manfaat.ikon" class="size-7 text-th-navy" aria-hidden="true" />
        <div><h3 class="text-xl font-extrabold tracking-tight">{{ manfaat.judul }}</h3><p class="mt-2 text-base leading-relaxed text-th-navy/80">{{ manfaat.isi }}</p></div>
      </article>
    </section>

    <section id="daftar-member" class="th-container scroll-mt-24" aria-labelledby="judul-daftar-member">
      <div class="rounded-[24px] bg-th-peach px-6 py-8 sm:px-10 sm:py-10 lg:px-12 lg:py-12">
        <div class="max-w-4xl">
          <h2 id="judul-daftar-member" class="text-3xl leading-tight font-extrabold tracking-tight text-th-navy sm:text-4xl">Jadi member Trans Family, gratis.</h2>
          <p class="mt-3 max-w-2xl text-base leading-relaxed text-th-orange sm:text-lg">Daftar pakai email, dapat kode member digital. Sebutkan kodenya ke kasir, poin belanja langsung tercatat.</p>
          <div class="mt-5 inline-flex w-full flex-col gap-3 sm:w-auto sm:flex-row">
            <Button to="/daftar" size="lg" class="w-full sm:min-w-52 sm:w-auto">Daftar sekarang</Button>
            <TombolGoogle label="Masuk dengan Google" class="!h-14 !w-full !px-6 !text-base sm:!w-[19.25rem]" />
          </div>
        </div>
      </div>
    </section>

    <section class="th-container grid gap-8 lg:grid-cols-[0.8fr_1.2fr] lg:gap-16" aria-labelledby="judul-faq-member">
      <div class="flex flex-col items-start gap-4"><span class="th-eyebrow">Pertanyaan umum</span><h2 id="judul-faq-member" class="text-3xl leading-tight font-extrabold tracking-tight sm:text-4xl">Masih ada yang ingin ditanyakan?</h2><WhatsAppButton :href="generalLink()" class="mt-2" label="Tanya via WhatsApp" /></div>
      <div class="flex flex-col gap-3"><details v-for="item in FAQ" :key="item.tanya" class="th-card group p-5 sm:p-6"><summary class="flex cursor-pointer list-none items-center justify-between gap-4 text-base font-extrabold"><span>{{ item.tanya }}</span><CircleHelp class="size-5 shrink-0 text-th-orange" aria-hidden="true" /></summary><p class="mt-3 text-base leading-relaxed text-th-muted">{{ item.jawab }}</p></details></div>
    </section>
  </main>
</template>
