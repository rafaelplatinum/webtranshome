<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Eye } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import ArtikelLainnya from '@/components/konten/ArtikelLainnya.vue'
import IsiArtikel from '@/components/konten/IsiArtikel.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { getArticle } from '@/services/kontenService'
import { pesanError } from '@/services/errors'
import { TIPE_ARTIKEL, daftarAsal } from '@/utils/artikel'
import { formatTanggal } from '@/utils/format'

const AJAKAN = {
  PROMO: { judul: 'Mau ambil promo ini?', isi: 'Tanyakan stok dan syarat promonya ke tim kami lewat WhatsApp.' },
  EVENT: { judul: 'Tertarik ikut acara ini?', isi: 'Tanyakan jadwal dan tempatnya ke tim kami lewat WhatsApp.' },
  LAINNYA: { judul: 'Masih ada yang ingin ditanyakan?', isi: 'Tim kami bantu pilih bahan yang pas untuk rumahmu.' },
}

const route = useRoute()
const { kontenLink } = useWhatsApp()

const status = ref('loading') // loading | ready | tidak-ada | error
const artikel = ref(null)
const pesan = ref('')
let pengendali = null

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    artikel.value = await getArticle(route.params.slug, route.query.preview, { signal: pengendali.signal })
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    // Draf tanpa token pratinjau yang benar juga 404, sama seperti artikel yang tidak ada.
    if (error?.response?.status === 404) {
      status.value = 'tidak-ada'
      return
    }
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(
  [() => route.params.slug, () => route.query.preview],
  ([slug]) => {
    if (slug) muat()
  },
  { immediate: true },
)
onBeforeUnmount(() => pengendali?.abort())

const tipe = computed(() => TIPE_ARTIKEL[artikel.value?.type] ?? { label: artikel.value?.type, kelas: 'bg-subtle' })
const ajakan = computed(() => AJAKAN[artikel.value?.type] ?? AJAKAN.LAINNYA)
const halamanTetap = computed(() => artikel.value?.type === 'HALAMAN')

const remah = computed(() => {
  const items = [{ label: 'Beranda', to: '/' }]
  if (artikel.value && !halamanTetap.value) items.push(daftarAsal(artikel.value.type))
  items.push({ label: artikel.value?.title ?? 'Artikel' })
  return items
})

useSeo({
  title: () => {
    if (status.value === 'tidak-ada') return 'Artikel tidak ditemukan'
    return status.value === 'ready' ? artikel.value.title : null
  },
  description: () => (status.value === 'ready' ? artikel.value.excerpt : null),
})
</script>

<template>
  <NotFoundState
    v-if="status === 'tidak-ada'"
    title="Artikel tidak ditemukan"
    description="Artikel yang kamu cari tidak ada atau belum diterbitkan."
    action-label="Lihat artikel lain"
    action-to="/artikel"
  />

  <div v-else class="mx-auto max-w-7xl px-4 pt-4 pb-12 sm:px-6 lg:pt-6 lg:pb-16">
    <div v-if="status === 'loading'" class="mx-auto flex max-w-3xl flex-col gap-4" aria-busy="true">
      <span class="sr-only">Memuat artikel…</span>
      <Skeleton class="h-5 w-56 max-w-full" />
      <Skeleton class="h-6 w-20" />
      <Skeleton class="h-10 w-4/5" />
      <Skeleton class="aspect-video w-full rounded-xl" />
      <Skeleton v-for="n in 4" :key="n" class="h-5 w-full" />
    </div>

    <ErrorState v-else-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />

    <template v-else-if="artikel">
      <article class="mx-auto flex max-w-3xl flex-col gap-6" data-artikel>
        <Breadcrumb :items="remah" />
        <p
          v-if="artikel.status === 'DRAFT'"
          role="status"
          class="flex items-start gap-2 rounded-md bg-warning-soft px-4 py-3 text-sm text-warning-ink"
          data-pratinjau
        >
          <Eye class="mt-0.5 size-4 shrink-0" aria-hidden="true" />
          <span><strong>Pratinjau draf.</strong> Halaman ini belum terbit, jadi hanya bisa dibuka lewat tautan pratinjau.</span>
        </p>
        <header class="flex flex-col gap-3">
          <span v-if="!halamanTetap" class="self-start rounded px-2 py-0.5 text-xs font-bold" :class="tipe.kelas">{{ tipe.label }}</span>
          <h1 class="text-3xl leading-tight font-extrabold sm:text-4xl">{{ artikel.title }}</h1>
          <p class="text-sm text-muted">
            <time v-if="artikel.published_at" :datetime="artikel.published_at">Terbit {{ formatTanggal(artikel.published_at) }}</time>
            <span v-else>Belum terbit</span>
          </p>
        </header>
        <img
          v-if="artikel.thumbnail_url"
          :src="artikel.thumbnail_url"
          alt=""
          width="1200"
          height="675"
          class="aspect-video w-full rounded-xl bg-subtle object-cover"
        />
        <IsiArtikel :isi="artikel.content" />
        <aside
          aria-labelledby="judul-tanya-konten"
          class="flex flex-col gap-4 rounded-xl border border-border bg-surface p-5 sm:flex-row sm:items-center sm:justify-between"
        >
          <div class="flex flex-col gap-1">
            <h2 id="judul-tanya-konten" class="text-lg font-bold">{{ ajakan.judul }}</h2>
            <p class="text-muted">{{ ajakan.isi }}</p>
          </div>
          <WhatsAppButton class="self-start sm:self-auto" :href="kontenLink(artikel)" />
        </aside>
      </article>
      <ArtikelLainnya v-if="!halamanTetap" :key="artikel.id" class="mt-12" :artikel="artikel" />
    </template>
  </div>
</template>
