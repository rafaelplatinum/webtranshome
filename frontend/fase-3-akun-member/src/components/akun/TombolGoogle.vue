<script setup>
import { useRoute } from 'vue-router'
import logoGoogle from '@/assets/logo-google.svg'
import { urlGoogle } from '@/services/authService'
import { useUiStore } from '@/stores/ui'

defineProps({
  label: { type: String, default: 'Masuk dengan Google' },
})

const route = useRoute()
const ui = useUiStore()

// OAuth dikerjakan backend; setelah selesai kembali ke /masuk/google membawa token (atau error bila dibatalkan).
function mulai() {
  const kembali = new URL('/masuk/google', window.location.origin)
  kembali.searchParams.set('asal', route.path)
  const redirect = route.query.redirect
  if (typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')) kembali.searchParams.set('redirect', redirect)
  const url = urlGoogle(kembali.href)
  if (!url) {
    ui.tampilkanToast({ pesan: 'Masuk dengan Google belum tersedia di mode ini.', jenis: 'warning' })
    return
  }
  window.location.assign(url)
}
</script>

<template>
  <button
    type="button"
    class="inline-flex h-11 w-full items-center justify-center gap-3 rounded-md border border-border-strong bg-surface px-4 text-sm font-semibold transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
    @click="mulai"
  >
    <img :src="logoGoogle" alt="" width="18" height="18" class="size-4.5" />
    {{ label }}
  </button>
</template>
