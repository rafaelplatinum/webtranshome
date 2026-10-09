<script setup>
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Table from '@/components/ui/Table.vue'
import { formatTanggalJam } from '@/utils/format'
import { labelAksi, subjekLog } from '@/utils/log'

defineProps({
  // Item GET /admin/activity-logs
  logs: { type: Array, required: true },
})

// `detail`: buka drawer perbandingan data untuk log itu.
const emit = defineEmits(['detail'])

const KOLOM = [
  { key: 'waktu', label: 'Waktu', class: 'whitespace-nowrap' },
  { key: 'admin', label: 'Admin', class: 'min-w-40' },
  { key: 'aksi', label: 'Aksi', class: 'min-w-56' },
  { key: 'data', label: 'Data', class: 'min-w-48' },
  { key: 'detail', label: 'Detail', align: 'right' },
]

const adaDetail = (log) => Boolean(log.before_data || log.after_data)
</script>

<template>
  <Table :columns="KOLOM" :rows="logs" caption="Log aktivitas admin" bare min-width="min-w-200">
    <template #cell-waktu="{ row }">
      <time :datetime="row.created_at" class="tabular-nums">{{ formatTanggalJam(row.created_at) }}</time>
    </template>
    <template #cell-admin="{ row }">{{ row.admin_name }}</template>
    <template #cell-aksi="{ row }">
      <span class="flex flex-col items-start gap-1">
        <span class="font-semibold">{{ labelAksi(row.action) }}</span>
        <span class="flex flex-wrap items-center gap-1.5">
          <Badge>{{ row.module }}</Badge>
          <code class="font-mono text-xs text-faint">{{ row.action }}</code>
        </span>
      </span>
    </template>
    <template #cell-data="{ row }">
      <span v-if="subjekLog(row)" class="break-words">{{ subjekLog(row) }}</span>
      <span v-else class="text-faint">—</span>
    </template>
    <template #cell-detail="{ row }">
      <Button v-if="adaDetail(row)" variant="outline" size="sm" :data-log="row.id" @click="emit('detail', row)">
        Lihat detail <span class="sr-only">{{ labelAksi(row.action) }}, {{ formatTanggalJam(row.created_at) }}</span>
      </Button>
      <span v-else class="text-sm text-faint">Tanpa detail</span>
    </template>
  </Table>
</template>
