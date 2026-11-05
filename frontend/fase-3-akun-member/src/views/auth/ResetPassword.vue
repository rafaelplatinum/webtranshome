<script setup>
import { nextTick, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { BadgeCheck } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import { resetPassword } from '@/services/authService'
import { ambilTokenReset, bersihkanPasswordReset } from '@/services/passwordResetFlow'
import { errorField, pesanError } from '@/services/errors'
import { PASSWORD_MIN, PESAN } from '@/utils/validasi'

const token = ambilTokenReset()
const tidakBerlaku = ref(!token)
const berhasil = ref(false)
const form = reactive({ password: '', ulangi: '' })
const errors = reactive({ password: '', ulangi: '' })
const galat = ref('')
const menyimpan = ref(false)
const formEl = ref(null)

async function simpan() {
  if (menyimpan.value) return
  galat.value = ''
  errors.password = !form.password ? PESAN.passwordKosong : form.password.length < PASSWORD_MIN ? PESAN.passwordPendek : ''
  errors.ulangi = form.ulangi === form.password ? '' : PESAN.passwordBeda
  if (errors.password || errors.ulangi) {
    await nextTick()
    formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
    return
  }
  menyimpan.value = true
  try {
    await resetPassword({ token, password: form.password })
    bersihkanPasswordReset()
    berhasil.value = true
  } catch (error) {
    if (error?.response?.data?.message === 'TOKEN_TIDAK_VALID') tidakBerlaku.value = true
    else if (error?.response?.status === 422) errors.password = errorField(error).password ?? PESAN.passwordPendek
    else galat.value = pesanError(error)
  } finally { menyimpan.value = false }
}
</script>

<template>
  <KartuAuth v-if="berhasil" judul="Password berhasil diubah" deskripsi="Sekarang kamu bisa masuk menggunakan password baru.">
    <div role="status" class="flex items-center gap-3 rounded-lg bg-success-soft p-4 text-sm text-success-ink">
      <BadgeCheck class="size-6 shrink-0" aria-hidden="true" />
      <p>Password baru sudah aktif. Demi keamanan, silakan masuk kembali.</p>
    </div>
    <Button to="/masuk" block>Masuk ke akun</Button>
  </KartuAuth>

  <KartuAuth v-else-if="tidakBerlaku" judul="Sesi reset tidak tersedia" deskripsi="Mulai kembali untuk meminta kode verifikasi baru.">
    <Button to="/lupa-password" block>Minta kode baru</Button>
  </KartuAuth>

  <KartuAuth v-else judul="Buat password baru" deskripsi="Gunakan setidaknya 8 karakter untuk mengamankan akunmu.">
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
    <form ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="simpan">
      <PasswordInput v-model="form.password" label="Password baru" required autocomplete="new-password" :hint="`Minimal ${PASSWORD_MIN} karakter.`" :error="errors.password" />
      <PasswordInput v-model="form.ulangi" label="Ulangi password baru" required autocomplete="new-password" :error="errors.ulangi" />
      <Button type="submit" block :loading="menyimpan">Simpan password baru</Button>
    </form>
  </KartuAuth>
</template>
