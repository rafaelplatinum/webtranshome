<script setup>
const model = defineModel({ type: Boolean, default: false })

defineProps({
  label: { type: String, required: true },
  hideLabel: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})
</script>

<template>
  <button
    type="button"
    role="switch"
    :aria-checked="model ? 'true' : 'false'"
    :disabled="disabled"
    class="inline-flex min-h-11 items-center gap-3 rounded-md text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
    @click="model = !model"
  >
    <span
      class="relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors duration-150"
      :class="model ? 'bg-success' : 'bg-border-strong'"
      aria-hidden="true"
    >
      <span
        class="inline-block size-5 rounded-full bg-surface shadow-sm transition-transform duration-150"
        :class="model ? 'translate-x-5.5' : 'translate-x-0.5'"
      />
    </span>
    <span :class="hideLabel ? 'sr-only' : 'text-sm font-semibold'">{{ label }}</span>
  </button>
</template>
