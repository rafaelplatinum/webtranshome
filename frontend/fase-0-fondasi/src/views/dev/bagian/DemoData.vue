<script setup>
import { ref } from 'vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import Table from '@/components/ui/Table.vue'

const halamanSedikit = ref(1)
const halamanBanyak = ref(6)

const KOLOM = [
  { key: 'name', label: 'Produk' },
  { key: 'category', label: 'Kategori' },
  { key: 'price', label: 'Harga umum', align: 'right' },
  { key: 'stock', label: 'Stok' },
]

const BARIS = [
  { id: 1, name: 'Keramik lantai abu doff 40×40', sku: 'KRM-4040-01', category: 'Keramik lantai', price: '89500.00', unit: 'dus', stock: 'TERSEDIA' },
  { id: 2, name: 'Kloset duduk monoblok putih', sku: 'SAN-KLS-02', category: 'Kloset', price: '3250000.00', unit: 'unit', stock: 'SISA_STOK' },
  { id: 3, name: 'Engsel pintu 4 inch', sku: 'HRD-ENG-04', category: 'Engsel & kunci', price: '32000.00', unit: 'pasang', stock: 'HABIS' },
]
</script>

<template>
  <div class="flex flex-col gap-6">
    <Table :columns="KOLOM" :rows="BARIS" caption="Contoh tabel produk">
      <template #cell-name="{ row }">
        <span class="flex flex-col">
          <strong>{{ row.name }}</strong>
          <span class="font-mono text-xs text-faint">{{ row.sku }}</span>
        </span>
      </template>
      <template #cell-price="{ row }"><PriceTag :price="row.price" :unit="row.unit" size="sm" /></template>
      <template #cell-stock="{ row }"><StockBadge :status="row.stock" /></template>
    </Table>
    <div class="flex flex-col gap-3">
      <p class="text-sm font-semibold text-muted">Paginasi 5 halaman (sekarang {{ halamanSedikit }}) dan 12 halaman (sekarang {{ halamanBanyak }})</p>
      <PaginationBar v-model:page="halamanSedikit" :total-pages="5" />
      <PaginationBar v-model:page="halamanBanyak" :total-pages="12" />
    </div>
  </div>
</template>
