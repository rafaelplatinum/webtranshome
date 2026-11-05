<script setup>
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Mail, MessageCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import KartuAuth from '@/components/akun/KartuAuth.vue'
import { useJeda } from '@/composables/useJeda'
import { requestPasswordOTP } from '@/services/authService'
import { mulaiPasswordReset } from '@/services/passwordResetFlow'
import { pesanError } from '@/services/errors'
import { cekEmail, nomorHpValid, PESAN } from '@/utils/validasi'

const router = useRouter()
const jeda = useJeda(60)
const channel = ref('whatsapp')
const tujuan = ref('')
const galatTujuan = ref('')
const galat = ref('')
const mengirim = ref(false)

async function kirim() {
  galat.value = ''
  galatTujuan.value = channel.value === 'email'
    ? cekEmail(tujuan.value)
    : (!tujuan.value.trim() ? 'Nomor WhatsApp wajib diisi' : nomorHpValid(tujuan.value) ? '' : PESAN.hpFormat)
  if (galatTujuan.value) return
  mengirim.value = true
  try {
    const destination = channel.value === 'email' ? tujuan.value.trim().toLowerCase() : tujuan.value.trim()
    await requestPasswordOTP({ channel: channel.value, destination })
    mulaiPasswordReset(channel.value, destination)
    jeda.mulai()
    await router.push('/lupa-password/verifikasi')
  } catch (error) {
    galat.value = pesanError(error)
  } finally {
    mengirim.value = false
  }
}
</script>

<template>
  <KartuAuth judul="Lupa password" deskripsi="Pilih cara menerima kode untuk membuat password baru.">
    <div class="flex rounded-lg border border-border bg-muted/10 p-1" role="group" aria-label="Cara menerima kode">
      <button type="button" class="flex min-h-11 flex-1 items-center justify-center gap-2 rounded-md px-3 text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" :class="channel === 'whatsapp' ? 'bg-surface text-primary shadow-sm' : 'text-muted hover:text-foreground'" :aria-pressed="channel === 'whatsapp'" @click="channel = 'whatsapp'; tujuan = ''; galatTujuan = ''">
        <MessageCircle class="size-4" aria-hidden="true" /> WhatsApp
      </button>
      <button type="button" class="flex min-h-11 flex-1 items-center justify-center gap-2 rounded-md px-3 text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" :class="channel === 'email' ? 'bg-surface text-primary shadow-sm' : 'text-muted hover:text-foreground'" :aria-pressed="channel === 'email'" @click="channel = 'email'; tujuan = ''; galatTujuan = ''">
        <Mail class="size-4" aria-hidden="true" /> Email
      </button>
    </div>
    <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
    <form class="flex flex-col gap-4" novalidate @submit.prevent="kirim">
      <Input v-if="channel === 'whatsapp'" v-model="tujuan" label="Nomor WhatsApp" type="tel" autocomplete="tel" inputmode="tel" placeholder="0812 3456 7890" required :error="galatTujuan" />
      <Input v-else v-model="tujuan" label="Email" type="email" autocomplete="email" inputmode="email" placeholder="nama@email.com" required :error="galatTujuan" />
      <p class="text-xs leading-relaxed text-muted">Jika data cocok dengan akun aktif dan terverifikasi, kode 6 digit akan dikirim. Kode berlaku 5 menit.</p>
      <Button type="submit" block :loading="mengirim" :disabled="jeda.aktif.value">{{ jeda.aktif.value ? `Tunggu ${jeda.sisa.value} detik` : 'Kirim kode verifikasi' }}</Button>
    </form>
    <RouterLink to="/masuk" class="self-center rounded-sm text-sm font-semibold text-primary transition-colors hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Kembali ke halaman masuk</RouterLink>
  </KartuAuth>
</template>
