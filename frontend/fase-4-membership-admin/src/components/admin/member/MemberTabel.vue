<script setup>
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Table from '@/components/ui/Table.vue'
import StatusSync from '@/components/admin/StatusSync.vue'
import { formatAngka, formatNomorHp, formatTanggal } from '@/utils/format'

defineProps({
  // Item GET /admin/members
  member: { type: Array, required: true },
})

const KOLOM = [
  { key: 'member', label: 'Member', class: 'min-w-52' },
  { key: 'email', label: 'Email' },
  { key: 'hp', label: 'Nomor HP', class: 'whitespace-nowrap' },
  { key: 'saldo', label: 'Saldo poin', align: 'right', class: 'whitespace-nowrap' },
  { key: 'consent', label: 'Info promo', class: 'whitespace-nowrap' },
  { key: 'sync', label: 'Sync Qontak', class: 'whitespace-nowrap' },
  { key: 'terdaftar', label: 'Terdaftar', class: 'whitespace-nowrap' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

// Email panjang hanya turun baris sebelum "@" (nama.panjang / @example.com), tidak di tengah kata.
function pecahEmail(email) {
  const i = String(email ?? '').lastIndexOf('@')
  return i > 0 ? [email.slice(0, i), email.slice(i)] : [email, '']
}
</script>

<template>
  <Table :columns="KOLOM" :rows="member" caption="Daftar member" bare min-width="min-w-240">
    <template #cell-member="{ row }">
      <span class="flex flex-col items-start gap-0.5">
        <strong class="font-semibold">{{ row.full_name }}</strong>
        <span class="font-mono text-xs text-faint">{{ row.member_code }}</span>
        <Badge v-if="!row.is_active" variant="danger">Nonaktif</Badge>
      </span>
    </template>
    <template #cell-email="{ row }">
      <span class="flex flex-col items-start gap-0.5">
        <span>{{ pecahEmail(row.email)[0] }}<wbr />{{ pecahEmail(row.email)[1] }}</span>
        <Badge v-if="!row.email_verified" variant="warning">Belum verifikasi</Badge>
      </span>
    </template>
    <template #cell-hp="{ row }">
      <span v-if="row.phone_number" class="tabular-nums">{{ formatNomorHp(row.phone_number) }}</span>
      <span v-else class="text-faint">Belum diisi</span>
    </template>
    <template #cell-saldo="{ row }">
      <strong class="tabular-nums">{{ formatAngka(row.balance) }}</strong>
    </template>
    <template #cell-consent="{ row }">
      <Badge :variant="row.communication_consent ? 'success' : 'neutral'">{{ row.communication_consent ? 'Setuju' : 'Tidak' }}</Badge>
    </template>
    <template #cell-sync="{ row }">
      <StatusSync :status="row.sync_status" />
    </template>
    <template #cell-terdaftar="{ row }">{{ formatTanggal(row.created_at) }}</template>
    <template #cell-aksi="{ row }">
      <Button variant="outline" size="sm" :to="`/admin/member/${row.id}`">
        Lihat <span class="sr-only">{{ row.full_name }}</span>
      </Button>
    </template>
  </Table>
</template>
