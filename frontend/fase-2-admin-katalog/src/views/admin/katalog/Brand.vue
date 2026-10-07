<script setup>
import { computed, nextTick, ref } from 'vue'
import { Plus, Search, SearchX, Tag } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import Table from '@/components/ui/Table.vue'
import MasterFormModal from '@/components/admin/MasterFormModal.vue'
import { useMasterData } from '@/composables/useMasterData'
import { usePermission } from '@/composables/usePermission'
import { deleteBrand } from '@/services/adminMasterService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'

const m = useMasterData('brands')
const { can } = usePermission()
const ui = useUiStore()

// Cari di browser: daftar brand admin sudah dimuat lengkap (termasuk yang nonaktif).
const cari = ref('')
const kotakCari = ref(null)
const semua = computed(() => [...m.daftar.value].sort((a, b) => a.name.localeCompare(b.name, 'id')))
const baris = computed(() => {
  const q = cari.value.trim().toLowerCase()
  return q ? semua.value.filter((b) => b.name.toLowerCase().includes(q) || b.slug.includes(q)) : semua.value
})

const KOLOM = [
  { key: 'nama', label: 'Brand', class: 'min-w-56' },
  { key: 'produk', label: 'Produk', align: 'right' },
  { key: 'aktif', label: 'Aktif', align: 'center' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

const pesanNonaktif = computed(() => {
  const b = m.konfirmasi.value
  if (!b) return ''
  return `Halaman dan filter brand “${b.name}” tidak tampil di website.${b.product_count ? ` ${b.product_count} produknya tetap tersimpan.` : ''}`
})

// Hapus permanen: hanya brand nonaktif yang tidak dipakai produk (server juga menolak: 422 / 409).
const akanDihapus = ref(null)
const menghapus = ref(false)
const bisaDihapus = (b) => !b.is_active && !b.product_count

async function hapus() {
  const b = akanDihapus.value
  menghapus.value = true
  try {
    await deleteBrand(b.id)
    m.daftar.value = m.daftar.value.filter((x) => x.id !== b.id)
    akanDihapus.value = null
    ui.tampilkanToast({ pesan: `Brand ${b.name} dihapus.`, jenis: 'success' })
    await nextTick()
    kotakCari.value?.querySelector('input')?.focus()
  } catch (error) {
    akanDihapus.value = null
    // Brand keburu dipakai produk (admin lain): perbarui jumlah produknya.
    if (error?.response?.status === 409) b.product_count = error.response.data?.product_count ?? b.product_count
    if (error?.response?.status !== 401) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    menghapus.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Brand untuk filter katalog dan halaman /brand. Brand juga bisa ditambah langsung dari form produk.</p>
      <Button v-if="can('brand.create')" @click="m.bukaForm()">
        <template #icon><Plus class="size-4.5" aria-hidden="true" /></template>
        Tambah brand
      </Button>
    </div>

    <div v-if="m.status.value === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4" aria-busy="true">
      <span class="sr-only">Memuat brand…</span>
      <Skeleton v-for="n in 6" :key="n" class="h-10" />
    </div>
    <ErrorState v-else-if="m.status.value === 'error'" :message="m.pesan.value" @retry="m.muat" />
    <EmptyState v-else-if="!semua.length" title="Belum ada brand" description="Brand yang ditambahkan akan muncul di sini.">
      <template #icon><Tag class="size-10" aria-hidden="true" /></template>
    </EmptyState>

    <template v-else>
      <div ref="kotakCari" class="flex flex-wrap items-end justify-between gap-3">
        <div class="w-full sm:w-80">
          <Input v-model="cari" label="Cari brand" type="search" placeholder="Nama atau slug" autocomplete="off">
            <template #akhir><Search class="mr-2.5 size-4.5 text-muted" aria-hidden="true" /></template>
          </Input>
        </div>
        <p class="text-sm text-muted" aria-live="polite">Menampilkan {{ baris.length }} dari {{ semua.length }} brand</p>
      </div>

      <EmptyState v-if="!baris.length" :title="`Brand “${cari.trim()}” tidak ditemukan`" description="Periksa ejaan, atau tambahkan brand baru.">
        <template #icon><SearchX class="size-10" aria-hidden="true" /></template>
        <template #actions><Button variant="outline" @click="cari = ''">Tampilkan semua</Button></template>
      </EmptyState>

      <Table v-else :columns="KOLOM" :rows="baris" caption="Daftar brand">
        <template #cell-nama="{ row }">
          <div class="flex items-center gap-3">
            <span class="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-md border border-border bg-surface text-sm font-extrabold text-muted">
              <img v-if="row.logo_url" :src="row.logo_url" alt="" width="40" height="40" class="size-full object-contain p-1" />
              <span v-else aria-hidden="true">{{ row.name.charAt(0).toUpperCase() }}</span>
            </span>
            <span class="flex min-w-0 flex-col">
              <span class="font-semibold">{{ row.name }}</span>
              <span class="font-mono text-xs text-faint">/brand/{{ row.slug }}</span>
            </span>
          </div>
        </template>
        <template #cell-produk="{ row }"><span class="tabular-nums">{{ row.product_count }}</span></template>
        <template #cell-aktif="{ row }">
          <Switch
            :model-value="row.is_active"
            :label="`Tampil di website: ${row.name}`"
            hide-label
            :disabled="!can('brand.update') || row.sibuk"
            @update:model-value="(v) => m.ubahAktif(row, v)"
          />
        </template>
        <template #cell-aksi="{ row }">
          <span class="flex flex-col items-end gap-1">
            <span class="flex justify-end gap-2">
              <Button v-if="can('brand.update')" variant="outline" size="sm" @click="m.bukaForm(row)">Ubah <span class="sr-only">{{ row.name }}</span></Button>
              <Button
                v-if="can('brand.delete') && !row.is_active"
                variant="ghost"
                size="sm"
                class="text-danger-ink"
                :disabled="!bisaDihapus(row)"
                :data-hapus-brand="row.id"
                @click="akanDihapus = row"
              >
                Hapus <span class="sr-only">{{ row.name }}</span>
              </Button>
            </span>
            <span v-if="can('brand.delete') && !row.is_active && row.product_count" class="text-xs text-muted">Dipakai {{ row.product_count }} produk</span>
          </span>
        </template>
      </Table>
    </template>

    <MasterFormModal v-model:open="m.formBuka.value" jenis="brands" :baris="m.dipilih.value" @tersimpan="m.tersimpan" />
    <ConfirmModal
      :open="Boolean(m.konfirmasi.value)"
      title="Nonaktifkan brand?"
      :message="pesanNonaktif"
      confirm-label="Nonaktifkan"
      :loading="m.sedangNonaktif.value"
      @update:open="(buka) => !buka && (m.konfirmasi.value = null)"
      @confirm="m.nonaktifkan"
    />
    <ConfirmModal
      :open="Boolean(akanDihapus)"
      :title="`Hapus brand ${akanDihapus?.name ?? ''} permanen?`"
      message="Brand ini sudah nonaktif dan tidak dipakai produk. Setelah dihapus tidak bisa dikembalikan; brand dengan nama yang sama harus ditambah ulang."
      confirm-label="Ya, hapus"
      danger
      :loading="menghapus"
      @update:open="(buka) => !buka && (akanDihapus = null)"
      @confirm="hapus"
    />
  </div>
</template>
