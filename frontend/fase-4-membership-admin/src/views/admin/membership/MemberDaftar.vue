<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Coins, UserSearch } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import MemberFilterBar from '@/components/admin/member/MemberFilterBar.vue'
import MemberTabel from '@/components/admin/member/MemberTabel.vue'
import { useFilter } from '@/composables/useFilter'
import { usePermission } from '@/composables/usePermission'
import { ingatDaftar } from '@/composables/useTujuanKembali'
import { listMembers } from '@/services/adminMemberService'
import { pesanError } from '@/services/errors'
import { perilakuGulir } from '@/utils/gerak'

// Filter disimpan di URL (kembali dari detail member → filter & halaman sama).
const KUNCI_FILTER = ['q', 'consent', 'hp', 'sync', 'active']

const { nilai, terapkan, reset } = useFilter()
const { can } = usePermission()
ingatDaftar('member')

const status = ref('loading') // loading | ready | error
const member = ref([])
const meta = ref(null)
const pesan = ref('')
const atasTabel = ref(null)
let pengendali = null

const params = computed(() => {
  const hasil = { per_page: 20, page: nilai.value.page || undefined }
  for (const k of KUNCI_FILTER) if (nilai.value[k]) hasil[k] = nilai.value[k]
  return hasil
})
const adaFilter = computed(() => KUNCI_FILTER.some((k) => nilai.value[k]))
const rentang = computed(() => {
  if (!meta.value?.total) return ''
  const awal = (meta.value.page - 1) * meta.value.per_page + 1
  return `Menampilkan ${awal}–${awal + member.value.length - 1} dari ${meta.value.total} member`
})

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listMembers(params.value, { signal: pengendali.signal })
    member.value = hasil.data
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
      <p class="text-sm text-muted">Semua member Trans Family. Buka detail untuk melihat riwayat poin, status sync, dan mengelola akun.</p>
      <Button v-if="can('point.create')" to="/admin/poin/input">
        <template #icon><Coins class="size-4.5" aria-hidden="true" /></template>
        Input poin
      </Button>
    </div>

    <MemberFilterBar :nilai="nilai" @ubah="terapkan" @reset="reset()" />

    <section ref="atasTabel" aria-label="Daftar member" class="scroll-mt-24" :aria-busy="status === 'loading' ? 'true' : 'false'">
      <div v-if="status === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
        <span class="sr-only">Memuat member…</span>
        <div v-for="n in 8" :key="n" class="flex items-center gap-3">
          <Skeleton class="h-4 w-40" />
          <Skeleton class="h-4 flex-1" />
          <Skeleton class="h-4 w-24" />
          <Skeleton class="h-4 w-16" />
        </div>
      </div>

      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

      <EmptyState
        v-else-if="!member.length"
        :title="adaFilter ? 'Tidak ada member yang cocok' : 'Belum ada member'"
        :description="adaFilter ? 'Periksa lagi kode member (contoh TF-000123), email, atau nomor HP, atau kurangi filter.' : 'Member yang mendaftar di website akan muncul di sini.'"
      >
        <template #icon><UserSearch class="size-10" aria-hidden="true" /></template>
        <template v-if="adaFilter" #actions>
          <Button variant="outline" @click="reset()">Reset filter</Button>
        </template>
      </EmptyState>

      <div v-else class="overflow-hidden rounded-xl border border-border bg-surface">
        <MemberTabel :member="member" />
        <div class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span class="text-sm text-muted" aria-live="polite">{{ rentang }}</span>
          <PaginationBar v-if="meta.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="gantiHalaman" />
        </div>
      </div>
    </section>
  </div>
</template>
