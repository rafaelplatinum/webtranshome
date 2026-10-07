<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Newspaper, PenLine } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Table from '@/components/ui/Table.vue'
import ArtikelFilterBar from '@/components/admin/konten/ArtikelFilterBar.vue'
import StatusArtikel from '@/components/admin/konten/StatusArtikel.vue'
import { useFilter } from '@/composables/useFilter'
import { usePermission } from '@/composables/usePermission'
import { ingatDaftar } from '@/composables/useTujuanKembali'
import { listAdminArticles } from '@/services/adminKontenService'
import { pesanError } from '@/services/errors'
import { TIPE_ARTIKEL, tautanArtikel } from '@/utils/artikel'
import { formatTanggal } from '@/utils/format'
import { perilakuGulir } from '@/utils/gerak'

// Filter disimpan di URL (kembali dari editor → filter & halaman sama).
const KUNCI_FILTER = ['q', 'type', 'status']

const { nilai, terapkan, reset } = useFilter()
const { can } = usePermission()
ingatDaftar('artikel')

const status = ref('loading') // loading | ready | error
const artikel = ref([])
const meta = ref(null)
const pesan = ref('')
const atasTabel = ref(null)
let pengendali = null

const KOLOM = [
  { key: 'judul', label: 'Judul', class: 'min-w-72' },
  { key: 'tipe', label: 'Tipe' },
  { key: 'status', label: 'Status' },
  { key: 'diubah', label: 'Terakhir diubah', class: 'whitespace-nowrap' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

const params = computed(() => {
  const hasil = { per_page: 20, page: nilai.value.page || undefined }
  for (const k of KUNCI_FILTER) if (nilai.value[k]) hasil[k] = nilai.value[k]
  return hasil
})
const adaFilter = computed(() => KUNCI_FILTER.some((k) => nilai.value[k]))
const rentang = computed(() => {
  if (!meta.value?.total) return ''
  const awal = (meta.value.page - 1) * meta.value.per_page + 1
  return `Menampilkan ${awal}–${awal + artikel.value.length - 1} dari ${meta.value.total}`
})

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listAdminArticles(params.value, { signal: pengendali.signal })
    artikel.value = hasil.data
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
  atasTabel.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Artikel, promo, dan event yang terbit tampil di website. Draf hanya bisa dilihat lewat pratinjau.</p>
      <Button v-if="can('article.create')" to="/admin/artikel/baru">
        <template #icon><PenLine class="size-4.5" aria-hidden="true" /></template>
        Tulis artikel
      </Button>
    </div>

    <ArtikelFilterBar :nilai="nilai" @ubah="terapkan" @reset="reset()" />

    <section ref="atasTabel" aria-label="Daftar artikel" class="scroll-mt-24" :aria-busy="status === 'loading' ? 'true' : 'false'">
      <div v-if="status === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
        <span class="sr-only">Memuat artikel…</span>
        <Skeleton v-for="n in 6" :key="n" class="h-12" />
      </div>

      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

      <EmptyState
        v-else-if="!artikel.length"
        :title="adaFilter ? 'Tidak ada artikel yang cocok' : 'Belum ada artikel'"
        :description="adaFilter ? 'Coba ubah kata kunci atau filter.' : 'Artikel, promo, dan event yang ditulis akan muncul di sini.'"
      >
        <template #icon><Newspaper class="size-10" aria-hidden="true" /></template>
        <template #actions>
          <Button v-if="adaFilter" variant="outline" @click="reset()">Reset filter</Button>
          <Button v-else-if="can('article.create')" to="/admin/artikel/baru">Tulis artikel</Button>
        </template>
      </EmptyState>

      <div v-else class="overflow-hidden rounded-xl border border-border bg-surface">
        <Table :columns="KOLOM" :rows="artikel" caption="Daftar artikel" bare min-width="min-w-200">
          <template #cell-judul="{ row }">
            <span class="flex min-w-0 flex-col">
              <strong class="font-semibold">{{ row.title }}</strong>
              <span class="font-mono text-xs text-faint">/artikel/{{ row.slug }}</span>
            </span>
          </template>
          <template #cell-tipe="{ row }">
            <span class="rounded px-2 py-0.5 text-xs font-semibold whitespace-nowrap" :class="TIPE_ARTIKEL[row.type]?.kelas">{{ TIPE_ARTIKEL[row.type]?.label ?? row.type }}</span>
          </template>
          <template #cell-status="{ row }">
            <span class="flex flex-col items-start gap-1">
              <StatusArtikel :status="row.status" />
              <span v-if="row.published_at" class="text-xs whitespace-nowrap text-muted">{{ formatTanggal(row.published_at) }}</span>
            </span>
          </template>
          <template #cell-diubah="{ row }">
            <span class="flex flex-col">
              <span>{{ formatTanggal(row.updated_at) }}</span>
              <span v-if="row.author_name" class="text-xs text-muted">Penulis: {{ row.author_name }}</span>
            </span>
          </template>
          <template #cell-aksi="{ row }">
            <div class="flex items-center justify-end gap-1.5 whitespace-nowrap">
              <Button v-if="row.status === 'PUBLISHED'" variant="ghost" size="sm" :href="tautanArtikel(row.slug)">
                Lihat <span class="sr-only">{{ row.title }} (tab baru)</span>
              </Button>
              <Button v-if="can('article.update')" variant="outline" size="sm" :to="`/admin/artikel/${row.id}`">
                Ubah <span class="sr-only">{{ row.title }}</span>
              </Button>
            </div>
          </template>
        </Table>
        <div class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span class="text-sm text-muted" aria-live="polite">{{ rentang }}</span>
          <PaginationBar v-if="meta.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="gantiHalaman" />
        </div>
      </div>
    </section>
  </div>
</template>
