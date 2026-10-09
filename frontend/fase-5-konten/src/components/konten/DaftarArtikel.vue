<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Newspaper } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { listArticles } from '@/services/kontenService'
import { pesanError } from '@/services/errors'
import KartuArtikel from './KartuArtikel.vue'

const props = defineProps({
  // Tipe untuk API, mis. 'ARTIKEL' atau 'PROMO,EVENT'.
  tipe: { type: String, required: true },
  // Judul daftar untuk pembaca layar (h2 di antara h1 halaman dan judul kartu h3).
  judul: { type: String, required: true },
  judulKosong: { type: String, default: 'Belum ada artikel' },
  deskripsiKosong: { type: String, default: '' },
})

const PER_HALAMAN = 9

const status = ref('loading') // loading | ready | error
const items = ref([])
const meta = ref(null)
const pesan = ref('')
const memuatLagi = ref(false)
let pengendali = null

const masihAda = computed(() => meta.value && meta.value.page < meta.value.total_pages)

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listArticles({ type: props.tipe, page: 1, per_page: PER_HALAMAN }, { signal: pengendali.signal })
    items.value = hasil.data
    meta.value = hasil.meta
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

// "Muat lebih banyak": halaman berikutnya ditambahkan ke daftar; tombol hilang bila habis.
async function muatBerikutnya() {
  memuatLagi.value = true
  try {
    const hasil = await listArticles({ type: props.tipe, page: meta.value.page + 1, per_page: PER_HALAMAN })
    items.value = [...items.value, ...hasil.data]
    meta.value = hasil.meta
  } catch (error) {
    pesan.value = pesanError(error)
  } finally {
    memuatLagi.value = false
  }
}

watch(() => props.tipe, muat, { immediate: true })
onBeforeUnmount(() => pengendali?.abort())
</script>

<template>
  <section aria-labelledby="judul-daftar-artikel" :aria-busy="status === 'loading' ? 'true' : 'false'">
    <h2 id="judul-daftar-artikel" class="sr-only">{{ judul }}</h2>
    <div v-if="status === 'loading'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <span class="sr-only">Memuat…</span>
      <Skeleton v-for="n in 6" :key="n" class="h-72 rounded-xl" />
    </div>
    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
    <EmptyState v-else-if="!items.length" :title="judulKosong" :description="deskripsiKosong">
      <template #icon><Newspaper class="size-10" aria-hidden="true" /></template>
    </EmptyState>
    <div v-else class="flex flex-col gap-6">
      <ul class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" data-daftar-artikel>
        <li v-for="(a, i) in items" :key="a.id" v-reveal="i % 3">
          <KartuArtikel :artikel="a" />
        </li>
      </ul>
      <div class="flex flex-col items-center gap-2">
        <p class="text-sm text-muted" aria-live="polite">Menampilkan {{ items.length }} dari {{ meta.total }}</p>
        <Button v-if="masihAda" variant="outline" :loading="memuatLagi" @click="muatBerikutnya">Muat lebih banyak</Button>
      </div>
    </div>
  </section>
</template>
