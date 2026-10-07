<script setup>
import { Inbox } from 'lucide-vue-next'

defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>

<template>
  <div class="flex flex-col items-center gap-3 rounded-xl border border-dashed border-border-strong bg-surface px-6 py-12 text-center">
    <div class="text-faint">
      <slot name="icon"><Inbox class="size-10" aria-hidden="true" /></slot>
    </div>
    <p class="text-lg font-bold">{{ title }}</p>
    <p v-if="description" class="max-w-md text-muted">{{ description }}</p>
    <div v-if="$slots.actions" class="mt-2 flex flex-wrap justify-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
