<script setup>
import { nextTick, reactive, ref, useId } from 'vue'
import { BadgeCheck, Plus } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { errorField, pesanError } from '@/services/errors'
import { updateProfile } from '@/services/memberService'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'
import { formatNomorHp, normalisasiNomorHp } from '@/utils/format'
import { nomorHpValid, PESAN } from '@/utils/validasi'

// Nama dan HP bisa diubah; email hanya-baca (dipakai untuk masuk dan reset password).
const auth = useAuthStore()
const ui = useUiStore()
const idJudul = useId()

const mode = ref('lihat') // lihat | ubah
const form = reactive({ full_name: '', phone_number: '' })
const errors = ref({})
const galat = ref('')
const menyimpan = ref(false)
const formEl = ref(null)

async function ubah({ fokusHp = false } = {}) {
  form.full_name = auth.user?.full_name ?? ''
  form.phone_number = auth.user?.phone_number ? formatNomorHp(auth.user.phone_number) : ''
  errors.value = {}
  galat.value = ''
  mode.value = 'ubah'
  await nextTick()
  formEl.value?.querySelectorAll('input')[fokusHp ? 1 : 0]?.focus()
}

async function batal() {
  mode.value = 'lihat'
  await nextTick()
  document.querySelector('[data-tombol-ubah]')?.focus()
}

async function simpan() {
  if (menyimpan.value) return
  galat.value = ''
  const e = {}
  if (!form.full_name.trim()) e.full_name = PESAN.namaKosong
  if (form.phone_number.trim() && !nomorHpValid(form.phone_number)) e.phone_number = PESAN.hpFormat
  errors.value = e
  if (Object.keys(e).length) {
    await nextTick()
    formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
    return
  }
  menyimpan.value = true
  try {
    const user = await updateProfile({
      full_name: form.full_name.trim(),
      phone_number: form.phone_number.trim() ? normalisasiNomorHp(form.phone_number) : null,
    })
    auth.perbaruiUser(user)
    ui.tampilkanToast({ pesan: 'Profil diperbarui', jenis: 'success' })
    await batal()
  } catch (error) {
    const status = error?.response?.status
    if (status === 409 || status === 422) {
      errors.value = { ...errorField(error), ...(status === 409 && { phone_number: 'Nomor HP sudah dipakai akun lain' }) }
      await nextTick()
      formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
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
    <h1 class="text-2xl font-extrabold sm:text-3xl">Profil</h1>

    <section :aria-labelledby="idJudul" class="flex max-w-2xl flex-col gap-5 rounded-xl border border-border bg-surface p-5 sm:p-6">
      <div class="flex items-center justify-between gap-3">
        <h2 :id="idJudul" class="text-lg font-extrabold">Data akun</h2>
        <Button v-if="mode === 'lihat'" variant="outline" size="sm" data-tombol-ubah @click="ubah()">Ubah</Button>
      </div>

      <dl v-if="mode === 'lihat'" class="flex flex-col gap-4">
        <div class="flex flex-col gap-0.5">
          <dt class="text-xs font-semibold text-muted">Nama lengkap</dt>
          <dd class="font-semibold">{{ auth.user?.full_name }}</dd>
        </div>
        <div class="flex flex-col gap-0.5">
          <dt class="text-xs font-semibold text-muted">Email</dt>
          <dd class="flex flex-wrap items-center gap-2">
            <span class="[overflow-wrap:anywhere]">{{ auth.user?.email }}</span>
            <span v-if="auth.user?.email_verified" class="inline-flex items-center gap-1 rounded bg-success-soft px-2 py-0.5 text-xs font-semibold text-success-ink">
              <BadgeCheck class="size-3.5" aria-hidden="true" />Terverifikasi
            </span>
          </dd>
        </div>
        <div class="flex flex-col gap-0.5">
          <dt class="text-xs font-semibold text-muted">Nomor HP (opsional)</dt>
          <dd v-if="auth.user?.phone_number" class="tabular-nums" data-nomor-hp>{{ formatNomorHp(auth.user.phone_number) }}</dd>
          <dd v-else>
            <button
              type="button"
              class="inline-flex min-h-11 items-center gap-1 rounded-sm text-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-8"
              @click="ubah({ fokusHp: true })"
            >
              <Plus class="size-4" aria-hidden="true" />
              Tambah nomor HP untuk notifikasi poin via WhatsApp
            </button>
          </dd>
        </div>
        <div class="flex flex-col gap-0.5">
          <dt class="text-xs font-semibold text-muted">Akun Google</dt>
          <dd>{{ auth.user?.google_connected ? 'Terhubung' : 'Tidak terhubung' }}</dd>
        </div>
        <div class="flex flex-col gap-0.5">
          <dt class="text-xs font-semibold text-muted">Kode member</dt>
          <dd class="font-mono font-semibold">{{ auth.user?.member_code }}</dd>
        </div>
      </dl>

      <form v-else ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="simpan">
        <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
        <Input v-model="form.full_name" label="Nama lengkap" required autocomplete="name" :error="errors.full_name ?? ''" />
        <Input
          v-model="form.phone_number"
          label="Nomor HP (opsional)"
          type="tel"
          inputmode="tel"
          autocomplete="tel"
          placeholder="0812 3456 7890"
          hint="Untuk notifikasi poin lewat WhatsApp. Kosongkan bila tidak ingin."
          :error="errors.phone_number ?? ''"
        />
        <div class="flex flex-col gap-0.5">
          <span class="text-sm font-semibold">Email</span>
          <span class="[overflow-wrap:anywhere] text-muted">{{ auth.user?.email }}</span>
          <span class="text-xs text-muted">Email tidak bisa diubah sendiri. Hubungi CS bila perlu mengganti email.</span>
        </div>
        <div class="flex flex-wrap justify-end gap-2">
          <Button variant="outline" :disabled="menyimpan" @click="batal">Batal</Button>
          <Button type="submit" :loading="menyimpan">Simpan</Button>
        </div>
      </form>
    </section>
  </div>
</template>
