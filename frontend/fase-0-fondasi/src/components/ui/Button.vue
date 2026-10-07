<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { LoaderCircle } from 'lucide-vue-next'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (v) => ['primary', 'secondary', 'outline', 'ghost', 'whatsapp', 'danger'].includes(v),
  },
  size: { type: String, default: 'md', validator: (v) => ['sm', 'md', 'lg'].includes(v) },
  type: { type: String, default: 'button' },
  // `to` → RouterLink; `href` → tautan luar di tab baru; tanpa keduanya → <button>.
  to: { type: [String, Object], default: null },
  href: { type: String, default: null },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  block: { type: Boolean, default: false },
  // true: label boleh turun ke baris kedua (mis. tombol di kartu produk 2 kolom di HP).
  wrap: { type: Boolean, default: false },
})

const VARIAN = {
  primary: 'bg-primary text-primary-foreground hover:bg-primary-hover',
  secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary-hover',
  outline: 'border border-border-strong bg-surface text-foreground hover:border-foreground',
  ghost: 'text-foreground hover:bg-subtle',
  whatsapp: 'bg-whatsapp text-whatsapp-foreground hover:bg-whatsapp-hover',
  danger: 'bg-danger text-danger-foreground hover:bg-danger-ink',
}

const UKURAN = {
  sm: { tinggi: 'h-10', minimal: 'min-h-10', px: 'px-3.5', teks: 'text-sm' },
  md: { tinggi: 'h-11', minimal: 'min-h-11', px: 'px-4.5', teks: 'text-sm' },
  lg: { tinggi: 'h-14', minimal: 'min-h-14', px: 'px-6', teks: 'text-base' },
}

const nonaktif = computed(() => props.disabled || props.loading)

// Tautan yang nonaktif dirender sebagai <button disabled> agar tidak bisa diklik.
const tag = computed(() => {
  if (nonaktif.value) return 'button'
  if (props.to) return RouterLink
  if (props.href) return 'a'
  return 'button'
})

const atribut = computed(() => {
  if (tag.value === RouterLink) return { to: props.to }
  if (tag.value === 'a') return { href: props.href, target: '_blank', rel: 'noopener' }
  return { type: props.type, disabled: nonaktif.value }
})

const kelas = computed(() => {
  const ukuran = UKURAN[props.size]
  return [
    'inline-flex shrink-0 select-none items-center justify-center gap-2 rounded-md font-semibold',
    'transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
    'disabled:cursor-not-allowed disabled:opacity-60',
    props.wrap ? `${ukuran.minimal} py-2 text-center leading-tight` : `${ukuran.tinggi} whitespace-nowrap`,
    // Tombol selebar wadah tidak butuh padding samping lebar; ruangnya untuk label (mis. kartu produk 4 kolom).
    props.block ? 'px-3' : ukuran.px,
    ukuran.teks,
    VARIAN[props.variant],
    props.block && 'w-full',
  ]
})
</script>

<template>
  <!-- data-tombol: efek tekan (scale) di situs publik, lihat main.css. -->
  <component :is="tag" v-bind="atribut" :class="kelas" :aria-busy="loading ? 'true' : undefined" data-tombol>
    <LoaderCircle v-if="loading" class="size-4 animate-spin" aria-hidden="true" />
    <slot v-else name="icon" />
    <slot />
  </component>
</template>
