<script setup>
import { MessageCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

defineProps({
  // Link dari useWhatsApp(); null selama nomor toko belum dimuat → tombol nonaktif.
  href: { type: String, default: null },
  label: { type: String, default: 'Tanya via WhatsApp' },
  size: { type: String, default: 'md' },
  block: { type: Boolean, default: false },
  wrap: { type: Boolean, default: false },
})
</script>

<template>
  <Button
    variant="whatsapp"
    :href="href"
    :disabled="!href"
    :size="size"
    :block="block"
    :wrap="wrap"
    :title="href ? undefined : 'Nomor WhatsApp toko belum tersedia'"
  >
    <template #icon><MessageCircle class="size-4.5" aria-hidden="true" /></template>
    {{ label }}
  </Button>
</template>
