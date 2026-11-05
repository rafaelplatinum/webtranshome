<script setup>
import { MessageCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { useWhatsAppContactStore } from '@/stores/whatsappContact'

const props = defineProps({
  // Context untuk mengisi topik dan detail produk pada form bersama.
  context: { type: Object, default: () => ({}) },
  label: { type: String, default: 'Tanya via WhatsApp' },
  variant: { type: String, default: 'whatsapp' },
  size: { type: String, default: 'md' },
  block: { type: Boolean, default: false },
  wrap: { type: Boolean, default: false },
})
const form = useWhatsAppContactStore()

function bukaForm() { form.buka(props.context) }
</script>

<template>
  <Button
    :variant="variant"
    :size="size"
    :block="block"
    :wrap="wrap"
    @click="bukaForm"
  >
    <template #icon><MessageCircle class="size-4.5" aria-hidden="true" /></template>
    {{ label }}
  </Button>
</template>
