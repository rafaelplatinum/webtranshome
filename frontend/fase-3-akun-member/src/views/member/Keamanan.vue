<script setup>
import { computed, nextTick, reactive, ref, useId } from 'vue'
import Button from '@/components/ui/Button.vue'
import Checkbox from '@/components/ui/Checkbox.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import { errorField, pesanError } from '@/services/errors'
import { changePassword, updateConsent } from '@/services/memberService'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'
import { PASSWORD_MIN, PESAN } from '@/utils/validasi'

const auth = useAuthStore()
const ui = useUiStore()
const idPreferensi = useId()
const idPassword = useId()

// Persetujuan: optimistic — langsung berubah, dikembalikan + toast bila gagal.
const setuju = ref(Boolean(auth.user?.communication_consent))

async function ubahPersetujuan(nilai) {
  const lama = setuju.value
  setuju.value = nilai
  try {
    auth.perbaruiUser(await updateConsent(nilai))
    ui.tampilkanToast({ pesan: nilai ? 'Kamu akan menerima info promo dan poin.' : 'Kamu tidak akan menerima info promo lagi.', jenis: 'success' })
  } catch (error) {
    setuju.value = lama
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  }
}

// Akun yang dibuat lewat Google belum punya password: tampil sebagai "Buat password".
const punyaPassword = computed(() => auth.user?.has_password !== false)
const form = reactive({ lama: '', baru: '', ulangi: '' })
const errors = reactive({ lama: '', baru: '', ulangi: '' })
const galat = ref('')
const menyimpan = ref(false)
const formEl = ref(null)

async function fokusGalat() {
  await nextTick()
  formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
}

async function simpanPassword() {
  if (menyimpan.value) return
  galat.value = ''
  errors.lama = punyaPassword.value && !form.lama ? 'Password lama wajib diisi' : ''
  errors.baru = !form.baru ? PESAN.passwordKosong : form.baru.length < PASSWORD_MIN ? PESAN.passwordPendek : ''
  errors.ulangi = form.ulangi === form.baru ? '' : PESAN.passwordBeda
  if (errors.lama || errors.baru || errors.ulangi) return fokusGalat()
  menyimpan.value = true
  const membuat = !punyaPassword.value
  try {
    await changePassword(membuat ? { new_password: form.baru } : { old_password: form.lama, new_password: form.baru })
    auth.perbaruiUser({ has_password: true })
    Object.assign(form, { lama: '', baru: '', ulangi: '' })
    ui.tampilkanToast({
      pesan: membuat ? 'Password dibuat. Sekarang kamu juga bisa masuk dengan email dan password.' : 'Password diperbarui. Sesi di perangkat lain sudah dikeluarkan.',
      jenis: 'success',
    })
  } catch (error) {
    if (error?.response?.status === 422) {
      const galatField = errorField(error)
      errors.lama = galatField.old_password ?? ''
      errors.baru = galatField.new_password ?? ''
      fokusGalat()
    } else {
      galat.value = pesanError(error)
    }
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <h1 class="text-2xl font-extrabold sm:text-3xl">Keamanan</h1>

    <section :aria-labelledby="idPreferensi" class="flex max-w-2xl flex-col gap-3 rounded-xl border border-border bg-surface p-5 sm:p-6">
      <h2 :id="idPreferensi" class="text-lg font-extrabold">Preferensi</h2>
      <Checkbox
        :model-value="setuju"
        label="Terima info promo dan poin"
        description="Lewat email dan WhatsApp. Bisa dimatikan kapan saja."
        @update:model-value="ubahPersetujuan"
      />
    </section>

    <section :aria-labelledby="idPassword" class="flex max-w-2xl flex-col gap-4 rounded-xl border border-border bg-surface p-5 sm:p-6">
      <div class="flex flex-col gap-1">
        <h2 :id="idPassword" class="text-lg font-extrabold">{{ punyaPassword ? 'Ubah password' : 'Buat password' }}</h2>
        <p class="text-sm text-muted">
          {{
            punyaPassword
              ? 'Setelah diubah, akun ini keluar dari perangkat lain.'
              : 'Akun kamu masuk lewat Google. Buat password supaya bisa masuk dengan email juga.'
          }}
        </p>
      </div>
      <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
      <form ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="simpanPassword">
        <PasswordInput v-if="punyaPassword" v-model="form.lama" label="Password lama" required autocomplete="current-password" :error="errors.lama" />
        <PasswordInput v-model="form.baru" label="Password baru" required autocomplete="new-password" :hint="`Minimal ${PASSWORD_MIN} karakter.`" :error="errors.baru" />
        <PasswordInput v-model="form.ulangi" label="Ulangi password baru" required autocomplete="new-password" :error="errors.ulangi" />
        <Button type="submit" class="self-start" :loading="menyimpan">{{ punyaPassword ? 'Ubah password' : 'Buat password' }}</Button>
      </form>
    </section>
  </div>
</template>
