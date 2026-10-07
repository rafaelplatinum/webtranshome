<script setup>
import { computed, onMounted, ref } from 'vue'
import { ScrollText } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import SkeletonTabel from '@/components/admin/SkeletonTabel.vue'
import DetailLog from '@/components/admin/audit/DetailLog.vue'
import LogFilterBar from '@/components/admin/audit/LogFilterBar.vue'
import LogTabel from '@/components/admin/audit/LogTabel.vue'
import { useDaftar } from '@/composables/useDaftar'
import { useFilter } from '@/composables/useFilter'
import { getActivityLogOptions, listActivityLogs } from '@/services/adminLogService'
import { perilakuGulir } from '@/utils/gerak'

// Filter disimpan di URL: bisa dibagikan dan tetap ada setelah refresh.
const KUNCI_FILTER = ['admin_id', 'module', 'from', 'to']

const { nilai, terapkan, reset } = useFilter()
const opsi = ref({ admins: [], modules: [] })
const opsiGagal = ref(false)
const atasTabel = ref(null)
const terpilih = ref(null)
const drawerBuka = ref(false)

const errorTanggal = computed(() => (nilai.value.from && nilai.value.to && nilai.value.to < nilai.value.from ? 'Tanggal akhir tidak boleh sebelum tanggal awal' : ''))
const params = computed(() => {
  if (errorTanggal.value) return null
  const hasil = { per_page: 20, page: nilai.value.page || undefined }
  for (const k of KUNCI_FILTER) if (nilai.value[k]) hasil[k] = nilai.value[k]
  return hasil
})
const adaFilter = computed(() => KUNCI_FILTER.some((k) => nilai.value[k]))

const { status, baris: logs, meta, pesan, muat, rentang } = useDaftar(listActivityLogs, params, { satuan: 'aktivitas' })

// Pilihan admin & modul hanya pelengkap: bila gagal dimuat, daftar log tetap bisa dipakai.
onMounted(async () => {
  try {
    opsi.value = await getActivityLogOptions()
  } catch {
    opsiGagal.value = true
  }
})

function lihat(log) {
  terpilih.value = log
  drawerBuka.value = true
}

async function gantiHalaman(nomor) {
  await terapkan({ page: nomor > 1 ? nomor : null })
  atasTabel.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <p class="text-sm text-muted">Semua perubahan yang dilakukan admin, terbaru di atas. Buka detail untuk melihat data sebelum dan sesudah.</p>

    <LogFilterBar :nilai="nilai" :opsi="opsi" :error-tanggal="errorTanggal" @ubah="terapkan" @reset="reset()" />
    <p v-if="opsiGagal" class="text-sm text-muted">Pilihan admin dan modul gagal dimuat. Filter tanggal tetap bisa dipakai.</p>

    <section ref="atasTabel" aria-label="Daftar log aktivitas" class="scroll-mt-24" :aria-busy="status === 'loading' ? 'true' : 'false'">
      <SkeletonTabel v-if="status === 'loading'" label="Memuat log aktivitas…" />

      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

      <EmptyState
        v-else-if="!logs.length"
        :title="adaFilter ? 'Tidak ada aktivitas yang cocok' : 'Belum ada aktivitas'"
        :description="adaFilter ? 'Coba ganti admin, modul, atau rentang tanggal.' : 'Perubahan dari panel admin akan tercatat di sini.'"
      >
        <template #icon><ScrollText class="size-10" aria-hidden="true" /></template>
        <template v-if="adaFilter" #actions>
          <Button variant="outline" @click="reset()">Reset filter</Button>
        </template>
      </EmptyState>

      <div v-else class="overflow-hidden rounded-xl border border-border bg-surface">
        <LogTabel :logs="logs" @detail="lihat" />
        <div class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span class="text-sm text-muted" aria-live="polite">{{ rentang }}</span>
          <PaginationBar v-if="meta?.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="gantiHalaman" />
        </div>
      </div>
    </section>

    <DetailLog v-model:open="drawerBuka" :log="terpilih" />
  </div>
</template>
