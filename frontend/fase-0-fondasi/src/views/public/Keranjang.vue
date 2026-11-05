<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ImageOff, Minus, Plus, ShoppingCart, Trash2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import { useAuthStore } from '@/stores/auth'
import { useCartStore } from '@/stores/cart'
import { useSettingsStore } from '@/stores/settings'
import { buatLinkWhatsApp } from '@/utils/whatsapp'

const cart = useCartStore()
const auth = useAuthStore()
const settings = useSettingsStore()
const pesan = computed(() => {
  const daftar = cart.items.map((item, index) => `${index + 1}. ${item.name} SKU: ${item.sku} ${new URL(`/produk/${item.slug}`, window.location.origin).href}${item.quantity > 1 ? ` (jumlah: ${item.quantity})` : ''}`)
  return [
    'Halo Transhome, saya ingin bertanya tentang produk ini:',
    ...daftar,
    '',
    `Nama : ${auth.user?.full_name || '[Nama customer]'}`,
    `Nomor telepon : ${auth.user?.phone_number || '[Nomor telepon]'}`,
  ].join('\n')
})
const linkChat = computed(() => settings.waNumber && cart.items.length ? buatLinkWhatsApp(settings.waNumber, pesan.value) : null)
</script>

<template>
  <div class="mx-auto flex w-full max-w-4xl flex-col gap-6 px-4 py-8 sm:px-6 sm:py-12">
    <div class="flex items-center gap-3">
      <ShoppingCart class="size-7 text-primary" aria-hidden="true" />
      <div><h1 class="text-2xl font-extrabold">Keranjang belanja</h1><p class="text-sm text-muted">{{ cart.jumlah }} produk</p></div>
    </div>

    <EmptyState v-if="!cart.items.length" title="Keranjang masih kosong" description="Tambahkan produk yang ingin ditanyakan ke Transhome.">
      <template #actions><Button to="/katalog">Lihat katalog</Button></template>
    </EmptyState>
    <template v-else>
      <ul class="divide-y divide-border rounded-xl border border-border bg-surface">
        <li v-for="item in cart.items" :key="item.sku" class="flex flex-col gap-4 p-4 sm:flex-row sm:items-center">
          <RouterLink :to="`/produk/${item.slug}`" class="grid aspect-square w-full max-w-28 shrink-0 place-items-center overflow-hidden rounded-lg bg-subtle sm:size-24" :aria-label="`Lihat ${item.name}`">
            <img v-if="item.image" :src="item.image" :alt="item.name" loading="lazy" class="size-full object-cover" />
            <ImageOff v-else class="size-8 text-faint" aria-hidden="true" />
          </RouterLink>
          <div class="min-w-0 flex-1 space-y-1">
            <span v-if="item.brand" class="block text-xs font-bold tracking-wide text-muted uppercase">{{ item.brand }}</span>
            <RouterLink :to="`/produk/${item.slug}`" class="font-semibold leading-snug hover:text-primary">{{ item.name }}</RouterLink>
            <p class="font-mono text-xs text-muted">SKU: {{ item.sku }}</p>
            <PriceTag size="sm" :price="item.price" :show="item.showPrice" :unit="item.unit" />
            <StockBadge :status="item.stockStatus" :keterangan="item.stockLabel" />
          </div>
          <div class="flex items-center gap-2 self-end sm:self-center">
            <Button variant="outline" size="sm" :aria-label="`Kurangi ${item.name}`" @click="cart.ubahJumlah(item.sku, item.quantity - 1)"><Minus class="size-4" /></Button>
            <span class="min-w-6 text-center">{{ item.quantity }}</span>
            <Button variant="outline" size="sm" :aria-label="`Tambah ${item.name}`" @click="cart.ubahJumlah(item.sku, item.quantity + 1)"><Plus class="size-4" /></Button>
            <Button variant="ghost" size="sm" :aria-label="`Hapus ${item.name}`" @click="cart.hapus(item.sku)"><Trash2 class="size-4 text-danger" /></Button>
          </div>
        </li>
      </ul>

      <div class="rounded-xl border border-border bg-surface p-4 sm:p-5">
        <h2 class="font-bold">Data untuk chat</h2>
        <p v-if="auth.user" class="mt-1 text-sm text-muted">Nama dan nomor telepon diambil dari profil akun Anda.</p>
        <p v-else class="mt-1 text-sm text-muted">Masuk agar nama dan nomor telepon dari akun terisi otomatis.</p>
        <p v-if="!auth.user" class="mt-3 text-sm">Nama dan telepon akan ditampilkan sebagai isian di pesan WhatsApp.</p>
        <Button v-if="auth.user" class="mt-4" :disabled="!linkChat" :href="linkChat" target="_blank" rel="noopener">Lanjut ke WhatsApp</Button>
        <Button v-else class="mt-4" to="/masuk">Masuk untuk lanjut</Button>
      </div>
    </template>
  </div>
</template>
