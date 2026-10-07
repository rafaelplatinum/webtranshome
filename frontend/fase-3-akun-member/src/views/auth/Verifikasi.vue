<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { CircleCheck, LoaderCircle, MailCheck } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import PetunjukMock from '@/components/akun/PetunjukMock.vue'
import { useJeda } from '@/composables/useJeda'
import { resendVerification, verifyEmail } from '@/services/authService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'
import { cekEmail } from '@/utils/validasi'

// Dua pintu masuk: setelah daftar (/verifikasi?email=) dan dari link email (/verifikasi?token=).
const route = useRoute()
const ui = useUiStore()
const jeda = useJeda(60)

const token = typeof route.query.token === 'string' ? route.query.token : ''
const status = ref(token ? 'memverifikasi' : 'menunggu') // memverifikasi | berhasil | gagal | menunggu
const email = ref(typeof route.query.email === 'string' ? route.query.email : '')
const galatEmail = ref('')
const mengirim = ref(false)

onMounted(async () => {
  if (!token) return
  try {
    const hasil = await verifyEmail(token)
    email.value = hasil?.email ?? ''
    status.value = 'berhasil'
  } catch (error) {
    if (error?.response?.status !== 422) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    status.value = 'gagal'
  }
})

async function kirimUlang() {
  galatEmail.value = cekEmail(email.value)
  if (galatEmail.value) return
  mengirim.value = true
  try {
    await resendVerification(email.value.trim())
    ui.tampilkanToast({ pesan: 'Link verifikasi baru sudah dikirim. Cek email kamu.', jenis: 'success' })
    jeda.mulai()
  } catch (error) {
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    mengirim.value = false
  }
}

const kelasTautan = 'rounded-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <KartuAuth v-if="status === 'memverifikasi'" judul="Memverifikasi email">
    <p class="flex items-center gap-2 text-muted" aria-busy="true">
      <LoaderCircle class="size-5 animate-spin" aria-hidden="true" />
      Sebentar, link verifikasi sedang dicek…
    </p>
  </KartuAuth>

  <KartuAuth v-else-if="status === 'berhasil'" judul="Email terverifikasi" deskripsi="Akun Trans Family kamu sudah aktif. Silakan masuk untuk melihat kode member.">
    <CircleCheck class="size-10 text-success" aria-hidden="true" />
    <Button :to="{ path: '/masuk', query: email ? { email } : {} }" block>Masuk sekarang</Button>
  </KartuAuth>

  <KartuAuth
    v-else
    :judul="status === 'gagal' ? 'Link verifikasi tidak berlaku' : 'Cek email kamu'"
    :deskripsi="status === 'gagal' ? 'Link ini tidak valid atau sudah dipakai. Minta link baru di bawah.' : ''"
  >
    <div v-if="status === 'menunggu' && route.query.email" class="flex items-start gap-3 rounded-md bg-subtle p-3 text-sm">
      <MailCheck class="mt-0.5 size-5 shrink-0 text-primary" aria-hidden="true" />
      <p>
        Kami mengirim link verifikasi ke <strong class="[overflow-wrap:anywhere]">{{ route.query.email }}</strong>. Buka link itu untuk mengaktifkan akun, lalu masuk.
      </p>
    </div>
    <form class="flex flex-col gap-3" novalidate @submit.prevent="kirimUlang">
      <Input v-if="!route.query.email" v-model="email" label="Email akun" type="email" autocomplete="email" :error="galatEmail" />
      <p v-else class="text-sm text-muted">Belum menerima email? Cek folder spam, atau kirim ulang.</p>
      <Button type="submit" variant="outline" block :loading="mengirim" :disabled="jeda.aktif.value">
        {{ jeda.aktif.value ? `Kirim ulang (${jeda.sisa.value})` : 'Kirim ulang email verifikasi' }}
      </Button>
    </form>
    <PetunjukMock />
    <p class="text-center text-sm text-muted">
      Sudah verifikasi? <RouterLink to="/masuk" :class="kelasTautan">Masuk</RouterLink>
      · Salah email? <RouterLink to="/daftar" :class="kelasTautan">Daftar lagi</RouterLink>
    </p>
  </KartuAuth>
</template>
