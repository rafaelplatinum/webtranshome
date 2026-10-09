<script setup>
import { computed } from 'vue'
import { ImageOff, LayoutDashboard, Plus } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import Table from '@/components/ui/Table.vue'
import MasterFormModal from '@/components/admin/MasterFormModal.vue'
import TombolUrut from '@/components/admin/TombolUrut.vue'
import { useMasterData } from '@/composables/useMasterData'
import { usePermission } from '@/composables/usePermission'

const m = useMasterData('rooms')
const { can } = usePermission()

const baris = computed(() => [...m.daftar.value].sort((a, b) => a.sort_order - b.sort_order || a.id - b.id))

const KOLOM = [
  { key: 'nama', label: 'Ruangan', class: 'min-w-56' },
  { key: 'produk', label: 'Produk', align: 'right' },
  { key: 'urutan', label: 'Urutan', align: 'center' },
  { key: 'aktif', label: 'Aktif', align: 'center' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

const pesanNonaktif = computed(() => {
  const r = m.konfirmasi.value
  return r ? `Ruangan “${r.name}” tidak tampil di website.${r.product_count ? ` ${r.product_count} produknya tetap tersimpan.` : ''}` : ''
})
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Urutan di sini menjadi urutan kartu di halaman Belanja per ruangan.</p>
      <Button v-if="can('room.create')" @click="m.bukaForm()">
        <template #icon><Plus class="size-4.5" aria-hidden="true" /></template>
        Tambah ruangan
      </Button>
    </div>

    <div v-if="m.status.value === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4" aria-busy="true">
      <span class="sr-only">Memuat ruangan…</span>
      <Skeleton v-for="n in 5" :key="n" class="h-12" />
    </div>
    <ErrorState v-else-if="m.status.value === 'error'" :message="m.pesan.value" @retry="m.muat" />
    <EmptyState v-else-if="!baris.length" title="Belum ada ruangan" description="Tambahkan ruangan beserta foto cover-nya.">
      <template #icon><LayoutDashboard class="size-10" aria-hidden="true" /></template>
    </EmptyState>

    <Table v-else :columns="KOLOM" :rows="baris" caption="Daftar ruangan">
      <template #cell-nama="{ row }">
        <div class="flex items-center gap-3">
          <span class="flex h-12 w-16 shrink-0 items-center justify-center overflow-hidden rounded-md bg-subtle text-faint">
            <img v-if="row.image_cover" :src="row.image_cover" alt="" width="64" height="48" class="size-full object-cover" />
            <ImageOff v-else class="size-4.5" aria-hidden="true" />
          </span>
          <span class="flex min-w-0 flex-col">
            <span class="font-semibold">{{ row.name }}</span>
            <span class="font-mono text-xs text-faint">/ruangan/{{ row.slug }}</span>
          </span>
        </div>
      </template>
      <template #cell-produk="{ row }"><span class="tabular-nums">{{ row.product_count }}</span></template>
      <template #cell-urutan="{ row }">
        <TombolUrut
          :id="row.id"
          :nama="row.name"
          :pertama="baris[0] === row"
          :terakhir="baris.at(-1) === row"
          :disabled="!can('room.update')"
          @geser="(arah) => m.geser(row, arah, baris)"
        />
      </template>
      <template #cell-aktif="{ row }">
        <Switch
          :model-value="row.is_active"
          :label="`Tampil di website: ${row.name}`"
          hide-label
          :disabled="!can('room.update') || row.sibuk"
          @update:model-value="(v) => m.ubahAktif(row, v)"
        />
      </template>
      <template #cell-aksi="{ row }">
        <div class="flex items-center justify-end gap-1.5 whitespace-nowrap">
          <Button v-if="can('room.update')" variant="outline" size="sm" @click="m.bukaForm(row)">Ubah <span class="sr-only">{{ row.name }}</span></Button>
          <Button v-if="can('product.view')" variant="ghost" size="sm" :to="{ path: '/admin/produk', query: { room: row.id } }">
            Lihat produk <span class="sr-only">{{ row.name }}</span>
          </Button>
        </div>
      </template>
    </Table>

    <MasterFormModal v-model:open="m.formBuka.value" jenis="rooms" :baris="m.dipilih.value" @tersimpan="m.tersimpan" />
    <ConfirmModal
      :open="Boolean(m.konfirmasi.value)"
      title="Nonaktifkan ruangan?"
      :message="pesanNonaktif"
      confirm-label="Nonaktifkan"
      :loading="m.sedangNonaktif.value"
      @update:open="(buka) => !buka && (m.konfirmasi.value = null)"
      @confirm="m.nonaktifkan"
    />
  </div>
</template>
