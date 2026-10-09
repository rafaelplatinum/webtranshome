<script setup>
import { computed, reactive, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import Select from '@/components/ui/Select.vue'
import { createAdmin } from '@/services/adminAksesService'
import { errorField, pesanError } from '@/services/errors'

const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  // meta.roles dari GET /admin/users: [{ code, name, description }]
  roles: { type: Array, default: () => [] },
})

const emit = defineEmits(['tersimpan'])

const POLA_EMAIL = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const isian = reactive({ full_name: '', email: '', role: '', password: '' })
const errors = ref({})
const pesan = ref('')
const menyimpan = ref(false)
const acak = ref('')

watch(open, (buka) => {
  if (!buka) return
  Object.assign(isian, { full_name: '', email: '', role: '', password: '' })
  errors.value = {}
  pesan.value = ''
  acak.value = ''
})

// Error satu field hilang begitu isiannya diubah (keterangan role tampil lagi).
watch(
  () => ({ ...isian }),
  (baru, lama) => {
    const sisa = { ...errors.value }
    for (const k of Object.keys(baru)) if (baru[k] !== lama[k]) delete sisa[k]
    errors.value = sisa
  },
)

const opsiRole = computed(() => props.roles.map((r) => ({ value: r.code, label: r.name })))
const keteranganRole = computed(() => props.roles.find((r) => r.code === isian.role)?.description ?? 'Izin tiap role diatur di halaman Role dan izin.')

// Huruf dan angka tanpa yang mirip (0/O, 1/l/I) supaya mudah didiktekan.
function buatAcak() {
  const huruf = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789'
  const angka = crypto.getRandomValues(new Uint32Array(12))
  isian.password = [...angka].map((n) => huruf[n % huruf.length]).join('')
  acak.value = isian.password
}

function periksa() {
  const e = {}
  if (!isian.full_name.trim()) e.full_name = 'Nama wajib diisi'
  if (!isian.email.trim()) e.email = 'Email wajib diisi'
  else if (!POLA_EMAIL.test(isian.email.trim())) e.email = 'Format email belum benar, contoh nama@transhome.id'
  if (!isian.role) e.role = 'Pilih role'
  if (isian.password.length < 8) e.password = 'Password sementara minimal 8 karakter'
  errors.value = e
  return !Object.keys(e).length
}

async function simpan() {
  if (menyimpan.value) return
  pesan.value = ''
  if (!periksa()) return
  menyimpan.value = true
  try {
    const admin = await createAdmin({ ...isian, full_name: isian.full_name.trim(), email: isian.email.trim() })
    emit('tersimpan', admin)
    open.value = false
  } catch (error) {
    errors.value = errorField(error)
    if (!Object.keys(errors.value).length) pesan.value = pesanError(error)
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <Modal v-model:open="open" title="Tambah admin" description="Admin baru langsung bisa masuk di /admin/masuk dengan password sementara ini." :dismissible="!menyimpan">
    <form id="form-tambah-admin" class="flex flex-col gap-4" novalidate @submit.prevent="simpan">
      <p v-if="pesan" role="alert" class="rounded-md bg-danger-soft px-3 py-2 text-sm font-medium text-danger-ink">{{ pesan }}</p>
      <Input v-model="isian.full_name" label="Nama" required autocomplete="off" :error="errors.full_name" />
      <Input v-model="isian.email" label="Email" type="email" required autocomplete="off" inputmode="email" :error="errors.email" />
      <Select v-model="isian.role" label="Role" required placeholder="Pilih role" :options="opsiRole" :error="errors.role" :hint="keteranganRole" />
      <div class="flex flex-col gap-2">
        <PasswordInput v-model="isian.password" label="Password sementara" required autocomplete="new-password" :error="errors.password" hint="Minimal 8 karakter. Sampaikan langsung ke admin itu, jangan lewat grup chat." />
        <div class="flex flex-wrap items-center gap-2">
          <Button variant="ghost" size="sm" @click="buatAcak">Buat password acak</Button>
          <template v-if="acak && acak === isian.password">
            <code class="rounded bg-subtle px-2 py-1 font-mono text-sm" data-password-acak>{{ acak }}</code>
            <CopyButton :text="acak" label="Salin password" />
          </template>
        </div>
      </div>
    </form>
    <template #footer>
      <Button variant="outline" :disabled="menyimpan" @click="open = false">Batal</Button>
      <Button type="submit" form="form-tambah-admin" :loading="menyimpan">Tambah admin</Button>
    </template>
  </Modal>
</template>
