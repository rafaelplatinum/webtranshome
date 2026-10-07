<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { MessageCircle } from 'lucide-vue-next'
import { useWhatsApp } from '@/composables/useWhatsApp'

// Tombol melayang ke CS (FR-R). Kalau sudah masuk, pesan diawali nama + member_code (FR-S).
// Halaman dengan bar WhatsApp sendiri di HP (meta.barAksiHp) hanya menampilkannya di layar lebar.
const route = useRoute()
const { generalLink } = useWhatsApp()
const link = computed(() => generalLink())
const hanyaDesktop = computed(() => Boolean(route.meta.barAksiHp))
</script>

<template>
  <a
    v-if="link"
    :href="link"
    target="_blank"
    rel="noopener"
    aria-label="Chat WhatsApp CS"
    class="[view-transition-name:wa-melayang] fixed right-4 bottom-[calc(1rem+env(safe-area-inset-bottom))] z-40 h-14 items-center gap-2 rounded-full bg-whatsapp px-4 font-semibold text-whatsapp-foreground shadow-lg transition-colors duration-150 hover:bg-whatsapp-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 sm:right-6 sm:bottom-[calc(1.5rem+env(safe-area-inset-bottom))] sm:pr-5"
    :class="hanyaDesktop ? 'hidden lg:inline-flex' : 'inline-flex'"
  >
    <MessageCircle class="size-6" aria-hidden="true" />
    <span class="hidden sm:inline" aria-hidden="true">Chat CS</span>
  </a>
</template>
