<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import TombolGoogle from '@/components/akun/TombolGoogle.vue'
import { useAuth } from '@/composables/useAuth'
import { useJeda } from '@/composables/useJeda'
import { resendVerification } from '@/services/authService'
import { pesanError } from '@/services/errors'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'
import { cekEmail, PESAN } from '@/utils/validasi'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const ui = useUiStore()
const { login } = useAuth()
const jeda = useJeda(60)

const form = reactive({ email: typeof route.query.email === 'string' ? route.query.email : '', password: '' })
const errors = reactive({ email: '', password: '' })
const galat = ref('')
const belumVerifikasi = ref(false)
const memuat = ref(false)
const mengirimUlang = ref(false)
const formEl = ref(null)

// Hanya kembali ke halaman di situs ini (bukan admin; admin punya halaman masuk sendiri).
const tujuan = computed(() => {
  const r = route.query.redirect
  return typeof r === 'string' && r.startsWith('/') && !r.startsWith('//') && !r.startsWith('/admin') ? r : null
})

async function masuk() {
  if (memuat.value) return
  galat.value = ''
  belumVerifikasi.value = false
  errors.email = cekEmail(form.email)
  errors.password = form.password ? '' : PESAN.passwordKosong
  if (errors.email || errors.password) {
    await nextTick()
    formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
    return
  }
  memuat.value = true
  try {
    await login({ email: form.email.trim(), password: form.password })
    ui.tampilkanToast({ pesan: `Selamat datang, ${auth.user?.full_name ?? 'member'}.`, jenis: 'success' })
    await router.replace(auth.hasRole('CUSTOMER') ? (tujuan.value ?? '/akun') : '/admin')
  } catch (error) {
    const status = error?.response?.status
    if (status === 401) galat.value = 'Email atau password salah.'
    else if (status === 403 && error.response.data?.code === 'EMAIL_BELUM_VERIFIKASI') belumVerifikasi.value = true
    else galat.value = pesanError(error)
  } finally {
    memuat.value = false
  }
}

async function kirimUlang() {
  mengirimUlang.value = true
  try {
    await resendVerification(form.email.trim())
    ui.tampilkanToast({ pesan: 'Link verifikasi baru sudah dikirim. Cek email kamu.', jenis: 'success' })
    jeda.mulai()
  } catch (error) {
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    mengirimUlang.value = false
  }
}

const kelasTautan = 'rounded-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <KartuAuth judul="Masuk" deskripsi="Masuk untuk melihat kode member dan poin Trans Family kamu.">
    <p v-if="tujuan && !galat && !belumVerifikasi" class="rounded-md bg-info-soft px-3 py-2.5 text-sm font-medium text-info-ink">
      Masuk dulu untuk melanjutkan. Setelah masuk, kamu kembali ke halaman tadi.
    </p>
    <div v-if="belumVerifikasi" role="alert" class="flex flex-col items-start gap-2 rounded-md bg-warning-soft px-3 py-3 text-sm text-warning-ink">
      <p class="font-semibold">Email kamu belum diverifikasi.</p>
      <p>Buka link verifikasi yang kami kirim ke {{ form.email.trim() }}, atau minta link baru.</p>
      <Button variant="outline" size="sm" :loading="mengirimUlang" :disabled="jeda.aktif.value" @click="kirimUlang">
        {{ jeda.aktif.value ? `Kirim ulang (${jeda.sisa.value})` : 'Kirim ulang verifikasi' }}
      </Button>
    </div>
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>

    <form ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="masuk">
      <Input v-model="form.email" label="Email" type="email" autocomplete="email" inputmode="email" required :error="errors.email" />
      <PasswordInput v-model="form.password" autocomplete="current-password" required :error="errors.password" />
      <RouterLink to="/lupa-password" :class="[kelasTautan, '-mt-1 inline-flex min-h-11 items-center self-end text-sm md:min-h-8']">Lupa password?</RouterLink>
      <Button type="submit" block :loading="memuat">Masuk</Button>
    </form>

    <div class="flex items-center gap-3 text-xs text-muted" aria-hidden="true">
      <span class="h-px flex-1 bg-border" />atau<span class="h-px flex-1 bg-border" />
    </div>
    <TombolGoogle label="Masuk dengan Google" />

    <p class="text-center text-sm text-muted">
      Belum punya akun? <RouterLink to="/daftar" :class="kelasTautan">Daftar member</RouterLink>
    </p>
  </KartuAuth>
</template>
