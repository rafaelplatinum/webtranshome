<script setup>
import { computed } from 'vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'

// Pita ajakan WhatsApp di atas footer (desain dari redesign Dylan, 6 Nov 2026). Tombol langsung ke WhatsApp;
// `href` boleh diisi tautan khusus, bawaannya pesan umum ke CS.
const props = defineProps({
  heading: { type: String, required: true },
  description: { type: String, required: true },
  buttonLabel: { type: String, default: 'Chat via WhatsApp' },
  tone: { type: String, default: 'orange' },
  href: { type: String, default: null },
})

const { generalLink } = useWhatsApp()
const tautan = computed(() => props.href ?? generalLink())
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
    <WhatsAppButton class="shrink-0" size="lg" :href="tautan" :label="buttonLabel" />
  </section>
</template>
