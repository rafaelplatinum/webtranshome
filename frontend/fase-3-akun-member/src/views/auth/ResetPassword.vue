<script setup>
import { nextTick, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import { resetPassword } from '@/services/authService'
import { errorField, pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'
import { PASSWORD_MIN, PESAN } from '@/utils/validasi'

// Link reset hanya bisa dipakai sekali; yang tidak berlaku diarahkan untuk minta link baru.
const route = useRoute()
const router = useRouter()
const ui = useUiStore()

const token = typeof route.query.token === 'string' ? route.query.token : ''
const tidakBerlaku = ref(!token)
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
    ui.tampilkanToast({ pesan: 'Password diperbarui. Silakan masuk dengan password baru.', jenis: 'success' })
    await router.replace('/masuk')
  } catch (error) {
    if (error?.response?.data?.code === 'TOKEN_TIDAK_VALID') tidakBerlaku.value = true
    else if (error?.response?.status === 422) errors.password = errorField(error).password ?? PESAN.passwordPendek
    else galat.value = pesanError(error)
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <KartuAuth v-if="tidakBerlaku" judul="Link tidak berlaku" deskripsi="Link reset ini tidak valid, sudah dipakai, atau sudah kedaluwarsa.">
    <Button to="/lupa-password" block>Minta link baru</Button>
  </KartuAuth>

  <KartuAuth v-else judul="Buat password baru" deskripsi="Password baru langsung berlaku untuk masuk berikutnya.">
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
    <form ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="simpan">
      <PasswordInput v-model="form.password" label="Password baru" required autocomplete="new-password" :hint="`Minimal ${PASSWORD_MIN} karakter.`" :error="errors.password" />
      <PasswordInput v-model="form.ulangi" label="Ulangi password baru" required autocomplete="new-password" :error="errors.ulangi" />
      <Button type="submit" block :loading="menyimpan">Simpan password baru</Button>
    </form>
  </KartuAuth>
</template>
