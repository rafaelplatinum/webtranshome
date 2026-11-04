<script setup>
import { computed } from 'vue'
import { Check, ShoppingCart } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { MAKS_PRODUK_KERANJANG, jumlahMinimal, useKeranjangStore } from '@/stores/keranjang'
import { useUiStore } from '@/stores/ui'
import { formatAngka } from '@/utils/format'

// "Tambah ke keranjang" (KEPUTUSAN #23). Produk sudah di keranjang → jumlahnya ditambah. Stok HABIS → nonaktif
// (tetap bisa "Tanya ketersediaan" lewat WhatsApp). `ikon`: tombol bulat untuk pojok foto kartu produk.
const props = defineProps({
  produk: { type: Object, required: true },
  // Jumlah yang ditambahkan; null = minimal pembelian produk.
  jumlah: { type: Number, default: null },
  ikon: { type: Boolean, default: false },
  variant: { type: String, default: 'secondary' },
  size: { type: String, default: 'md' },
  block: { type: Boolean, default: false },
})

const keranjang = useKeranjangStore()
const ui = useUiStore()

const habis = computed(() => props.produk.stock_status === 'HABIS')
const diKeranjang = computed(() => keranjang.cari(props.produk.id))
const satuan = (jumlah) => `${formatAngka(jumlah)}${props.produk.unit_sale ? ` ${props.produk.unit_sale}` : ''}`
// Nama aksesibel diawali teks yang terlihat ("Tambah ke keranjang"), lalu nama produk.
const label = computed(() => {
  if (habis.value) return `Stok habis: ${props.produk.name} tidak bisa ditambah ke keranjang`
  const ada = diKeranjang.value ? ` (sudah ada ${satuan(diKeranjang.value.jumlah)} di keranjang)` : ''
  return `Tambah ke keranjang: ${props.produk.name}${ada}`
})
const LIHAT = { label: 'Lihat keranjang', to: '/keranjang' }

function tambah() {
  const { hasil, jumlah } = keranjang.tambah(props.produk, props.jumlah ?? jumlahMinimal(props.produk))
  const nama = props.produk.name
  if (hasil === 'baru') {
    ui.tampilkanToast({ pesan: `${nama} masuk keranjang (${satuan(jumlah)}).`, jenis: 'success', aksi: LIHAT })
  } else if (hasil === 'tambah') {
    ui.tampilkanToast({ pesan: `Jumlah ${nama} di keranjang sekarang ${satuan(jumlah)}.`, jenis: 'success', aksi: LIHAT })
  } else if (hasil === 'penuh') {
    ui.tampilkanToast({
      pesan: `Keranjang sudah penuh (maksimal ${MAKS_PRODUK_KERANJANG} produk). Kirim pesanan lewat WhatsApp dulu, atau hapus produk lain.`,
      jenis: 'warning',
      durasi: 6000,
      aksi: LIHAT,
    })
  }
}
</script>

<template>
  <button
    v-if="ikon"
    type="button"
    class="relative inline-flex size-11 cursor-pointer items-center justify-center rounded-full border border-border bg-surface text-foreground transition-colors duration-150 hover:border-primary hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:border-border disabled:hover:bg-surface disabled:hover:text-foreground"
    :disabled="habis"
    :aria-label="label"
    :title="habis ? 'Stok habis' : 'Tambah ke keranjang'"
    data-tambah-keranjang
    :data-di-keranjang="diKeranjang ? '' : undefined"
    @click="tambah"
  >
    <ShoppingCart class="size-5" aria-hidden="true" />
    <span v-if="diKeranjang" class="absolute -top-1 -right-1 flex size-5 items-center justify-center rounded-full bg-primary text-primary-foreground ring-2 ring-surface" aria-hidden="true">
      <Check class="size-3" :stroke-width="3" />
    </span>
  </button>
  <Button
    v-else
    :variant="variant"
    :size="size"
    :block="block"
    :disabled="habis"
    :aria-label="label"
    data-tambah-keranjang
    :data-di-keranjang="diKeranjang ? '' : undefined"
    @click="tambah"
  >
    <template #icon><ShoppingCart class="size-4.5" aria-hidden="true" /></template>
    {{ habis ? 'Stok habis' : 'Tambah ke keranjang' }}
  </Button>
</template>
