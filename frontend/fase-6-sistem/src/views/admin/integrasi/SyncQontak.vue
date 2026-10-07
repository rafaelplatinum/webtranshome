<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { RefreshCw, UserSearch } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import SkeletonTabel from '@/components/admin/SkeletonTabel.vue'
import KontakFilterBar from '@/components/admin/integrasi/KontakFilterBar.vue'
import KontakTabel from '@/components/admin/integrasi/KontakTabel.vue'
import LogSyncDrawer from '@/components/admin/integrasi/LogSyncDrawer.vue'
import { useDaftar } from '@/composables/useDaftar'
import { useFilter } from '@/composables/useFilter'
import { useSyncQontak } from '@/composables/useSyncQontak'
import { pesanError } from '@/services/errors'
import { listContacts, retryFailed } from '@/services/qontakService'
import { useUiStore } from '@/stores/ui'
import { perilakuGulir } from '@/utils/gerak'

const JEDA_POLLING = 5000

const { nilai, terapkan, reset } = useFilter()
const { cobaLagi } = useSyncQontak()
const ui = useUiStore()

const params = computed(() => {
  const hasil = { per_page: 20, page: nilai.value.page || undefined }
  if (nilai.value.status) hasil.status = nilai.value.status
  if (nilai.value.q) hasil.q = nilai.value.q
  return hasil
})
const { status, baris: kontak, meta, pesan, muat, rentang } = useDaftar(listContacts, params, { satuan: 'kontak' })

const atasTabel = ref(null)
const konfirmasi = ref(false)
const mengulang = ref(false)
const terpilih = ref(null)
const drawerBuka = ref(false)
const jumlahGagal = computed(() => meta.value?.counts?.FAILED ?? 0)
const adaFilter = computed(() => Boolean(nilai.value.status || nilai.value.q))

// Status berubah otomatis: selama ada kontak PENDING, daftar dimuat ulang diam-diam tiap 5 detik.
// Ditunda selama "Coba lagi" per baris berjalan (baris itu sedang dipantau sendiri) atau tab tidak terlihat.
let jeda = null
let ulangBerjalan = 0
function jadwalkan() {
  clearTimeout(jeda)
  if (!kontak.value.some((k) => k.sync_status === 'PENDING')) return
  jeda = setTimeout(() => {
    if (document.hidden || ulangBerjalan) jadwalkan()
    else muat({ diam: true })
  }, JEDA_POLLING)
}
watch(kontak, jadwalkan)
onBeforeUnmount(() => clearTimeout(jeda))

async function ulangSatu(row) {
  ulangBerjalan++
  try {
    await cobaLagi(row, row.name)
  } finally {
    ulangBerjalan--
  }
  // Angka per status ikut diperbarui; fokus kembali ke baris itu (tombol Coba lagi hilang bila berhasil).
  await muat({ diam: true })
  await nextTick()
  const tujuan = document.querySelector(`[data-retry-kontak="${row.id}"]`) ?? document.querySelector(`[data-log-kontak="${row.id}"]`)
  ;(tujuan ?? atasTabel.value)?.focus()
}

async function ulangSemua() {
  mengulang.value = true
  try {
    const jumlah = await retryFailed()
    konfirmasi.value = false
    ui.tampilkanToast({ pesan: jumlah ? `${jumlah} kontak dikirim ulang. Status diperbarui otomatis.` : 'Tidak ada kontak yang gagal.', jenis: jumlah ? 'success' : 'info' })
    await muat({ diam: true })
  } catch (error) {
    konfirmasi.value = false
    // 5xx/jaringan/403 sudah ditoast interceptor.
    const kode = error?.response?.status
    if (kode && kode < 500 && kode !== 401 && kode !== 403) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    mengulang.value = false
  }
}

function bukaLog(row) {
  terpilih.value = row
  drawerBuka.value = true
}

async function gantiHalaman(nomor) {
  await terapkan({ page: nomor > 1 ? nomor : null })
  atasTabel.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="max-w-2xl text-sm text-muted">
        Kontak member dikirim ke Mekari Qontak di background; pendaftaran tidak menunggu Qontak. Kontak yang gagal bisa dicoba lagi di sini.
      </p>
      <Button v-if="jumlahGagal" variant="secondary" :disabled="mengulang" @click="konfirmasi = true">
        <template #icon><RefreshCw class="size-4.5" aria-hidden="true" /></template>
        Coba lagi semua yang gagal ({{ jumlahGagal }})
      </Button>
    </div>

    <KontakFilterBar :nilai="nilai" :counts="meta?.counts ?? null" @ubah="terapkan" />

    <section ref="atasTabel" tabindex="-1" aria-label="Daftar kontak Qontak" class="scroll-mt-24 focus:outline-none" :aria-busy="status === 'loading' ? 'true' : 'false'">
      <SkeletonTabel v-if="status === 'loading'" label="Memuat status sync…" />

      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

      <EmptyState
        v-else-if="!kontak.length"
        :title="adaFilter ? 'Tidak ada kontak yang cocok' : 'Belum ada kontak'"
        :description="adaFilter ? 'Coba status lain atau periksa kata pencarian.' : 'Kontak muncul di sini setelah member mendaftar.'"
      >
        <template #icon><UserSearch class="size-10" aria-hidden="true" /></template>
        <template v-if="adaFilter" #actions>
          <Button variant="outline" @click="reset()">Tampilkan semua</Button>
        </template>
      </EmptyState>

      <div v-else class="overflow-hidden rounded-xl border border-border bg-surface">
        <KontakTabel :kontak="kontak" @coba-lagi="ulangSatu" @log="bukaLog" />
        <div class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span class="text-sm text-muted" aria-live="polite">{{ rentang }}</span>
          <PaginationBar v-if="meta?.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="gantiHalaman" />
        </div>
      </div>
    </section>

    <ConfirmModal
      v-model:open="konfirmasi"
      :title="`Coba lagi ${jumlahGagal} kontak yang gagal?`"
      message="Kontak dikirim ulang ke Qontak di background. Statusnya diperbarui otomatis di halaman ini."
      confirm-label="Ya, coba lagi"
      :loading="mengulang"
      @confirm="ulangSemua"
    />
    <LogSyncDrawer v-model:open="drawerBuka" :kontak="terpilih" />
  </div>
</template>
