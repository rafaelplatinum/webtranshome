<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ImageOff, Trash2 } from 'lucide-vue-next'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import PengaturJumlah from './PengaturJumlah.vue'
import { MAKS_JUMLAH, jumlahMinimal, useKeranjangStore } from '@/stores/keranjang'
import { formatAngka, formatRupiah } from '@/utils/format'
import { hargaDisembunyikan } from '@/utils/whatsapp'

// Satu produk di /keranjang: foto, nama, harga satuan, stok, jumlah, subtotal, hapus.
const props = defineProps({
  // { produk, jumlah, tidakDijual } dari store keranjang.
  item: { type: Object, required: true },
})
const emit = defineEmits(['hapus'])

const keranjang = useKeranjangStore()
const produk = computed(() => props.item.produk)
const tautan = computed(() => `/produk/${produk.value.slug}`)
const minimal = computed(() => jumlahMinimal(produk.value))
const jumlah = computed({
  get: () => props.item.jumlah,
  set: (nilai) => keranjang.ubahJumlah(produk.value.id, nilai),
})
const subtotal = computed(() => (hargaDisembunyikan(produk.value) ? null : formatRupiah(Number(produk.value.price_general) * props.item.jumlah)))
</script>

<template>
  <li class="flex gap-3 p-3 sm:gap-4 sm:p-4" :class="item.tidakDijual && 'bg-subtle'" data-baris-keranjang>
    <RouterLink :to="tautan" class="size-20 shrink-0 overflow-hidden rounded-lg bg-subtle sm:size-24" tabindex="-1" aria-hidden="true">
      <img v-if="produk.primary_image_url" :src="produk.primary_image_url" :alt="produk.name" loading="lazy" width="96" height="96" class="size-full object-cover" :class="item.tidakDijual && 'opacity-50 grayscale'" />
      <span v-else class="flex size-full items-center justify-center text-faint"><ImageOff class="size-6" aria-hidden="true" /></span>
    </RouterLink>

    <div class="flex min-w-0 flex-1 flex-col gap-2">
      <div class="flex items-start gap-2">
        <div class="flex min-w-0 flex-1 flex-col gap-0.5">
          <span v-if="produk.brand" class="text-xs font-bold tracking-wide text-muted uppercase">{{ produk.brand.name }}</span>
          <RouterLink
            :to="tautan"
            class="line-clamp-2 rounded-sm font-semibold leading-snug transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            {{ produk.name }}
          </RouterLink>
          <span class="font-mono text-xs text-faint">SKU {{ produk.sku }}</span>
        </div>
        <button
          type="button"
          class="-mt-1.5 -mr-1.5 inline-flex size-11 shrink-0 cursor-pointer items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-danger-soft hover:text-danger-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          :aria-label="`Hapus ${produk.name} dari keranjang`"
          :data-hapus-keranjang="produk.id"
          @click="emit('hapus', item)"
        >
          <Trash2 class="size-5" aria-hidden="true" />
        </button>
      </div>

      <p v-if="item.tidakDijual" class="text-sm font-semibold text-danger-ink" data-tidak-dijual>
        Produk ini sudah tidak dijual, jadi tidak ikut dikirim. Hapus dari keranjang.
      </p>
      <template v-else>
        <div class="flex flex-wrap items-center gap-x-3 gap-y-1.5">
          <PriceTag size="sm" :price="produk.price_general" :show="produk.show_price !== false" :unit="produk.unit_sale" />
          <StockBadge :status="produk.stock_status" />
        </div>
        <div class="flex flex-wrap items-end justify-between gap-x-4 gap-y-2">
          <PengaturJumlah v-model="jumlah" :min="minimal" :max="MAKS_JUMLAH" :satuan="produk.unit_sale" hide-label :nama="produk.name" />
          <p class="ml-auto flex flex-col items-end text-right" data-subtotal>
            <span class="text-xs text-muted">Subtotal</span>
            <strong v-if="subtotal" class="text-lg tabular-nums">{{ subtotal }}</strong>
            <span v-else class="text-sm font-bold text-muted">Harga via WhatsApp</span>
          </p>
        </div>
        <p v-if="minimal > 1" class="text-xs text-muted">Minimal pembelian {{ formatAngka(minimal) }} {{ produk.unit_sale }}</p>
        <p v-if="produk.stock_status === 'HABIS'" class="text-sm text-danger-ink">Stok sedang habis. Tetap ikut dikirim supaya tim kami bisa memberi info ketersediaannya.</p>
      </template>
    </div>
  </li>
</template>
