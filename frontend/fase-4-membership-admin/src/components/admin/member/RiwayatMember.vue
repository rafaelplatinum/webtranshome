<script setup>
import { computed, onBeforeUnmount, ref, useId, watch } from 'vue'
import { History } from 'lucide-vue-next'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import Select from '@/components/ui/Select.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import RiwayatPoinAdmin from './RiwayatPoinAdmin.vue'
import { useFilter } from '@/composables/useFilter'
import { listPoints } from '@/services/adminMemberService'
import { pesanError } from '@/services/errors'
import { perilakuGulir } from '@/utils/gerak'

const props = defineProps({
  userId: { type: [Number, String], required: true },
})

// Riwayat poin di detail member. Filter jenis (?jenis=) dan halaman disimpan di URL.
const OPSI_JENIS = [
  { value: '', label: 'Semua jenis' },
  { value: 'EARN', label: 'Belanja' },
  { value: 'REDEEM', label: 'Tukar poin' },
  { value: 'ADJUST', label: 'Penyesuaian' },
]

const { nilai, terapkan } = useFilter()
const idJudul = useId()
const status = ref('loading') // loading | ready | error
const items = ref([])
const meta = ref(null)
const pesan = ref('')
const atas = ref(null)
let pengendali = null

const params = computed(() => ({ user_id: props.userId, type: nilai.value.jenis || undefined, page: nilai.value.page || undefined, per_page: 10 }))

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listPoints(params.value, { signal: pengendali.signal })
    items.value = hasil.data
    meta.value = hasil.meta
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(() => JSON.stringify(params.value), muat, { immediate: true })
onBeforeUnmount(() => pengendali?.abort())

async function gantiHalaman(nomor) {
  await terapkan({ page: nomor > 1 ? nomor : null })
  atas.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
}
</script>

<template>
  <section
    ref="atas"
    :aria-labelledby="idJudul"
    class="flex scroll-mt-24 flex-col gap-3 rounded-xl border border-border bg-surface p-5"
    :aria-busy="status === 'loading' ? 'true' : 'false'"
  >
    <div class="flex flex-wrap items-end justify-between gap-3">
      <h2 :id="idJudul" class="text-base font-bold">Riwayat poin</h2>
      <div class="w-full sm:w-52">
        <Select :model-value="nilai.jenis ?? ''" label="Jenis transaksi" hide-label :options="OPSI_JENIS" @update:model-value="(v) => terapkan({ jenis: v || null })" />
      </div>
    </div>

    <div v-if="status === 'loading'" class="flex flex-col gap-2">
      <span class="sr-only">Memuat riwayat poin…</span>
      <Skeleton v-for="n in 4" :key="n" class="h-10" />
    </div>
    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
    <EmptyState
      v-else-if="!items.length"
      :title="nilai.jenis ? 'Belum ada transaksi jenis ini' : 'Belum ada transaksi poin'"
      :description="nilai.jenis ? 'Pilih “Semua jenis” untuk melihat transaksi lain.' : 'Poin dari nota belanja yang diinput admin akan muncul di sini.'"
    >
      <template #icon><History class="size-10" aria-hidden="true" /></template>
    </EmptyState>
    <template v-else>
      <RiwayatPoinAdmin :items="items" />
      <div class="flex flex-wrap items-center justify-between gap-3">
        <span class="text-sm text-muted" aria-live="polite">{{ meta.total }} transaksi</span>
        <PaginationBar v-if="meta.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="gantiHalaman" />
      </div>
    </template>
  </section>
</template>
