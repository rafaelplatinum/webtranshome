<script setup>
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import PetunjukMock from '@/components/akun/PetunjukMock.vue'
import { useJeda } from '@/composables/useJeda'
import { forgotPassword } from '@/services/authService'
import { pesanError } from '@/services/errors'
import { cekEmail } from '@/utils/validasi'

// Pesan selalu sama, email terdaftar atau tidak, supaya orang lain tidak bisa menebak email yang punya akun.
const PESAN_TERKIRIM = 'Jika email terdaftar, link reset sudah dikirim.'

const jeda = useJeda(60)
const email = ref('')
const galatEmail = ref('')
const galat = ref('')
const terkirim = ref(false)
const mengirim = ref(false)

async function kirim() {
  galat.value = ''
  galatEmail.value = cekEmail(email.value)
  if (galatEmail.value) return
  mengirim.value = true
  try {
    await forgotPassword(email.value.trim())
    terkirim.value = true
    jeda.mulai()
  } catch (error) {
    galat.value = pesanError(error)
  } finally {
    mengirim.value = false
  }
}
</script>

<template>
  <KartuAuth judul="Lupa password" deskripsi="Masukkan email akun kamu. Kami kirim link untuk membuat password baru.">
    <div v-if="terkirim" role="status" class="flex flex-col gap-1 rounded-md bg-success-soft px-3 py-3 text-sm text-success-ink">
      <p class="font-semibold" data-pesan-reset>{{ PESAN_TERKIRIM }}</p>
      <p>Cek kotak masuk dan folder spam. Link berlaku 1 jam.</p>
    </div>
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>

    <form class="flex flex-col gap-4" novalidate @submit.prevent="kirim">
      <Input v-model="email" label="Email" type="email" required autocomplete="email" inputmode="email" :error="galatEmail" />
      <Button type="submit" block :loading="mengirim" :disabled="jeda.aktif.value">
        {{ jeda.aktif.value ? `Kirim ulang (${jeda.sisa.value})` : terkirim ? 'Kirim ulang link reset' : 'Kirim link reset' }}
      </Button>
    </form>
    <PetunjukMock />
    <RouterLink
      to="/masuk"
      class="self-center rounded-sm text-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
    >
      Kembali ke halaman masuk
    </RouterLink>
  </KartuAuth>
</template>
