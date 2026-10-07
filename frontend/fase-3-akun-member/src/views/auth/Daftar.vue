<script setup>
import { nextTick, reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Checkbox from '@/components/ui/Checkbox.vue'
import Input from '@/components/ui/Input.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import TombolGoogle from '@/components/akun/TombolGoogle.vue'
import { register } from '@/services/authService'
import { errorField, pesanError } from '@/services/errors'
import { normalisasiNomorHp } from '@/utils/format'
import { cekEmail, nomorHpValid, PASSWORD_MIN, PESAN } from '@/utils/validasi'

// Email wajib, HP opsional, persetujuan komunikasi tidak dicentang dari awal (CLAUDE.md §7, KEPUTUSAN #11).
const router = useRouter()
const form = reactive({ full_name: '', email: '', password: '', phone_number: '', communication_consent: false })
const errors = ref({})
const galat = ref('')
const emailTerdaftar = ref(false)
const memuat = ref(false)
const formEl = ref(null)

function validasi() {
  const e = {}
  if (!form.full_name.trim()) e.full_name = PESAN.namaKosong
  const galatEmail = cekEmail(form.email)
  if (galatEmail) e.email = galatEmail
  if (!form.password) e.password = PESAN.passwordKosong
  else if (form.password.length < PASSWORD_MIN) e.password = PESAN.passwordPendek
  if (form.phone_number.trim() && !nomorHpValid(form.phone_number)) e.phone_number = PESAN.hpFormat
  errors.value = e
  return !Object.keys(e).length
}

async function fokusGalat() {
  await nextTick()
  formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
}

async function daftar() {
  if (memuat.value) return
  galat.value = ''
  emailTerdaftar.value = false
  if (!validasi()) return fokusGalat()
  memuat.value = true
  try {
    const email = form.email.trim().toLowerCase()
    await register({
      full_name: form.full_name.trim(),
      email,
      password: form.password,
      phone_number: form.phone_number.trim() ? normalisasiNomorHp(form.phone_number) : null,
      communication_consent: form.communication_consent,
    })
    await router.push({ path: '/verifikasi', query: { email } })
  } catch (error) {
    const status = error?.response?.status
    if (status === 409 || status === 422) {
      const galatField = errorField(error)
      emailTerdaftar.value = status === 409 && Boolean(galatField.email)
      errors.value = {
        ...galatField,
        ...(emailTerdaftar.value && { email: 'Email sudah terdaftar.' }),
        ...(status === 409 && galatField.phone_number && { phone_number: 'Nomor HP sudah dipakai akun lain' }),
      }
      fokusGalat()
    } else {
      galat.value = pesanError(error)
    }
  } finally {
    memuat.value = false
  }
}

const kelasTautan = 'rounded-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <KartuAuth judul="Daftar member" deskripsi="Gratis. Kumpulkan poin dari setiap belanja di toko Transhome.">
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>

    <form ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="daftar">
      <Input v-model="form.full_name" label="Nama lengkap" required autocomplete="name" :error="errors.full_name ?? ''" />
      <Input v-model="form.email" label="Email" type="email" required autocomplete="email" inputmode="email" :error="errors.email ?? ''">
        <template #bawah>
          <RouterLink v-if="emailTerdaftar" :to="{ path: '/masuk', query: { email: form.email.trim() } }" :class="[kelasTautan, 'self-start text-sm']">
            Masuk dengan email ini
          </RouterLink>
        </template>
      </Input>
      <PasswordInput v-model="form.password" required autocomplete="new-password" :hint="`Minimal ${PASSWORD_MIN} karakter.`" :error="errors.password ?? ''" />
      <Input
        v-model="form.phone_number"
        label="Nomor HP (opsional)"
        type="tel"
        inputmode="tel"
        autocomplete="tel"
        placeholder="0812 3456 7890"
        hint="Untuk notifikasi poin lewat WhatsApp. Bisa ditambah nanti."
        :error="errors.phone_number ?? ''"
      />
      <Checkbox
        v-model="form.communication_consent"
        label="Kirimi saya info promo dan poin"
        description="Lewat email dan WhatsApp. Bisa diubah kapan saja di halaman akun."
      />
      <Button type="submit" block :loading="memuat">Daftar</Button>
    </form>

    <div class="flex items-center gap-3 text-xs text-muted" aria-hidden="true">
      <span class="h-px flex-1 bg-border" />atau<span class="h-px flex-1 bg-border" />
    </div>
    <TombolGoogle label="Daftar dengan Google" />

    <p class="text-center text-sm text-muted">
      Sudah punya akun? <RouterLink to="/masuk" :class="kelasTautan">Masuk</RouterLink>
    </p>
  </KartuAuth>
</template>
