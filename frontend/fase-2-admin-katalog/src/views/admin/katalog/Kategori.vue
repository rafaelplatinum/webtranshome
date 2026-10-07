<script setup>
import { computed } from 'vue'
import { CornerDownRight, FolderTree, Plus } from 'lucide-vue-next'
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

const m = useMasterData('categories')
const { can } = usePermission()

const urut = (a, b) => a.sort_order - b.sort_order || a.id - b.id
const induk = computed(() => m.daftar.value.filter((k) => k.parent_id == null).sort(urut))
const anakDari = (id) => m.daftar.value.filter((k) => k.parent_id === id).sort(urut)
// Kategori utama diikuti subkategorinya; urutan hanya ditukar di antara saudara.
const baris = computed(() => induk.value.flatMap((p) => [p, ...anakDari(p.id)]))
const saudara = (k) => (k.parent_id == null ? induk.value : anakDari(k.parent_id))

const KOLOM = [
  { key: 'nama', label: 'Kategori', class: 'min-w-56' },
  { key: 'produk', label: 'Produk', align: 'right' },
  { key: 'urutan', label: 'Urutan', align: 'center' },
  { key: 'aktif', label: 'Aktif', align: 'center' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

const pesanNonaktif = computed(() => {
  const k = m.konfirmasi.value
  if (!k) return ''
  const bagian = [`“${k.name}” tidak akan tampil di menu dan katalog website.`]
  if (k.parent_id == null && anakDari(k.id).length) bagian.push('Subkategori di bawahnya ikut tersembunyi.')
  if (k.product_count) bagian.push(`Kategori ini masih punya ${k.product_count} produk; produknya tetap tersimpan.`)
  return bagian.join(' ')
})
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Kategori utama dan subkategori untuk mega menu, filter, dan halaman katalog.</p>
      <Button v-if="can('category.create')" @click="m.bukaForm()">
        <template #icon><Plus class="size-4.5" aria-hidden="true" /></template>
        Tambah kategori
      </Button>
    </div>

    <div v-if="m.status.value === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4" aria-busy="true">
      <span class="sr-only">Memuat kategori…</span>
      <Skeleton v-for="n in 8" :key="n" class="h-6" />
    </div>
    <ErrorState v-else-if="m.status.value === 'error'" :message="m.pesan.value" @retry="m.muat" />
    <EmptyState v-else-if="!baris.length" title="Belum ada kategori" description="Tambahkan kategori utama dulu, lalu subkategorinya.">
      <template #icon><FolderTree class="size-10" aria-hidden="true" /></template>
    </EmptyState>

    <Table v-else :columns="KOLOM" :rows="baris" caption="Daftar kategori">
      <template #cell-nama="{ row }">
        <div class="flex items-center gap-2" :class="row.parent_id != null && 'pl-6'">
          <CornerDownRight v-if="row.parent_id != null" class="size-4 shrink-0 text-faint" aria-hidden="true" />
          <span class="flex min-w-0 flex-col">
            <span :class="row.parent_id == null ? 'font-bold' : 'font-semibold'">
              {{ row.name }}<span v-if="row.parent_id != null" class="sr-only"> (subkategori)</span>
            </span>
            <span class="font-mono text-xs text-faint">/katalog/{{ row.slug }}</span>
          </span>
        </div>
      </template>
      <template #cell-produk="{ row }"><span class="tabular-nums">{{ row.product_count }}</span></template>
      <template #cell-urutan="{ row }">
        <TombolUrut
          :id="row.id"
          :nama="row.name"
          :pertama="saudara(row)[0] === row"
          :terakhir="saudara(row).at(-1) === row"
          :disabled="!can('category.update')"
          @geser="(arah) => m.geser(row, arah, saudara(row))"
        />
      </template>
      <template #cell-aktif="{ row }">
        <Switch
          :model-value="row.is_active"
          :label="`Tampil di website: ${row.name}`"
          hide-label
          :disabled="!can('category.update') || row.sibuk"
          @update:model-value="(v) => m.ubahAktif(row, v)"
        />
      </template>
      <template #cell-aksi="{ row }">
        <Button v-if="can('category.update')" variant="outline" size="sm" @click="m.bukaForm(row)">Ubah <span class="sr-only">{{ row.name }}</span></Button>
      </template>
    </Table>

    <MasterFormModal v-model:open="m.formBuka.value" jenis="categories" :baris="m.dipilih.value" :induk="induk" @tersimpan="m.tersimpan" />
    <ConfirmModal
      :open="Boolean(m.konfirmasi.value)"
      title="Nonaktifkan kategori?"
      :message="pesanNonaktif"
      confirm-label="Nonaktifkan"
      :loading="m.sedangNonaktif.value"
      @update:open="(buka) => !buka && (m.konfirmasi.value = null)"
      @confirm="m.nonaktifkan"
    />
  </div>
</template>
