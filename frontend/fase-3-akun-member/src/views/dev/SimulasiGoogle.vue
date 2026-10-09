<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { UserRound } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { googleMock } from '@/services/authService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'

// Pengganti layar pilih akun Google di mode mock. Hasilnya sama dengan backend asli:
// kembali ke return_to#token=… (berhasil) atau return_to?error=access_denied (batal).
const AKUN = [
  { email: 'rina.lestari@gmail.example', full_name: 'Rina Lestari', keterangan: 'Belum punya akun, dibuatkan member baru' },
  { email: 'member@example.com', full_name: 'Budi Santoso', keterangan: 'Sudah member, langsung masuk' },
]

const route = useRoute()
const ui = useUiStore()
const sedang = ref('')

// Hanya kembali ke situs ini.
function alamatKembali() {
  const nilai = typeof route.query.return_to === 'string' ? route.query.return_to : '/masuk/google'
  const url = new URL(nilai, window.location.origin)
  return url.origin === window.location.origin ? url : new URL('/masuk/google', window.location.origin)
}

async function pilih(akun) {
  sedang.value = akun.email
  try {
    const hasil = await googleMock({ email: akun.email, full_name: akun.full_name })
    const url = alamatKembali()
    url.hash = `token=${encodeURIComponent(hasil.token)}${hasil.new ? '&new=1' : ''}`
    window.location.assign(url.href)
  } catch (error) {
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    sedang.value = ''
  }
}

function batal() {
  const url = alamatKembali()
  url.searchParams.set('error', 'access_denied')
  window.location.assign(url.href)
}
</script>

<template>
  <section class="mx-auto flex max-w-md flex-col gap-5 px-4 py-12">
    <div class="flex flex-col gap-1">
      <p class="text-sm font-semibold text-muted">Mode mock</p>
      <h1 class="text-2xl font-extrabold">Simulasi masuk dengan Google</h1>
      <p class="text-sm text-muted">Pilih akun Google contoh. Di API asli, layar ini milik Google.</p>
    </div>
    <ul class="flex flex-col gap-2">
      <li v-for="akun in AKUN" :key="akun.email">
        <button
          type="button"
          class="flex w-full items-center gap-3 rounded-xl border border-border bg-surface p-4 text-left transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="Boolean(sedang)"
          @click="pilih(akun)"
        >
          <span class="grid size-10 shrink-0 place-items-center rounded-full bg-primary-soft text-primary-hover"><UserRound class="size-5" aria-hidden="true" /></span>
          <span class="flex min-w-0 flex-col">
            <span class="font-semibold">{{ akun.full_name }}</span>
            <span class="truncate text-sm text-muted">{{ akun.email }}</span>
            <span class="text-xs text-faint">{{ akun.keterangan }}</span>
          </span>
        </button>
      </li>
    </ul>
    <Button variant="outline" :disabled="Boolean(sedang)" @click="batal">Batal</Button>
  </section>
</template>
