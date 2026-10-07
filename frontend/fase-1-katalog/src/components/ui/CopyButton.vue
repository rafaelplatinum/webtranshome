<script setup>
import { onBeforeUnmount, ref } from 'vue'
import { Check, Copy } from 'lucide-vue-next'
import { salinTeks } from '@/utils/clipboard'

const props = defineProps({
  text: { type: String, required: true },
  label: { type: String, default: 'Salin' },
  copiedLabel: { type: String, default: 'Disalin' },
  // true: di atas latar gelap (kartu member navy).
  gelap: { type: Boolean, default: false },
})

// `failed`: browser menolak clipboard; induk bisa memilih teksnya agar disalin manual.
const emit = defineEmits(['copied', 'failed'])

const tersalin = ref(false)
let timer = null

async function salin() {
  if (!(await salinTeks(props.text))) {
    emit('failed')
    return
  }
  tersalin.value = true
  emit('copied')
  clearTimeout(timer)
  timer = setTimeout(() => (tersalin.value = false), 2000)
}

onBeforeUnmount(() => clearTimeout(timer))
</script>

<template>
  <button
    type="button"
    class="inline-flex h-11 items-center gap-1.5 rounded-md border px-3 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 md:h-9"
    :class="
      gelap
        ? 'border-on-dark-muted bg-transparent text-secondary-foreground hover:border-secondary-foreground focus-visible:ring-primary-on-dark'
        : 'border-border-strong bg-surface hover:border-foreground focus-visible:ring-ring'
    "
    @click="salin"
  >
    <component :is="tersalin ? Check : Copy" class="size-4" aria-hidden="true" />
    <span aria-live="polite">{{ tersalin ? copiedLabel : label }}</span>
  </button>
</template>
