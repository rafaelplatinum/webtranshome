<script setup>
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Switch from '@/components/ui/Switch.vue'
import Table from '@/components/ui/Table.vue'
import { usePermission } from '@/composables/usePermission'
import { formatTanggalJam } from '@/utils/format'

defineProps({
  // Item GET /admin/users (+ `sedang` selama aksi berjalan)
  admins: { type: Array, required: true },
  // users.id admin yang sedang masuk: akunnya sendiri tidak bisa dinonaktifkan atau diganti rolenya di sini.
  idSaya: { type: Number, default: null },
})

// `aktif`: (admin, nilaiBaru). `role` / `reset`: buka modal untuk admin itu.
const emit = defineEmits(['aktif', 'role', 'reset'])

const { can } = usePermission()

const KOLOM = [
  { key: 'admin', label: 'Admin', class: 'min-w-60' },
  { key: 'role', label: 'Role', class: 'min-w-40' },
  { key: 'aktif', label: 'Aktif' },
  { key: 'masuk', label: 'Masuk terakhir', class: 'whitespace-nowrap' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]
</script>

<template>
  <Table :columns="KOLOM" :rows="admins" caption="Daftar pengguna admin" bare min-width="min-w-220">
    <template #cell-admin="{ row }">
      <span class="flex flex-col items-start gap-0.5">
        <span class="flex flex-wrap items-center gap-1.5">
          <strong class="font-semibold">{{ row.full_name ?? row.email }}</strong>
          <Badge v-if="row.id === idSaya" variant="info">Kamu</Badge>
          <Badge v-if="!row.is_active" variant="danger">Nonaktif</Badge>
        </span>
        <span class="text-sm text-muted">{{ row.email }}</span>
      </span>
    </template>
    <template #cell-role="{ row }">
      <span class="flex flex-wrap gap-1">
        <Badge v-for="r in row.roles" :key="r.code" :variant="r.code === 'SUPER_ADMIN' ? 'primary' : 'neutral'">{{ r.name }}</Badge>
      </span>
    </template>
    <template #cell-aktif="{ row }">
      <span class="flex flex-col items-start gap-0.5">
        <Switch
          :model-value="row.is_active"
          :label="`Aktif: ${row.full_name ?? row.email}`"
          hide-label
          :disabled="!can('user.update') || row.id === idSaya || Boolean(row.sedang)"
          :data-aktif-admin="row.id"
          @update:model-value="emit('aktif', row, $event)"
        />
        <span v-if="row.id === idSaya" class="text-xs text-muted">Akunmu sendiri</span>
      </span>
    </template>
    <template #cell-masuk="{ row }">
      <time v-if="row.last_login_at" :datetime="row.last_login_at" class="tabular-nums">{{ formatTanggalJam(row.last_login_at) }}</time>
      <span v-else class="text-faint">Belum pernah</span>
    </template>
    <template #cell-aksi="{ row }">
      <span v-if="can('user.update')" class="flex justify-end gap-2">
        <Button v-if="row.id !== idSaya" variant="outline" size="sm" :data-role-admin="row.id" @click="emit('role', row)">
          Ubah role <span class="sr-only">{{ row.full_name ?? row.email }}</span>
        </Button>
        <Button variant="ghost" size="sm" :disabled="!row.is_active" :data-reset-admin="row.id" @click="emit('reset', row)">
          Reset password <span class="sr-only">{{ row.full_name ?? row.email }}</span>
        </Button>
      </span>
      <span v-else class="text-sm text-faint">—</span>
    </template>
  </Table>
</template>
