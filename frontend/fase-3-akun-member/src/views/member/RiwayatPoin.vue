<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { History } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import DaftarPoin from '@/components/akun/DaftarPoin.vue'
import { pesanError } from '@/services/errors'
import { getPoints } from '@/services/memberService'

// Filter disimpan di URL (?jenis=didapat|ditukar) supaya tetap setelah refresh dan bisa dibagikan.
const FILTER = [
  { value: '', label: 'Semua', type: undefined, kosong: 'Belum ada transaksi poin.' },
  { value: 'didapat', label: 'Didapat', type: 'EARN', kosong: 'Belum ada poin yang didapat.' },
  { value: 'ditukar', label: 'Ditukar', type: 'REDEEM', kosong: 'Belum ada poin yang ditukar.' },
]
const PER_HALAMAN = 10

const route = useRoute()
const router = useRouter()
const status = ref('loading') // loading | ready | error
const items = ref([])
const meta = ref(null)
const pesan = ref('')
const memuatLagi = ref(false)
let pengendali = null

const filter = computed(() => FILTER.find((f) => f.value === route.query.jenis) ?? FILTER[0])
const masihAda = computed(() => meta.value && meta.value.page < meta.value.total_pages)
const angka = new Intl.NumberFormat('id-ID')

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await getPoints({ type: filter.value.type, page: 1, per_page: PER_HALAMAN }, { signal: pengendali.signal })
    items.value = hasil.data
    meta.value = hasil.meta
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

// "Muat lebih banyak": halaman berikutnya ditambahkan ke daftar; tombol hilang bila data habis.
async function muatBerikutnya() {
  memuatLagi.value = true
  try {
    const hasil = await getPoints({ type: filter.value.type, page: meta.value.page + 1, per_page: PER_HALAMAN })
    items.value = [...items.value, ...hasil.data]
    meta.value = hasil.meta
  } catch (error) {
    pesan.value = pesanError(error)
  } finally {
    memuatLagi.value = false
  }
}

function pilih(f) {
  router.replace({ query: f.value ? { jenis: f.value } : {} })
}

watch(() => filter.value.value, muat, { immediate: true })
onBeforeUnmount(() => pengendali?.abort())
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col gap-1">
      <h1 class="text-2xl font-extrabold sm:text-3xl">Riwayat poin</h1>
      <p v-if="meta" class="text-muted">Saldo sekarang <strong class="text-foreground tabular-nums">{{ angka.format(meta.balance) }} poin</strong></p>
    </div>

    <section aria-label="Riwayat transaksi poin" class="flex flex-col gap-4 rounded-xl border border-border bg-surface p-5 sm:p-6">
      <div role="group" aria-label="Filter jenis transaksi" class="flex flex-wrap gap-2">
        <button
          v-for="f in FILTER"
          :key="f.value"
          type="button"
          :aria-pressed="filter.value === f.value ? 'true' : 'false'"
          class="inline-flex h-11 items-center rounded-full border px-4 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
          :class="filter.value === f.value ? 'border-secondary bg-secondary text-secondary-foreground' : 'border-border-strong bg-surface hover:border-foreground'"
          @click="pilih(f)"
        >
          {{ f.label }}
        </button>
      </div>

      <div v-if="status === 'loading'" class="flex flex-col gap-2" aria-busy="true">
        <span class="sr-only">Memuat riwayat poin…</span>
        <Skeleton v-for="n in 5" :key="n" class="h-11" />
      </div>
      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
      <EmptyState v-else-if="!items.length" :title="filter.kosong" description="Poin bertambah setiap belanja di toko dengan menyebutkan kode member ke kasir.">
        <template #icon><History class="size-10" aria-hidden="true" /></template>
      </EmptyState>
      <template v-else>
        <DaftarPoin :items="items" />
        <div class="flex flex-col items-center gap-2">
          <p class="text-sm text-muted" aria-live="polite">Menampilkan {{ items.length }} dari {{ meta.total }} transaksi</p>
          <Button v-if="masihAda" variant="outline" :loading="memuatLagi" @click="muatBerikutnya">Muat lebih banyak</Button>
        </div>
      </template>
    </section>
  </div>
</template>
