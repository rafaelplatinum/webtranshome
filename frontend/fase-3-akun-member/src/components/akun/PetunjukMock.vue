<script setup>
import { RouterLink } from 'vue-router'

// Mode mock tidak mengirim email sungguhan; petunjuk ini hanya tampil saat development + mock.
const tampil = import.meta.env.DEV && import.meta.env.VITE_USE_MOCK === 'true'
</script>

<template>
  <p v-if="tampil" class="rounded-md bg-subtle px-3 py-2.5 text-xs text-muted">
    Mode mock: email tidak benar-benar dikirim. Buka
    <RouterLink to="/dev/email" class="rounded-sm font-semibold text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">kotak email tiruan</RouterLink>
    untuk membuka link-nya.
  </p>
</template>
