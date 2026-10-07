<script setup>
import { useId } from 'vue'

const model = defineModel({ type: Boolean, default: false })

defineProps({
  label: { type: String, required: true },
  description: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
})

const id = useId()
</script>

<template>
  <label :for="id" class="flex min-h-11 cursor-pointer items-start gap-3 py-2" :class="disabled && 'cursor-not-allowed opacity-60'">
    <input
      :id="id"
      v-model="model"
      type="checkbox"
      :disabled="disabled"
      :aria-describedby="description ? `${id}-desc` : undefined"
      class="mt-0.5 size-5 shrink-0 cursor-pointer accent-primary disabled:cursor-not-allowed"
    />
    <span class="flex flex-col gap-0.5">
      <span class="text-sm font-semibold">{{ label }}</span>
      <span v-if="description" :id="`${id}-desc`" class="text-sm text-muted">{{ description }}</span>
    </span>
  </label>
</template>
