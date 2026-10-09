<script setup>
import { computed, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Modal from '@/components/ui/Modal.vue'
import Select from '@/components/ui/Select.vue'
import { updateAdminRoles } from '@/services/adminAksesService'
import { errorField, pesanError } from '@/services/errors'

const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  // Baris admin yang diubah: { id, full_name, email, roles: [{ code, name }] }
  admin: { type: Object, default: null },
  roles: { type: Array, default: () => [] },
})

const emit = defineEmits(['tersimpan'])

const role = ref('')
const error = ref('')
const menyimpan = ref(false)

watch(open, (buka) => {
  if (!buka) return
  role.value = props.admin?.roles[0]?.code ?? ''
  error.value = ''
})

const opsi = computed(() => props.roles.map((r) => ({ value: r.code, label: r.name })))
const keterangan = computed(() => props.roles.find((r) => r.code === role.value)?.description ?? '')
const nama = computed(() => props.admin?.full_name ?? props.admin?.email ?? '')
const berubah = computed(() => props.admin && role.value !== props.admin.roles.map((r) => r.code).join())

async function simpan() {
  if (menyimpan.value) return
  if (!role.value) {
    error.value = 'Pilih role'
    return
  }
  menyimpan.value = true
  try {
    emit('tersimpan', await updateAdminRoles(props.admin.id, [role.value]))
    open.value = false
  } catch (e) {
    error.value = errorField(e).roles ?? pesanError(e)
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <Modal v-model:open="open" :title="`Ubah role ${nama}`" description="Izin admin ini ikut role barunya setelah ia memuat ulang halaman." size="sm" :dismissible="!menyimpan">
    <form id="form-ubah-role" class="flex flex-col gap-3" novalidate @submit.prevent="simpan">
      <Select v-model="role" label="Role" required :options="opsi" :error="error" :hint="keterangan" />
    </form>
    <template #footer>
      <Button variant="outline" :disabled="menyimpan" @click="open = false">Batal</Button>
      <Button type="submit" form="form-ubah-role" :loading="menyimpan" :disabled="!berubah">Simpan role</Button>
    </template>
  </Modal>
</template>
