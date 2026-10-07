<script setup>
import { nextTick, reactive, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import { createMaster } from '@/services/adminMasterService'
import { errorField } from '@/services/errors'
import { buatSlug } from '@/utils/slug'

// "+ Tambah brand" di form produk: brand baru langsung terpilih tanpa meninggalkan form.
const open = defineModel('open', { type: Boolean, default: false })
const emit = defineEmits(['dibuat'])

const form = reactive({ name: '' })
const errors = reactive({ name: '', slug: '' })
const menyimpan = ref(false)
const formEl = ref(null)

watch(open, (buka) => {
  if (!buka) return
  form.name = ''
  errors.name = ''
  errors.slug = ''
})

async function simpan() {
  if (menyimpan.value) return
  errors.name = form.name.trim() ? '' : 'Nama brand wajib diisi'
  errors.slug = ''
  if (errors.name) return
  menyimpan.value = true
  try {
    const brand = await createMaster('brands', { name: form.name.trim(), slug: buatSlug(form.name), is_active: true })
    emit('dibuat', brand)
    open.value = false
  } catch (error) {
    const galat = errorField(error)
    // Slug dibuat dari nama, jadi bentrok slug berarti nama brand sudah ada.
    errors.name = galat.name ?? (galat.slug ? 'Brand dengan nama ini sudah ada' : '')
  } finally {
    menyimpan.value = false
    await nextTick()
    if (errors.name) formEl.value?.querySelector('input')?.focus()
  }
}
</script>

<template>
  <Modal v-model:open="open" title="Tambah brand" description="Brand baru langsung dipilih untuk produk ini." size="sm" :dismissible="!menyimpan">
    <form :id="`form-brand-cepat`" ref="formEl" novalidate @submit.prevent="simpan">
      <Input v-model="form.name" label="Nama brand" required autocomplete="off" :error="errors.name" />
    </form>
    <template #footer>
      <Button variant="outline" :disabled="menyimpan" @click="open = false">Batal</Button>
      <Button type="submit" form="form-brand-cepat" :loading="menyimpan">Simpan brand</Button>
    </template>
  </Modal>
</template>
