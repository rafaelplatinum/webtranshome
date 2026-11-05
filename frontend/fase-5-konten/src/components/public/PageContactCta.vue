<script setup>
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'

defineProps({
  heading: { type: String, required: true },
  description: { type: String, required: true },
  buttonLabel: { type: String, default: 'Chat via WhatsApp' },
  tone: { type: String, default: 'orange' },
  context: { type: Object, default: () => ({}) },
})
</script>

<template>
  <section
    class="flex flex-col gap-5 px-5 py-8 text-white sm:flex-row sm:items-center sm:justify-between sm:px-8 lg:px-10 lg:py-12"
    :class="tone === 'navy' ? 'bg-th-navy' : 'bg-th-orange'"
  >
    <div class="max-w-2xl">
      <h2 class="text-2xl leading-tight font-extrabold tracking-tight sm:text-3xl">{{ heading }}</h2>
      <p class="mt-2 text-base leading-relaxed text-white/85">{{ description }}</p>
    </div>
    <WhatsAppButton class="shrink-0" size="lg" :context="context" :label="buttonLabel" />
  </section>
</template>
