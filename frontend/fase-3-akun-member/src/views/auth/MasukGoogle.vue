<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { LoaderCircle } from 'lucide-vue-next'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

// Callback "Masuk/Daftar dengan Google". Backend kembali ke sini dengan #token=… (berhasil, token di fragment
// supaya tidak terkirim ke server mana pun) atau ?error=… (dibatalkan/gagal).
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const ui = useUiStore()

const asal = route.query.asal === '/daftar' ? '/daftar' : '/masuk'
const redirect = typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/') && !route.query.redirect.startsWith('//') ? route.query.redirect : null

function kembaliKeAsal(pesan) {
  ui.tampilkanToast({ pesan, jenis: 'warning' })
  router.replace({ path: asal, query: redirect ? { redirect } : {} })
}

onMounted(async () => {
  const fragmen = new URLSearchParams(window.location.hash.slice(1))
  const token = fragmen.get('token')
  // Token jangan tertinggal di riwayat browser.
  window.history.replaceState(window.history.state, '', window.location.pathname + window.location.search)
  if (route.query.error || !token) {
    kembaliKeAsal(route.query.error === 'access_denied' ? 'Masuk dengan Google dibatalkan.' : 'Masuk dengan Google gagal. Coba lagi.')
    return
  }
  if (!(await auth.masukDenganToken(token))) {
    kembaliKeAsal('Masuk dengan Google gagal. Coba lagi.')
    return
  }
  const baru = fragmen.get('new') === '1'
  ui.tampilkanToast({ pesan: baru ? `Akun member kamu sudah dibuat. Kode member: ${auth.user?.member_code}.` : `Selamat datang, ${auth.user?.full_name ?? 'member'}.`, jenis: 'success' })
  router.replace(auth.hasRole('CUSTOMER') ? (redirect && !redirect.startsWith('/admin') ? redirect : '/akun') : '/admin')
})
</script>

<template>
  <KartuAuth judul="Masuk dengan Google">
    <p class="flex items-center gap-2 text-muted" aria-busy="true">
      <LoaderCircle class="size-5 animate-spin" aria-hidden="true" />
      Sebentar, akun Google kamu sedang disambungkan…
    </p>
  </KartuAuth>
</template>
