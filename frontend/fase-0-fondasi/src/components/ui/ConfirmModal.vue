<script setup>
import Button from './Button.vue'
import Modal from './Modal.vue'

// Konfirmasi untuk aksi destruktif atau ringkasan sebelum simpan (CLAUDE.md §11).
const open = defineModel('open', { type: Boolean, default: false })

defineProps({
  title: { type: String, required: true },
  message: { type: String, default: '' },
  confirmLabel: { type: String, default: 'Ya, lanjutkan' },
  cancelLabel: { type: String, default: 'Batal' },
  danger: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['confirm', 'cancel'])

function batal() {
  open.value = false
  emit('cancel')
}
</script>

<template>
  <!-- Tombol Batal ada lebih dulu di DOM sehingga menerima fokus pertama (pilihan aman). -->
  <Modal v-model:open="open" :title="title" :description="message" size="sm" hide-close :dismissible="!loading" @close="emit('cancel')">
    <template #footer>
      <Button variant="outline" :disabled="loading" @click="batal">{{ cancelLabel }}</Button>
      <Button :variant="danger ? 'danger' : 'primary'" :loading="loading" @click="emit('confirm')">{{ confirmLabel }}</Button>
    </template>
  </Modal>
</template>
