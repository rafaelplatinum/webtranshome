<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ImageOff } from 'lucide-vue-next'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'

const props = defineProps({
  // Item dari GET /products: { sku, name, slug, price_general, unit_sale, stock_status, stock_qty_label, brand, primary_image_url }
  product: { type: Object, required: true },
})

const { productLink } = useWhatsApp()
const tautan = computed(() => `/produk/${props.product.slug}`)
const linkWhatsApp = computed(() => productLink(props.product))
const habis = computed(() => props.product.stock_status === 'HABIS')
</script>

<template>
  <article class="flex flex-col overflow-hidden rounded-xl border border-border bg-surface transition-colors duration-150 hover:border-border-strong">
    <RouterLink :to="tautan" class="block aspect-square bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset" tabindex="-1" aria-hidden="true">
      <img
        v-if="product.primary_image_url"
        :src="product.primary_image_url"
        :alt="product.name"
        loading="lazy"
        width="400"
        height="400"
        class="size-full object-cover"
      />
      <span v-else class="flex size-full flex-col items-center justify-center gap-2 text-faint">
        <ImageOff class="size-8" aria-hidden="true" />
        <span class="text-xs">Foto belum ada</span>
      </span>
    </RouterLink>

    <div class="flex flex-1 flex-col gap-1.5 p-3 sm:p-4">
      <span v-if="product.brand" class="text-xs font-bold tracking-wide text-muted uppercase">{{ product.brand.name }}</span>
      <RouterLink
        :to="tautan"
        class="line-clamp-2 rounded-sm font-semibold leading-snug transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        {{ product.name }}
      </RouterLink>
      <span class="font-mono text-xs text-faint">SKU {{ product.sku }}</span>
      <PriceTag class="mt-1" size="card" :price="product.price_general" :unit="product.unit_sale" />
      <StockBadge :status="product.stock_status" :keterangan="product.stock_qty_label ?? ''" />
      <WhatsAppButton class="mt-auto" :href="linkWhatsApp" :label="habis ? 'Tanya ketersediaan' : 'Tanya via WhatsApp'" block wrap />
    </div>
  </article>
</template>
