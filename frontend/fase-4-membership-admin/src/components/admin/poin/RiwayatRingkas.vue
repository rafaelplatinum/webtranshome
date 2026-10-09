<script setup>
import { onBeforeUnmount, ref, useId, watch } from 'vue'
import { RouterLink } from 'vue-router'
import Skeleton from '@/components/ui/Skeleton.vue'
import { usePermission } from '@/composables/usePermission'
import { listPoints } from '@/services/adminMemberService'
import { pesanError } from '@/services/errors'
import { formatRupiah, formatTanggal } from '@/utils/format'
import { JENIS_TRANSAKSI, tandaPoin } from '@/utils/poin'

const props = defineProps({
  userId: { type: Number, required: true },
})

const JUMLAH = 5

const { can } = usePermission()
const idJudul = useId()
const status = ref('loading') // loading | ready | error
const items = ref([])
const pesan = ref('')
let pengendali = null

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listPoints({ user_id: props.userId, per_page: JUMLAH }, { signal: pengendali.signal })
    items.value = hasil.data
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(() => props.userId, muat, { immediate: true })
onBeforeUnmount(() => pengendali?.abort())

/** Transaksi yang baru disimpan langsung tampil paling atas tanpa memuat ulang. */
function tambah(transaksi) {
  items.value = [transaksi, ...items.value].slice(0, JUMLAH)
}

// Baris kedua: "Rp 12.800.000 · 5 Okt 2026 · oleh Admin" (catatan penukaran sudah di baris pertama).
const keterangan = (t) =>
  [t.purchase_amount && formatRupiah(t.purchase_amount), formatTanggal(t.created_at), t.created_by_name && `oleh ${t.created_by_name}`].filter(Boolean).join(' · ')

defineExpose({ tambah })
</script>

<template>
  <section :aria-labelledby="idJudul" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-5" :aria-busy="status === 'loading' ? 'true' : 'false'">
    <h2 :id="idJudul" class="text-base font-bold">Riwayat poin member ini</h2>

    <div v-if="status === 'loading'" class="flex flex-col gap-2">
      <span class="sr-only">Memuat riwayat poin…</span>
      <Skeleton v-for="n in 4" :key="n" class="h-11" />
    </div>
    <p v-else-if="status === 'error'" role="alert" class="flex flex-wrap items-center gap-2 text-sm text-muted">
      {{ pesan }}
      <button
        type="button"
        class="inline-flex min-h-11 items-center rounded-sm font-bold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
        @click="muat"
      >
        Coba lagi
      </button>
    </p>
    <p v-else-if="!items.length" class="text-sm text-muted">Belum ada transaksi poin.</p>
    <ul v-else class="flex flex-col divide-y divide-subtle" data-riwayat-ringkas>
      <li v-for="t in items" :key="t.id" class="flex items-center gap-3 py-2.5">
        <span class="rounded px-2 py-0.5 text-xs font-bold whitespace-nowrap" :class="JENIS_TRANSAKSI[t.type]?.kelas">{{ JENIS_TRANSAKSI[t.type]?.label ?? t.type }}</span>
        <span class="flex min-w-0 flex-1 flex-col">
          <span class="text-sm [overflow-wrap:anywhere]" :class="t.reference_no && 'font-mono text-xs'">{{ t.reference_no ?? t.note ?? '—' }}</span>
          <span class="text-xs text-muted">{{ keterangan(t) }}</span>
        </span>
        <strong class="shrink-0 text-base tabular-nums" :class="t.points > 0 ? 'text-success-ink' : 'text-danger-ink'">{{ tandaPoin(t.points) }}</strong>
      </li>
    </ul>

    <RouterLink
      v-if="can('member.view') && status === 'ready' && items.length"
      :to="`/admin/member/${userId}`"
      class="inline-flex min-h-11 items-center self-start rounded-sm text-sm font-bold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
    >
      Lihat semua riwayat
    </RouterLink>
  </section>
</template>
