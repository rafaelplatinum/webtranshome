<script setup>
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { useAuth } from '@/composables/useAuth'
import { pesanError } from '@/services/errors'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

// Akun contoh dari services/mock/data/akun.js, hanya untuk menguji header dan route guard.
const AKUN = [
  { label: 'Member', email: 'member@example.com', admin: false },
  { label: 'Admin Katalog', email: 'katalog@example.com', admin: true },
  { label: 'Super Admin', email: 'superadmin@example.com', admin: true },
]
const UJI = ['/akun', '/akun/poin', '/admin', '/admin/produk', '/admin/poin/input', '/admin/role', '/admin/sync-qontak', '/admin/masuk', '/masuk']

const pakaiMock = import.meta.env.VITE_USE_MOCK === 'true'
const auth = useAuthStore()
const ui = useUiStore()
const { login, logout } = useAuth()
const sedangMasuk = ref('')

// Mode mock: semua permintaan bertoken ditolak 401 sampai login lagi (menguji alur sesi habis di form admin).
function simulasikanSesiHabis() {
  try {
    localStorage.setItem('transhome.mockSesiHabis', '1')
    ui.tampilkanToast({ pesan: 'Sesi dianggap habis. Permintaan berikutnya ke API akan ditolak (401).', jenis: 'warning' })
  } catch {
    ui.tampilkanToast({ pesan: 'localStorage diblokir browser.', jenis: 'danger' })
  }
}

// Mode mock: Qontak dianggap tidak bisa dihubungi (503). Pendaftaran tetap berhasil, kontaknya FAILED dan bisa di-retry.
const KUNCI_QONTAK_MATI = 'transhome.mockQontakMati'
function bacaQontakMati() {
  try {
    return localStorage.getItem(KUNCI_QONTAK_MATI) === '1'
  } catch {
    return false
  }
}
const qontakMati = ref(bacaQontakMati())

function gantiQontak() {
  try {
    if (qontakMati.value) localStorage.removeItem(KUNCI_QONTAK_MATI)
    else localStorage.setItem(KUNCI_QONTAK_MATI, '1')
    qontakMati.value = bacaQontakMati()
    ui.tampilkanToast({ pesan: qontakMati.value ? 'Qontak dimatikan: sync berikutnya gagal (503).' : 'Qontak hidup lagi: sync berikutnya berhasil.', jenis: 'info' })
  } catch {
    ui.tampilkanToast({ pesan: 'localStorage diblokir browser.', jenis: 'danger' })
  }
}

// Hapus perubahan dari panel admin yang tersimpan di localStorage, lalu muat ulang dengan data contoh awal.
async function resetDataContoh() {
  const { resetDb } = await import('@/services/mock/db')
  resetDb()
  window.location.reload()
}

async function masukSebagai(akun) {
  sedangMasuk.value = akun.email
  try {
    await login({ email: akun.email, password: 'rahasia123' }, { admin: akun.admin })
    ui.tampilkanToast({ pesan: `Masuk sebagai ${akun.label}.`, jenis: 'success' })
  } catch (error) {
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    sedangMasuk.value = ''
  }
}
</script>

<template>
  <div v-if="pakaiMock" class="flex flex-col gap-5">
    <p class="text-sm">
      Sesi sekarang:
      <strong v-if="auth.user">{{ auth.user.full_name }} ({{ auth.roles.join(', ') }})</strong>
      <strong v-else>belum masuk</strong>
    </p>
    <div class="flex flex-wrap gap-2">
      <Button v-for="akun in AKUN" :key="akun.email" variant="outline" size="sm" :loading="sedangMasuk === akun.email" @click="masukSebagai(akun)">
        Masuk sebagai {{ akun.label }}
      </Button>
      <Button v-if="auth.isLoggedIn" variant="ghost" size="sm" @click="logout">Keluar</Button>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Uji data mock</p>
      <div class="flex flex-wrap gap-2">
        <Button variant="outline" size="sm" @click="simulasikanSesiHabis">Simulasikan sesi habis</Button>
        <Button variant="outline" size="sm" @click="gantiQontak">
          {{ qontakMati ? 'Hidupkan Qontak lagi' : 'Matikan Qontak' }}
        </Button>
        <Button variant="outline" size="sm" @click="resetDataContoh">Reset data contoh</Button>
      </div>
      <p class="text-sm text-muted">Perubahan dari panel admin disimpan di browser ini (localStorage) agar sama di semua tab.</p>
    </div>
    <div class="flex flex-col gap-2">
      <p class="text-sm font-semibold text-muted">Uji route guard (buka setelah masuk/keluar)</p>
      <ul class="flex flex-wrap gap-x-5 gap-y-2 text-sm">
        <li v-for="path in UJI" :key="path">
          <RouterLink :to="path" class="rounded-sm font-mono font-semibold text-primary hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">{{ path }}</RouterLink>
        </li>
      </ul>
    </div>
  </div>
  <p v-else class="text-sm text-muted">Bagian ini hanya aktif saat <code class="font-mono">VITE_USE_MOCK=true</code>.</p>
</template>
