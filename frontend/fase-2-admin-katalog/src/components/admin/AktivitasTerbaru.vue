<script setup>
import { useId } from 'vue'
import { RouterLink } from 'vue-router'
import { usePermission } from '@/composables/usePermission'
import { formatTanggal, formatWaktuRelatif } from '@/utils/format'
import { labelAksi } from '@/utils/log'

defineProps({
  // logs dari GET /admin/dashboard: [{ id, admin_name, action, module, created_at }]
  logs: { type: Array, required: true },
  // 'all' (punya log.view) atau 'mine' (hanya aktivitas sendiri)
  scope: { type: String, default: 'all' },
})

const { can } = usePermission()
const id = useId()
</script>

<template>
  <section :aria-labelledby="id" class="flex min-w-0 flex-col gap-3.5 rounded-xl border border-border bg-surface p-5">
    <div class="flex items-baseline justify-between gap-3">
      <h2 :id="id" class="text-lg font-extrabold">{{ scope === 'mine' ? 'Aktivitas kamu' : 'Aktivitas terbaru' }}</h2>
      <RouterLink
        v-if="can('log.view')"
        to="/admin/log-aktivitas"
        class="inline-flex min-h-11 items-center rounded-sm text-sm font-bold text-primary hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
      >
        Log lengkap
      </RouterLink>
    </div>
    <p v-if="!logs.length" class="text-sm text-muted">Belum ada aktivitas.</p>
    <div v-else class="overflow-x-auto">
      <table class="w-full min-w-105 border-collapse text-sm">
        <caption class="sr-only">{{ scope === 'mine' ? 'Aktivitas kamu' : 'Aktivitas admin terbaru' }}</caption>
        <thead>
          <tr class="text-left text-xs font-semibold tracking-wide text-muted uppercase">
            <th scope="col" class="border-b border-border px-2 py-2">Admin</th>
            <th scope="col" class="border-b border-border px-2 py-2">Aksi</th>
            <th scope="col" class="border-b border-border px-2 py-2">Modul</th>
            <th scope="col" class="border-b border-border px-2 py-2">Waktu</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id">
            <td class="border-b border-subtle px-2 py-2.5">{{ log.admin_name }}</td>
            <td class="border-b border-subtle px-2 py-2.5">
              {{ labelAksi(log.action) }}
              <code class="block font-mono text-xs text-faint">{{ log.action }}</code>
            </td>
            <td class="border-b border-subtle px-2 py-2.5">{{ log.module }}</td>
            <td class="border-b border-subtle px-2 py-2.5 whitespace-nowrap text-muted">
              <time :datetime="log.created_at" :title="formatTanggal(log.created_at)">{{ formatWaktuRelatif(log.created_at) }}</time>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
