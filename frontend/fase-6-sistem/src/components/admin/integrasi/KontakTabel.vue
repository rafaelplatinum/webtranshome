<script setup>
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Table from '@/components/ui/Table.vue'
import StatusSync from '@/components/admin/StatusSync.vue'
import { usePermission } from '@/composables/usePermission'
import { formatNomorHp, formatTanggalJam } from '@/utils/format'

defineProps({
  // Item GET /admin/qontak/contacts (+ `sedang` selama Coba lagi berjalan)
  kontak: { type: Array, required: true },
})

// `coba-lagi`: induk menjalankan retry; `log`: buka drawer sync_logs.
const emit = defineEmits(['coba-lagi', 'log'])

const { can } = usePermission()

const KOLOM = [
  { key: 'member', label: 'Member', class: 'min-w-56' },
  { key: 'status', label: 'Status', class: 'whitespace-nowrap' },
  { key: 'percobaan', label: 'Percobaan', align: 'right', class: 'whitespace-nowrap' },
  { key: 'respons', label: 'Respons terakhir', class: 'whitespace-nowrap' },
  { key: 'diperbarui', label: 'Diperbarui', class: 'whitespace-nowrap' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

// Kode HTTP dari Qontak; null = belum pernah dikirim.
const teksRespons = (k) => (k.status_code == null ? 'Belum dikirim' : k.status_code === 200 ? '200 OK' : `${k.status_code} Gagal`)
</script>

<template>
  <Table :columns="KOLOM" :rows="kontak" caption="Status sync kontak ke Qontak" bare min-width="min-w-230">
    <template #cell-member="{ row }">
      <span class="flex flex-col items-start gap-0.5">
        <RouterLink
          v-if="can('member.view')"
          :to="`/admin/member/${row.user_id}`"
          class="rounded-sm font-semibold hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          {{ row.name }}
        </RouterLink>
        <strong v-else class="font-semibold">{{ row.name }}</strong>
        <span class="font-mono text-xs text-faint">{{ row.member_code }}</span>
        <span class="text-xs text-muted">{{ row.email }}<template v-if="row.phone_number"> · {{ formatNomorHp(row.phone_number) }}</template></span>
        <span v-if="!row.phone_number" class="text-xs text-muted">Tanpa HP, dikirim dengan email saja</span>
      </span>
    </template>
    <template #cell-status="{ row }">
      <StatusSync :status="row.sync_status" :memproses="Boolean(row.sedang)" />
    </template>
    <template #cell-percobaan="{ row }"><span class="tabular-nums">{{ row.retry_count }}</span></template>
    <template #cell-respons="{ row }">
      <span :class="row.status_code && row.status_code !== 200 ? 'font-semibold text-danger-ink' : 'text-muted'" class="tabular-nums">{{ teksRespons(row) }}</span>
    </template>
    <template #cell-diperbarui="{ row }">
      <time :datetime="row.updated_at" class="tabular-nums">{{ formatTanggalJam(row.updated_at) }}</time>
    </template>
    <template #cell-aksi="{ row }">
      <span class="flex justify-end gap-2">
        <!-- Tetap tampil selama proses supaya fokus tidak hilang. -->
        <Button
          v-if="row.sync_status === 'FAILED' || row.sedang"
          size="sm"
          :loading="Boolean(row.sedang)"
          :data-retry-kontak="row.id"
          @click="emit('coba-lagi', row)"
        >
          Coba lagi <span class="sr-only">sync {{ row.name }}</span>
        </Button>
        <Button variant="outline" size="sm" :data-log-kontak="row.id" @click="emit('log', row)">
          Lihat log <span class="sr-only">sync {{ row.name }}</span>
        </Button>
      </span>
    </template>
  </Table>
</template>
