<script setup>
import { formatRupiah, formatTanggal } from '@/utils/format'

defineProps({
  // Item GET /me/points: { id, type, points, reference_no, purchase_amount, note, created_at }
  items: { type: Array, required: true },
})

// Jenis dibedakan dengan label teks, bukan warna saja.
const JENIS = {
  EARN: { label: 'Didapat', kelas: 'bg-success-soft text-success-ink' },
  REDEEM: { label: 'Ditukar', kelas: 'bg-primary-soft text-primary-hover' },
  ADJUST: { label: 'Penyesuaian', kelas: 'bg-info-soft text-info-ink' },
  EXPIRE: { label: 'Kedaluwarsa', kelas: 'bg-subtle text-muted' },
}

const tandaPoin = (n) => (n > 0 ? `+${n}` : `−${Math.abs(n)}`)
const keterangan = (t) => t.reference_no ?? t.note ?? '—'
</script>

<template>
  <!-- HP: daftar bertumpuk. Mulai 640px: tabel. -->
  <div>
    <ul class="flex flex-col divide-y divide-subtle sm:hidden">
      <li v-for="t in items" :key="t.id" class="flex items-start justify-between gap-3 py-3">
        <div class="flex min-w-0 flex-col gap-1">
          <span class="flex flex-wrap items-center gap-2">
            <span class="rounded px-2 py-0.5 text-xs font-bold" :class="JENIS[t.type]?.kelas">{{ JENIS[t.type]?.label ?? t.type }}</span>
            <span class="text-xs text-muted">{{ formatTanggal(t.created_at) }}</span>
          </span>
          <span class="text-sm" :class="t.reference_no && 'font-mono'">{{ keterangan(t) }}</span>
          <span v-if="t.purchase_amount" class="text-xs text-muted">Belanja {{ formatRupiah(t.purchase_amount) }}</span>
        </div>
        <span class="shrink-0 text-base font-extrabold tabular-nums" :class="t.points > 0 ? 'text-success-ink' : 'text-danger-ink'">{{ tandaPoin(t.points) }}</span>
      </li>
    </ul>

    <div class="hidden overflow-x-auto sm:block">
      <table class="w-full border-collapse text-sm">
        <thead>
          <tr class="text-left text-xs font-bold tracking-wide text-muted uppercase">
            <th scope="col" class="border-b border-border px-3 py-2.5">Tanggal</th>
            <th scope="col" class="border-b border-border px-3 py-2.5">Jenis</th>
            <th scope="col" class="border-b border-border px-3 py-2.5">No. nota / keterangan</th>
            <th scope="col" class="border-b border-border px-3 py-2.5 text-right">Total belanja</th>
            <th scope="col" class="border-b border-border px-3 py-2.5 text-right">Poin</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in items" :key="t.id">
            <td class="border-b border-subtle px-3 py-3 whitespace-nowrap">{{ formatTanggal(t.created_at) }}</td>
            <td class="border-b border-subtle px-3 py-3">
              <span class="rounded px-2 py-0.5 text-xs font-bold whitespace-nowrap" :class="JENIS[t.type]?.kelas">{{ JENIS[t.type]?.label ?? t.type }}</span>
            </td>
            <td class="border-b border-subtle px-3 py-3" :class="t.reference_no && 'font-mono text-xs'">{{ keterangan(t) }}</td>
            <td class="border-b border-subtle px-3 py-3 text-right whitespace-nowrap tabular-nums">{{ t.purchase_amount ? formatRupiah(t.purchase_amount) : '—' }}</td>
            <td class="border-b border-subtle px-3 py-3 text-right font-extrabold tabular-nums" :class="t.points > 0 ? 'text-success-ink' : 'text-danger-ink'">
              {{ tandaPoin(t.points) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
