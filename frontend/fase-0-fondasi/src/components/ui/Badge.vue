<script setup>
defineProps({
  variant: {
    type: String,
    default: 'neutral',
    validator: (v) => ['neutral', 'primary', 'success', 'warning', 'info', 'danger', 'accent'].includes(v),
  },
})

const VARIAN = {
  neutral: 'bg-subtle text-foreground',
  primary: 'bg-primary-soft text-primary-hover',
  success: 'bg-success-soft text-success-ink',
  warning: 'bg-warning-soft text-warning-ink',
  info: 'bg-info-soft text-info-ink',
  danger: 'bg-danger-soft text-danger-ink',
  accent: 'bg-accent text-accent-foreground',
}
</script>

<template>
  <span class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs font-semibold whitespace-nowrap" :class="VARIAN[variant]">
    <slot />
  </span>
</template>
