<script setup>
import { RouterView } from 'vue-router'
import Toast from '@/components/ui/Toast.vue'
import WhatsAppContactForm from '@/components/layout/WhatsAppContactForm.vue'
</script>

<template>
  <RouterView />
  <Toast />
  <WhatsAppContactForm />
</template>
