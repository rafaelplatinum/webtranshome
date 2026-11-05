<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import { useJeda } from '@/composables/useJeda'
import { requestPasswordOTP, verifyPasswordOTP } from '@/services/authService'
import { dataPasswordReset, simpanTokenReset } from '@/services/passwordResetFlow'
import { pesanError } from '@/services/errors'

const router = useRouter()
const challenge = dataPasswordReset()
const jeda = useJeda(60)
const digits = ref(['', '', '', '', '', ''])
const inputs = ref([])
const galat = ref('')
const mengirim = ref(false)
const memverifikasi = ref(false)
const tujuanTersamar = computed(() => {
  if (!challenge) return ''
  if (challenge.channel === 'email') {
    const [nama, domain] = challenge.destination.split('@')
    return `${nama.slice(0, 1)}${'*'.repeat(Math.max(2, nama.length - 1))}@${domain}`
  }
  const phone = challenge.destination.replace(/\D/g, '')
  return `${phone.slice(0, 3)}${'*'.repeat(Math.max(4, phone.length - 5))}${phone.slice(-2)}`
})

onMounted(() => {
  if (!challenge) router.replace('/lupa-password')
  else { jeda.mulai(); nextTick(() => inputs.value[0]?.focus()) }
})

function setInputRef(el, index) { if (el) inputs.value[index] = el }
function inputDigit(index, event) {
  const value = event.target.value.replace(/\D/g, '').slice(-1)
  digits.value[index] = value
  event.target.value = value
  if (value && index < 5) inputs.value[index + 1]?.focus()
}
function keyDigit(index, event) {
  if (event.key === 'Backspace' && !digits.value[index] && index > 0) inputs.value[index - 1]?.focus()
  if (event.key === 'ArrowLeft' && index > 0) inputs.value[index - 1]?.focus()
  if (event.key === 'ArrowRight' && index < 5) inputs.value[index + 1]?.focus()
}
function pasteCode(event) {
  const value = event.clipboardData?.getData('text')?.replace(/\D/g, '').slice(0, 6) ?? ''
  if (!value) return
  event.preventDefault()
  digits.value = Array.from({ length: 6 }, (_, i) => value[i] ?? '')
  nextTick(() => inputs.value[Math.min(value.length, 5)]?.focus())
}

async function verifikasi() {
  if (!challenge || memverifikasi.value) return
  galat.value = ''
  const code = digits.value.join('')
  if (code.length !== 6) { galat.value = 'Masukkan semua 6 digit kode.'; return }
  memverifikasi.value = true
  try {
    const result = await verifyPasswordOTP({ ...challenge, code })
    simpanTokenReset(result.reset_token)
    await router.replace('/reset-password')
  } catch (error) {
    galat.value = pesanError(error)
    digits.value = ['', '', '', '', '', '']
    await nextTick(() => inputs.value[0]?.focus())
  } finally { memverifikasi.value = false }
}

async function kirimUlang() {
  if (!challenge || jeda.aktif.value || mengirim.value) return
  mengirim.value = true
  galat.value = ''
  try {
    await requestPasswordOTP(challenge)
    digits.value = ['', '', '', '', '', '']
    jeda.mulai()
    await nextTick(() => inputs.value[0]?.focus())
  } catch (error) { galat.value = pesanError(error) }
  finally { mengirim.value = false }
}
</script>

<template>
  <KartuAuth judul="Masukkan kode verifikasi" :deskripsi="`Kode 6 digit dikirim ke ${tujuanTersamar}.`">
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
    <form class="flex flex-col gap-5" @submit.prevent="verifikasi">
      <fieldset class="flex flex-col gap-2">
        <legend class="text-sm font-semibold">Kode OTP</legend>
        <div class="flex justify-between gap-2 sm:justify-start sm:gap-3" @paste="pasteCode">
          <input v-for="(_, index) in digits" :key="index" :ref="(el) => setInputRef(el, index)" :value="digits[index]" type="text" inputmode="numeric" autocomplete="one-time-code" maxlength="1" :aria-label="`Digit ${index + 1} dari 6`" class="h-12 min-w-0 flex-1 rounded-md border border-border-strong bg-surface text-center text-xl font-bold tabular-nums focus-visible:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:w-12 sm:flex-none" @input="inputDigit(index, $event)" @keydown="keyDigit(index, $event)" />
        </div>
        <p class="text-xs text-muted">Kode berlaku 5 menit dan hanya dapat dipakai sekali.</p>
      </fieldset>
      <Button type="submit" block :loading="memverifikasi">Verifikasi kode</Button>
    </form>
    <div class="flex flex-col items-center gap-2 text-sm">
      <p class="text-muted">Tidak menerima kode?</p>
      <button type="button" class="min-h-10 rounded-sm font-semibold text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:text-muted" :disabled="jeda.aktif.value || mengirim" @click="kirimUlang">{{ jeda.aktif.value ? `Kirim ulang dalam ${jeda.sisa.value} detik` : mengirim ? 'Mengirim kode…' : 'Kirim ulang kode' }}</button>
      <RouterLink to="/lupa-password" class="rounded-sm text-muted hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Ganti email atau nomor WhatsApp</RouterLink>
    </div>
  </KartuAuth>
</template>
